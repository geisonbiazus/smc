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
		compiler, err := compileFSM("& a:b {}", &bytes.Buffer{})
		assertContainsError(t, compiler,
			parser.SyntaxError{
				Type: parser.ErrorSyntax, LineNumber: 1, Position: 1,
			},
		)
		assert.Equal(t, CompileError, err)
	})

	t.Run("Collect parse errors", func(t *testing.T) {
		compiler, err := compileFSM("a:b:c {}", &bytes.Buffer{})
		assertContainsError(t, compiler,
			parser.SyntaxError{
				Type: parser.ErrorParse, LineNumber: 1, Position: 4, Msg: "HEADER|COLON",
			},
		)
		assert.Equal(t, CompileError, err)
	})

	t.Run("Collect semantic errors", func(t *testing.T) {
		compiler, err := compileFSM("a:b {}", &bytes.Buffer{})
		assertContainsError(t, compiler,
			semantic.Error{Type: semantic.ErrorNoFSM, Element: "FSM"},
		)
		assert.Equal(t, CompileError, err)
	})

	t.Run("Write the compiled output", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		_, err := compileFSM("FSM: fsm Initial: state { state event state action }", buffer)

		assert.Equal(t, compiledFSM, buffer.String())
		assert.Nil(t, err)
	})
}

func compileFSM(input string, output *bytes.Buffer) (*Compiler, error) {
	compiler := NewCompiler(bytes.NewBufferString(input), output)
	err := compiler.Compile()
	return compiler, err
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
