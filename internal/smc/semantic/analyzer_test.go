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
			assertContainsError(t,
				analizeSemantically("{}"),
				Error{ErrorNoFSM, "FSM"},
			)

			assertNotContainsError(t,
				analizeSemantically("FSM:a{}"),
				Error{ErrorNoFSM, "FSM"},
			)

			assertContainsError(t,
				analizeSemantically("{}"),
				Error{ErrorNoInitial, "Initial"},
			)

			assertNotContainsError(t,
				analizeSemantically("Initial:a{}"),
				Error{ErrorNoFSM, "Initial"},
			)

			assertContainsError(t,
				analizeSemantically("Actions:a {}"),
				Error{ErrorNoFSM, "FSM"}, Error{ErrorNoInitial, "Initial"},
			)

			assertContainsError(t,
				analizeSemantically("a:b {}"),
				Error{ErrorInvalidHeader, "a"},
			)

			assertNotContainsError(t,
				analizeSemantically("Actions:a FSM:b Initial:c {}"),
				Error{ErrorNoFSM, "FSM"},
				Error{ErrorNoInitial, "Initial"},
				Error{ErrorInvalidHeader, "Actions"},
			)

			assertNotContainsError(t,
				analizeSemantically("actions:a fsm:b initial:c {}"),
				Error{ErrorNoFSM, "FSM"},
				Error{ErrorNoInitial, "Initial"},
				Error{ErrorInvalidHeader, "actions"},
			)

			assertContainsError(t,
				analizeSemantically("Actions:a Actions:b {}"),
				Error{ErrorDuplicateHeader, "Actions"},
			)

			assertContainsError(t,
				analizeSemantically("FSM:a fsm:b {}"),
				Error{ErrorDuplicateHeader, "FSM"},
			)

			assertContainsError(t,
				analizeSemantically("Initial:b Initial:c {}"),
				Error{ErrorDuplicateHeader, "Initial"},
			)

			assertNotContainsError(t,
				analizeSemantically("Actions:a FSM:b Initial:c {}"),
				Error{ErrorDuplicateHeader, "FSM"},
				Error{ErrorDuplicateHeader, "Initial"},
				Error{ErrorDuplicateHeader, "Actions"},
			)
		})
	})

	t.Run("States analysis", func(t *testing.T) {
		t.Run("Values", func(t *testing.T) {
			t.Run("One state without event", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a - - -}")
				assert.Equal(t, &State{Name: "a", Used: true}, semanticFSM.States["a"])
			})

			t.Run("Self referenced state", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a b - -}")
				stateA := &State{Name: "a", Used: true}
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
				stateC := &State{Name: "c", Used: true}
				stateA := &State{Name: "a", Transitions: []Transition{
					{Event: "b", NextState: stateC, Actions: []string{"d", "e"}},
				}}
				assert.Equal(t, stateA, semanticFSM.States["a"])
			})

			t.Run("Super state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{(a) b c d \n c:a - - -}")
				stateC := &State{Name: "c", Used: true}
				stateA := &State{
					Name:     "a",
					Abstract: true,
					Used:     true,
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
				stateA := &State{Name: "a", Used: true}
				stateA.Transitions = []Transition{
					{Event: "b", NextState: stateA, Actions: []string{"c"}},
					{Event: "d", NextState: stateA, Actions: []string{"e"}},
				}
				assert.Equal(t, stateA, semanticFSM.States["a"])

				semanticFSM = analizeSemantically("{a { b - c \n d - e}}")
				assert.Equal(t, stateA, semanticFSM.States["a"])
			})

			t.Run("Undefined states are not added", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial: a{b c d -}")
				assert.Len(t, semanticFSM.States, 1)
				assert.Contains(t, semanticFSM.States, "b")
			})
		})

		t.Run("Errors", func(t *testing.T) {
			assertContainsError(t,
				analizeSemantically("Initial: a{}"),
				Error{ErrorUndefinedState, "a"},
			)

			assertNotContainsError(t,
				analizeSemantically("Initial: a{a - - -}"),
				Error{ErrorUndefinedState, "a"},
			)

			assertContainsError(t,
				analizeSemantically("{a b c -}"),
				Error{ErrorUndefinedState, "c"},
			)

			assertNotContainsError(t,
				analizeSemantically("{a b c - c - - -}"),
				Error{ErrorUndefinedState, "c"},
			)

			assertContainsError(t,
				analizeSemantically("{a >b - - - \n a >c - - -}"),
				Error{ErrorEntryActionsAlreadyDefined, "a"},
			)

			assertNotContainsError(t,
				analizeSemantically("{a >b - - - \n b >c - - -}"),
				Error{ErrorEntryActionsAlreadyDefined, "a"},
			)

			assertContainsError(t,
				analizeSemantically("{a <b - - - \n a <c - - -}"),
				Error{ErrorExitActionsAlreadyDefined, "a"},
			)

			assertNotContainsError(t,
				analizeSemantically("{a <b - - - \n b <c - - -}"),
				Error{ErrorExitActionsAlreadyDefined, "a"},
			)

			assertContainsError(t,
				analizeSemantically("{(a) - - - \n a - - -}"),
				Error{ErrorAbstractStateRedefinedAsNonAbstract, "a"},
			)

			assertContainsError(t,
				analizeSemantically("{a - - - \n (a) - - -}"),
				Error{ErrorAbstractStateRedefinedAsNonAbstract, "a"},
			)

			assertNotContainsError(t,
				analizeSemantically("{(a) - - - \n (a) - - -}"),
				Error{ErrorAbstractStateRedefinedAsNonAbstract, "a"},
			)

			assertContainsError(t,
				analizeSemantically("{a:b - - -}"),
				Error{ErrorUndefinedSuperState, "b"},
			)

			assertNotContainsError(t,
				analizeSemantically("{a:b - - - b - - - }"),
				Error{ErrorUndefinedSuperState, "b"},
			)

			assertContainsError(t,
				analizeSemantically("{a b c - (c) - - -}"),
				Error{ErrorAbstractStateUsedAsNextState, "c"},
			)

			assertNotContainsError(t,
				analizeSemantically("{a b c - c - - -}"),
				Error{ErrorAbstractStateUsedAsNextState, "c"},
			)
		})

		t.Run("Warnings", func(t *testing.T) {
			assertContainsWarning(t,
				analizeSemantically("{a b c d}"),
				Error{ErrorUnusedState, "a"},
			)

			assertNotContainsWarning(t,
				analizeSemantically("{a b a d}"),
				Error{ErrorUnusedState, "a"},
			)

			assertNotContainsWarning(t,
				analizeSemantically("{a b - d}"),
				Error{ErrorUnusedState, "a"},
			)

			assertNotContainsWarning(t,
				analizeSemantically("Initial: a{a b c d}"),
				Error{ErrorUnusedState, "a"},
			)

			assertNotContainsWarning(t,
				analizeSemantically("{a b c d c:a - - -}"),
				Error{ErrorUnusedState, "a"},
			)
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

func assertContainsError(t *testing.T, semanticFSM *FSM, errors ...Error) {
	t.Helper()
	for _, err := range errors {
		assert.Contains(t, semanticFSM.Errors, err)
	}
}

func assertNotContainsError(t *testing.T, semanticFSM *FSM, errors ...Error) {
	t.Helper()
	for _, err := range errors {
		assert.NotContains(t, semanticFSM.Errors, err)
	}
}

func assertContainsWarning(t *testing.T, semanticFSM *FSM, errors ...Error) {
	t.Helper()
	for _, err := range errors {
		assert.Contains(t, semanticFSM.Warnings, err)
	}
}

func assertNotContainsWarning(t *testing.T, semanticFSM *FSM, errors ...Error) {
	t.Helper()
	for _, err := range errors {
		assert.NotContains(t, semanticFSM.Warnings, err)
	}
}
