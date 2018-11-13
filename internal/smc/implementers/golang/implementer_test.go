package golang

import (
	"bytes"
	"testing"

	"github.com/geisonbiazus/smc/internal/smc/generator/statepattern"
	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/optimizer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
	"github.com/stretchr/testify/assert"
)

func TestImplementer(t *testing.T) {
	assertImplementedFSM(t,
		"FSM: fsm Initial: state { state event state action }",
		``,
	)
}

func assertImplementedFSM(t *testing.T, input, expected string) {
	node := implementFSM(input)
	assert.Equal(t, expected, node)
}

func implementFSM(input string) string {
	builder := parser.NewSyntaxBuilder()
	psr := parser.NewParser(builder)
	lxr := lexer.NewLexer(psr)
	lxr.Lex(bytes.NewBufferString(input))

	parsedFSM := builder.FSM()

	analyzer := semantic.NewAnalyzer()
	semanticFSM := analyzer.Analyze(parsedFSM)

	opt := optimizer.New()
	optimizedFSM := opt.Optimize(semanticFSM)

	gen := statepattern.NewNodeGenerator()
	node := gen.Generate(optimizedFSM)

	implementer := NewImplementer()

	return implementer.Implement(node)
}
