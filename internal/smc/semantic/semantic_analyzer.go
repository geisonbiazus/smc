package semantic

import (
	"strings"

	"github.com/geisonbiazus/smc/internal/smc/parser"
)

type SemanticAnalyzer struct {
	semanticFSM *SemanticFSM
	fsm         parser.FSMSyntax
}

func NewSemanticAnalyzer() *SemanticAnalyzer {
	return &SemanticAnalyzer{}
}

func (a *SemanticAnalyzer) Analyze(fsm parser.FSMSyntax) *SemanticFSM {
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
		switch strings.ToLower(header.Name) {
		case "fsm":
			if !a.isDuplicate(a.semanticFSM.Name, ErrorDuplicateHeader) {
				a.semanticFSM.Name = header.Value
			}
		case "actions":
			if !a.isDuplicate(a.semanticFSM.Actions, ErrorDuplicateHeader) {
				a.semanticFSM.Actions = header.Value
			}
		case "initial":
			if !a.isDuplicate(a.semanticFSM.Initial, ErrorDuplicateHeader) {
				a.semanticFSM.Initial = header.Value
			}
		default:
			a.addError(ErrorInvalidHeader)
		}
	}
}

func (a *SemanticAnalyzer) isDuplicate(value string, errorType ErrorType) bool {
	if value == "" {
		a.addError(errorType)
		return false
	}
	return true
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
