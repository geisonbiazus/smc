package semantic

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/stretchr/testify/assert"
)

func TestAnalyzer(t *testing.T) {
	t.Run("Header analysis", func(t *testing.T) {
		t.Run("Values", func(t *testing.T) {
			semanticFSM := analizeSemantically("Actions:a FSM:b Initial:c {}")
			assert.Equal(t, "a", semanticFSM.ActionsClass)
			assert.Equal(t, "b", semanticFSM.Name)
			assert.Equal(t, "c", semanticFSM.InitialState.Name)
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
			assertContainsErrors(t, analizeSemantically("Actions:a Actions:b {}"), ErrorDuplicateHeader)
			assertContainsErrors(t, analizeSemantically("FSM:a fsm:b {}"), ErrorDuplicateHeader)
			assertContainsErrors(t, analizeSemantically("Initial:b Initial:c {}"), ErrorDuplicateHeader)
			assertNotContainsErrors(t, analizeSemantically("Actions:a FSM:b Initial:c {}"), ErrorDuplicateHeader)
		})
	})

	t.Run("States analysis", func(t *testing.T) {
		t.Run("Values", func(t *testing.T) {
			t.Run("One state", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a - - -}")
				assert.Equal(t, createSelfReferencedState("a"), semanticFSM.States["a"])
			})

			t.Run("More than one state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a - - - b - - -}")
				assert.Equal(t, createSelfReferencedState("a"), semanticFSM.States["a"])
				assert.Equal(t, createSelfReferencedState("b"), semanticFSM.States["b"])
			})

			t.Run("State with event, next state, and actions", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a b c {d e} c - - -}")
				stateA := createSelfReferencedState("a")
				stateC := createSelfReferencedState("c")
				stateA.Transitions[0] = Transition{
					Event: "b", NextState: stateC, Actions: []string{"d", "e"},
				}
				assert.Equal(t, stateA, semanticFSM.States["a"])
			})

			t.Run("Super state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{(a) b c d \n c:a - - -}")
				stateA := createSelfReferencedState("a")
				stateC := createSelfReferencedState("c")
				stateA.Abstract = true
				stateA.Transitions[0] = Transition{
					Event: "b", NextState: stateC, Actions: []string{"d"},
				}
				stateC.SuperStates = append(stateC.SuperStates, stateA)
				assert.Equal(t, stateC, semanticFSM.States["c"])
			})
		})
	})
}

func createSelfReferencedState(name string) *State {
	state := &State{Name: name}
	state.Transitions = append(state.Transitions, Transition{NextState: state, Actions: []string{}})
	return state
}

func analizeSemantically(input string) *FSM {
	builder := parser.NewSyntaxBuilder()
	parser := parser.NewParser(builder)
	lexer := lexer.NewLexer(parser)
	lexer.Lex(bytes.NewBufferString(input))

	fsm := builder.FSM()
	analyzer := NewAnalyzer()
	return analyzer.Analyze(fsm)
}

func assertContainsErrors(t *testing.T, semanticFSM *FSM, errorTypes ...ErrorType) {
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

func assertNotContainsErrors(t *testing.T, semanticFSM *FSM, errorTypes ...ErrorType) {
	t.Helper()
	for _, errorType := range errorTypes {
		for _, e := range semanticFSM.Errors {
			if e.Type == errorType {
				t.Errorf("\n Expected: %v \n To not contain %v", semanticFSM, errorType)
			}
		}
	}
}
