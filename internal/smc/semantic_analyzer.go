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

	a.validateRequiredHeader(ErrorNoFSM, "FSM")
	a.validateRequiredHeader(ErrorNoInitial, "Initial")
	a.validateAvailableHeaders(ErrorInvalidHeader, "FSM", "Initial", "Actions")

	return a.semanticFSM
}

func (a *SemanticAnalyzer) validateRequiredHeader(errorType ErrorType, name string) {
	for _, header := range a.fsm.Headers {
		if header.Name == name {
			return
		}
	}
	a.semanticFSM.Errors = append(a.semanticFSM.Errors, SemanticError{Type: errorType})
}

func (a *SemanticAnalyzer) validateAvailableHeaders(errorType ErrorType, availableHeaders ...string) {
	for _, header := range a.fsm.Headers {
		found := false
		for _, available := range availableHeaders {
			if header.Name == available {
				found = true
			}
		}
		if !found {
			a.semanticFSM.Errors = append(a.semanticFSM.Errors, SemanticError{Type: errorType})
		}
	}
}
