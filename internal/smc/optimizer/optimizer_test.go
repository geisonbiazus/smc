package optimizer

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
	"github.com/stretchr/testify/assert"
)

func TestOptimizer(t *testing.T) {
	t.Run("Simple FSM", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: a
	    Initial: c
	    {
	      d e f g
	      h i j {k l}
	    }
			`,
			&FSM{
				Name:         "a",
				InitialState: "c",
				States: []*State{
					{Name: "d", Transitions: []*Transition{
						{Event: "e", NextState: "f", Actions: []string{"g"}},
					}},
					{Name: "h", Transitions: []*Transition{
						{Event: "i", NextState: "j", Actions: []string{"k", "l"}},
					}},
				},
				Events:  []string{"e", "i"},
				Actions: []string{"g", "k", "l"},
			},
		)
	})

	t.Run("With abstract superclass", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				(a) b c d
				e:a - - -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "e", Transitions: []*Transition{
						{Event: "b", NextState: "c", Actions: []string{"d"}},
					}},
				},
				Events:  []string{"b"},
				Actions: []string{"d"},
			},
		)
	})

	t.Run("With not abstract superclass", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				a b c d
				e:a - - -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "a", Transitions: []*Transition{
						{Event: "b", NextState: "c", Actions: []string{"d"}},
					}},
					{Name: "e", Transitions: []*Transition{
						{Event: "b", NextState: "c", Actions: []string{"d"}},
					}},
				},
				Events:  []string{"b"},
				Actions: []string{"d"},
			},
		)
	})

	t.Run("With deep inheritance", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				(a) Ea - Aa
				(b):a Eb Nb {Ab1 Ab2}
				c:b Ec Nc Ac
				d Ed Nd Ad
				e:c:d Ee Ne Ae
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "c", Transitions: []*Transition{
						{Event: "Ec", NextState: "Nc", Actions: []string{"Ac"}},
						{Event: "Eb", NextState: "Nb", Actions: []string{"Ab1", "Ab2"}},
						{Event: "Ea", NextState: "", Actions: []string{"Aa"}},
					}},
					{Name: "d", Transitions: []*Transition{
						{Event: "Ed", NextState: "Nd", Actions: []string{"Ad"}},
					}},
					{Name: "e", Transitions: []*Transition{
						{Event: "Ee", NextState: "Ne", Actions: []string{"Ae"}},
						{Event: "Ec", NextState: "Nc", Actions: []string{"Ac"}},
						{Event: "Eb", NextState: "Nb", Actions: []string{"Ab1", "Ab2"}},
						{Event: "Ea", NextState: "", Actions: []string{"Aa"}},
						{Event: "Ed", NextState: "Nd", Actions: []string{"Ad"}},
					}},
				},
				Events:  []string{"Ea", "Eb", "Ec", "Ed", "Ee"},
				Actions: []string{"Aa", "Ab1", "Ab2", "Ac", "Ad", "Ae"},
			},
		)
	})

	t.Run("With transition override", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				(a) E - Aa
				(b):a E Nb {Ab1 Ab2}
				c:b E Nc Ac
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "c", Transitions: []*Transition{
						{Event: "E", NextState: "Nc", Actions: []string{"Ac"}},
					}},
				},
				Events:  []string{"E"},
				Actions: []string{"Aa", "Ab1", "Ab2", "Ac"},
			},
		)
	})

	t.Run("With entry actions", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				S1 >EA1 >EA2 E1 S2 -
				S2 E2 S1 A2
				S3 E3 S1 -
				S4 E4 S2 -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S1", Transitions: []*Transition{
						{Event: "E1", NextState: "S2", Actions: []string{}},
					}},
					{Name: "S2", Transitions: []*Transition{
						{Event: "E2", NextState: "S1", Actions: []string{"A2", "EA1", "EA2"}},
					}},
					{Name: "S3", Transitions: []*Transition{
						{Event: "E3", NextState: "S1", Actions: []string{"EA1", "EA2"}},
					}},
					{Name: "S4", Transitions: []*Transition{
						{Event: "E4", NextState: "S2", Actions: []string{}},
					}},
				},
				Events:  []string{"E1", "E2", "E3", "E4"},
				Actions: []string{"EA1", "EA2", "A2"},
			},
		)
	})

	t.Run("Entry actions are not executed with no transition", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				S1 >EA1 >EA2 E1 S1 -
				S2 >EA1 >EA2 E1 - -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S1", Transitions: []*Transition{
						{Event: "E1", NextState: "S1", Actions: []string{}},
					}},
					{Name: "S2", Transitions: []*Transition{
						{Event: "E1", NextState: "", Actions: []string{}},
					}},
				},
				Events:  []string{"E1"},
				Actions: []string{"EA1", "EA2"},
			},
		)
	})

	t.Run("With entry inheritance", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				(S1) >EA1 >EA2 E1 S3 -
				S2:S1 - - -
				S3 E3 S2 -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S2", Transitions: []*Transition{
						{Event: "E1", NextState: "S3", Actions: []string{}},
					}},
					{Name: "S3", Transitions: []*Transition{
						{Event: "E3", NextState: "S2", Actions: []string{"EA1", "EA2"}},
					}},
				},
				Events:  []string{"E1", "E3"},
				Actions: []string{"EA1", "EA2"},
			},
		)
	})

	t.Run("Inherited duplciated entry actions are ignored", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
			Initial: initial
			{
				(S1) >EA1 >EA3 - - -
				(S2):S1 >EA1 >EA2 - - -
				S3:S2 - - -
				S4 E1 S3 {A1 EA2}
			}
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S3"},
					{Name: "S4", Transitions: []*Transition{
						{Event: "E1", NextState: "S3", Actions: []string{"A1", "EA2", "EA1", "EA3"}},
					}},
				},
				Events:  []string{"E1"},
				Actions: []string{"EA1", "EA3", "EA2", "A1"},
			},
		)
	})

	t.Run("With exit actions", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				S1 <EA1 <EA2 {
					E1 S2 A1
					E2 S3 -
				}
				S2 E2 S1 A2
				S3 E3 S1 -
				S4 E4 S2 -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S1", Transitions: []*Transition{
						{Event: "E1", NextState: "S2", Actions: []string{"A1", "EA1", "EA2"}},
						{Event: "E2", NextState: "S3", Actions: []string{"EA1", "EA2"}},
					}},
					{Name: "S2", Transitions: []*Transition{
						{Event: "E2", NextState: "S1", Actions: []string{"A2"}},
					}},
					{Name: "S3", Transitions: []*Transition{
						{Event: "E3", NextState: "S1", Actions: []string{}},
					}},
					{Name: "S4", Transitions: []*Transition{
						{Event: "E4", NextState: "S2", Actions: []string{}},
					}},
				},
				Events:  []string{"E1", "E2", "E3", "E4"},
				Actions: []string{"EA1", "EA2", "A1", "A2"},
			},
		)
	})

	t.Run("With exit inheritance", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
			Initial: initial
			{
				(S1) <EA1 <EA2 E1 S3 -
				S2:S1 - - -
				S3 E3 S2 -
			}
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S2", Transitions: []*Transition{
						{Event: "E1", NextState: "S3", Actions: []string{"EA1", "EA2"}},
					}},
					{Name: "S3", Transitions: []*Transition{
						{Event: "E3", NextState: "S2", Actions: []string{}},
					}},
				},
				Events:  []string{"E1", "E3"},
				Actions: []string{"EA1", "EA2"},
			},
		)
	})

	t.Run("Inherited duplciated exit actions are ignored", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
			Initial: initial
			{
				(S1) <EA1 <EA3 - - -
				(S2):S1 <EA1 <EA2 - - -
				S3:S2 E1 S4 {A1 EA2}
				S4 - - -
			}
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S3", Transitions: []*Transition{
						{Event: "E1", NextState: "S4", Actions: []string{"A1", "EA2", "EA1", "EA3"}},
					}},
					{Name: "S4"},
				},
				Events:  []string{"E1"},
				Actions: []string{"EA1", "EA3", "EA2", "A1"},
			},
		)
	})

	t.Run("Exit actions are not executed with no transition", func(t *testing.T) {
		assertOptimizedFSM(t, `
			FSM: fsm
	    Initial: initial
	    {
				S1 <EA1 <EA2 E1 S1 -
				S2 <EA1 <EA2 E1 - -
	    }
			`,
			&FSM{
				Name:         "fsm",
				InitialState: "initial",
				States: []*State{
					{Name: "S1", Transitions: []*Transition{
						{Event: "E1", NextState: "S1", Actions: []string{}},
					}},
					{Name: "S2", Transitions: []*Transition{
						{Event: "E1", NextState: "", Actions: []string{}},
					}},
				},
				Events:  []string{"E1"},
				Actions: []string{"EA1", "EA2"},
			},
		)
	})

	t.Run("Acceptance tests", func(t *testing.T) {
		assertOptimizedFSM(t, `
			Actions: Turnstile
			FSM: OneCoinTurnstile
			Initial: Locked
			{
			  Locked	Coin	Unlocked	{alarmOff unlock}
			  Locked 	Pass	Locked		alarmOn
			  Unlocked	Coin	Unlocked	thankyou
			  Unlocked	Pass	Locked		lock
			}
			`,
			&FSM{
				Name:         "OneCoinTurnstile",
				InitialState: "Locked",
				States: []*State{
					{Name: "Locked", Transitions: []*Transition{
						{Event: "Coin", NextState: "Unlocked", Actions: []string{"alarmOff", "unlock"}},
						{Event: "Pass", NextState: "Locked", Actions: []string{"alarmOn"}},
					}},
					{Name: "Unlocked", Transitions: []*Transition{
						{Event: "Coin", NextState: "Unlocked", Actions: []string{"thankyou"}},
						{Event: "Pass", NextState: "Locked", Actions: []string{"lock"}},
					}},
				},
				Events:  []string{"Coin", "Pass"},
				Actions: []string{"alarmOff", "unlock", "alarmOn", "thankyou", "lock"},
			},
		)

		assertOptimizedFSM(t, `
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
			`,
			&FSM{
				Name:         "TwoCoinTurnstile",
				InitialState: "Locked",
				States: []*State{
					{Name: "Locked", Transitions: []*Transition{
						{Event: "Pass", NextState: "Alarming", Actions: []string{"alarmOn"}},
						{Event: "Coin", NextState: "FirstCoin", Actions: []string{}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock", "alarmOff"}},
					}},
					{Name: "Alarming", Transitions: []*Transition{
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock", "alarmOff"}},
					}},
					{Name: "FirstCoin", Transitions: []*Transition{
						{Event: "Pass", NextState: "Alarming", Actions: []string{}},
						{Event: "Coin", NextState: "Unlocked", Actions: []string{"unlock"}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock", "alarmOff"}},
					}},
					{Name: "Unlocked", Transitions: []*Transition{
						{Event: "Pass", NextState: "Locked", Actions: []string{"lock"}},
						{Event: "Coin", NextState: "", Actions: []string{"thankyou"}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock", "alarmOff"}},
					}},
				},
				Events:  []string{"Pass", "Coin", "Reset"},
				Actions: []string{"alarmOn", "lock", "alarmOff", "unlock", "thankyou"},
			},
		)

		assertOptimizedFSM(t, `
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
			`,
			&FSM{
				Name:         "TwoCoinTurnstile",
				InitialState: "Locked",
				States: []*State{
					{Name: "Locked", Transitions: []*Transition{
						{Event: "Pass", NextState: "Alarming", Actions: []string{"alarmOn"}},
						{Event: "Coin", NextState: "FirstCoin", Actions: []string{}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"alarmOff", "lock"}},
					}},
					{Name: "Alarming", Transitions: []*Transition{
						{Event: "Reset", NextState: "Locked", Actions: []string{"alarmOff", "lock"}},
					}},
					{Name: "FirstCoin", Transitions: []*Transition{
						{Event: "Pass", NextState: "Alarming", Actions: []string{}},
						{Event: "Coin", NextState: "Unlocked", Actions: []string{"unlock"}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"alarmOff", "lock"}},
					}},
					{Name: "Unlocked", Transitions: []*Transition{
						{Event: "Pass", NextState: "Locked", Actions: []string{"lock"}},
						{Event: "Coin", NextState: "", Actions: []string{"thankyou"}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"alarmOff", "lock"}},
					}},
				},
				Events:  []string{"Reset", "Pass", "Coin"},
				Actions: []string{"alarmOff", "lock", "alarmOn", "unlock", "thankyou"},
			},
		)

		assertOptimizedFSM(t, `
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
			`,
			&FSM{
				Name:         "TwoCoinTurnstile",
				InitialState: "Locked",
				States: []*State{
					{Name: "Locked", Transitions: []*Transition{
						{Event: "Pass", NextState: "Alarming", Actions: []string{"alarmOn"}},
						{Event: "Coin", NextState: "FirstCoin", Actions: []string{}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock"}},
					}},
					{Name: "Alarming", Transitions: []*Transition{
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock", "alarmOff"}},
					}},
					{Name: "FirstCoin", Transitions: []*Transition{
						{Event: "Pass", NextState: "Alarming", Actions: []string{"alarmOn"}},
						{Event: "Coin", NextState: "Unlocked", Actions: []string{"unlock"}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock"}},
					}},
					{Name: "Unlocked", Transitions: []*Transition{
						{Event: "Pass", NextState: "Locked", Actions: []string{"lock"}},
						{Event: "Coin", NextState: "", Actions: []string{"thankyou"}},
						{Event: "Reset", NextState: "Locked", Actions: []string{"lock"}},
					}},
				},
				Events:  []string{"Reset", "Pass", "Coin"},
				Actions: []string{"lock", "alarmOn", "alarmOff", "unlock", "thankyou"},
			},
		)
	})
}

func optimizeFSM(input string) *FSM {
	builder := parser.NewSyntaxBuilder()
	parser := parser.NewParser(builder)
	lexer := lexer.NewLexer(parser)
	lexer.Lex(bytes.NewBufferString(input))

	analyzer := semantic.NewAnalyzer()
	semanticFSM := analyzer.Analyze(builder.FSM())
	opt := New()

	return opt.Optimize(semanticFSM)
}

func assertOptimizedFSM(t *testing.T, input string, expected *FSM) {
	assert.Equal(t, expected, optimizeFSM(input))
}
