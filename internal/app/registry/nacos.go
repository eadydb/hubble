package registry

import "github.com/eadydb/hubble/internal/app"

type Nacos struct{}

func NewNacos() *Nacos {
	return &Nacos{}
}

func (n *Nacos) Info(name string) *app.App {
	return nil
}
