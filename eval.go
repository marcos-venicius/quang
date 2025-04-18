package quang

import (
	"fmt"
	"regexp"
)

type evaluator_t struct {
	symbols    map[string]variable_t
	expression *expression_t
}

func createEvaluator(expression *expression_t) evaluator_t {
	return evaluator_t{
		symbols:    make(map[string]variable_t),
		expression: expression,
	}
}

func (e *evaluator_t) addStringVar(name, value string) {
	e.symbols[name] = variable_t{
		dtype:  dtype_string,
		string: value,
	}
}

func (e *evaluator_t) addIntegerVar(name string, value integer_t) {
	e.symbols[name] = variable_t{
		dtype:   dtype_integer,
		integer: value,
	}
}

func (e *evaluator_t) addFloatVar(name string, value float_t) {
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

	return expr, nil
}

// TODO: eval expressions between float|integer and float|integer
func (e *evaluator_t) evaluateExpression(expr *expression_t) (bool, error) {
	if expr == nil {
		return false, nil
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

					return false, fmt.Errorf("you cannot do such operation '%s %s %s'", ek_to_string[left.kind], bo_to_string[op], ek_to_string[right.kind])
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

func cmpIntegerToInteger(left integer_t, op binary_operator_t, right integer_t) (bool, error) {
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

func cmpFloatToFloat(left float_t, op binary_operator_t, right float_t) (bool, error) {
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
