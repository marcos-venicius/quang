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

		e := createEvaluator(expr)

		result, err := e.eval()

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

		e := createEvaluator(expr)

		result, err := e.eval()

		assert.Nil(t, err)
		assert.Equal(t, expected, result)
	}
}

func TestEvaluatingFloatExpressions(t *testing.T) {
	// TODO: add support for (<number>) syntax
	tests := map[string]bool{
		"1.5 eq 1.5":    true,
		"1.5 eq 10.309": false,
		"2. eq 1.9":     false,
		"(1. eq 1.)":    true,
		"((1. eq 1.0))": true,
		"1. ne 1.0001":  true,
		"1. ne 2.":      true,
		"2. ne 1.":      true,
		"10. gt 5.":     true,
		"10. gt 15.":    false,
		"15. gt 10.":    true,
		"10. gt 10.":    false,

		"10. lt 5.":  false,
		"10. lt 15.": true,
		"15. lt 10.": false,
		"10. lt 10.": false,

		"10. gte 10.": true,
		"10. gte 11.": false,
		"11. gte 10.": true,

		"10. lte 10.":   true,
		"10. lte 11.":   true,
		"(10. lte 11.)": true,
		"11. lte 10.":   false,
	}

	type test_case struct {
		op     binary_operator_t
		expr   string
		err    string
		result bool
	}

	fail_tests := []test_case{
		{
			op:     bo_reg,
			expr:   "10. reg 'dsflksjdf'",
			err:    "error: you cannot do such operation 'float reg string'",
			result: false,
		},
		{
			op:     bo_reg,
			expr:   "'dsflksjdf' reg 10.",
			err:    "error: you cannot do such operation 'string reg float'",
			result: false,
		},
	}

	for test, expected := range tests {
		l := createLexer(test)

		assert.Nil(t, l.lex())

		p := createParser(l.tokens)

		expr, err := p.parseExpression()

		assert.Nil(t, err)

		e := createEvaluator(expr)

		result, err := e.eval()

		assert.Nil(t, err)
		assert.Equal(t, expected, result, "test: %s", test)
	}

	for _, test := range fail_tests {
		l := createLexer(test.expr)

		assert.Nil(t, l.lex())

		p := createParser(l.tokens)

		expr, err := p.parseExpression()

		assert.Nil(t, err)

		e := createEvaluator(expr)

		result, err := e.eval()

		assert.NotNil(t, err)
		assert.Equal(t, test.err, err.Error())
		assert.Equal(t, test.result, result, "test: %s", test)
	}
}

func TestEvaluatingStringExpressions(t *testing.T) {
	tests := map[string]bool{
		"'hello world' eq 'hello world'":       true,
		"'hello world ' eq 'hello world'":      false,
		"'hello world ' ne 'hello world'":      true,
		"'hello world' ne 'hello world'":       false,
		"'hello \\'world' eq 'hello \\'world'": true,
		"'z' gt 'a'":                           true,
		"'a' lt 'z'":                           true,
		"'a' eq 'a'":                           true,
		"'/test/3e7f0bb3-d315-46ec-a92f-9bd694e5e281/fake' reg '^/test/[a-z0-9]{8}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{4}-[a-z0-9]{12}/fake$'": true,
	}

	for test, expected := range tests {
		l := createLexer(test)

		assert.Nil(t, l.lex())

		p := createParser(l.tokens)

		expr, err := p.parseExpression()

		assert.Nil(t, err)

		e := createEvaluator(expr)

		result, err := e.eval()

		assert.Nil(t, err)
		assert.Equal(t, expected, result, "test: %s", test)
	}
}

func TestEvaluatingLazySymbols(t *testing.T) {
	test := "size gt 40"

	l := createLexer(test)

	assert.Nil(t, l.lex())

	p := createParser(l.tokens)

	expr, err := p.parseExpression()

	assert.Nil(t, err)

	e := createEvaluator(expr)

	result, err := e.eval()

	assert.NotNil(t, err)
	assert.Equal(t, false, result)
	assert.Equal(t, "error: the variable 'size' does not exist", err.Error())

	e.addStringVar("size", "anything")

	result, err = e.eval()

	assert.NotNil(t, err)
	assert.Equal(t, false, result)
	assert.Equal(t, "error: you cannot do such operation 'string gt integer'", err.Error())

	e.addIntegerVar("size", 41)

	result, err = e.eval()

	assert.Nil(t, err)
	assert.Equal(t, true, result)

	e.addIntegerVar("size", 38)

	result, err = e.eval()

	assert.Nil(t, err)
	assert.Equal(t, false, result)
}

func TestEvaluatingLazySymbolsStrings(t *testing.T) {
	test := "agent reg this"

	l := createLexer(test)

	assert.Nil(t, l.lex())

	p := createParser(l.tokens)

	expr, err := p.parseExpression()

	assert.Nil(t, err)

	e := createEvaluator(expr)

	e.addStringVar("agent", "hello world")
	e.addStringVar("this", "^\\w+\\s\\w+$")

	result, err := e.eval()

	assert.Nil(t, err)
	assert.Equal(t, true, result)
}

func TestEvaluatingLazyAtoms(t *testing.T) {
	test := "method eq :get"

	l := createLexer(test)

	assert.Nil(t, l.lex())

	p := createParser(l.tokens)

	expr, err := p.parseExpression()

	assert.Nil(t, err)

	e := createEvaluator(expr)

	result, err := e.eval()

	assert.NotNil(t, err)
	assert.Equal(t, false, result)
	assert.Equal(t, "error: the variable 'method' does not exist", err.Error())

	e.addAtomVar("method", 0)

	result, err = e.eval()

	assert.NotNil(t, err)
	assert.Equal(t, false, result)
	assert.Equal(t, "error: the atom ':get' does not exist", err.Error())

	e.setAtomValue(":get", 1)

	result, err = e.eval()

	assert.Nil(t, err)
	assert.Equal(t, false, result)

	e.setAtomValue(":get", 0)

	result, err = e.eval()

	assert.Nil(t, err)
	assert.Equal(t, true, result)

	e.addStringVar("method", "get")

	result, err = e.eval()

	assert.NotNil(t, err)
	assert.Equal(t, false, result)
	assert.Equal(t, "error: you cannot do such operation 'string eq atom'", err.Error())
}

func TestEvaluatingLazyAtomsNeOperator(t *testing.T) {
	test := "method ne :get"

	l := createLexer(test)

	assert.Nil(t, l.lex())

	p := createParser(l.tokens)

	expr, err := p.parseExpression()

	assert.Nil(t, err)

	e := createEvaluator(expr)

	e.addAtomVar("method", 0)
	e.setAtomValue(":get", 0)

	result, err := e.eval()

	assert.Nil(t, err)
	assert.Equal(t, false, result)

	e.setAtomValue(":get", 1)

	result, err = e.eval()

	assert.Nil(t, err)
	assert.Equal(t, true, result)
}
