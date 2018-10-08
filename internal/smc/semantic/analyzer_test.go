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
				assert.Equal(t, &State{Name: "a", Used: true}, findState(semanticFSM, "a"))
				assert.Empty(t, semanticFSM.Events)
				assert.Empty(t, semanticFSM.Actions)
			})

			t.Run("Empty next state", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a b - -}")
				stateA := &State{Name: "a", Used: true}
				stateA.Transitions = []Transition{
					{Event: "b", NextState: nil, Actions: []string{}},
				}
				assert.Equal(t, stateA, findState(semanticFSM, "a"))
				assert.Equal(t, []string{"b"}, semanticFSM.Events)
				assert.Empty(t, semanticFSM.Actions)
			})

			t.Run("Self referenced state", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial:a{a b a -}")
				stateA := &State{Name: "a", Used: true}
				stateA.Transitions = []Transition{
					{Event: "b", NextState: stateA, Actions: []string{}},
				}

				assert.Equal(t, stateA, findState(semanticFSM, "a"))
				assert.Equal(t, []string{"b"}, semanticFSM.Events)
				assert.Empty(t, semanticFSM.Actions)
			})

			t.Run("More than one state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a - - - b - - -}")
				assert.Equal(t, &State{Name: "a"}, findState(semanticFSM, "a"))
				assert.Equal(t, &State{Name: "b"}, findState(semanticFSM, "b"))
				assert.Empty(t, semanticFSM.Events)
				assert.Empty(t, semanticFSM.Actions)
			})

			t.Run("State with event, next state, and actions", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a b c {d e} c - - -}")
				stateC := &State{Name: "c", Used: true}
				stateA := &State{Name: "a", Transitions: []Transition{
					{Event: "b", NextState: stateC, Actions: []string{"d", "e"}},
				}}
				assert.Equal(t, stateA, findState(semanticFSM, "a"))
				assert.Equal(t, []string{"b"}, semanticFSM.Events)
				assert.Equal(t, []string{"d", "e"}, semanticFSM.Actions)
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
				assert.Equal(t, stateC, findState(semanticFSM, "c"))
				assert.Equal(t, []string{"b"}, semanticFSM.Events)
				assert.Equal(t, []string{"d"}, semanticFSM.Actions)
			})

			t.Run("Entry and exit action", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a >b <c - - -}")
				stateA := &State{
					Name:         "a",
					EntryActions: []string{"b"},
					ExitActions:  []string{"c"},
				}
				assert.Equal(t, stateA, findState(semanticFSM, "a"))
				assert.Equal(t, []string{"b", "c"}, semanticFSM.Actions)
			})

			t.Run("Two transitions for the same state", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a b - c \n a d - e}")
				stateA := &State{Name: "a", Used: true}
				stateA.Transitions = []Transition{
					{Event: "b", NextState: nil, Actions: []string{"c"}},
					{Event: "d", NextState: nil, Actions: []string{"e"}},
				}
				assert.Equal(t, stateA, findState(semanticFSM, "a"))
				assert.Equal(t, []string{"b", "d"}, semanticFSM.Events)
				assert.Equal(t, []string{"c", "e"}, semanticFSM.Actions)

				semanticFSM = analizeSemantically("{a { b - c \n d - e}}")
				assert.Equal(t, stateA, findState(semanticFSM, "a"))
				assert.Equal(t, []string{"b", "d"}, semanticFSM.Events)
				assert.Equal(t, []string{"c", "e"}, semanticFSM.Actions)
			})

			t.Run("Undefined states are not added", func(t *testing.T) {
				semanticFSM := analizeSemantically("Initial: a{b c d -}")
				assert.Len(t, semanticFSM.States, 1)
				assert.NotNil(t, findState(semanticFSM, "b"))
				assert.Equal(t, []string{"c"}, semanticFSM.Events)
				assert.Empty(t, semanticFSM.Actions)
			})

			t.Run("Duplicate events are not added", func(t *testing.T) {
				semanticFSM := analizeSemantically("{a b - - c d - - e b - -}")
				assert.Equal(t, []string{"b", "d"}, semanticFSM.Events)
			})

			t.Run("Duplicate actions are not added", func(t *testing.T) {
				semanticFSM := analizeSemantically(`
					  {
							s1 >a1 <a1 <a3 e1 - {a1 a1}
							s2 e2 - a2
						}`)

				assert.Equal(t, []string{"a1", "a3", "a2"}, semanticFSM.Actions)
			})

			t.Run("Duplicate state definitions are merged", func(t *testing.T) {
				semanticFSM := analizeSemantically(`
					  {
							s1 e1 - a1
							s1 e2 - a2
						}`)

				assert.Len(t, semanticFSM.States, 1)
				assert.Equal(t,
					&State{Name: "s1", Used: true, Transitions: []Transition{
						{Event: "e1", Actions: []string{"a1"}},
						{Event: "e2", Actions: []string{"a2"}},
					}},
					findState(semanticFSM, "s1"),
				)
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

			assertContainsError(t,
				analizeSemantically("{a { b c - \n b d - }"),
				Error{ErrorDuplicateTransition, "a:b"},
			)

			assertNotContainsError(t,
				analizeSemantically("{a { b c - \n e d - }"),
				Error{ErrorDuplicateTransition, "a:b"},
			)

			assertContainsError(t,
				analizeSemantically(`
					{
						(a) e1 c -
						(b) e1 d -
						c:a:b
					}
					`),
				Error{ErrorConflictingSuperStates, "c:e1"},
			)

			assertNotContainsError(t,
				analizeSemantically(`
					{
						(a) e1 c -
						(b) e2 d -
						c:a:b
					}
					`),
				Error{ErrorConflictingSuperStates, "c:e1"},
				Error{ErrorConflictingSuperStates, "c:e2"},
			)

			t.Run("State can be overriden", func(t *testing.T) {
				assertNotContainsError(t,
					analizeSemantically(`
						{
							(a) e1 c -
							b:a e1 d -
						}
						`),
					Error{ErrorConflictingSuperStates, "c:e1"},
				)
			})
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

		t.Run("Acceptance tests", func(t *testing.T) {
			assertValid(t, `
					Actions: Turnstile
					FSM: OneCoinTurnstile
					Initial: Locked
					{
						Locked	Coin	Unlocked	{alarmOff unlock}
						Locked 	Pass	Locked		alarmOn
						Unlocked	Coin	Unlocked	thankyou
						Unlocked	Pass	Locked		lock
					}
				`)

			assertValid(t, `
					Actions: Turnstile
					FSM: TwoCoinTurnstile
					Initial: Locked
					{
					  Locked {
					    Pass  Alarming   alarmOn
					    Coin  FirstCoin  -
					    Reset Locked     {lock alarmOff}
					  }

					  Alarming  Reset  Locked  {lock alarmOff}

					  FirstCoin {
					    Pass  Alarming  -
					    Coin  Unlocked  unlock
					    Reset Locked    {lock alarmOff}
					  }

					  Unlocked {
					    Pass  Locked  lock
					    Coin  -       thankyou
					    Reset Locked  {lock alarmOff}
					  }
					}
				`)

			assertValid(t, `
				Actions: Turnstile
				FSM: TwoCoinTurnstile
				Initial: Locked
				{
				  (Base)  Reset  Locked  {alarmOff lock}

				  Locked : Base {
				    Pass  Alarming  alarmOn
				    Coin  FirstCoin -
				  }

				  Alarming : Base  -  -  -

				  FirstCoin : Base {
				    Pass  Alarming  -
				    Coin  Unlocked  unlock
				  }

				  Unlocked : Base {
				    Pass  Locked  lock
				    Coin  -       thankyou
				  }
				}
			`)

			assertValid(t, `
				Actions: Turnstile
				FSM: TwoCoinTurnstile
				Initial: Locked
				{
				  (Base)  Reset  Locked  lock

				  Locked : Base {
				    Pass  Alarming   -
				    Coin  FirstCoin  -
				  }

				  Alarming : Base  >alarmOn <alarmOff {
				    - - -
				  }

				  FirstCoin : Base {
				    Pass  Alarming  -
				    Coin  Unlocked  unlock
				  }

				  Unlocked : Base {
				    Pass  Locked  lock
				    Coin  -       thankyou
				  }
				}
			`)
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

func assertValid(t *testing.T, input string) {
	fsm := analizeSemantically(input)
	assert.Empty(t, fsm.Errors)
	assert.Empty(t, fsm.Warnings)
}

func findState(fsm *FSM, stateName string) *State {
	for _, state := range fsm.States {
		if state.Name == stateName {
			return state
		}
	}
	return nil
}
