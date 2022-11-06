package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/eadydb/hubble/pkg/kubernetes/client"
	"github.com/eadydb/hubble/pkg/output/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type PodWatcher interface {
	Register(receiver chan<- PodEvent)
	Deregister(receiver chan<- PodEvent)
	Start(ctx context.Context, kubeContext string, ns []string) (func(), error)
}

// podWatcher is a pod watcher for multiple namespaces.
type podWatcher struct {
	podSelector  PodSelector
	receivers    map[chan<- PodEvent]bool
	receiverLock sync.RWMutex
}

type PodEvent struct {
	Type watch.EventType
	Pod  *v1.Pod
}

func NewPodWatcher(podSelector PodSelector) PodWatcher {
	return &podWatcher{
		podSelector: podSelector,
		receivers:   make(map[chan<- PodEvent]bool),
	}
}

func (w *podWatcher) Register(receiver chan<- PodEvent) {
	w.receiverLock.Lock()
	w.receivers[receiver] = true
	w.receiverLock.Unlock()
}

func (w *podWatcher) Deregister(receiver chan<- PodEvent) {
	w.receiverLock.Lock()
	w.receivers[receiver] = false
	w.receiverLock.Unlock()
}

func (w *podWatcher) Start(ctx context.Context, kubeContext string, namespaces []string) (func(), error) {
	if len(w.receivers) == 0 {
		return func() {}, errors.New("no receiver was registered")
	}

	var watchers []watch.Interface
	stopWatchers := func() {
		for _, w := range watchers {
			w.Stop()
		}
	}

	kubeclient, err := client.Client(kubeContext)
	if err != nil {
		return func() {}, fmt.Errorf("getting k8s client: %w", err)
	}

	var forever int64 = 3600 * 24 * 365 * 100

	for _, ns := range namespaces {
		watcher, err := kubeclient.CoreV1().Pods(ns).Watch(context.Background(), metav1.ListOptions{
			TimeoutSeconds: &forever,
		})
		if err != nil {
			stopWatchers()
			return func() {}, fmt.Errorf("initializing pod watcher for %q : %w", ns, err)
		}

		watchers = append(watchers, watcher)

		go func() {
			l := log.Entry(ctx)
			defer l.Tracef("podWatcher: cease waiting")
			l.Tracef("podWatcher: waiting")
			for {
				select {
				case <-ctx.Done():
					l.Tracef("podWatcher: context canceled, returing")
					return
				case evt, ok := <-watcher.ResultChan():
					if !ok {
						l.Tracef("podWatcher: channel closed, returning")
						return
					}
					// If the event's type is "ERROR", log and continue.
					if evt.Type == watch.Error {
						l.Debugf("podWatcher: got unexpected event of type %s: %v", evt.Type, evt.Object)
						continue
					}

					// Grab thd pod from the event
					pod, ok := evt.Object.(*v1.Pod)
					if !ok {
						continue
					}

					if !w.podSelector.Select(pod) {
						continue
					}

					if log.IsTraceLevelEnabled() {
						st := fmt.Sprintf("podWatcher[%s/%s:%v] phase: %v", pod.Namespace, pod.Name, evt.Type, pod.Status.Phase)
						if len(pod.Status.Reason) > 0 {
							st += fmt.Sprintf("reason:%s", pod.Status.Reason)
						}
						for _, c := range append(pod.Status.InitContainerStatuses, pod.Status.ContainerStatuses...) {
							switch {
							case c.State.Waiting != nil:
								st += fmt.Sprintf("%s<waiting> ", c.Name)
							case c.State.Running != nil:
								st += fmt.Sprintf("%s<running> ", c.Name)
							case c.State.Terminated != nil:
								st += fmt.Sprintf("%s<terminated> ", c.Name)
							}
						}
						l.Trace(st)
					}

					l.Tracef("podWatcher: sending to all receivers")
					w.receiverLock.RLock()
					for receiver, open := range w.receivers {
						if open {
							l.Tracef("podWatcher: sending event type %v pod name %v namespace %v", evt.Type, pod.GetName(), pod.GetNamespace())
							receiver <- PodEvent{
								Type: evt.Type,
								Pod:  pod,
							}
						}
					}
					w.receiverLock.RUnlock()
					l.Tracef("podWatcher: done sending to all receivers")
				}
			}
		}()
	}
	return stopWatchers, nil
}
