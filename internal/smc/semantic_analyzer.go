package smc

type SemanticAnalyzer struct {
	semanticFSM *SemanticFSM
	fsm         FSMSyntax
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{}
}

func (a *SemanticAnalyzer) Analyze(fsm FSMSyntax) *SemanticFSM {
	a.semanticFSM = &SemanticFSM{}
	a.fsm = fsm

	a.validateHeader("FSM")

	return a.semanticFSM
}

func (a *SemanticAnalyzer) validateHeader(name string) {
	for _, header := range a.fsm.Headers {
		if header.Name == name {
			return
		}
	}
	a.semanticFSM.Errors = append(a.semanticFSM.Errors, SemanticError{Type: ErrorNoFSM})
}
