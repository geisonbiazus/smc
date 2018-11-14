package statepattern

import "github.com/geisonbiazus/smc/internal/smc/optimizer"

type NodeGenerator struct {
	fsm *optimizer.FSM
}

func NewNodeGenerator() *NodeGenerator {
	return &NodeGenerator{}
}

func (g *NodeGenerator) Generate(fsm *optimizer.FSM) Node {
	g.fsm = fsm
	return CompositeNode(
		[]Node{
			g.stateInterfaceNode(),
			g.actionsInterfaceNode(),
			g.fsmClassNode(),
			g.baseStateClassNode(),
			g.stateClassNodes(),
		},
	)
}

func (g *NodeGenerator) stateInterfaceNode() Node {
	return StateInterfaceNode{
		FSMClassName: g.fsm.Name,
		Events:       g.fsm.Events,
	}
}

func (g *NodeGenerator) fsmClassNode() Node {
	return FSMClassNode{
		ClassName:    g.fsm.Name,
		InitialState: g.fsm.InitialState,
		EventMethods: g.eventMethodNodes(),
	}
}

func (g *NodeGenerator) eventMethodNodes() []Node {
	nodes := []Node{}
	for _, event := range g.fsm.Events {
		eventNode := EventMethodNode{ClassName: g.fsm.Name, EventName: event}
		nodes = append(nodes, eventNode)
	}
	return nodes
}

func (g *NodeGenerator) actionsInterfaceNode() Node {
	return ActionsInterfaceNode{
		Actions: g.fsm.Actions,
	}
}

func (g *NodeGenerator) baseStateClassNode() Node {
	return BaseStateClassNode{
		Events: g.fsm.Events,
	}
}

func (g *NodeGenerator) stateClassNodes() Node {
	return CompositeNode(g.stateClassNodeList())
}

func (g *NodeGenerator) stateClassNodeList() []Node {
	nodes := []Node{}

	for _, state := range g.fsm.States {
		nodes = append(nodes, g.stateClassNode(state))
	}

	return nodes
}

func (g *NodeGenerator) stateClassNode(state *optimizer.State) Node {
	return StateClassNode{
		StateName:         state.Name,
		StateEventMethods: g.stateEventMethodNodes(state),
	}
}

func (g *NodeGenerator) stateEventMethodNodes(state *optimizer.State) []Node {
	nodes := []Node{}
	for _, transition := range state.Transitions {
		nodes = append(nodes, g.stateEventMethodNode(state, transition))
	}
	return nodes
}

func (g *NodeGenerator) stateEventMethodNode(
	state *optimizer.State, transition *optimizer.Transition,
) Node {
	return StateEventMethodNode{
		FSMClassName: g.fsm.Name,
		StateName:    state.Name,
		EventName:    transition.Event,
		NextState:    transition.NextState,
		Actions:      transition.Actions,
	}
}
