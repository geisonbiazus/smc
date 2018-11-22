package golang

import (
	"strings"

	"github.com/geisonbiazus/smc/internal/smc/generator/statepattern"
)

type Implementer struct {
	pkg    string
	result string
}

func NewImplementer(pkg string) *Implementer {
	return &Implementer{
		pkg: pkg,
	}
}

func (i *Implementer) Implement(node statepattern.Node) string {
	i.result = ""

	if i.pkg != "" {
		i.result += "package " + i.pkg + "\n"
	}

	node.Accept(i)
	return i.result
}

func (i *Implementer) VisitStateInterfaceNode(node statepattern.StateInterfaceNode) {
	i.result += "\n"
	i.result += "type State interface {\n"

	for _, event := range node.Events {
		i.result += "  " + title(event) + "(fsm *" + title(node.FSMClassName) + ")\n"
	}

	i.result += "}\n"
}

func (i *Implementer) VisitActionsInterfaceNode(node statepattern.ActionsInterfaceNode) {
	i.result += "\n"
	i.result += "type Actions interface {\n"

	for _, action := range node.Actions {
		i.result += "  " + title(action) + "()\n"
	}

	i.result += "  UnhandledTransition(state string, event string)\n"
	i.result += "}\n"
}

func (i *Implementer) VisitFSMClassNode(node statepattern.FSMClassNode) {
	className := title(node.ClassName)

	i.result += "\n"
	i.result += "type " + className + " struct {\n"
	i.result += "  State State\n"
	i.result += "  Actions Actions\n"
	i.result += "}\n"
	i.result += "\n"
	i.result += "func New" + className + "(actions Actions) *" + className + " {\n"
	i.result += "  return &" + className + "{\n"
	i.result += "    Actions: actions,\n"
	i.result += "    State:   NewState" + title(node.InitialState) + "(),\n"
	i.result += "  }\n"
	i.result += "}\n"

	for _, methodNode := range node.EventMethods {
		methodNode.Accept(i)
	}
}

func (i *Implementer) VisitEventMethodNode(node statepattern.EventMethodNode) {
	i.result += "\n"
	i.result += "func (f *" + title(node.ClassName) + ") " + title(node.EventName) + "() {\n"
	i.result += "  f.State." + title(node.EventName) + "(f)\n"
	i.result += "}\n"
}

func (i *Implementer) VisitBaseStateClassNode(node statepattern.BaseStateClassNode) {
	i.result += "\n"
	i.result += "type BaseState struct {"
	i.result += "  StateName string\n"
	i.result += "}\n"

	for _, event := range node.Events {
		i.result += "\n"
		i.result += "func (b BaseState) " + title(event) + "(fsm *" + title(node.FSMClassName) + ") {\n"
		i.result += "  fsm.Actions.UnhandledTransition(b.StateName, \"" + event + "\")\n"
		i.result += "}\n"
	}
}

func (i *Implementer) VisitStateClassNode(node statepattern.StateClassNode) {
	i.result += "\n"
	i.result += "type State" + title(node.StateName) + " struct {\n"
	i.result += "  BaseState\n"
	i.result += "}\n"
	i.result += "\n"
	i.result += "func NewState" + title(node.StateName) + "() State" + title(node.StateName) + " {\n"
	i.result += "  return State" + title(node.StateName) + "{BaseState{StateName: \"" + node.StateName + "\"}}\n"
	i.result += "}\n"

	for _, method := range node.StateEventMethods {
		method.Accept(i)
	}
}
func (i *Implementer) VisitStateEventMethodNode(node statepattern.StateEventMethodNode) {
	i.result += "\n"
	i.result += "func (s State" + title(node.StateName) + ") " + title(node.EventName) + "(fsm *" + title(node.FSMClassName) + ") {\n"

	if node.NextState != "" {
		i.result += "  fsm.State = NewState" + title(node.NextState) + "()\n"
	}

	for _, action := range node.Actions {
		i.result += "  fsm.Actions." + title(action) + "()\n"
	}

	i.result += "}\n"
}

func title(s string) string {
	return strings.Title(s)
}
