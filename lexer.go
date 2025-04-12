package quang

import "fmt"

type token_kind_t int

type token_t struct {
	value string
	kind  token_kind_t
}

type lexer_t struct {
	content     string
	cursor, bot int
	tokens      []token_t
}

const (
	tk_open_paren token_kind_t = iota
	tk_close_paren
	tk_and_keyword
	tk_or_keyword
	tk_true_keyword
	tk_false_keyword
	tk_nil_keyword
	tk_eq_keyword
	tk_ne_keyword
	tk_gt_keyword
	tk_lt_keyword
	tk_gte_keyword
	tk_lte_keyword
	tk_reg_keyword
	tk_symbol
	tk_integer
	tk_atom
	tk_string
	/* tk_float   token_kind_t = iota */
)

var keywords = map[string]token_kind_t{
	"true":  tk_true_keyword,
	"false": tk_false_keyword,
	"nil":   tk_nil_keyword,
	"and":   tk_and_keyword,
	"or":    tk_or_keyword,
	"reg":   tk_reg_keyword,
	"eq":    tk_eq_keyword,
	"ne":    tk_ne_keyword,
	"gt":    tk_gt_keyword,
	"lt":    tk_lt_keyword,
	"gte":   tk_gte_keyword,
	"lte":   tk_lte_keyword,
}

func createLexer(content string) lexer_t {
	return lexer_t{
		content: content,
		cursor:  0,
		bot:     0,
		tokens:  make([]token_t, 0, 512),
	}
}

func (l lexer_t) isEmpty() bool {
	return l.cursor >= len(l.content)
}

func (l lexer_t) isEmptyAhead() bool {
	return l.cursor+1 >= len(l.content)
}

func (l lexer_t) char() byte {
	return l.content[l.cursor]
}

func (l lexer_t) charAhead() byte {
	return l.content[l.cursor+1]
}

func (l *lexer_t) forward() {
	l.cursor++
}

func (l *lexer_t) trimWhitespaces() {
	for !l.isEmpty() && l.content[l.cursor] == ' ' {
		l.forward()
	}
}

func (l *lexer_t) lexSingleChar(kind token_kind_t) {
	l.forward()

	token := token_t{
		value: l.content[l.bot:l.cursor],
		kind:  kind,
	}

	l.tokens = append(l.tokens, token)
}

func (l *lexer_t) lexNumber() {
	for !l.isEmpty() && isDigit(l.char()) {
		l.forward()
	}

	// TODO: identify floats

	token := token_t{
		kind:  tk_integer,
		value: l.content[l.bot:l.cursor],
	}

	l.tokens = append(l.tokens, token)
}

func (l *lexer_t) lexSymbolOrKeyword() {
	for !l.isEmpty() && isSymbol(l.char()) {
		l.forward()
	}

	token := token_t{
		kind:  tk_symbol,
		value: l.content[l.bot:l.cursor],
	}

	if kind, ok := keywords[token.value]; ok {
		token.kind = kind
	}

	l.tokens = append(l.tokens, token)
}

func (l *lexer_t) lexAtom() error {
	l.forward()

	atomNameSize := 0

	for !l.isEmpty() && isSymbol(l.char()) {
		l.forward()
		atomNameSize++
	}

	if atomNameSize == 0 {
		return fmt.Errorf("error: missing atom name at position %d", l.cursor)
	}

	token := token_t{
		kind:  tk_atom,
		value: l.content[l.bot:l.cursor],
	}

	l.tokens = append(l.tokens, token)

	return nil
}

func (l *lexer_t) lexString() error {
	l.forward()

	for !l.isEmpty() && l.char() != '\'' {
		if l.char() == '\\' {
			if l.isEmptyAhead() {
				return fmt.Errorf("error: unterminated string literal at position %d", l.bot+1)
			}

			switch l.charAhead() {
			case '\'':
				l.forward()
			default:
				return fmt.Errorf("error: invalid scape sequence at position %d", l.cursor+1)
			}
		}

		l.forward()
	}

	if l.isEmpty() {
		return fmt.Errorf("error: unterminated string literal at position %d", l.bot+1)
	}

	token := token_t{
		kind:  tk_string,
		value: l.content[l.bot+1 : l.cursor],
	}

	l.tokens = append(l.tokens, token)

	l.forward()

	return nil
}

func (l *lexer_t) lex() error {
	for l.cursor < len(l.content) {
		l.trimWhitespaces()

		if l.isEmpty() {
			break
		}

		l.bot = l.cursor

		char := l.char()

		switch char {
		case '\'':
			if err := l.lexString(); err != nil {
				return err
			}
		case '(':
			l.lexSingleChar(tk_open_paren)
		case ')':
			l.lexSingleChar(tk_close_paren)
		case ':':
			if err := l.lexAtom(); err != nil {
				return err
			}
		default:
			if isDigit(char) {
				l.lexNumber()
			} else if isSymbol(char) {
				l.lexSymbolOrKeyword()
			} else {
				return fmt.Errorf("error: unexpected character \"%c\" at position %d", char, l.cursor+1)
			}
		}
	}

	return nil
}
