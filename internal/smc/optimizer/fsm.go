package optimizer

type FSM struct {
	Name         string
	ActionsClass string
	InitialState string
	Events       []string
	Actions      []string
	States       []State
}

type State struct {
	Name        string
	Transitions []Transition
}

type Transition struct {
	Event     string
	NextState string
	Actions   []string
}
