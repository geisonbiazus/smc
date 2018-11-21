package smc

import (
	"io"

	"github.com/geisonbiazus/smc/internal/smc/lexer"
	"github.com/geisonbiazus/smc/internal/smc/parser"
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
	fsm := c.parseFSM()
	for _, err := range fsm.Errors {
		c.Errors = append(c.Errors, err)
	}
}

func (c *Compiler) parseFSM() parser.FSMSyntax {
	builder := parser.NewSyntaxBuilder()
	psr := parser.NewParser(builder)
	lxr := lexer.NewLexer(psr)
	lxr.Lex(c.input)
	return builder.FSM()
}
