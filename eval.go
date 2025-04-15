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
