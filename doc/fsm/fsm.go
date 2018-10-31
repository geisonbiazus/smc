package turnstile

type State interface {
	Pass(fsm *TurnstileFSM)
	Coin(fsm *TurnstileFSM)
}

type Actions interface {
	AlarmOn()
	AlarmOff()
	Lock()
	Unlock()
	Thankyou()
	UnhandledTransition(state string, event string)
}

type TurnstileFSM struct {
	State   State
	Actions Actions
}

func NewTurnstileFSM(actions Actions) *TurnstileFSM {
	return &TurnstileFSM{
		Actions: actions,
		State:   NewStateLocked(),
	}
}

func (t *TurnstileFSM) Pass() {
	t.State.Pass(t)
}

func (t *TurnstileFSM) Coin() {
	t.State.Coin(t)
}

type BaseState struct {
	StateName string
}

func (b BaseState) Pass(fsm *TurnstileFSM) {
	fsm.Actions.UnhandledTransition(b.StateName, "Pass")
}

func (b BaseState) Coin(fsm *TurnstileFSM) {
	fsm.Actions.UnhandledTransition(b.StateName, "Coin")
}

type StateLocked struct {
	BaseState
}

func NewStateLocked() StateLocked {
	return StateLocked{BaseState{StateName: "Locked"}}
}

func (s StateLocked) Pass(fsm *TurnstileFSM) {
	fsm.Actions.AlarmOn()
}

func (s StateLocked) Coin(fsm *TurnstileFSM) {
	fsm.State = NewStateUnlocked()
	fsm.Actions.AlarmOff()
	fsm.Actions.Unlock()
}

type StateUnlocked struct {
	BaseState
}

func NewStateUnlocked() StateLocked {
	return StateLocked{BaseState{StateName: "Unlocked"}}
}

func (s StateUnlocked) Pass(fsm *TurnstileFSM) {
	fsm.State = NewStateLocked()
	fsm.Actions.Lock()
}

func (s StateUnlocked) Coin(fsm *TurnstileFSM) {
	fsm.Actions.Thankyou()
}
