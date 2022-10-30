package registry

import "github.com/eadydb/hubble/internal/app"

type KubePod struct{}

func NewKubePod() *KubePod {
	return &KubePod{}
}

func (k *KubePod) Info(name string) *app.App {
	return nil
}
