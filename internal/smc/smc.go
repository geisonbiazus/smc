package smc

import (
	"errors"
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
	input          io.Reader
	output         io.Writer
	Errors         []Error
	parsedFSM      parser.FSMSyntax
	semanticFSM    *semantic.FSM
	optimizedFSM   *optimizer.FSM
	node           statepattern.Node
	implementedFSM string
}

func NewCompiler(input io.Reader, output io.Writer) *Compiler {
	return &Compiler{
		input:  input,
		output: output,
	}
}

func (c *Compiler) Compile() error {
	if !c.parseFSM() {
		return CompileError
	}

	if !c.analyzeFSM() {
		return CompileError
	}

	c.optimizeFSM()
	c.generateFSM()
	c.implementFSM()
	c.writeImplementation()
	return nil
}

func (c *Compiler) parseFSM() bool {
	builder := parser.NewSyntaxBuilder()
	psr := parser.NewParser(builder)
	lxr := lexer.NewLexer(psr)
	lxr.Lex(c.input)

	c.parsedFSM = builder.FSM()
	c.collectParseErrors()
	return len(c.Errors) == 0
}

func (c *Compiler) collectParseErrors() {
	for _, err := range c.parsedFSM.Errors {
		c.Errors = append(c.Errors, err)
	}
}

func (c *Compiler) analyzeFSM() bool {
	analyzer := semantic.NewAnalyzer()
	c.semanticFSM = analyzer.Analyze(c.parsedFSM)
	c.collectSemanticErrors()
	return len(c.Errors) == 0
}

func (c *Compiler) collectSemanticErrors() {
	for _, err := range c.semanticFSM.Errors {
		c.Errors = append(c.Errors, err)
	}
}

func (c *Compiler) optimizeFSM() {
	opt := optimizer.New()
	c.optimizedFSM = opt.Optimize(c.semanticFSM)
}

func (c *Compiler) generateFSM() {
	generator := statepattern.NewNodeGenerator()
	c.node = generator.Generate(c.optimizedFSM)
}

func (c *Compiler) implementFSM() {
	impl := golang.NewImplementer("fsm")
	c.implementedFSM = impl.Implement(c.node)
}

func (c *Compiler) writeImplementation() {
	fmt.Fprint(c.output, c.implementedFSM)
}

var CompileError = errors.New("Compile error")
