package smc

type SemanticFSM struct {
	Errors []SemanticError
}

type SemanticError struct {
	Type ErrorType
}
