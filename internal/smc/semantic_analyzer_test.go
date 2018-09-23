package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/testing/assert"
)

func TestSemanticAnalyzer(t *testing.T) {
	t.Run("Header analysis", func(t *testing.T) {
		t.Run("Values", func(t *testing.T) {
			semanticFSM := analizeSemantically("Actions:a FSM:b Initial:c {}")
			assert.Equal(t, "a", semanticFSM.Actions)
			assert.Equal(t, "b", semanticFSM.Name)
			assert.Equal(t, "c", semanticFSM.Initial)
		})

		t.Run("Errors", func(t *testing.T) {
			assertContainsErrors(t, analizeSemantically("{}"), ErrorNoFSM)
			assertNotContainsErrors(t, analizeSemantically("FSM:a{}"), ErrorNoFSM)
			assertContainsErrors(t, analizeSemantically("{}"), ErrorNoInitial)
			assertNotContainsErrors(t, analizeSemantically("Initial:a{}"), ErrorNoInitial)
			assertContainsErrors(t, analizeSemantically("Actions:a {}"), ErrorNoFSM, ErrorNoInitial)
			assertContainsErrors(t, analizeSemantically("a:b {}"), ErrorInvalidHeader)
			assertNotContainsErrors(t, analizeSemantically("Actions:a FSM:b Initial:c {}"), ErrorNoFSM, ErrorNoInitial, ErrorInvalidHeader)
			assertNotContainsErrors(t, analizeSemantically("actions:a fsm:b initial:c {}"), ErrorNoFSM, ErrorNoInitial, ErrorInvalidHeader)
			assertNotContainsErrors(t, analizeSemantically("FSM:b Initial:c {}"), ErrorNoFSM, ErrorNoInitial, ErrorInvalidHeader)
			assertContainsErrors(t, analizeSemantically("FSM:a FSM:b {}"), ErrorDuplicateHeader)
			assertContainsErrors(t, analizeSemantically("Actions:a Actions:b {}"), ErrorDuplicateHeader)
			assertContainsErrors(t, analizeSemantically("Initial:a Initial:b {}"), ErrorDuplicateHeader)
		})
	})

	t.Run("Logic analysis", func(t *testing.T) {
		// semanticFSM := analizeSemantically("{}")
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

func assertContainsErrors(t *testing.T, semanticFSM *SemanticFSM, errorTypes ...ErrorType) {
	t.Helper()
	for _, errorType := range errorTypes {
		found := false
		for _, e := range semanticFSM.Errors {
			if e.Type == errorType {
				found = true
			}
		}
		if !found {
			t.Errorf("\n Expected: %v \n To contain %v", semanticFSM, errorType)
		}
	}
}

func assertNotContainsErrors(t *testing.T, semanticFSM *SemanticFSM, errorTypes ...ErrorType) {
	t.Helper()
	for _, errorType := range errorTypes {
		for _, e := range semanticFSM.Errors {
			if e.Type == errorType {
				t.Errorf("\n Expected: %v \n To not contain %v", semanticFSM, errorType)
			}
		}
	}
}
