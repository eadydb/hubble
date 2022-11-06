package kubernetes

import (
	"sync"

	v1 "k8s.io/api/core/v1"
)

// PodSelector is used to choose which pods to log.
type PodSelector interface {
	Select(pod *v1.Pod) bool
}

// ImageList implements PodSelector based on a list of images names.
type ImageList struct {
	sync.RWMutex
	names map[string]bool
}

// NewImageList creates a new ImageList.
func NewImageList() *ImageList {
	return &ImageList{
		names: make(map[string]bool),
	}
}

// Add adds an image to the list.
func (l *ImageList) Add(image string) {
	l.Lock()
	l.names[image] = true
	l.Unlock()
}

// Select returns true if one of the pod's images is in the list.
func (l *ImageList) Select(pod *v1.Pod) bool {
	l.RLock()
	defer l.RUnlock()

	for _, container := range append(pod.Spec.InitContainers, pod.Spec.Containers...) {
		if l.names[container.Image] {
			return true
		}
	}

	return false
}

// TODO(nkubala): remove this when podSelector is moved entirely into Deployer
type ImageListMux []*ImageList

func (l ImageListMux) Select(pod *v1.Pod) bool {
	for _, selector := range l {
		if selector.Select(pod) {
			return true
		}
	}
	return false
}
