package smc

import (
	"bytes"
	"testing"
)

func TestSemanticAnalyzer(t *testing.T) {

	t.Run("Semantic Errors", func(t *testing.T) {
		assertContainsError(t, analizeSemantically("{}"), ErrorNoFSM)
		assertNotContainsError(t, analizeSemantically("FSM:a{}"), ErrorNoFSM)
	})
}

func analizeSemantically(input string) *SemanticFSM {
	builder := NewSyntaxBuilder()
	parser := NewParser(builder)
	lexer := NewLexer(parser)
	lexer.Lex(bytes.NewBufferString(input))

	fsm := builder.FSM()
	analyzer := NewSemanticAnalyzer()
	return analyzer.Analyze(fsm)
}

func assertContainsError(t *testing.T, semanticFSM *SemanticFSM, errorType ErrorType) {
	t.Helper()
	for _, e := range semanticFSM.Errors {
		if e.Type == errorType {
			return
		}
	}
	t.Errorf("\n Expected: %v \n To contain %v", semanticFSM, errorType)
}

func assertNotContainsError(t *testing.T, semanticFSM *SemanticFSM, errorType ErrorType) {
	t.Helper()
	for _, e := range semanticFSM.Errors {
		if e.Type == errorType {
			t.Errorf("\n Expected: %v \n To not contain %v", semanticFSM, errorType)
		}
	}
}
