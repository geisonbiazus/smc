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

type State struct {
	Name         string
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
	Type ErrorType
}

type ErrorType string

const (
	ErrorNoFSM           ErrorType = "NO_FSM"
	ErrorNoInitial       ErrorType = "NO_INITIAL"
	ErrorInvalidHeader   ErrorType = "INVALID_HEADER"
	ErrorDuplicateHeader ErrorType = "DUPLICATE_HEADER"
)