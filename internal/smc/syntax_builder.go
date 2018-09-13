package smc

import "fmt"

type SyntaxBuilder struct {
	fsm         FSMSyntax
	currentName string
}

func NewSyntaxBuilder() *SyntaxBuilder {
	return &SyntaxBuilder{}
}

func (b *SyntaxBuilder) FSM() FSMSyntax {
	return b.fsm
}

func (b *SyntaxBuilder) SetName(name string) {
	b.currentName = name
}

func (b *SyntaxBuilder) NewHeader() {
	b.fsm.Headers = append(b.fsm.Headers, Header{Name: b.currentName})
}

func (b *SyntaxBuilder) AddHeaderValue() {
	b.fsm.Headers[len(b.fsm.Headers)-1].Value = b.currentName
}

func (b *SyntaxBuilder) HeaderError(s State, e Event, line, pos int) {
	b.addError(ErrorHeader, s, e, line, pos)
}

func (b *SyntaxBuilder) addError(t ErrorType, s State, e Event, line, pos int) {
	msg := fmt.Sprintf("%s|%s", s, e)
	b.fsm.Errors = append(
		b.fsm.Errors,
		SyntaxError{Type: ErrorHeader, Msg: msg, LineNumber: line, Position: pos},
	)
}
