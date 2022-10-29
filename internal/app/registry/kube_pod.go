package registry

import "github.com/eadydb/probe-pilot/internal/app"

type KubePod struct{}

func NewKubePod() *KubePod {
	return &KubePod{}
}

func (k *KubePod) Info(name string) *app.App {
	return nil
}
