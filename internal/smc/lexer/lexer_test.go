package lexer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestLexer(t *testing.T) {
	t.Run("Captures the tokens and the positions", func(t *testing.T) {
		assertLexResult(t, "{", "OB:1/1.")
		assertLexResult(t, "}", "CB:1/1.")
		assertLexResult(t, ":", "C:1/1.")
		assertLexResult(t, "(", "OP:1/1.")
		assertLexResult(t, ")", "CP:1/1.")
		assertLexResult(t, "<", "OA:1/1.")
		assertLexResult(t, ">", "CA:1/1.")
		assertLexResult(t, "-", "D:1/1.")
		assertLexResult(t, "-", "D:1/1.")
		assertLexResult(t, ".", "E:1/1.") // Error
		assertLexResult(t, "&", "E:1/1.") // Error
		assertLexResult(t, "*", "E:1/1.") // Error
		assertLexResult(t, "name", "#name#:1/1.")
		assertLexResult(t, "Name", "#Name#:1/1.")
		assertLexResult(t, "Complex_Name", "#Complex_Name#:1/1.")
		assertLexResult(t, "{}", "OB:1/1,CB:1/2.")
		assertLexResult(t, "{-}<>&:", "OB:1/1,D:1/2,CB:1/3,OA:1/4,CA:1/5,E:1/6,C:1/7.")
		assertLexResult(t, "{name}", "OB:1/1,#name#:1/2,CB:1/6.")
		assertLexResult(t, "{name}asd:fgh>", "OB:1/1,#name#:1/2,CB:1/6,#asd#:1/7,C:1/10,#fgh#:1/11,CA:1/14.")
		assertLexResult(t, "{ name }", "OB:1/1,#name#:1/3,CB:1/8.")
		assertLexResult(t, "{\n  name\n}", "OB:1/1,#name#:2/3,CB:3/1.")
		assertLexResult(t, "FSM: fsm {\n name : >asd &      \n\n  }\n", "#FSM#:1/1,C:1/4,#fsm#:1/6,OB:1/10,#name#:2/2,C:2/7,CA:2/9,#asd#:2/10,E:2/14,CB:4/3.")
		assertLexEndPosition(t, "\n\n\na:b", 5, 1)
	})
}

func assertLexResult(t *testing.T, input, expected string) {
	t.Helper()
	collector := NewTokenCollectorSpy()
	lexer := NewLexer(collector)
	lexer.Lex(bytes.NewBufferString(input))
	assert.Equal(t, expected, collector.Result)
}

func assertLexEndPosition(t *testing.T, input string, line, pos int) {
	t.Helper()
	collector := NewTokenCollectorSpy()
	lexer := NewLexer(collector)
	lexer.Lex(bytes.NewBufferString(input))
	assert.Equal(t, line, collector.EndLine)
	assert.Equal(t, pos, collector.EndPos)
}

type TokenCollectorSpy struct {
	Result  string
	EndLine int
	EndPos  int
}

func NewTokenCollectorSpy() *TokenCollectorSpy {
	return &TokenCollectorSpy{}
}

func (c *TokenCollectorSpy) addToken(token string, line, pos int) {
	if c.Result != "" {
		c.Result += ","
	}
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

func (c *TokenCollectorSpy) End(line, pos int) {
	c.EndLine = line
	c.EndPos = pos
	c.Result += "."
}
