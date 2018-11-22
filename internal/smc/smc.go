package smc

import (
	"fmt"
	"io"

	"github.com/geisonbiazus/smc/internal/smc/generator/statepattern"
	"github.com/geisonbiazus/smc/internal/smc/implementers/golang"
	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/optimizer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
	"github.com/geisonbiazus/smc/internal/smc/semantic"
)

type Error interface {
	String() string
}

type Compiler struct {
	input  io.Reader
	output io.Writer
	Errors []Error
}

func NewCompiler(input io.Reader, output io.Writer) *Compiler {
	return &Compiler{
		input:  input,
		output: output,
	}
}

func (c *Compiler) Compile() {
	parsedFSM := c.parseFSM()
	c.collectParseErrors(parsedFSM)

	if len(parsedFSM.Errors) > 0 {
		return
	}

	semanticFSM := c.analyzeFSM(parsedFSM)
	c.collectSemanticErrors(semanticFSM)

	if len(semanticFSM.Errors) > 0 {
		return
	}

	optimizedFSM := c.optimizeFSM(semanticFSM)
	node := c.generateFSM(optimizedFSM)
	c.implementFSM(node)
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

func (c *Compiler) optimizeFSM(semanticFSM *semantic.FSM) *optimizer.FSM {
	opt := optimizer.New()
	return opt.Optimize(semanticFSM)
}

func (c *Compiler) generateFSM(optimizedFSM *optimizer.FSM) statepattern.Node {
	generator := statepattern.NewNodeGenerator()
	return generator.Generate(optimizedFSM)
}

func (c *Compiler) implementFSM(node statepattern.Node) {
	impl := golang.NewImplementer("fsm")
	result := impl.Implement(node)
	fmt.Fprint(c.output, result)
}
