package smc

import (
	"io"

	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
)

type Error interface {
	String() string
}

type Compiler struct {
	input  io.Reader
	Errors []Error
}

func NewCompiler(input io.Reader) *Compiler {
	return &Compiler{
		input: input,
	}
}

func (c *Compiler) Compile() {
	parsedFSM := c.parseFSM()
	c.collectParseErrors(parsedFSM)

	semanticFSM := c.analyzeFSM(parsedFSM)
	c.collectSemanticErrors(semanticFSM)
}

func (c *Compiler) parseFSM() parser.FSMSyntax {
	builder := parser.NewSyntaxBuilder()
	psr := parser.NewParser(builder)
	lxr := lexer.NewLexer(psr)
	lxr.Lex(c.input)
	return builder.FSM()
}

func (c *Compiler) analyzeFSM(parsedFSM parser.FSMSyntax) *semantic.FSM {
	analyzer := semantic.NewAnalyzer()
	return analyzer.Analyze(parsedFSM)
}

func (c *Compiler) collectParseErrors(parsedFSM parser.FSMSyntax) {
	for _, err := range parsedFSM.Errors {
		c.Errors = append(c.Errors, err)
	}
}

func (c *Compiler) collectSemanticErrors(semanticFSM *semantic.FSM) {
	for _, err := range semanticFSM.Errors {
		c.Errors = append(c.Errors, err)
	}
}
