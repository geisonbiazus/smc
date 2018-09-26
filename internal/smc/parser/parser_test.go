package parser

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Run("Incremental Tests", func(t *testing.T) {
		assertParserResult(t,
			"a:b c:d {}",
			FSMSyntax{
				Headers: []Header{
					{Name: "a", Value: "b"},
					{Name: "c", Value: "d"},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b{c d e f}",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{StateSpec{Name: "c"}, []SubTransition{{"d", "e", []string{"f"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b{c d e {f g} \n h i j k}",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{StateSpec{Name: "c"}, []SubTransition{{"d", "e", []string{"f", "g"}}}},
					{StateSpec{Name: "h"}, []SubTransition{{"i", "j", []string{"k"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b { c d e - \n f g h i }",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{StateSpec{Name: "c"}, []SubTransition{{"d", "e", []string{}}}},
					{StateSpec{Name: "f"}, []SubTransition{{"g", "h", []string{"i"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b { c d - e }",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{StateSpec{Name: "c"}, []SubTransition{{"d", "", []string{"e"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b { c - d e }",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{StateSpec{Name: "c"}, []SubTransition{{"", "d", []string{"e"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b { c { d e f \n g h i }}",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{
						StateSpec{Name: "c"}, []SubTransition{
							{"d", "e", []string{"f"}},
							{"g", "h", []string{"i"}},
						},
					},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b { c { - - - } g { h i j } }",
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{StateSpec{Name: "c"}, []SubTransition{{"", "", []string{}}}},
					{StateSpec{Name: "g"}, []SubTransition{{"h", "i", []string{"j"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			`a:b {
				c {
					d e {f g}
					h i j
				}
			}`,
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{
						StateSpec{Name: "c"}, []SubTransition{
							{"d", "e", []string{"f", "g"}},
							{"h", "i", []string{"j"}},
						},
					},
				},
				Done: true,
			})

		assertParserResult(t,
			`a:b {
				(c) d e f
				(g) h i -
				j : c : g - - -
			}`,
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{
						StateSpec{Name: "c", AbstractState: true}, []SubTransition{
							{"d", "e", []string{"f"}},
						},
					},
					{
						StateSpec{Name: "g", AbstractState: true}, []SubTransition{
							{"h", "i", []string{}},
						},
					},
					{
						StateSpec{Name: "j", SuperStates: []string{"c", "g"}}, []SubTransition{
							{"", "", []string{}},
						},
					},
				},
				Done: true,
			})

		assertParserResult(t,
			`a:b {
				c >d >e <f <g h i j
			}`,
			FSMSyntax{
				Headers: []Header{{Name: "a", Value: "b"}},
				Logic: []Transition{
					{
						StateSpec{
							Name:         "c",
							EntryActions: []string{"d", "e"},
							ExitActions:  []string{"f", "g"},
						}, []SubTransition{
							{"h", "i", []string{"j"}},
						},
					},
				},
				Done: true,
			})
	})

	t.Run("Error tests", func(t *testing.T) {
		assertParserResult(t,
			"a:b . {}",
			FSMSyntax{
				Headers: []Header{
					{Name: "a", Value: "b"},
				},
				Errors: []SyntaxError{
					{Type: ErrorSyntax, LineNumber: 1, Position: 5, Msg: ""},
				},
				Done: true,
			})

		assertParserResult(t,
			`a:b:c:d {
				e f { g h
			}`,
			FSMSyntax{
				Headers: []Header{
					{Name: "a", Value: "b"},
					{Name: "c", Value: "d"},
				},
				Logic: []Transition{
					{StateSpec{Name: "e"}, []SubTransition{{"f", "g", []string{"h"}}}},
				},
				Errors: []SyntaxError{
					{Type: ErrorParse, LineNumber: 1, Position: 4, Msg: "HEADER|COLON"},
					{Type: ErrorParse, LineNumber: 2, Position: 9, Msg: "SINGLE_EVENT|OPEN_BRACE"},
				},
				Done: true,
			})

		assertParserResult(t,
			"a:b {",
			FSMSyntax{
				Headers: []Header{
					{Name: "a", Value: "b"},
				},
				Errors: []SyntaxError{
					{Type: ErrorParse, LineNumber: 2, Position: 1, Msg: "TRANSITION_GROUP|END"},
				},
				Done: true,
			})
	})

	t.Run("Acceptance tests", func(t *testing.T) {
		assertParserResult(t,
			`Actions: Turnstile
			FSM: OneCoinTurnstile
			Initial: Locked
			{
			  Locked	Coin	Unlocked	{alarmOff unlock}
			  Locked 	Pass	Locked		alarmOn
			  Unlocked	Coin	Unlocked	thankyou
			  Unlocked	Pass	Locked		lock
			}`,
			FSMSyntax{
				Headers: []Header{
					{Name: "Actions", Value: "Turnstile"},
					{Name: "FSM", Value: "OneCoinTurnstile"},
					{Name: "Initial", Value: "Locked"},
				},
				Logic: []Transition{
					{StateSpec{Name: "Locked"}, []SubTransition{{"Coin", "Unlocked", []string{"alarmOff", "unlock"}}}},
					{StateSpec{Name: "Locked"}, []SubTransition{{"Pass", "Locked", []string{"alarmOn"}}}},
					{StateSpec{Name: "Unlocked"}, []SubTransition{{"Coin", "Unlocked", []string{"thankyou"}}}},
					{StateSpec{Name: "Unlocked"}, []SubTransition{{"Pass", "Locked", []string{"lock"}}}},
				},
				Done: true,
			})

		assertParserResult(t,
			`Actions: Turnstile
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
			}`,
			FSMSyntax{
				Headers: []Header{
					{Name: "Actions", Value: "Turnstile"},
					{Name: "FSM", Value: "TwoCoinTurnstile"},
					{Name: "Initial", Value: "Locked"},
				},
				Logic: []Transition{
					{StateSpec{Name: "Locked"}, []SubTransition{
						{"Pass", "Alarming", []string{"alarmOn"}},
						{"Coin", "FirstCoin", []string{}},
						{"Reset", "Locked", []string{"lock", "alarmOff"}},
					}},
					{StateSpec{Name: "Alarming"}, []SubTransition{
						{"Reset", "Locked", []string{"lock", "alarmOff"}},
					}},
					{StateSpec{Name: "FirstCoin"}, []SubTransition{
						{"Pass", "Alarming", []string{}},
						{"Coin", "Unlocked", []string{"unlock"}},
						{"Reset", "Locked", []string{"lock", "alarmOff"}},
					}},
					{StateSpec{Name: "Unlocked"}, []SubTransition{
						{"Pass", "Locked", []string{"lock"}},
						{"Coin", "", []string{"thankyou"}},
						{"Reset", "Locked", []string{"lock", "alarmOff"}},
					}},
				},
				Done: true,
			})

		assertParserResult(t,
			`Actions: Turnstile
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
			}`,
			FSMSyntax{
				Headers: []Header{
					{Name: "Actions", Value: "Turnstile"},
					{Name: "FSM", Value: "TwoCoinTurnstile"},
					{Name: "Initial", Value: "Locked"},
				},
				Logic: []Transition{
					{StateSpec{Name: "Base", AbstractState: true}, []SubTransition{
						{"Reset", "Locked", []string{"alarmOff", "lock"}},
					}},
					{StateSpec{Name: "Locked", SuperStates: []string{"Base"}}, []SubTransition{
						{"Pass", "Alarming", []string{"alarmOn"}},
						{"Coin", "FirstCoin", []string{}},
					}},
					{StateSpec{Name: "Alarming", SuperStates: []string{"Base"}}, []SubTransition{
						{"", "", []string{}},
					}},
					{StateSpec{Name: "FirstCoin", SuperStates: []string{"Base"}}, []SubTransition{
						{"Pass", "Alarming", []string{}},
						{"Coin", "Unlocked", []string{"unlock"}},
					}},
					{StateSpec{Name: "Unlocked", SuperStates: []string{"Base"}}, []SubTransition{
						{"Pass", "Locked", []string{"lock"}},
						{"Coin", "", []string{"thankyou"}},
					}},
				},
				Done: true,
			})

		assertParserResult(t,
			`Actions: Turnstile
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
			}`,
			FSMSyntax{
				Headers: []Header{
					{Name: "Actions", Value: "Turnstile"},
					{Name: "FSM", Value: "TwoCoinTurnstile"},
					{Name: "Initial", Value: "Locked"},
				},
				Logic: []Transition{
					{StateSpec{Name: "Base", AbstractState: true}, []SubTransition{
						{"Reset", "Locked", []string{"lock"}},
					}},
					{StateSpec{Name: "Locked", SuperStates: []string{"Base"}}, []SubTransition{
						{"Pass", "Alarming", []string{}},
						{"Coin", "FirstCoin", []string{}},
					}},
					{StateSpec{
						Name:         "Alarming",
						SuperStates:  []string{"Base"},
						EntryActions: []string{"alarmOn"},
						ExitActions:  []string{"alarmOff"},
					}, []SubTransition{
						{"", "", []string{}},
					}},
					{StateSpec{Name: "FirstCoin", SuperStates: []string{"Base"}}, []SubTransition{
						{"Pass", "Alarming", []string{}},
						{"Coin", "Unlocked", []string{"unlock"}},
					}},
					{StateSpec{Name: "Unlocked", SuperStates: []string{"Base"}}, []SubTransition{
						{"Pass", "Locked", []string{"lock"}},
						{"Coin", "", []string{"thankyou"}},
					}},
				},
				Done: true,
			})
	})
}

func assertParserResult(t *testing.T, input string, expected FSMSyntax) {
	t.Helper()
	builder := NewSyntaxBuilder()
	parser := NewParser(builder)
	lexer := lexer.NewLexer(parser)

	lexer.Lex(bytes.NewBufferString(input))
	assert.Equal(t, expected, builder.FSM())
}
