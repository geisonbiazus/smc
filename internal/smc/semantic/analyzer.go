package semantic

import (
	"strings"

	"github.com/geisonbiazus/smc/internal/smc/parser"
)

type Analyzer struct {
	semanticFSM *FSM
	parsedFSM   parser.FSMSyntax
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(parsedFSM parser.FSMSyntax) *FSM {
	a.semanticFSM = NewFSM()
	a.parsedFSM = parsedFSM

	a.setAndValidateHeaders()
	a.setStates()

	return a.semanticFSM
}

func (a *Analyzer) setAndValidateHeaders() {
	a.setHeaders()
	a.validateRequiredHeaders()
}

func (a *Analyzer) setHeaders() {
	for _, header := range a.parsedFSM.Headers {
		switch strings.ToLower(header.Name) {
		case "fsm":
			a.setName(header.Value)
		case "actions":
			a.setActionsClass(header.Value)
		case "initial":
			a.setInitialState(header.Value)
		default:
			a.addError(ErrorInvalidHeader)
		}
	}
}

func (a *Analyzer) setName(value string) {
	if !a.isDuplicate(a.semanticFSM.Name, ErrorDuplicateHeader) {
		a.semanticFSM.Name = value
	}
}

func (a *Analyzer) setActionsClass(value string) {
	if !a.isDuplicate(a.semanticFSM.ActionsClass, ErrorDuplicateHeader) {
		a.semanticFSM.ActionsClass = value
	}
}

func (a *Analyzer) setInitialState(value string) {
	if !a.isDuplicateState(a.semanticFSM.InitialState, ErrorDuplicateHeader) {
		a.semanticFSM.InitialState = a.findOrCreateState(value)
	}
}

func (a *Analyzer) isDuplicate(value string, errorType ErrorType) bool {
	if value != "" {
		a.addError(errorType)
		return true
	}
	return false
}

func (a *Analyzer) isDuplicateState(value *State, errorType ErrorType) bool {
	if value != nil {
		a.addError(errorType)
		return true
	}
	return false
}

func (a *Analyzer) validateRequiredHeaders() {
	if a.semanticFSM.Name == "" {
		a.addError(ErrorNoFSM)
	}

	if a.semanticFSM.InitialState == nil {
		a.addError(ErrorNoInitial)
	}
}

func (a *Analyzer) setStates() {
	for _, parsedTransition := range a.parsedFSM.Logic {
		a.setState(parsedTransition)
	}
}

func (a *Analyzer) setState(t parser.Transition) {
	state := a.findOrCreateState(t.StateSpec.Name)
	state.Abstract = t.StateSpec.AbstractState
	state.EntryActions = t.StateSpec.EntryActions
	state.ExitActions = t.StateSpec.ExitActions
	a.setSuperStates(state, t)
	a.setTransitions(state, t)
}

func (a *Analyzer) setSuperStates(state *State, t parser.Transition) {
	for _, name := range t.StateSpec.SuperStates {
		state.SuperStates = append(state.SuperStates, a.findOrCreateState(name))
	}
}

func (a *Analyzer) setTransitions(state *State, t parser.Transition) {
	for _, sub := range t.SubTransitions {
		a.setTransition(state, sub)
	}
}

func (a *Analyzer) setTransition(state *State, sub parser.SubTransition) {
	transition := Transition{
		Event:     sub.Event,
		NextState: a.resolveNextState(state, sub.NextState),
		Actions:   sub.Actions,
	}
	state.Transitions = append(state.Transitions, transition)
}

func (a *Analyzer) resolveNextState(state *State, nextStateName string) *State {
	nextState := state
	if nextStateName != "" {
		nextState = a.findOrCreateState(nextStateName)
	}
	return nextState
}

func (a *Analyzer) findOrCreateState(name string) *State {
	state, ok := a.semanticFSM.States[name]
	if !ok {
		state = &State{Name: name}
		a.semanticFSM.States[name] = state
	}
	return state
}

func (a *Analyzer) addError(errorType ErrorType) {
	a.semanticFSM.Errors = append(a.semanticFSM.Errors, Error{Type: errorType})
}
