package registry

import "github.com/eadydb/hubble/internal/app"

type Eureka struct{}

func NewEureka() *Eureka {
	return &Eureka{}
}

func (e *Eureka) Info(name string) *app.App {
	return nil
}
