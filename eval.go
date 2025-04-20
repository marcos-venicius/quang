package quang

import (
	"fmt"
	"regexp"
)

type data_type_t int
type variable_t struct {
	dtype data_type_t

	bool    bool
	float   FloatType
	integer IntegerType
	atom    AtomType
	string  string
}

type evaluator_t struct {
	symbols    map[string]variable_t
	atoms      map[string]AtomType
	expression *expression_t
}

const (
	dtype_integer data_type_t = iota
	dtype_float
	dtype_string
	dtype_bool
	dtype_atom
	dtype_nil
)

var dtype_to_string = map[data_type_t]string{
	dtype_integer: "integer",
	dtype_float:   "float",
	dtype_string:  "string",
	dtype_bool:    "bool",
	dtype_atom:    "atom",
	dtype_nil:     "nil",
}

func createEvaluator(expression *expression_t) evaluator_t {
	return evaluator_t{
		symbols:    make(map[string]variable_t),
		atoms:      make(map[string]AtomType),
		expression: expression,
	}
}

func (e *evaluator_t) addStringVar(name, value string) {
	e.symbols[name] = variable_t{
		dtype:  dtype_string,
		string: value,
	}
}

func (e *evaluator_t) addIntegerVar(name string, value IntegerType) {
	e.symbols[name] = variable_t{
		dtype:   dtype_integer,
		integer: value,
	}
}

func (e *evaluator_t) addFloatVar(name string, value FloatType) {
	e.symbols[name] = variable_t{
		dtype: dtype_float,
		float: value,
	}
}

func (e *evaluator_t) addBoolVar(name string, value bool) {
	e.symbols[name] = variable_t{
		dtype: dtype_bool,
		bool:  value,
	}
}

func (e *evaluator_t) addAtomVar(name string, value AtomType) {
	e.symbols[name] = variable_t{
		dtype: dtype_atom,
		atom:  value,
	}
}

func (e *evaluator_t) setAtomValue(name string, value AtomType) error {
	l := createLexer(name)

	if err := l.lex(); err != nil {
		return err
	}

	if len(l.tokens) == 0 {
		return fmt.Errorf("error: missing atom name")
	}

	if l.tokens[0].kind != tk_atom {
		return fmt.Errorf("error: invalid atom name")
	}

	e.atoms[l.tokens[0].value] = value

	return nil
}

func (e *evaluator_t) lazyEvalVar(expr *expression_t) (*expression_t, error) {
	if expr.kind == ek_lazy_symbol {
		if variable, ok := e.symbols[expr.symbolName]; ok {
			switch variable.dtype {
			case dtype_string:
				return &expression_t{
					kind:   ek_string,
					string: variable.string,
				}, nil
			case dtype_integer:
				return &expression_t{
					kind:    ek_integer,
					integer: variable.integer,
				}, nil
			case dtype_float:
				return &expression_t{
					kind:  ek_float,
					float: variable.float,
				}, nil
			case dtype_bool:
				return &expression_t{
					kind: ek_bool,
					bool: variable.bool,
				}, nil
			case dtype_atom:
				return &expression_t{
					kind: ek_atom,
					atom: variable.atom,
				}, nil
			default:
				return nil, fmt.Errorf("error: could not lazy evaluate type %s", dtype_to_string[variable.dtype])
			}
		} else {
			return nil, fmt.Errorf("error: the variable '%s' does not exist", expr.symbolName)
		}
	}

	if expr.kind == ek_lazy_atom {
		if atom, ok := e.atoms[expr.symbolName]; ok {
			return &expression_t{
				kind: ek_atom,
				atom: atom,
			}, nil
		} else {
			return nil, fmt.Errorf("error: the atom '%s' does not exist", expr.symbolName)
		}
	}

	return expr, nil
}

