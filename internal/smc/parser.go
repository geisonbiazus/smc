package smc

type Builder interface {
	SetName(name string)
	NewHeader()
	AddHeaderValue()
	AddNewTransition()
	AddNewAbstractTransition()
	AddSuperState()
	AddEntryAction()
	AddExitAction()
	AddEmptyEvent()
	AddEvent()
	AddNextState()
	AddAction()
	Done()
	SyntaxError(line, pos int)
	ParseError(s State, e Event, line, pos int)
}

type Parser struct {
	Builder Builder
	state   State
}

func NewParser(builder Builder) *Parser {
	return &Parser{Builder: builder, state: StateHeader}
}

func (p *Parser) OpenBrace(line, pos int) {
	p.HandleEvent(EventOpenBrace, line, pos)
}

func (p *Parser) ClosedBrace(line, pos int) {
	p.HandleEvent(EventClosedBrace, line, pos)
}

func (p *Parser) Colon(line, pos int) {
	p.HandleEvent(EventColon, line, pos)
}

func (p *Parser) OpenParen(line, pos int) {
	p.HandleEvent(EventOpenParen, line, pos)
}

func (p *Parser) ClosedParen(line, pos int) {
	p.HandleEvent(EventClosedParen, line, pos)
}

func (p *Parser) OpenAngle(line, pos int) {
	p.HandleEvent(EventOpenAngle, line, pos)
}

func (p *Parser) ClosedAngle(line, pos int) {
	p.HandleEvent(EventClosedAngle, line, pos)
}

func (p *Parser) Dash(line, pos int) {
	p.HandleEvent(EventDash, line, pos)
}

func (p *Parser) Name(name string, line, pos int) {
	p.Builder.SetName(name)
	p.HandleEvent(EventName, line, pos)
}

func (p *Parser) Error(line, pos int) {
	p.Builder.SyntaxError(line, pos)
}

func (p *Parser) End(line, pos int) {
	p.Builder.Done()
	p.HandleEvent(EventEnd, line, pos)
}

type State string
type Event string

type transition struct {
	currentState State
	event        Event
	newState     State
	action       func(Builder)
}

var transitions = []transition{
	{StateHeader, EventName, StateHeaderColon, func(b Builder) { b.NewHeader() }},
	{StateHeader, EventOpenBrace, StateTransitionGroup, NoAction},
	{StateHeaderColon, EventColon, StateHeaderValue, NoAction},
	{StateHeaderValue, EventName, StateHeader, func(b Builder) { b.AddHeaderValue() }},

	{StateTransitionGroup, EventName, StateNewTransition, func(b Builder) { b.AddNewTransition() }},
	{StateTransitionGroup, EventClosedBrace, StateEnd, NoAction},
	{StateTransitionGroup, EventOpenParen, StateSuperState, NoAction},
	{StateSuperState, EventName, StateSuperStateName, func(b Builder) { b.AddNewAbstractTransition() }},
	{StateSuperStateName, EventClosedParen, StateNewTransition, NoAction},
	{StateNewTransition, EventName, StateSingleEvent, func(b Builder) { b.AddEvent() }},
	{StateNewTransition, EventDash, StateSingleEvent, func(b Builder) { b.AddEmptyEvent() }},
	{StateNewTransition, EventOpenBrace, StateSubTransitionGroup, NoAction},
	{StateNewTransition, EventColon, StateStateBase, NoAction},
	{StateNewTransition, EventOpenAngle, StateEntryAction, NoAction},
	{StateNewTransition, EventClosedAngle, StateExitAction, NoAction},
	{StateStateBase, EventName, StateNewTransition, func(b Builder) { b.AddSuperState() }},
	{StateEntryAction, EventName, StateNewTransition, func(b Builder) { b.AddEntryAction() }},
	{StateExitAction, EventName, StateNewTransition, func(b Builder) { b.AddExitAction() }},
	{StateEnd, EventEnd, StateEnd, NoAction},

	{StateSingleEvent, EventName, StateNextState, func(b Builder) { b.AddNextState() }},
	{StateSingleEvent, EventDash, StateNextState, NoAction},
	{StateNextState, EventName, StateTransitionGroup, func(b Builder) { b.AddAction() }},
	{StateNextState, EventOpenBrace, StateActionGroup, NoAction},
	{StateNextState, EventDash, StateTransitionGroup, NoAction},
	{StateActionGroup, EventName, StateActionGroup, func(b Builder) { b.AddAction() }},
	{StateActionGroup, EventClosedBrace, StateTransitionGroup, NoAction},

	{StateSubTransitionGroup, EventClosedBrace, StateTransitionGroup, NoAction},
	{StateSubTransitionGroup, EventName, StateSubTransitionEvent, func(b Builder) { b.AddEvent() }},
	{StateSubTransitionGroup, EventDash, StateSubTransitionEvent, func(b Builder) { b.AddEmptyEvent() }},
	{StateSubTransitionEvent, EventName, StateSubTransitionNextState, func(b Builder) { b.AddNextState() }},
	{StateSubTransitionEvent, EventDash, StateSubTransitionNextState, NoAction},
	{StateSubTransitionNextState, EventName, StateSubTransitionGroup, func(b Builder) { b.AddAction() }},
	{StateSubTransitionNextState, EventDash, StateSubTransitionGroup, NoAction},
	{StateSubTransitionNextState, EventOpenBrace, StateSubTransitionActionGroup, NoAction},
	{StateSubTransitionActionGroup, EventClosedBrace, StateSubTransitionGroup, NoAction},
	{StateSubTransitionActionGroup, EventName, StateSubTransitionActionGroup, func(b Builder) { b.AddAction() }},
}

