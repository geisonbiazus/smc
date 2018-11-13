package golang

import "github.com/geisonbiazus/smc/internal/smc/generator/statepattern"

type Implementer struct{}

func NewImplementer() *Implementer {
	return &Implementer{}
}

func (i *Implementer) Implement(node statepattern.Node) string {
	return ""
}
