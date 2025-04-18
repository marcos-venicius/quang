package quang

import (
	"fmt"
	"regexp"
)

// TODO: eval expressions between float|integer and float|integer
func evaluateExpression(expr *expression_t) (bool, error) {
	if expr == nil {
		return false, nil
	}

	switch expr.kind {
	case ek_binary:
		{
			op := expr.binary.operator
			left := expr.binary.left
			right := expr.binary.right

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
					leftValue, err := evaluateExpression(left)

					if err != nil {
						return false, err
					}

					rightValue, err := evaluateExpression(right)

					if err != nil {
						return false, err
					}

					return leftValue || rightValue, nil
				}
			case bo_and:
				{
					leftValue, err := evaluateExpression(left)

					if err != nil {
						return false, err
					}

					rightValue, err := evaluateExpression(right)

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

func cmpIntegerToInteger(left int64, op binary_operator_t, right int64) (bool, error) {
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

func cmpFloatToFloat(left float64, op binary_operator_t, right float64) (bool, error) {
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
