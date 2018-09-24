package semantic

type SemanticFSM struct {
	Errors  []SemanticError
	Name    string
	Actions string
	Initial string
}

type SemanticError struct {
	Type ErrorType
}

type ErrorType string

const (
	ErrorNoFSM           ErrorType = "NO_FSM"
	ErrorNoInitial       ErrorType = "NO_INITIAL"
	ErrorInvalidHeader   ErrorType = "INVALID_HEADER"
	ErrorDuplicateHeader ErrorType = "DUPLICATE_HEADER"
)
