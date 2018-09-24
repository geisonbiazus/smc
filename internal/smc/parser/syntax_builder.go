package parser

import "fmt"

type SyntaxBuilder struct {
	fsm           FSMSyntax
	currentName   string
	currentHeader Header
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
	b.lastHeader().Value = b.currentName
}

func (b *SyntaxBuilder) AddNewTransition() {
	b.fsm.Logic = append(
		b.fsm.Logic, Transition{StateSpec: StateSpec{Name: b.currentName}},
	)
}

func (b *SyntaxBuilder) AddNewAbstractTransition() {
	b.AddNewTransition()
	b.lastStateSpec().AbstractState = true
}

func (b *SyntaxBuilder) AddSuperState() {
	b.lastStateSpec().SuperStates = append(b.lastStateSpec().SuperStates, b.currentName)
}

func (b *SyntaxBuilder) AddEntryAction() {
	b.lastStateSpec().EntryActions = append(b.lastStateSpec().EntryActions, b.currentName)
}

func (b *SyntaxBuilder) AddExitAction() {
	b.lastStateSpec().ExitActions = append(b.lastStateSpec().ExitActions, b.currentName)
}

func (b *SyntaxBuilder) AddEmptyEvent() {
	b.lastTransition().SubTransitions = append(
		b.lastTransition().SubTransitions,
		SubTransition{Actions: []string{}},
	)
}

func (b *SyntaxBuilder) AddEvent() {
	b.AddEmptyEvent()
	b.lastSubTransition().Event = b.currentName
}

func (b *SyntaxBuilder) AddNextState() {
	b.lastSubTransition().NextState = b.currentName
}

func (b *SyntaxBuilder) AddAction() {
	b.lastSubTransition().Actions = append(
		b.lastSubTransition().Actions, b.currentName,
	)
}

func (b *SyntaxBuilder) Done() {
	b.fsm.Done = true
}

func (b *SyntaxBuilder) SyntaxError(line, pos int) {
	b.fsm.Errors = append(
		b.fsm.Errors,
		SyntaxError{Type: ErrorSyntax, LineNumber: line, Position: pos},
	)
}

func (b *SyntaxBuilder) ParseError(s State, e Event, line, pos int) {
	b.addError(ErrorParse, s, e, line, pos)
}

func (b *SyntaxBuilder) lastHeader() *Header {
	return &b.fsm.Headers[len(b.fsm.Headers)-1]
}

func (b *SyntaxBuilder) lastTransition() *Transition {
	return &b.fsm.Logic[len(b.fsm.Logic)-1]
}

func (b *SyntaxBuilder) lastStateSpec() *StateSpec {
	return &b.lastTransition().StateSpec
}

func (b *SyntaxBuilder) lastSubTransition() *SubTransition {
	return &b.lastTransition().
		SubTransitions[len(b.lastTransition().SubTransitions)-1]
}

func (b *SyntaxBuilder) addError(t ErrorType, s State, e Event, line, pos int) {
	msg := fmt.Sprintf("%s|%s", s, e)
	b.fsm.Errors = append(
		b.fsm.Errors,
		SyntaxError{Type: ErrorParse, Msg: msg, LineNumber: line, Position: pos},
	)
}
