package semantic

type FSM struct {
	Errors       []Error
	Warnings     []Error
	Name         string
	InitialState *State
	States       []*State
	Events       []string
	Actions      []string
}

type State struct {
	Name         string
	Abstract     bool
	Used         bool
	SuperStates  []*State
	EntryActions []string
	ExitActions  []string
	Transitions  []Transition
}

type Transition struct {
	Event     string
	NextState *State
	Actions   []string
}

type Error struct {
	Type    ErrorType
	Element string
}

type ErrorType string

const (
	ErrorNoFSM                               ErrorType = "NO_FSM"
	ErrorNoInitial                           ErrorType = "NO_INITIAL"
	ErrorInvalidHeader                       ErrorType = "INVALID_HEADER"
	ErrorDuplicateHeader                     ErrorType = "DUPLICATE_HEADER"
	ErrorNoTransitions                       ErrorType = "NO_TRANSITIONS"
	ErrorUndefinedState                      ErrorType = "UNDEFINED_STATE"
	ErrorUndefinedSuperState                 ErrorType = "UNDEFINED_SUPER_STATE"
	ErrorEntryActionsAlreadyDefined          ErrorType = "ENTRY_ACTIONS_ALREADY_DEFINED"
	ErrorExitActionsAlreadyDefined           ErrorType = "EXIT_ACTIONS_ALREADY_DEFINED"
	ErrorAbstractStateRedefinedAsNonAbstract ErrorType = "ABSTRACT_STATE_REDEFINED_AS_NON_ABSTRACT"
	ErrorAbstractStateUsedAsNextState        ErrorType = "ABSTRACT_STATE_USED_AS_NEXT_STATE"
	ErrorUnusedState                         ErrorType = "UNUSED_STATE"
	ErrorDuplicateTransition                 ErrorType = "DUPLICATE_TRANSITION"
	ErrorConflictingSuperStates              ErrorType = "CONFLICTING_SUPER_STATES"
)
