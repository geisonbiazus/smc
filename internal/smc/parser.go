package smc

type Builder interface {
	SetName(name string)
	NewHeader()
	AddHeaderValue()
	AddNewTransition()
	AddEvent()
	AddNextState()
	AddAction()
	Done()
	HeaderError(s State, e Event, line, pos int)
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
}

func (p *Parser) ClosedParen(line, pos int) {
}

func (p *Parser) OpenAngle(line, pos int) {
}

func (p *Parser) ClosedAngle(line, pos int) {
}

func (p *Parser) Dash(line, pos int) {
	p.HandleEvent(EventDash, line, pos)
}

func (p *Parser) Name(name string, line, pos int) {
	p.Builder.SetName(name)
	p.HandleEvent(EventName, line, pos)
}

func (p *Parser) Error(line, pos int) {
}

type State string
type Event string
type Action func(Builder)

type transition struct {
	currentState State
	event        Event
	newState     State
	action       Action
}

var transitions = []transition{
	{StateHeader, EventName, StateHeaderColon, func(b Builder) { b.NewHeader() }},
	{StateHeader, EventOpenBrace, StateTransitionGroup, NoAction},
	{StateHeaderColon, EventColon, StateHeaderValue, NoAction},
	{StateHeaderValue, EventName, StateHeader, func(b Builder) { b.AddHeaderValue() }},

	{StateTransitionGroup, EventName, StateNewTransition, func(b Builder) { b.AddNewTransition() }},
	{StateTransitionGroup, EventClosedBrace, StateEnd, func(b Builder) { b.Done() }},
	{StateNewTransition, EventName, StateSingleEvent, func(b Builder) { b.AddEvent() }},

	{StateSingleEvent, EventName, StateNextState, func(b Builder) { b.AddNextState() }},
	{StateNextState, EventName, StateTransitionGroup, func(b Builder) { b.AddAction() }},
	{StateNextState, EventOpenBrace, StateActionGroup, NoAction},
	{StateNextState, EventDash, StateTransitionGroup, NoAction},
	{StateActionGroup, EventName, StateActionGroup, func(b Builder) { b.AddAction() }},
	{StateActionGroup, EventClosedBrace, StateTransitionGroup, NoAction},
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
	p.Builder.HeaderError(p.state, event, line, pos)
}

const (
	StateHeader          State = "HEADER"
	StateHeaderColon     State = "HEADER_COLON"
	StateHeaderValue     State = "HEADER_VALUE"
	StateTransitionGroup State = "TRANSITION_GROUP"
	StateNewTransition   State = "NEW_TRANSITION"
	StateSingleEvent     State = "SINGLE_EVENT"
	StateNextState       State = "NEXT_STATE"
	StateActionGroup     State = "ACTION_GROUP"
	StateEnd             State = "END"

	EventName        Event = "NAME"
	EventColon       Event = "COLON"
	EventOpenBrace   Event = "OPEN_BRACE"
	EventClosedBrace Event = "CLOSED_BRACE"
	EventDash        Event = "DASH"
)

var NoAction Action = func(Builder) {}
