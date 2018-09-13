package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestParser(t *testing.T) {
	assertParserResult(t,
		"a:b c:d {}",
		FSMSyntax{
			Headers: []Header{
				{Name: "a", Value: "b"},
				{Name: "c", Value: "d"},
			},
		})

	assertParserResult(t,
		"a:b:c: {}",
		FSMSyntax{
			Headers: []Header{
				{Name: "a", Value: "b"},
				{Name: "c"},
			},
			Errors: []SyntaxError{
				{Type: ErrorHeader, LineNumber: 1, Position: 4, Msg: "HEADER|COLON"},
			},
		})
}

func assertParserResult(t *testing.T, input string, expected FSMSyntax) {
	t.Helper()
	builder := NewSyntaxBuilder()
	parser := NewParser(builder)
	lexer := NewLexer(parser)

	lexer.Lex(bytes.NewBufferString(input))
	assert.DeepEqual(t, expected, builder.FSM())
}
