package manifest

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
)

type Factory interface {
	// Get fetch a give resource.
	Get(gvr, path string, wait bool, sel labels.Selector) (runtime.Object, error)

	// List fetch a collection of resources.
	List(gvr, ns string, wait bool, sel labels.Selector) ([]runtime.Object, error)

	// ForResource fetch an informer for a given resource.
	ForResource(ns, gvr string) (informers.GenericInformer, error)

	// CanForResource fetch an informer for a given resource if authorized
	CanForResource(ns, gvr string, verbs []string) (informers.GenericInformer, error)

	// WaitForCacheSync synchronize the cache.
	WaitForCacheSync()
}
