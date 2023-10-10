package manifest

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Pod represents a pod resource.
type Pod struct {
	Resource
}

// IsHappy check for happy deployments.
func (p *Pod) IsHappy(po v1.Pod) bool {
	for _, c := range po.Status.Conditions {
		if c.Status == v1.ConditionFalse {
			return false
		}
	}
	return true
}

// Get returns a resource instance if found, else an error.
func (p *Pod) Get(ctx context.Context, path string) (runtime.Object, error) {
	return nil, nil
}

// List returns a collection of nodes.
func (p *Pod) List(ctx context.Context, ns string) ([]runtime.Object, error) {
	return nil, nil
}

// Containers returns all container names on pod.
func (p *Pod) Containers(path string, includeInit bool) ([]string, error) {
	return nil, nil
}

// Pod returns a pod victim by name.
func (p *Pod) Pod(fqn string) (string, error) {
	return fqn, nil
}

// GetInstance returns a pod instance.
func (p *Pod) GetInstance(fqn string) (*v1.Pod, error) {
	return nil, nil
}

// MetaFQN returns a fully qualified resource name.
func MetaFQN(m metav1.ObjectMeta) string {
	if m.Namespace == "" {
		return m.Name
	}

	return FQN(m.Namespace, m.Name)
}

// FQN returns a fully qualified resource name.
func FQN(ns, n string) string {
	if ns == "" {
		return n
	}
	return ns + "/" + n
}

func extractFQN(o runtime.Object) string {
	return ""
}

// GetPodSpec returns a pod spec given a resource.
func (p *Pod) GetPodSpec(path string) (*v1.PodSpec, error) {
	pod, err := p.GetInstance(path)
	if err != nil {
		return nil, err
	}
	podSpec := pod.Spec
	return &podSpec, nil
}

func (p *Pod) isControlled(path string) (string, bool, error) {
	pod, err := p.GetInstance(path)
	if err != nil {
		return "", false, err
	}
	references := pod.GetObjectMeta().GetOwnerReferences()
	if len(references) > 0 {
		return fmt.Sprintf("%s/%s", references[0].Kind, references[0].Name), true, nil
	}
	return "", false, nil
}

// GetDefaultContainer returns a container name if specified in an annotation.
func GetDefaultContainer(m metav1.ObjectMeta, spec v1.PodSpec) (string, bool) {
	return "", false
}
