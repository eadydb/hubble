package manifest

import (
	"context"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
)

type Resource struct {
	Generic
}

// List returns a collection of resources.
func (r *Resource) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	strLabel, _ := ctx.Value(KeyLabels).(string)
	lsel := labels.Everything()
	if strLabel != "" {
		if sel, err := labels.Parse(strLabel); err == nil {
			lsel = sel
		}
	}

	return r.GetFactory().List(r.gvr.String(), ns, false, lsel)
}

// Get returns a resource instance if found, else an error.
func (r *Resource) Get(_ context.Context, path string) (runtime.Object, error) {
	return r.GetFactory().Get(r.gvr.String(), path, true, labels.Everything())
}

// ToYAML returns a resource yaml.
func (r *Resource) ToYAML(path string, showManaged bool) (string, error) {
	return "", nil
}
