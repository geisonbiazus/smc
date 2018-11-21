package smc

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	t.Run("With syntax errors", func(t *testing.T) {
		assertContainsError(t, compileFSM("& a:b {}"),
			parser.SyntaxError{
				Type: parser.ErrorSyntax, LineNumber: 1, Position: 1,
			},
		)
	})

	t.Run("With parse errors", func(t *testing.T) {
		assertContainsError(t, compileFSM("a:b:c {}"),
			parser.SyntaxError{
				Type: parser.ErrorParse, LineNumber: 1, Position: 4, Msg: "HEADER|COLON",
			},
		)
	})
}

func compileFSM(input string) *Compiler {
	compiler := NewCompiler(bytes.NewBufferString(input))
	compiler.Compile()
	return compiler
}

func assertContainsError(t *testing.T, compiler *Compiler, err Error) {
	t.Helper()
	assert.Contains(t, compiler.Errors, err)
}