func (p *Parser) HandleEvent(event Event, line, pos int) {
	for _, t := range transitions {
		if t.currentState == p.state && t.event == event {
			p.state = t.newState
			t.action(p.Builder)
			return
		}
	}
	p.HandleEventError(event, line, pos)
}

func (p *Parser) HandleEventError(event Event, line, pos int) {
	p.Builder.ParseError(p.state, event, line, pos)
}

const (
	StateHeader                   State = "HEADER"
	StateHeaderColon              State = "HEADER_COLON"
	StateHeaderValue              State = "HEADER_VALUE"
	StateTransitionGroup          State = "TRANSITION_GROUP"
	StateNewTransition            State = "NEW_TRANSITION"
	StateSingleEvent              State = "SINGLE_EVENT"
	StateNextState                State = "NEXT_STATE"
	StateActionGroup              State = "ACTION_GROUP"
	StateSubTransitionGroup       State = "STATE_SUB_TRANSITION_GROUP"
	StateSubTransitionEvent       State = "STATE_SUB_TRANSITION_EVENT"
	StateSubTransitionNextState   State = "STATE_SUB_TRANSITION_NEXT_STATE"
	StateSubTransitionActionGroup State = "STATE_SUB_TRANSITION_ACTION_GROUP"
	StateSuperState               State = "SUPER_STATE"
	StateSuperStateName           State = "SUPER_STATE_NAME"
	StateStateBase                State = "STATE_BASE"
	StateEntryAction              State = "ENTRY_ACTION"
	StateExitAction               State = "EXIT_ACTION"
	StateEnd                      State = "END"

	EventName        Event = "NAME"
	EventColon       Event = "COLON"
	EventOpenBrace   Event = "OPEN_BRACE"
	EventClosedBrace Event = "CLOSED_BRACE"
	EventDash        Event = "DASH"
	EventOpenParen   Event = "OPEN_PAREN"
	EventClosedParen Event = "CLOSED_PAREN"
	EventOpenAngle   Event = "OPEN_ANGLE"
	EventClosedAngle Event = "CLOSED_ANGLE"
	EventEnd         Event = "END"
)

var NoAction = func(Builder) {}
