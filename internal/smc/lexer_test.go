package smc

import (
	"fmt"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestLexer(t *testing.T) {
	t.Run("Single Tokens", func(t *testing.T) {
		assertLexResult(t, "{", "OB:1/1")
		assertLexResult(t, "}", "CB:1/1")
		assertLexResult(t, ":", "C:1/1")
		assertLexResult(t, "(", "OP:1/1")
		assertLexResult(t, ")", "CP:1/1")
		assertLexResult(t, "<", "OA:1/1")
		assertLexResult(t, ">", "CA:1/1")
		assertLexResult(t, "-", "D:1/1")
		assertLexResult(t, "-", "D:1/1")
		assertLexResult(t, "name", "#name#:1/1")
		assertLexResult(t, ".", "E:1/1") // Error
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

func (c *TokenCollectorSpy) addToken(token string, line, pos int) {
	c.Result += fmt.Sprintf("%s:%d/%d", token, line, pos)
}

func (c *TokenCollectorSpy) OpenBrace(line, pos int) {
	c.addToken("OB", line, pos)
}

func (c *TokenCollectorSpy) ClosedBrace(line, pos int) {
	c.addToken("CB", line, pos)
}

func (c *TokenCollectorSpy) Colon(line, pos int) {
	c.addToken("C", line, pos)
}

func (c *TokenCollectorSpy) OpenParen(line, pos int) {
	c.addToken("OP", line, pos)
}

func (c *TokenCollectorSpy) ClosedParen(line, pos int) {
	c.addToken("CP", line, pos)
}

func (c *TokenCollectorSpy) OpenAngle(line, pos int) {
	c.addToken("OA", line, pos)
}

func (c *TokenCollectorSpy) ClosedAngle(line, pos int) {
	c.addToken("CA", line, pos)
}

func (c *TokenCollectorSpy) Dash(line, pos int) {
	c.addToken("D", line, pos)
}

func (c *TokenCollectorSpy) Name(name string, line, pos int) {
	c.addToken("#"+name+"#", line, pos)
}

func (c *TokenCollectorSpy) Error(line, pos int) {
	c.addToken("E", line, pos)
}
