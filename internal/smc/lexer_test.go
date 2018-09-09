package smc

import (
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestLexer(t *testing.T) {
	t.Run("Single Tokens", func(t *testing.T) {
		assertLexResult(t, "{", "OB")
		assertLexResult(t, "}", "CB")
	})
}

func assertLexResult(t *testing.T, input, expected string) {
	t.Helper()
	collector := NewTokenCollectorSpy()
	lexer := NewLexer(collector)
	lexer.Lex(input)
	assert.Equal(t, expected, collector.Result)
}

type TokenCollectorSpy struct {
	Result string
}

func NewTokenCollectorSpy() *TokenCollectorSpy {
	return &TokenCollectorSpy{}
}

func (c *TokenCollectorSpy) OpenBrace() {
	c.Result += "OB"
}

func (c *TokenCollectorSpy) CloseBrace() {
	c.Result += "CB"
}
