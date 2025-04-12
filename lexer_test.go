package quang

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexingParenthesis(t *testing.T) {
	l := createLexer("())(")

	err := l.lex()

	assert.Nil(t, err)
	assert.Equal(t, 4, len(l.tokens))
	assert.Equal(t, tk_open_paren, l.tokens[0].kind)
	assert.Equal(t, tk_close_paren, l.tokens[1].kind)
	assert.Equal(t, tk_close_paren, l.tokens[2].kind)
	assert.Equal(t, tk_open_paren, l.tokens[3].kind)
	assert.Equal(t, "(", l.tokens[0].value)
	assert.Equal(t, ")", l.tokens[1].value)
	assert.Equal(t, ")", l.tokens[2].value)
	assert.Equal(t, "(", l.tokens[3].value)
}

func TestLexingNumbers(t *testing.T) {
	l := createLexer("0 124")

	err := l.lex()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(l.tokens))
	assert.Equal(t, tk_integer, l.tokens[0].kind)
	assert.Equal(t, tk_integer, l.tokens[1].kind)
	assert.Equal(t, "0", l.tokens[0].value)
	assert.Equal(t, "124", l.tokens[1].value)
}

func TestLexingKeywords(t *testing.T) {
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

	keys := make([]string, 0, len(keywords))

	for key := range keywords {
		keys = append(keys, key)
	}

	l := createLexer(strings.Join(keys, " "))

	err := l.lex()

	assert.Nil(t, err)
	assert.Equal(t, len(keywords), len(l.tokens))

	for _, token := range l.tokens {
		keyword := keywords[token.value]

		assert.Equal(t, keyword, token.kind)
	}
}

func TestLexingSymbols(t *testing.T) {
	l := createLexer("hello guys")

	err := l.lex()

	assert.Nil(t, err)
	assert.Equal(t, 2, len(l.tokens))

	assert.Equal(t, "hello", l.tokens[0].value)
	assert.Equal(t, "guys", l.tokens[1].value)

	assert.Equal(t, tk_symbol, l.tokens[0].kind)
	assert.Equal(t, tk_symbol, l.tokens[1].kind)
}

func TestLexingAtoms(t *testing.T) {
	l := createLexer(":hello_world :_ :h :hi")

	err := l.lex()

	assert.Nil(t, err)
	assert.Equal(t, 4, len(l.tokens))

	assert.Equal(t, ":hello_world", l.tokens[0].value)
	assert.Equal(t, ":_", l.tokens[1].value)
	assert.Equal(t, ":h", l.tokens[2].value)
	assert.Equal(t, ":hi", l.tokens[3].value)

	assert.Equal(t, tk_atom, l.tokens[0].kind)
	assert.Equal(t, tk_atom, l.tokens[1].kind)
	assert.Equal(t, tk_atom, l.tokens[2].kind)
	assert.Equal(t, tk_atom, l.tokens[3].kind)

	l = createLexer(":")

	err = l.lex()

	assert.NotNil(t, err)
	assert.Equal(t, "error: missing atom name at position 1", err.Error())
}

func TestLexingString(t *testing.T) {
	l := createLexer("'Hello World'")

	err := l.lex()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(l.tokens))
	assert.Equal(t, tk_string, l.tokens[0].kind)
	assert.Equal(t, "Hello World", l.tokens[0].value)

	l = createLexer("'Hello \\'World\\''")

	err = l.lex()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(l.tokens))
	assert.Equal(t, tk_string, l.tokens[0].kind)
	assert.Equal(t, "Hello \\'World\\'", l.tokens[0].value)

	l = createLexer("'Hello \\")

	err = l.lex()

	assert.NotNil(t, err)
	assert.Equal(t, "error: unterminated string literal at position 1", err.Error())

	l = createLexer("'Hello \\'")

	err = l.lex()

	assert.NotNil(t, err)
	assert.Equal(t, "error: unterminated string literal at position 1", err.Error())

	l = createLexer("   'Hello s    sdflkjsdf sdlkjsdf\\'sdflksdfj")

	err = l.lex()

	assert.NotNil(t, err)
	assert.Equal(t, "error: unterminated string literal at position 4", err.Error())
}
