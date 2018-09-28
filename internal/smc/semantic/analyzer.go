package semantic

import (
	"strings"

	"github.com/geisonbiazus/smc/internal/smc/parser"
)

type Analyzer struct {
	semanticFSM   *FSM
	parsedFSM     parser.FSMSyntax
	definedStates map[string]bool
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{definedStates: make(map[string]bool)}
}

func (a *Analyzer) Analyze(parsedFSM parser.FSMSyntax) *FSM {
	a.semanticFSM = NewFSM()
	a.parsedFSM = parsedFSM

	a.addDefinedStates()
	a.setAndValidateHeaders()
	a.setAndValidateStates()

	return a.semanticFSM
}

func (a *Analyzer) addDefinedStates() {
	for _, t := range a.parsedFSM.Logic {
		a.addState(t.StateSpec.Name)
	}
}

func (a *Analyzer) addState(name string) {
	a.semanticFSM.States[name] = a.findOrCreateState(name)
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
			a.addError(ErrorInvalidHeader, header.Name)
		}
	}
}

func (a *Analyzer) setName(value string) {
	if !a.isDuplicate(a.semanticFSM.Name, ErrorDuplicateHeader, "FSM") {
		a.semanticFSM.Name = value
	}
}

func (a *Analyzer) setActionsClass(value string) {
	if !a.isDuplicate(a.semanticFSM.ActionsClass, ErrorDuplicateHeader, "Actions") {
		a.semanticFSM.ActionsClass = value
	}
}

func (a *Analyzer) setInitialState(value string) {
	if !a.isDuplicateState(a.semanticFSM.InitialState, ErrorDuplicateHeader, "Initial") {
		a.semanticFSM.InitialState = a.findAndValidateState(value)
	}
}

func (a *Analyzer) isDuplicate(value string, errorType ErrorType, element string) bool {
	if value != "" {
		a.addError(errorType, element)
		return true
	}
	return false
}

func (a *Analyzer) isDuplicateState(value *State, errorType ErrorType, element string) bool {
	if value != nil {
		a.addError(errorType, element)
		return true
	}
	return false
}

func (a *Analyzer) validateRequiredHeaders() {
	if a.semanticFSM.Name == "" {
		a.addError(ErrorNoFSM, "FSM")
	}

	if a.semanticFSM.InitialState == nil {
		a.addError(ErrorNoInitial, "Initial")
	}
}

func (a *Analyzer) setAndValidateStates() {
	for _, parsedTransition := range a.parsedFSM.Logic {
		a.setAndValidateState(parsedTransition)
	}
}

func (a *Analyzer) setAndValidateState(t parser.Transition) {
	state := a.semanticFSM.States[t.StateSpec.Name]
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
		if sub.Event != "" {
			a.setTransition(state, sub)
		}
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
		nextState = a.findAndValidateState(nextStateName)
	}
	return nextState
}

func (a *Analyzer) findAndValidateState(name string) *State {
	if _, ok := a.semanticFSM.States[name]; !ok {
		a.addError(ErrorUndefinedState, name)
	}
	return a.findOrCreateState(name)
}

func (a *Analyzer) findOrCreateState(name string) *State {
	state, ok := a.semanticFSM.States[name]
	if !ok {
		state = &State{Name: name}
	}
	return state
}

func (a *Analyzer) addError(errorType ErrorType, element string) {
	a.semanticFSM.Errors = append(
		a.semanticFSM.Errors,
		Error{Type: errorType, Element: element},
	)
}
