package parser

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
	ErrorParse  ErrorType = "PARSE"
	ErrorSyntax ErrorType = "SYNTAX"
)
