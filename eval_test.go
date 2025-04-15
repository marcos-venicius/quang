package quang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluate(t *testing.T) {
	tests := map[string]bool{
		"true or false":       true,
		"false or true":       true,
		"true or true":        true,
		"false or false":      false,
		"true and false":      false,
		"false and true":      false,
		"false and false":     false,
		"true and true":       true,
		"true":                true,
		"false":               false,
		"(true)":              true,
		"(true) or (true)":    true,
		"(true) or true":      true,
		"true or (true)":      true,
		"(false) or (false)":  false,
		"(true) and (true)":   true,
		"(false) and (false)": false,
		"(false and true) or ((true and false) or true)": true,
		"true and false or true":                         true,
	}

	for test, expected := range tests {
		l := createLexer(test)

		assert.Nil(t, l.lex())

		p := createParser(l.tokens)

		expr, err := p.parseExpression()

		assert.Nil(t, err)

		result, err := evaluateExpression(expr)

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	}
}