// TODO: eval expressions between float|integer and float|integer
func (e *evaluator_t) evaluateExpression(expr *expression_t) (bool, error) {
	if expr == nil {
		return true, nil
	}

	switch expr.kind {
	case ek_binary:
		{
			op := expr.binary.operator
			left, err := e.lazyEvalVar(expr.binary.left)

			if err != nil {
				return false, err
			}

			right, err := e.lazyEvalVar(expr.binary.right)

			if err != nil {
				return false, err
			}

			switch op {
			case bo_eq, bo_ne, bo_gt, bo_lt, bo_gte, bo_lte, bo_reg:
				{
					if left.kind == ek_integer && right.kind == ek_integer {
						return cmpIntegerToInteger(left.integer, op, right.integer)
					}

					if left.kind == ek_float && right.kind == ek_float {
						return cmpFloatToFloat(left.float, op, right.float)
					}

					if left.kind == ek_string && right.kind == ek_string {
						return cmpStringToString(left.string, op, right.string)
					}

					if left.kind == ek_atom && right.kind == ek_atom {
						return cmpAtomToAtom(left.atom, op, right.atom)
					}

					return false, fmt.Errorf("error: you cannot do such operation '%s %s %s'", ek_to_string[left.kind], bo_to_string[op], ek_to_string[right.kind])
				}
			case bo_or:
				{
					leftValue, err := e.evaluateExpression(left)

					if err != nil {
						return false, err
					}

					rightValue, err := e.evaluateExpression(right)

					if err != nil {
						return false, err
					}

					return leftValue || rightValue, nil
				}
			case bo_and:
				{
					leftValue, err := e.evaluateExpression(left)

					if err != nil {
						return false, err
					}

					rightValue, err := e.evaluateExpression(right)

					if err != nil {
						return false, err
					}

					return leftValue && rightValue, nil
				}
			}
		}
	case ek_bool:
		{
			return expr.bool, nil
		}
	}

	return false, fmt.Errorf("error: could not parse expression kind %s", ek_to_string[expr.kind])
}

func (e *evaluator_t) eval() (bool, error) {
	return e.evaluateExpression(e.expression)
}

func cmpIntegerToInteger(left IntegerType, op binary_operator_t, right IntegerType) (bool, error) {
	switch op {
	case bo_eq:
		return left == right, nil
	case bo_ne:
		return left != right, nil
	case bo_gt:
		return left > right, nil
	case bo_lt:
		return left < right, nil
	case bo_gte:
		return left >= right, nil
	case bo_lte:
		return left <= right, nil
	}

	return false, fmt.Errorf("you cannot do such operation 'integer %s integer'", bo_to_string[op])
}

func cmpFloatToFloat(left FloatType, op binary_operator_t, right FloatType) (bool, error) {
	switch op {
	case bo_eq:
		return left == right, nil
	case bo_ne:
		return left != right, nil
	case bo_gt:
		return left > right, nil
	case bo_lt:
		return left < right, nil
	case bo_gte:
		return left >= right, nil
	case bo_lte:
		return left <= right, nil
	}

	return false, fmt.Errorf("you cannot do such operation 'float %s float'", bo_to_string[op])
}

func cmpStringToString(left string, op binary_operator_t, right string) (bool, error) {
	switch op {
	case bo_eq:
		return left == right, nil
	case bo_ne:
		return left != right, nil
	case bo_gt:
		return left > right, nil
	case bo_lt:
		return left < right, nil
	case bo_gte:
		return left >= right, nil
	case bo_lte:
		return left <= right, nil
	case bo_reg:
		regex := regexp.MustCompile(right)

		return regex.MatchString(left), nil
	}

	return false, fmt.Errorf("you cannot do such operation 'string %s string'", bo_to_string[op])
}

func cmpAtomToAtom(left AtomType, op binary_operator_t, right AtomType) (bool, error) {
	switch op {
	case bo_eq:
		return left == right, nil
	case bo_ne:
		return left != right, nil
	}

	return false, fmt.Errorf("you cannot do such operation 'atom %s atom'", bo_to_string[op])
}
