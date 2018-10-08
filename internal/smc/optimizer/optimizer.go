package optimizer

import (
	"github.com/geisonbiazus/smc/internal/smc/semantic"
)

type Optimizer struct {
	optimizedFSM *FSM
	semanticFSM  *semantic.FSM
}

func New() *Optimizer {
	return &Optimizer{}
}

func (o *Optimizer) Optimize(fsm *semantic.FSM) *FSM {
	o.optimizedFSM = &FSM{}
	o.semanticFSM = fsm

	o.setEventsAndActions()
	o.setHeaders()
	o.optimizeStates()
	o.optimizeEntryActions()
	o.eliminateDuplicatedActions()

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
			o.optmizeState(s)
		}
	}
}

func (o *Optimizer) optmizeState(s *semantic.State) {
	state := &State{Name: s.Name}
	o.optimizeTransitions(state, s, make(map[string]bool))
	o.optimizeExitActions(state, s)
	o.optimizedFSM.States = append(o.optimizedFSM.States, state)
}

func (o *Optimizer) optimizeTransitions(
	state *State, semanticState *semantic.State, definedEvents map[string]bool,
) {

	for _, t := range semanticState.Transitions {
		if !definedEvents[t.Event] {
			o.addTransition(state, t)
			definedEvents[t.Event] = true
		}
	}

	for _, superState := range semanticState.SuperStates {
		o.optimizeTransitions(state, superState, definedEvents)
	}
}

func (o *Optimizer) addTransition(state *State, t semantic.Transition) {
	transition := &Transition{
		Event:     t.Event,
		NextState: o.resolveNextState(t),
		Actions:   t.Actions,
	}

	state.Transitions = append(state.Transitions, transition)
}

func (o *Optimizer) resolveNextState(t semantic.Transition) string {
	if t.NextState != nil {
		return t.NextState.Name
	}

	return ""
}

func (o *Optimizer) optimizeExitActions(state *State, semanticState *semantic.State) {
	actions := o.getExitActionsRecursively(semanticState)
	if len(actions) > 0 {
		for _, t := range state.Transitions {
			t.Actions = append(t.Actions, actions...)
		}
	}
}

func (o *Optimizer) getExitActionsRecursively(s *semantic.State) []string {
	actions := []string{}
	actions = append(actions, s.ExitActions...)
	for _, super := range s.SuperStates {
		actions = append(actions, o.getExitActionsRecursively(super)...)
	}
	return actions
}

func (o *Optimizer) optimizeEntryActions() {
	for _, semanticState := range o.semanticFSM.States {
		actions := o.getEntryActionsRecursively(semanticState)

		if len(actions) > 0 {
			o.optimizeEntryActionsOfState(semanticState, actions)
		}
	}
}

func (o *Optimizer) getEntryActionsRecursively(s *semantic.State) []string {
	actions := []string{}
	actions = append(actions, s.EntryActions...)
	for _, super := range s.SuperStates {
		actions = append(actions, o.getEntryActionsRecursively(super)...)
	}
	return actions
}

func (o *Optimizer) optimizeEntryActionsOfState(semanticState *semantic.State, actions []string) {
	for _, optmizedState := range o.optimizedFSM.States {
		for _, transition := range optmizedState.Transitions {
			if transition.NextState == semanticState.Name {
				transition.Actions = append(transition.Actions, actions...)
			}
		}
	}
}

func (o *Optimizer) eliminateDuplicatedActions() {
	for _, s := range o.optimizedFSM.States {
		for _, t := range s.Transitions {
			t.Actions = unique(t.Actions)
		}
	}
}

func unique(list []string) []string {
	result := []string{}
	cache := map[string]bool{}

	for _, item := range list {
		if !cache[item] {
			cache[item] = true
			result = append(result, item)
		}
	}
	return result
}
