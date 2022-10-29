package registry

import "github.com/eadydb/probe-pilot/internal/app"

type Nacos struct{}

func NewNacos() *Nacos {
	return &Nacos{}
}

func (n *Nacos) Info(name string) *app.App {
	return nil
}
