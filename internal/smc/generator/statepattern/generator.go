package statepattern

import "github.com/geisonbiazus/smc/internal/smc/optimizer"

type StatePattern struct {
	fsm *optimizer.FSM
}

func NewStatePattern() *StatePattern {
	return &StatePattern{}
}

func (g *StatePattern) Generate(fsm *optimizer.FSM) Node {
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

func (g *StatePattern) stateInterfaceNode() Node {
	return StateInterfaceNode{
		FSMClassName: g.fsm.Name,
		States:       g.stateNames(),
	}
}

func (g *StatePattern) stateNames() []string {
	states := []string{}
	for _, state := range g.fsm.States {
		states = append(states, state.Name)
	}
	return states
}

func (g *StatePattern) fsmClassNode() Node {
	return FSMClassNode{
		ClassName:    g.fsm.Name,
		InitialState: g.fsm.InitialState,
		EventMethods: g.eventMethodNodes(),
	}
}

func (g *StatePattern) eventMethodNodes() []Node {
	nodes := []Node{}
	for _, event := range g.fsm.Events {
		eventNode := EventMethodNode{ClassName: g.fsm.Name, EventName: event}
		nodes = append(nodes, eventNode)
	}
	return nodes
}

func (g *StatePattern) actionsInterfaceNode() Node {
	return ActionsInterfaceNode{
		Actions: g.fsm.Actions,
	}
}

func (g *StatePattern) baseStateClassNode() Node {
	return BaseStateClassNode{
		Events: g.fsm.Events,
	}
}

func (g *StatePattern) stateClassNodes() Node {
	return CompositeNode(g.stateClassNodeList())
}

func (g *StatePattern) stateClassNodeList() []Node {
	nodes := []Node{}

	for _, state := range g.fsm.States {
		nodes = append(nodes, g.stateClassNode(state))
	}

	return nodes
}

func (g *StatePattern) stateClassNode(state *optimizer.State) Node {
	return StateClassNode{
		StateName:         state.Name,
		StateEventMethods: g.stateEventMethodNodes(state),
	}
}

func (g *StatePattern) stateEventMethodNodes(state *optimizer.State) []Node {
	nodes := []Node{}
	for _, transition := range state.Transitions {
		nodes = append(nodes, g.stateEventMethodNode(state, transition))
	}
	return nodes
}

func (g *StatePattern) stateEventMethodNode(
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
