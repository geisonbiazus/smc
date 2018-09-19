package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestSemanticAnalyzer(t *testing.T) {

	t.Run("Semantic Errors", func(t *testing.T) {
		builder := NewSyntaxBuilder()
		parser := NewParser(builder)
		lexer := NewLexer(parser)
		lexer.Lex(bytes.NewBufferString("{}"))

		fsm := builder.FSM()
		analyzer := NewSemanticAnalyzer()
		semanticFSM := analyzer.Analyze(fsm)

		assert.DeepEqual(t, &SemanticFSM{}, semanticFSM)
	})

}
