package smc

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
