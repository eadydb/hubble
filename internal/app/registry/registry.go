package registry

import "github.com/eadydb/probe-pilot/internal/app"

type Registry interface {
	Info(name string) *app.App
}

type RegistryClient struct {
	Registry Registry
	Type     string
}

func NewRegistryClient(reg string) *RegistryClient {
	return &RegistryClient{
		Type: reg,
	}
}

func (r *RegistryClient) Info(name string) *app.App {
	var reg Registry
	switch r.Type {
	case "EUREKA":
		reg = NewEureka()
	case "KUBERNETE":
		reg = NewKubePod()
	case "NACOS":
		reg = NewNacos()
	}
	r.Registry = reg
	return r.Registry.Info(name)
}
