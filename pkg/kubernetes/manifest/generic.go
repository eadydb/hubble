package manifest

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
)

type Generic struct {
	NonResource
}

// List returns a collection of resources.
func (g *Generic) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	return nil, nil
}

// Get returns a given resource.
func (g *Generic) Get(ctx context.Context, path string) (runtime.Object, error) {
	return nil, nil
}

// Describe describes a resource.
func (g *Generic) Describe(path string) (string, error) {
	return "", nil
}

// ToYAML returns a resource yaml.
func (g *Generic) ToYAML(path string, showManaged bool) (string, error) {
	return "", nil
}
