package quang

import (
	"fmt"
	"strconv"
)

type expression_kind_t int
type binary_operator_t int
type data_type_t int
type atom_t int8

type binary_expression_t struct {
	operator binary_operator_t

	left  *expression_t
	right *expression_t
}

type expression_t struct {
	kind       expression_kind_t
	symbolName string

	bool    bool
	float   float64
	integer int64
	atom    atom_t
	string  string
	binary  *binary_expression_t
}

type variable_t struct {
	dtype data_type_t

	bool    bool
	float   float64
	integer int64
	atom    atom_t
	string  string
}

type parser_t struct {
	current_token int
	tokens        []token_t
}

const (
	ek_nil expression_kind_t = iota
	ek_integer
	ek_float
	ek_string
	ek_atom
	ek_bool

	ek_binary

	ek_lazy_atom
	ek_lazy_symbol
)

const (
	bo_eq binary_operator_t = iota
	bo_ne
	bo_gt
	bo_lt
	bo_gte
	bo_lte
	bo_reg
	bo_and
	bo_or
)

const (
	dtype_integer data_type_t = iota
	dtype_float
	dtype_string
	dtype_atom
	dtype_bool
	dtype_nil
)

func lexerTokenKindToBinaryOperator(kind token_kind_t) binary_operator_t {
	switch kind {
	case tk_reg_keyword:
		return bo_reg
	case tk_eq_keyword:
		return bo_eq
	case tk_ne_keyword:
		return bo_ne
	case tk_gt_keyword:
		return bo_gt
	case tk_lt_keyword:
		return bo_lt
	case tk_gte_keyword:
		return bo_gte
	case tk_lte_keyword:
		return bo_lte
	case tk_and_keyword:
		return bo_and
	case tk_or_keyword:
		return bo_or
	}

	panic("unreacheable: invalid token kind")
}

func parseInteger(n string) (int64, error) {
	v, err := strconv.ParseInt(n, 10, 64)

	return v, err
}

func parseFloat(n string) (float64, error) {
	v, err := strconv.ParseFloat(n, 64)

	return v, err
}

func parseBool(n string) bool {
	switch n {
	case "true":
		return true
	case "false":
		return false
	}

	panic("unreacheable: parsing boolean")
}

func unescapeString(s string) string {
	bytes := make([]byte, 0, len(s))

	i := 0
	size := 0

	for i < len(s) {
		if s[i] == '\\' {
			i++
		}

		bytes = append(bytes, s[i])

		size++
		i++
	}

	return string(bytes[:size+1])
}

func createParser(tokens []token_t) parser_t {
	return parser_t{
		tokens:        tokens,
		current_token: 0,
	}
}

func (p parser_t) isEmpty() bool {
	return p.current_token >= len(p.tokens)
}

func (p parser_t) token() token_t {
	return p.tokens[p.current_token]
}

func (p *parser_t) forward() {
	p.current_token++
}

func (p *parser_t) parsePrimary() (*expression_t, error) {
	if p.isEmpty() {
		return nil, fmt.Errorf("error: missing token")
	}

	current := p.token()

	p.forward()

	switch current.kind {
	case tk_integer:
		{
			integer, err := parseInteger(current.value)

			if err != nil {
				return nil, fmt.Errorf("error: could not parse \"%s\" as integer due to %s", current.value, err.Error())
			}

			return &expression_t{
				kind:       ek_integer,
				symbolName: "",
				integer:    integer,
			}, nil
		}
	case tk_float:
		{
			float, err := parseFloat(current.value)

			if err != nil {
				return nil, fmt.Errorf("error: could not parse \"%s\" as float due to %s", current.value, err.Error())
			}

			return &expression_t{
				kind:       ek_float,
				symbolName: "",
				float:      float,
			}, nil
		}
	case tk_true_keyword, tk_false_keyword:
		{
			return &expression_t{
				kind:       ek_bool,
				symbolName: "",
				bool:       parseBool(current.value),
			}, nil
		}
	case tk_atom:
		{
			return &expression_t{
				kind:       ek_lazy_atom,
				symbolName: current.value,
			}, nil
		}
	case tk_symbol:
		{
			return &expression_t{
				kind:       ek_lazy_symbol,
				symbolName: current.value,
			}, nil
		}
	case tk_nil_keyword:
		{
			return &expression_t{
				kind:       ek_nil,
				symbolName: "",
			}, nil
		}
	case tk_string:
		{
			return &expression_t{
				kind:       ek_string,
				symbolName: "",
				string:     unescapeString(current.value),
			}, nil
		}
	}

	return nil, fmt.Errorf("error: unexpected token \"%s\"", current.value)
}

func (p *parser_t) parseComparison() (*expression_t, error) {
	left, err := p.parsePrimary()

	if p.isEmpty() {
		return left, err
	}

	current := p.token()

	switch current.kind {

	case tk_eq_keyword, tk_ne_keyword, tk_gt_keyword, tk_lt_keyword, tk_gte_keyword, tk_lte_keyword, tk_reg_keyword:
		{
			p.forward()

			right, err := p.parsePrimary()

			if err != nil {
				return nil, err
			}

			return &expression_t{
				kind: ek_binary,
				binary: &binary_expression_t{
					operator: lexerTokenKindToBinaryOperator(current.kind),
					left:     left,
					right:    right,
				},
			}, nil
		}
	}

	return nil, fmt.Errorf("error: expected comparison operator after expression but got \"%s\"", current.value)
}

func (p *parser_t) parseFactor() (*expression_t, error) {
	current := p.token()

	if current.kind == tk_open_paren {
		p.forward()

		expr, err := p.parseExpression()

		if err != nil {
			return nil, err
		}

		if p.token().kind != tk_close_paren {
			return nil, fmt.Errorf("error: expected ')' but got \"%s\"", p.token().value)
		}

		p.forward()

		return expr, nil
	}

	return p.parseComparison()
}

func (p *parser_t) parseTerm() (*expression_t, error) {
	left, err := p.parseFactor()

	if err != nil {
		return nil, err
	}

	for !p.isEmpty() {
		current := p.token()

		if current.kind != tk_and_keyword {
			break
		}

		p.forward()

		right, err := p.parseFactor()

		if err != nil {
			return nil, err
		}

		left = &expression_t{
			kind: ek_binary,
			binary: &binary_expression_t{
				operator: bo_and,
				left:     left,
				right:    right,
			},
		}
	}

	return left, nil
}

func (p *parser_t) parseExpression() (*expression_t, error) {
	if p.isEmpty() {
		return nil, nil
	}

	left, err := p.parseTerm()

	if err != nil {
		return nil, err
	}

	for !p.isEmpty() {
		current := p.token()

		if current.kind != tk_or_keyword {
			break
		}

		p.forward()

		right, err := p.parseTerm()

		if err != nil {
			return nil, err
		}

		left = &expression_t{
			kind: ek_binary,
			binary: &binary_expression_t{
				operator: bo_or,
				left:     left,
				right:    right,
			},
		}
	}

	return left, nil
}
