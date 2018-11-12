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
		gen := NewStatePattern()

		syntax := "FSM: fsm Initial: a { a b a c }"

		builder := parser.NewSyntaxBuilder()
		psr := parser.NewParser(builder)
		lxr := lexer.NewLexer(psr)
		lxr.Lex(bytes.NewBufferString(syntax))

		parsedFSM := builder.FSM()

		analyzer := semantic.NewAnalyzer()
		semanticFSM := analyzer.Analyze(parsedFSM)

		opt := optimizer.New()
		optimizedFSM := opt.Optimize(semanticFSM)

		node := gen.Generate(optimizedFSM)

		expected := CompositeNode([]Node{
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
		})

		assert.Equal(t, expected, node)
	})
}
