package smc

import (
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestLexer(t *testing.T) {
	t.Run("Single Tokens", func(t *testing.T) {
		assertLexResult(t, "{", "OB")
		assertLexResult(t, "}", "CB")
		assertLexResult(t, ":", "C")
		assertLexResult(t, "(", "OP")
		assertLexResult(t, ")", "CP")
		assertLexResult(t, "<", "OA")
		assertLexResult(t, ">", "CA")
		assertLexResult(t, "-", "D")
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

func (c *TokenCollectorSpy) ClosedBrace() {
	c.Result += "CB"
}

func (c *TokenCollectorSpy) Colon() {
	c.Result += "C"
}

func (c *TokenCollectorSpy) OpenParen() {
	c.Result += "OP"
}

func (c *TokenCollectorSpy) ClosedParen() {
	c.Result += "CP"
}

func (c *TokenCollectorSpy) OpenAngle() {
	c.Result += "OA"
}

func (c *TokenCollectorSpy) ClosedAngle() {
	c.Result += "CA"
}

func (c *TokenCollectorSpy) Dash() {
	c.Result += "D"
}
