package quang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluatingBooleanExpressions(t *testing.T) {
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

func TestEvaluatingIntegerExpressions(t *testing.T) {
	// TODO: add support for (<number>) syntax
	tests := map[string]bool{
		"1 eq 1":     true,
		"1 eq 10":    false,
		"2 eq 1":     false,
		"(1 eq 1)":   true,
		"((1 eq 1))": true,
		"1 ne 1":     false,
		"1 ne 2":     true,
		"2 ne 1":     true,
		"10 gt 5":    true,
		"10 gt 15":   false,
		"15 gt 10":   true,
		"10 gt 10":   false,

		"10 lt 5":  false,
		"10 lt 15": true,
		"15 lt 10": false,
		"10 lt 10": false,

		"10 gte 10": true,
		"10 gte 11": false,
		"11 gte 10": true,

		"10 lte 10": true,
		"10 lte 11": true,
		"11 lte 10": false,
		"(true and false) or (1 gte 0 or 10 lte 5)": true,
	}
	// TODO: operator bellow should complaing about types
	/* "reg"
	"and"
	"or" */

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
