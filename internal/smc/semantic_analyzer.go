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

	a.setAndValidateHeaders()

	return a.semanticFSM
}

func (a *SemanticAnalyzer) setAndValidateHeaders() {
	a.setHeaders()
	a.validateRequiredHeaders()
}

func (a *SemanticAnalyzer) setHeaders() {
	for _, header := range a.fsm.Headers {
		switch header.Name {
		case "FSM":
			a.semanticFSM.Name = header.Value
		case "Actions":
			a.semanticFSM.Actions = header.Value
		case "Initial":
			a.semanticFSM.Initial = header.Value
		default:
			a.addError(ErrorInvalidHeader)
		}
	}
}

func (a *SemanticAnalyzer) validateRequiredHeaders() {
	if a.semanticFSM.Name == "" {
		a.addError(ErrorNoFSM)
	}

	if a.semanticFSM.Initial == "" {
		a.addError(ErrorNoInitial)
	}
}

func (a *SemanticAnalyzer) addError(errorType ErrorType) {
	a.semanticFSM.Errors = append(a.semanticFSM.Errors, SemanticError{Type: errorType})
}
