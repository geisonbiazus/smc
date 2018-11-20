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
	t.Run("Simple FSM", func(t *testing.T) {
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
	})

	t.Run("Complex FSM", func(t *testing.T) {
		assertImplementedFSM(t, `
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
			}`,
			`
			type State interface {
			  Reset(fsm *TwoCoinTurnstile)
			  Pass(fsm *TwoCoinTurnstile)
			  Coin(fsm *TwoCoinTurnstile)
			}

			type Actions interface {
			  Lock()
			  AlarmOn()
			  AlarmOff()
			  Unlock()
			  Thankyou()
			  UnhandledTransition(state string, event string)
			}

			type TwoCoinTurnstile struct {
			  State State
			  Actions Actions
			}

			func NewTwoCoinTurnstile(actions Actions) *TwoCoinTurnstile {
			  return &TwoCoinTurnstile{
			    Actions: actions,
			    State:   NewStateLocked(),
			  }
			}

			func (f *TwoCoinTurnstile) Reset() {
			  f.State.Reset(f)
			}

			func (f *TwoCoinTurnstile) Pass() {
			  f.State.Pass(f)
			}

			func (f *TwoCoinTurnstile) Coin() {
			  f.State.Coin(f)
			}

			type BaseState struct {  StateName string
			}

			func (b BaseState) Reset(fsm *TwoCoinTurnstile) {
			  fsm.Actions.UnhandledTransition(b.StateName, "Reset")
			}

			func (b BaseState) Pass(fsm *TwoCoinTurnstile) {
			  fsm.Actions.UnhandledTransition(b.StateName, "Pass")
			}

			func (b BaseState) Coin(fsm *TwoCoinTurnstile) {
			  fsm.Actions.UnhandledTransition(b.StateName, "Coin")
			}

			type StateLocked struct {
			  BaseState
			}

			func NewStateLocked() StateLocked {
			  return StateLocked{BaseState{StateName: "Locked"}}
			}

			func (s StateLocked) Pass(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateAlarming()
			  fsm.Actions.AlarmOn()
			}

			func (s StateLocked) Coin(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateFirstCoin()
			}

			func (s StateLocked) Reset(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateLocked()
			  fsm.Actions.Lock()
			}

			type StateAlarming struct {
			  BaseState
			}

			func NewStateAlarming() StateAlarming {
			  return StateAlarming{BaseState{StateName: "Alarming"}}
			}

			func (s StateAlarming) Reset(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateLocked()
			  fsm.Actions.Lock()
			  fsm.Actions.AlarmOff()
			}

			type StateFirstCoin struct {
			  BaseState
			}

			func NewStateFirstCoin() StateFirstCoin {
			  return StateFirstCoin{BaseState{StateName: "FirstCoin"}}
			}

			func (s StateFirstCoin) Pass(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateAlarming()
			  fsm.Actions.AlarmOn()
			}

			func (s StateFirstCoin) Coin(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateUnlocked()
			  fsm.Actions.Unlock()
			}

			func (s StateFirstCoin) Reset(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateLocked()
			  fsm.Actions.Lock()
			}

			type StateUnlocked struct {
			  BaseState
			}

			func NewStateUnlocked() StateUnlocked {
			  return StateUnlocked{BaseState{StateName: "Unlocked"}}
			}

			func (s StateUnlocked) Pass(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateLocked()
			  fsm.Actions.Lock()
			}

			func (s StateUnlocked) Coin(fsm *TwoCoinTurnstile) {
			  fsm.Actions.Thankyou()
			}

			func (s StateUnlocked) Reset(fsm *TwoCoinTurnstile) {
			  fsm.State = NewStateLocked()
			  fsm.Actions.Lock()
			}
			`,
		)
	})
}

func assertImplementedFSM(t *testing.T, input, expected string) {
	result := implementFSM(input)
	// fmt.Println(result)
	assert.Equal(t, removeSpacing(expected), removeSpacing(result))
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
