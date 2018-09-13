package smc

type FSMSyntax struct {
	Headers []Header
	Logic   []Transition
	Errors  []SyntaxError
	Done    bool
}

type Header struct {
	Name  string
	Value string
}

type Transition struct {
	StateSpec      StateSpec
	SubTransitions []SubTransition
}

type StateSpec struct {
	Name          string
	SuperStates   []string
	EntryActions  []string
	ExitActions   []string
	AbstractState bool
}

type SubTransition struct {
	Event     string
	NextState string
	Actions   []string
}

type SyntaxError struct {
	Type       ErrorType
	Msg        string
	LineNumber int
	Position   int
}

type ErrorType string

const (
	ErrorHeader          ErrorType = "HEADER"
	ErrorState           ErrorType = "STATE"
	ErrorTransition      ErrorType = "TRANSITION"
	ErrorTransitionGroup ErrorType = "TRANSITION_GROUP"
	ErrorEnd             ErrorType = "END"
	ErrorSyntax          ErrorType = "SYNTAX"
)
