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

func (n CompositeNode) Accept(v Visitor) {
	for _, node := range n {
		node.Accept(v)
	}
}

type StateInterfaceNode struct {
	Events       []string
	FSMClassName string
}

func (n StateInterfaceNode) Accept(v Visitor) {
	v.VisitStateInterfaceNode(n)
}

type ActionsInterfaceNode struct {
	Actions []string
}

func (n ActionsInterfaceNode) Accept(v Visitor) {
	v.VisitActionsInterfaceNode(n)
}

type FSMClassNode struct {
	InitialState string
	ClassName    string
	ActionsClass string
	EventMethods []Node
}

func (n FSMClassNode) Accept(v Visitor) {
	v.VisitFSMClassNode(n)
}

type EventMethodNode struct {
	ClassName string
	EventName string
}

func (n EventMethodNode) Accept(v Visitor) {
	v.VisitEventMethodNode(n)
}

type BaseStateClassNode struct {
	Events []string
}

func (n BaseStateClassNode) Accept(v Visitor) {
	v.VisitBaseStateClassNode(n)
}

type StateClassNode struct {
	StateName         string
	StateEventMethods []Node
}

func (n StateClassNode) Accept(v Visitor) {
	v.VisitStateClassNode(n)
}

type StateEventMethodNode struct {
	StateName    string
	FSMClassName string
	EventName    string
	NextState    string
	Actions      []string
}

func (n StateEventMethodNode) Accept(v Visitor) {
	v.VisitStateEventMethodNode(n)
}
