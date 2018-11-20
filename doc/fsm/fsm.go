package turnstile

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
	State   State
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

type BaseState struct {
	StateName string
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
