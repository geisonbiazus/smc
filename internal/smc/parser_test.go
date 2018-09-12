package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestParser(t *testing.T) {
	builder := NewSyntaxBuilder()
	parser := NewParser(builder)
	lexer := NewLexer(parser)

	lexer.Lex(bytes.NewBufferString("a:b{}"))
	expected := FSMSyntax{
		Headers: []Header{
			Header{Name: "a", Value: "b"},
		},
	}
	assert.DeepEqual(t, expected, builder.FSM())
}
