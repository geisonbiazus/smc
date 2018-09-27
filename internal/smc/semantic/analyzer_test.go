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
			t.Run("One state without event", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a - - -}")
				assert.Equal(t, &State{Name: "a"}, semanticFSM.States["a"])
			})

			t.Run("Self referenced state", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a b - -}")
				stateA := &State{Name: "a"}
				stateA.Transitions = []Transition{
					{Event: "b", NextState: stateA, Actions: []string{}},
				}
				assert.Equal(t, stateA, semanticFSM.States["a"])

				semanticFSM = analizeSemantically("Initial:a{a b a -}")

				assert.Equal(t, stateA, semanticFSM.States["a"])
			})

			t.Run("More than one state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a - - - b - - -}")
				assert.Equal(t, &State{Name: "a"}, semanticFSM.States["a"])
				assert.Equal(t, &State{Name: "b"}, semanticFSM.States["b"])
			})

			t.Run("State with event, next state, and actions", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a b c {d e} c - - -}")
				stateC := &State{Name: "c"}
				stateA := &State{Name: "a", Transitions: []Transition{
					{Event: "b", NextState: stateC, Actions: []string{"d", "e"}},
				}}
				assert.Equal(t, stateA, semanticFSM.States["a"])
			})

			t.Run("Super state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{(a) b c d \n c:a - - -}")
				stateC := &State{Name: "c"}
				stateA := &State{
					Name:     "a",
					Abstract: true,
					Transitions: []Transition{
						{Event: "b", NextState: stateC, Actions: []string{"d"}},
					},
				}
				stateC.SuperStates = append(stateC.SuperStates, stateA)
				assert.Equal(t, stateC, semanticFSM.States["c"])
			})

			t.Run("Entry and exit action", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a >b <c - - -}")
				stateA := &State{
					Name:         "a",
					EntryActions: []string{"b"},
					ExitActions:  []string{"c"},
				}
				assert.Equal(t, stateA, semanticFSM.States["a"])
			})

			t.Run("Two transitions for the same state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a b - c \n a d - e}")
				stateA := &State{Name: "a"}
				stateA.Transitions = []Transition{
					{Event: "b", NextState: stateA, Actions: []string{"c"}},
					{Event: "d", NextState: stateA, Actions: []string{"e"}},
				}
				assert.Equal(t, stateA, semanticFSM.States["a"])

				semanticFSM = analizeSemantically("{a { b - c \n d - e}}")
				assert.Equal(t, stateA, semanticFSM.States["a"])
			})
		})

		t.Run("Errors", func(t *testing.T) {
			assertContainsError(t, analizeSemantically("Initial: a{}"), ErrorUndefinedState, "a")
		})
	})
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

func assertContainsError(t *testing.T, semanticFSM *FSM, errorType ErrorType, element string) {
	t.Helper()
	for _, e := range semanticFSM.Errors {
		if e.Type == errorType && e.Element == element {
			return
		}
	}
	t.Errorf("\n Expected: %v \n To contain %v", semanticFSM, errorType)
}

func assertNotContainsError(t *testing.T, semanticFSM *FSM, errorType ErrorType, element string) {
	t.Helper()
	for _, e := range semanticFSM.Errors {
		if e.Type == errorType && e.Element == element {
			t.Errorf("\n Expected: %v \n To not contain %v", semanticFSM, errorType)
		}
	}
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
