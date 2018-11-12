package statepattern

type Visitor interface {
	VisitStateInterfaceNode(node StateInterfaceNode)
	VisitActionsInterfaceNode(node ActionsInterfaceNode)
	VisitFSMClassNode(node FSMClassNode)
	VisitEventMethodNode(node EventMethodNode)
	VisitBaseStateClassNode(node BaseStateClassNode)
	VisitStateClassNode(node StateClassNode)
	VisitStateEventMethodNode(node StateEventMethodNode)
}

type Node interface {
	Accept(v Visitor)
}

type CompositeNode []Node

func (n CompositeNode) Accept(v Visitor) {}

type StateInterfaceNode struct {
	States       []string
	FSMClassName string
}

func (n StateInterfaceNode) Accept(v Visitor) {}

type ActionsInterfaceNode struct {
	Actions []string
}

func (n ActionsInterfaceNode) Accept(v Visitor) {}

type FSMClassNode struct {
	InitialState string
	ClassName    string
	ActionsClass string
	EventMethods []Node
}

func (n FSMClassNode) Accept(v Visitor) {}

type EventMethodNode struct {
	ClassName string
	EventName string
}

func (n EventMethodNode) Accept(v Visitor) {}

type BaseStateClassNode struct {
	Events []string
}

func (n BaseStateClassNode) Accept(v Visitor) {}

type StateClassNode struct {
	StateName         string
	StateEventMethods []Node
}

func (n StateClassNode) Accept(v Visitor) {}

type StateEventMethodNode struct {
	StateName    string
	FSMClassName string
	EventName    string
	NextState    string
	Actions      []string
}

func (n StateEventMethodNode) Accept(v Visitor) {}
