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

	o.setHeaders()
	o.setStates()

	return o.optimizedFSM
}

func (o *Optimizer) setHeaders() {
	o.optimizedFSM.Name = o.semanticFSM.Name
	o.optimizedFSM.ActionsClass = o.semanticFSM.ActionsClass
	o.optimizedFSM.InitialState = o.semanticFSM.InitialState.Name
}

func (o *Optimizer) setStates() {
	for _, s := range o.semanticFSM.States {
		state := State{Name: s.Name}

		for _, t := range s.Transitions {
			transition := Transition{Event: t.Event, NextState: t.NextState.Name, Actions: t.Actions}
			state.Transitions = append(state.Transitions, transition)
			o.optimizedFSM.Events = append(o.optimizedFSM.Events, t.Event)

			for _, action := range t.Actions {
				o.optimizedFSM.Actions = append(o.optimizedFSM.Actions, action)
			}
		}

		o.optimizedFSM.States = append(o.optimizedFSM.States, state)
	}
}
