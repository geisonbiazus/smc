package semantic

import (
	"strings"

	"github.com/geisonbiazus/smc/internal/smc/parser"
)

type Analyzer struct {
	semanticFSM *FSM
	parsedFSM   parser.FSMSyntax
	stateCache  map[string]*State
	eventCache  map[string]bool
	actionCache map[string]bool
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(parsedFSM parser.FSMSyntax) *FSM {
	a.stateCache = map[string]*State{}
	a.eventCache = map[string]bool{}
	a.actionCache = map[string]bool{}
	a.semanticFSM = &FSM{}
	a.parsedFSM = parsedFSM

	a.addDefinedStates()
	a.setAndValidateHeaders()
	a.setAndValidateStates()
	a.checkForUnusedStates()
	a.checkForConflictingTransitions()

	return a.semanticFSM
}

func (a *Analyzer) addDefinedStates() {
	for _, t := range a.parsedFSM.Logic {
		a.addState(t.StateSpec)
	}
}

func (a *Analyzer) addState(spec parser.StateSpec) {
	state := a.findOrCreateState(spec.Name)
	state.Abstract = spec.AbstractState
	a.stateCache[spec.Name] = state
	a.semanticFSM.States = append(a.semanticFSM.States, state)
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
		a.semanticFSM.InitialState = markUsed(a.findAndValidateState(value))
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
	state := a.stateCache[t.StateSpec.Name]
	a.validateAbstractState(state, t)
	a.setEntryActions(state, t)
	a.setExitActions(state, t)
	a.setSuperStates(state, t)
	a.setTransitions(state, t)
}

func (a *Analyzer) validateAbstractState(state *State, t parser.Transition) {
	if state.Abstract != t.StateSpec.AbstractState {
		a.addError(ErrorAbstractStateRedefinedAsNonAbstract, state.Name)
	}
}

func (a *Analyzer) setEntryActions(state *State, t parser.Transition) {
	if len(state.EntryActions) > 0 {
		a.addError(ErrorEntryActionsAlreadyDefined, state.Name)
		return
	}
	state.EntryActions = t.StateSpec.EntryActions
	a.addActions(state.EntryActions)
}

func (a *Analyzer) setExitActions(state *State, t parser.Transition) {
	if len(state.ExitActions) > 0 {
		a.addError(ErrorExitActionsAlreadyDefined, state.Name)
		return
	}
	state.ExitActions = t.StateSpec.ExitActions
	a.addActions(state.ExitActions)
}

func (a *Analyzer) setSuperStates(state *State, t parser.Transition) {
	for _, name := range t.StateSpec.SuperStates {
		state.SuperStates = append(
			state.SuperStates,
			markUsed(a.findAndValidateSuperState(name)),
		)
	}
}

func (a *Analyzer) setTransitions(state *State, t parser.Transition) {
	for _, sub := range t.SubTransitions {
		if sub.Event != "" {
			a.addEvent(sub.Event)
			a.setTransition(state, sub)
			a.addActions(sub.Actions)
		}
	}
}

func (a *Analyzer) addEvent(name string) {
	if !a.eventCache[name] {
		a.eventCache[name] = true
		a.semanticFSM.Events = append(a.semanticFSM.Events, name)
	}
}

func (a *Analyzer) addActions(actions []string) {
	for _, action := range actions {
		if !a.actionCache[action] {
			a.actionCache[action] = true
			a.semanticFSM.Actions = append(a.semanticFSM.Actions, action)
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
		nextState = a.findAndValidateNextState(nextStateName)
	}
	return markUsed(nextState)
}

func (a *Analyzer) findAndValidateSuperState(name string) *State {
	if _, ok := a.stateCache[name]; !ok {
		a.addError(ErrorUndefinedSuperState, name)
	}
	return a.findOrCreateState(name)
}

func (a *Analyzer) findAndValidateNextState(name string) *State {
	state := a.findAndValidateState(name)

	if state.Abstract {
		a.addError(ErrorAbstractStateUsedAsNextState, name)
	}

	return state
}

func (a *Analyzer) findAndValidateState(name string) *State {
	if _, ok := a.stateCache[name]; !ok {
		a.addError(ErrorUndefinedState, name)
	}
	return a.findOrCreateState(name)
}

func (a *Analyzer) findOrCreateState(name string) *State {
	state, ok := a.stateCache[name]
	if !ok {
		state = &State{Name: name}
	}
	return state
}

func (a *Analyzer) checkForUnusedStates() {
	for _, state := range a.semanticFSM.States {
		if !state.Used {
			a.addWarning(ErrorUnusedState, state.Name)
		}
	}
}

func (a *Analyzer) checkForConflictingTransitions() {
	for _, state := range a.semanticFSM.States {
		a.checkForDuplicateTransition(state)
		a.checkForConflictingTransitionsOnSuperStates(state)
	}
}

func (a *Analyzer) checkForDuplicateTransition(state *State) {
	index := make(map[string]bool)
	for _, transition := range state.Transitions {
		if index[transition.Event] {
			a.addError(ErrorDuplicateTransition, state.Name+":"+transition.Event)
		}
		index[transition.Event] = true
	}
}

func (a *Analyzer) checkForConflictingTransitionsOnSuperStates(state *State) {
	index := make(map[string]bool)
	for _, super := range state.SuperStates {
		for _, transition := range super.Transitions {
			if index[transition.Event] {
				a.addError(ErrorConflictingSuperStates, state.Name+":"+transition.Event)
			}
			index[transition.Event] = true
		}
	}
}

func (a *Analyzer) addError(errorType ErrorType, element string) {
	a.semanticFSM.Errors = append(
		a.semanticFSM.Errors,
		Error{Type: errorType, Element: element},
	)
}

func (a *Analyzer) addWarning(errorType ErrorType, element string) {
	a.semanticFSM.Warnings = append(
		a.semanticFSM.Warnings,
		Error{Type: errorType, Element: element},
	)
}

func markUsed(s *State) *State {
	s.Used = true
	return s
}
