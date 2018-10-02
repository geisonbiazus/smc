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
	    Actions: b
	    Initial: c
	    {
	      d e f g
	      h i j {k l}
	    }
			`,
			FSM{
				Name:         "a",
				ActionsClass: "b",
				InitialState: "c",
				States: []State{
					{Name: "d", Transitions: []Transition{
						{Event: "e", NextState: "f", Actions: []string{"g"}},
					}},
					{Name: "h", Transitions: []Transition{
						{Event: "i", NextState: "j", Actions: []string{"k", "l"}},
					}},
				},
				Events:  []string{"e", "i"},
				Actions: []string{"g", "k", "l"},
			},
		)
	})
}

func optimizeFSM(input string) FSM {
	builder := parser.NewSyntaxBuilder()
	parser := parser.NewParser(builder)
	lexer := lexer.NewLexer(parser)
	lexer.Lex(bytes.NewBufferString(input))

	analyzer := semantic.NewAnalyzer()
	semanticFSM := analyzer.Analyze(builder.FSM())
	opt := New()

	return opt.Optimize(semanticFSM)
}

func assertOptimizedFSM(t *testing.T, input string, expected FSM) {
	assert.Equal(t, expected, optimizeFSM(input))
}
