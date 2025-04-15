package quang

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseComparison(t *testing.T) {
	data := "size eq 0"

	l := createLexer(data)
	err := l.lex()

	assert.Nil(t, err)

	p := createParser(l.tokens)
	expr, err := p.parseComparison()

	assert.Nil(t, err)

	assert.Equal(t, ek_binary, expr.kind)
	assert.Equal(t, bo_eq, expr.binary.operator)
	assert.Equal(t, ek_lazy_symbol, expr.binary.left.kind)
	assert.Equal(t, "size", expr.binary.left.symbolName)
	assert.Equal(t, ek_integer, expr.binary.right.kind)
}

func TestParseTerm(t *testing.T) {
	data := "size gte 10.4 and method eq :post"

	l := createLexer(data)
	err := l.lex()

	assert.Nil(t, err)

	p := createParser(l.tokens)
	expr, err := p.parseTerm()

	assert.Nil(t, err)
	assert.Equal(t, ek_binary, expr.kind)
	assert.Equal(t, bo_and, expr.binary.operator)

	left := expr.binary.left

	assert.Equal(t, ek_binary, left.kind)

	assert.Equal(t, ek_lazy_symbol, left.binary.left.kind)
	assert.Equal(t, bo_gte, left.binary.operator)
	assert.Equal(t, ek_float, left.binary.right.kind)

	right := expr.binary.right

	assert.Equal(t, ek_binary, right.kind)

	assert.Equal(t, ek_lazy_symbol, right.binary.left.kind)
	assert.Equal(t, bo_eq, right.binary.operator)
	assert.Equal(t, ek_lazy_atom, right.binary.right.kind)
	assert.Equal(t, ":post", right.binary.right.symbolName)
}

func TestParseExpression(t *testing.T) {
	data := "(method eq :get and size gt 0 and size lte 1024) or (method eq :post and status ne 204) and true eq false or 10.5 lte 23.567 and name eq nil"

	l := createLexer(data)
	err := l.lex()

	assert.Nil(t, err)

	p := createParser(l.tokens)
	expr, err := p.parseExpression()

	assert.NotNil(t, expr)
	assert.Nil(t, err)

	tests := []string{
		"true or true",
		"true or false",
		"false or true",
		"false or false",
		"true and true",
		"true and false",
		"false and true",
		"false and false",
	}

	for _, test := range tests {
		l := createLexer(test)
		err := l.lex()

		assert.Nil(t, err)

		p := createParser(l.tokens)
		expr, err := p.parseExpression()

		assert.NotNil(t, expr)
		assert.Nil(t, err)
	}
}
