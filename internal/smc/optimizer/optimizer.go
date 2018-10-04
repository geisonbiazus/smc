package optimizer

import "github.com/geisonbiazus/smc/internal/smc/semantic"

type Optimizer struct {
	optimizedFSM FSM
	semanticFSM  *semantic.FSM
}

func New() *Optimizer {
	return &Optimizer{}
}

func (o *Optimizer) Optimize(fsm *semantic.FSM) FSM {
	o.optimizedFSM = FSM{}
	o.semanticFSM = fsm

	o.setEventsAndActions()
	o.setHeaders()
	o.optimizeStates()

	return o.optimizedFSM
}

func (o *Optimizer) setEventsAndActions() {
	o.optimizedFSM.Events = o.semanticFSM.Events
	o.optimizedFSM.Actions = o.semanticFSM.Actions
}

func (o *Optimizer) setHeaders() {
	o.optimizedFSM.Name = o.semanticFSM.Name
	o.optimizedFSM.ActionsClass = o.semanticFSM.ActionsClass
	o.optimizedFSM.InitialState = o.semanticFSM.InitialState.Name
}

func (o *Optimizer) optimizeStates() {
	for _, s := range o.semanticFSM.States {
		if !s.Abstract {
			o.setState(s)
		}
	}
}

func (o *Optimizer) setState(s *semantic.State) {
	state := State{Name: s.Name}
	o.setTransitions(&state, s)
	o.optimizedFSM.States = append(o.optimizedFSM.States, state)
}

func (o *Optimizer) setTransitions(state *State, semanticState *semantic.State) {
	for _, t := range semanticState.Transitions {
		transition := Transition{Event: t.Event, NextState: t.NextState.Name, Actions: t.Actions}
		state.Transitions = append(state.Transitions, transition)
	}

	for _, superState := range semanticState.SuperStates {
		o.setTransitions(state, superState)
	}
}
