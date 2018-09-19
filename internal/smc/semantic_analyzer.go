package smc

type SemanticAnalyzer struct {
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{}
}

func (a *SemanticAnalyzer) Analyze(fsm FSMSyntax) *SemanticFSM {
	return &SemanticFSM{}
}
