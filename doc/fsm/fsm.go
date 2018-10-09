package fsm

type State interface {
	Pass(fsm *TurnstileFSM)
	Coin(fsm *TurnstileFSM)
}

type Actions interface {
	Alarm()
	Lock()
	Unlock()
	Thankyou()
}

type TurnstileFSM struct {
	State   State
	Actions Actions
}

func (t *TurnstileFSM) Pass() {
	t.State.Pass(t)
}

func (t *TurnstileFSM) Coin() {
	t.State.Coin(t)
}

type StateLocked struct{}

func (s StateLocked) Pass(fsm *TurnstileFSM) {
	fsm.Actions.Alarm()
}

func (s StateLocked) Coin(fsm *TurnstileFSM) {
	fsm.State = StateUnlocked{}
	fsm.Actions.Unlock()
}

type StateUnlocked struct{}

func (s StateUnlocked) Pass(fsm *TurnstileFSM) {
	fsm.State = StateLocked{}
	fsm.Actions.Lock()
}

func (s StateUnlocked) Coin(fsm *TurnstileFSM) {
	fsm.Actions.Thankyou()
}
