package golang

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/generator/statepattern"
	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/optimizer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
	"github.com/stretchr/testify/assert"
)

func TestImplementer(t *testing.T) {
	assertImplementedFSM(t,
		"FSM: fsm Initial: state { state event state action }",
		`
		type State interface {
			Event(fsm *Fsm)
		}

		type Actions interface {
			Action()
			UnhandledTransition(state string, event string)
		}

		type Fsm struct {
			State   State
			Actions Actions
		}

		func NewFsm(actions Actions) *Fsm {
			return &Fsm{
				Actions: actions,
				State:   NewStateState(),
			}
		}

		func (f *Fsm) Event() {
			f.State.Event(f)
		}

		type BaseState struct {
			StateName string
		}

		func (b BaseState) Event(fsm *Fsm) {
			fsm.Actions.UnhandledTransition(b.StateName, "event")
		}

		type StateState struct {
			BaseState
		}

		func NewStateState() StateState {
			return StateState{BaseState{StateName: "state"}}
		}

		func (s StateState) Event(fsm *Fsm) {
			fsm.State = NewStateState()
			fsm.Actions.Action()
		}
		`,
	)
}

func assertImplementedFSM(t *testing.T, input, expected string) {
	node := implementFSM(input)
	assert.Equal(t, removeSpacing(expected), removeSpacing(node))
}

func implementFSM(input string) string {
	builder := parser.NewSyntaxBuilder()
	psr := parser.NewParser(builder)
	lxr := lexer.NewLexer(psr)
	lxr.Lex(bytes.NewBufferString(input))

	parsedFSM := builder.FSM()

	analyzer := semantic.NewAnalyzer()
	semanticFSM := analyzer.Analyze(parsedFSM)

	opt := optimizer.New()
	optimizedFSM := opt.Optimize(semanticFSM)

	gen := statepattern.NewNodeGenerator()
	node := gen.Generate(optimizedFSM)

	implementer := NewImplementer()

	return implementer.Implement(node)
}

var whitespaceRegex = regexp.MustCompile("\\s+")

func removeSpacing(s string) string {
	return whitespaceRegex.ReplaceAllString(s, " ")
}
