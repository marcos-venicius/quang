package quang

import "fmt"

func evaluateExpression(expr *expression_t) (bool, error) {
	switch expr.kind {
	case ek_binary:
		{
			op := expr.binary.operator
			left := expr.binary.left
			right := expr.binary.right

			switch op {
			case bo_eq, bo_ne, bo_gt, bo_lt, bo_gte, bo_lte:
				{
					if left.kind == ek_integer && right.kind == ek_integer {
						return cmpIntegerToInteger(left.integer, op, right.integer)
					}

					return false, fmt.Errorf("cannot equalize types %s and %s", ek_to_string[left.kind], ek_to_string[right.kind])
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

	return false, fmt.Errorf("error: could not parse such expression")
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

	return false, fmt.Errorf("comparison operator not implemented yet")
}
