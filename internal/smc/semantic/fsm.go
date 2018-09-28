package semantic

type FSM struct {
	Errors       []Error
	Name         string
	ActionsClass string
	InitialState *State
	States       map[string]*State
	Events       map[string]bool
	Actions      map[string]bool
}

func NewFSM() *FSM {
	return &FSM{
		States: make(map[string]*State),
	}
}

type State struct {
	Name         string
	Abstract     bool
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
	ErrorEntryActionsAlreadyDefined          ErrorType = "ENTRY_ACTIONS_ALREADY_DEFINED"
	ErrorExitActionsAlreadyDefined           ErrorType = "EXIT_ACTIONS_ALREADY_DEFINED"
	ErrorAbstractStateRedefinedAsNonAbstract ErrorType = "ABSTRACT_STATE_REDEFINED_AS_NON_ABSTRACT"
)
