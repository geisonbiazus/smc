package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	t.Run("Collect syntax errors", func(t *testing.T) {
		assertContainsError(t, compileFSM("& a:b {}", &bytes.Buffer{}),
			parser.SyntaxError{
				Type: parser.ErrorSyntax, LineNumber: 1, Position: 1,
			},
		)
	})

	t.Run("Collect parse errors", func(t *testing.T) {
		assertContainsError(t, compileFSM("a:b:c {}", &bytes.Buffer{}),
			parser.SyntaxError{
				Type: parser.ErrorParse, LineNumber: 1, Position: 4, Msg: "HEADER|COLON",
			},
		)
	})

	t.Run("Collect semantic errors", func(t *testing.T) {
		assertContainsError(t, compileFSM("a:b {}", &bytes.Buffer{}),
			semantic.Error{Type: semantic.ErrorNoFSM, Element: "FSM"},
		)
	})

	t.Run("Write the compiled output", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		compileFSM("FSM: fsm Initial: state { state event state action }", buffer)

		assert.Equal(t, compiledFSM, buffer.String())
	})
}

func compileFSM(input string, output *bytes.Buffer) *Compiler {
	compiler := NewCompiler(bytes.NewBufferString(input), output)
	compiler.Compile()
	return compiler
}

func assertContainsError(t *testing.T, compiler *Compiler, err Error) {
	t.Helper()
	assert.Contains(t, compiler.Errors, err)
}

var compiledFSM = `package fsm

type State interface {
  Event(fsm *Fsm)
}

type Actions interface {
  Action()
  UnhandledTransition(state string, event string)
}

type Fsm struct {
  State State
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

type BaseState struct {  StateName string
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
`
