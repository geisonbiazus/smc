package statepattern

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/optimizer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
	"github.com/stretchr/testify/assert"
)

func TestStatePattern(t *testing.T) {
	t.Run("Single state", func(t *testing.T) {
		assertGeneratedFSM(t,
			"FSM: fsm Initial: a { a b a c }",
			CompositeNode([]Node{
				StateInterfaceNode{
					FSMClassName: "fsm",
					States:       []string{"a"},
				},
				ActionsInterfaceNode{
					Actions: []string{"c"},
				},
				FSMClassNode{
					ClassName:    "fsm",
					InitialState: "a",
					EventMethods: []Node{
						EventMethodNode{ClassName: "fsm", EventName: "b"},
					},
				},
				BaseStateClassNode{
					Events: []string{"b"},
				},
				CompositeNode([]Node{
					StateClassNode{
						StateName: "a",
						StateEventMethods: []Node{
							StateEventMethodNode{
								FSMClassName: "fsm",
								StateName:    "a",
								EventName:    "b",
								NextState:    "a",
								Actions:      []string{"c"},
							},
						},
					},
				}),
			}),
		)
	})

	t.Run("Full FSM", func(t *testing.T) {
		assertGeneratedFSM(t, `
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
			CompositeNode([]Node{
				StateInterfaceNode{
					FSMClassName: "TwoCoinTurnstile",
					States:       []string{"Locked", "Alarming", "FirstCoin", "Unlocked"},
				},
				ActionsInterfaceNode{
					Actions: []string{"lock", "alarmOn", "alarmOff", "unlock", "thankyou"},
				},
				FSMClassNode{
					ClassName:    "TwoCoinTurnstile",
					InitialState: "Locked",
					EventMethods: []Node{
						EventMethodNode{ClassName: "TwoCoinTurnstile", EventName: "Reset"},
						EventMethodNode{ClassName: "TwoCoinTurnstile", EventName: "Pass"},
						EventMethodNode{ClassName: "TwoCoinTurnstile", EventName: "Coin"},
					},
				},
				BaseStateClassNode{
					Events: []string{"Reset", "Pass", "Coin"},
				},
				CompositeNode([]Node{
					StateClassNode{
						StateName: "Locked",
						StateEventMethods: []Node{
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Locked",
								EventName:    "Pass",
								NextState:    "Alarming",
								Actions:      []string{"alarmOn"},
							},
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Locked",
								EventName:    "Coin",
								NextState:    "FirstCoin",
								Actions:      []string{},
							},
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Locked",
								EventName:    "Reset",
								NextState:    "Locked",
								Actions:      []string{"lock"},
							},
						},
					},
					StateClassNode{
						StateName: "Alarming",
						StateEventMethods: []Node{
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Alarming",
								EventName:    "Reset",
								NextState:    "Locked",
								Actions:      []string{"lock", "alarmOff"},
							},
						},
					},
					StateClassNode{
						StateName: "FirstCoin",
						StateEventMethods: []Node{
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "FirstCoin",
								EventName:    "Pass",
								NextState:    "Alarming",
								Actions:      []string{"alarmOn"},
							},
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "FirstCoin",
								EventName:    "Coin",
								NextState:    "Unlocked",
								Actions:      []string{"unlock"},
							},
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "FirstCoin",
								EventName:    "Reset",
								NextState:    "Locked",
								Actions:      []string{"lock"},
							},
						},
					},
					StateClassNode{
						StateName: "Unlocked",
						StateEventMethods: []Node{
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Unlocked",
								EventName:    "Pass",
								NextState:    "Locked",
								Actions:      []string{"lock"},
							},
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Unlocked",
								EventName:    "Coin",
								NextState:    "",
								Actions:      []string{"thankyou"},
							},
							StateEventMethodNode{
								FSMClassName: "TwoCoinTurnstile",
								StateName:    "Unlocked",
								EventName:    "Reset",
								NextState:    "Locked",
								Actions:      []string{"lock"},
							},
						},
					},
				}),
			}),
		)
	})
}

func assertGeneratedFSM(t *testing.T, input string, expected Node) {
	node := generateFSM(input)
	assert.Equal(t, expected, node)
}

func generateFSM(input string) Node {
	builder := parser.NewSyntaxBuilder()
	psr := parser.NewParser(builder)
	lxr := lexer.NewLexer(psr)
	lxr.Lex(bytes.NewBufferString(input))

	parsedFSM := builder.FSM()

	analyzer := semantic.NewAnalyzer()
	semanticFSM := analyzer.Analyze(parsedFSM)

	opt := optimizer.New()
	optimizedFSM := opt.Optimize(semanticFSM)

	gen := NewStatePattern()
	return gen.Generate(optimizedFSM)
}
