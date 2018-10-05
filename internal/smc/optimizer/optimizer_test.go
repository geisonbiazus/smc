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
			&FSM{
				Name:         "a",
				ActionsClass: "b",
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
	    Actions: actions
	    Initial: initial
	    {
				(a) b c d
				e:a - - -
	    }
			`,
			&FSM{
				Name:         "fsm",
				ActionsClass: "actions",
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
	    Actions: actions
	    Initial: initial
	    {
				a b c d
				e:a - - -
	    }
			`,
			&FSM{
				Name:         "fsm",
				ActionsClass: "actions",
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
	    Actions: actions
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
				ActionsClass: "actions",
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
	    Actions: actions
	    Initial: initial
	    {
				(a) E - Aa
				(b):a E Nb {Ab1 Ab2}
				c:b E Nc Ac
	    }
			`,
			&FSM{
				Name:         "fsm",
				ActionsClass: "actions",
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
	    Actions: actions
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
				ActionsClass: "actions",
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
