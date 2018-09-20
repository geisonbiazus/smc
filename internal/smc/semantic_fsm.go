package smc

type SemanticFSM struct {
	Errors []SemanticError
}

type SemanticError struct {
	Type ErrorType
}

const (
	ErrorNoFSM         ErrorType = "NO_FSM"
	ErrorNoInitial     ErrorType = "NO_INITIAL"
	ErrorInvalidHeader ErrorType = "INVALID_HEADER"
)
