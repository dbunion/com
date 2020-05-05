package k8s

import (
	"github.com/dbunion/com/scheduler"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/watch"
)

type k8sWatcher struct {
	watcher watch.Interface
	done    chan struct{}
}

// NewWatcher - create new k8s watcher
func NewWatcher(w watch.Interface) scheduler.Interface {
	return &k8sWatcher{
		watcher: w,
		done:    make(chan struct{}),
	}
}

// Stops watching. Will close the channel returned by ResultChan(). Releases
// any resources used by the watch.
func (w *k8sWatcher) Stop() {
	w.watcher.Stop()
	w.done <- struct{}{}
}

// Returns a chan which will receive all the events. If an error occurs
// or Stop() is called, this channel will be closed, in which case the
// watch should be completely cleaned up.
func (w *k8sWatcher) ResultChan() <-chan scheduler.WatchEvent {
	event := make(chan scheduler.WatchEvent)
	go func() {
		for {
			select {
			case e := <-w.watcher.ResultChan():
				if e.Type == watch.Error {
					continue
				}
				event <- w.processEvent(e)
			case <-w.done:
				return
			}
		}
	}()

	return event
}

func (w *k8sWatcher) processEvent(e watch.Event) scheduler.WatchEvent {
	var event scheduler.WatchEvent
	var object scheduler.Object
	event.Type = scheduler.EventType(e.Type)

	versionKind := e.Object.GetObjectKind().GroupVersionKind()
	switch versionKind.Kind {
	case "ConfigMap":
		cf, ok := e.Object.(*v1.ConfigMap)
		if ok {
			object = cf
		}
	case "Namespace":
		ns, ok := e.Object.(*v1.Namespace)
		if ok {
			object = ns
		}
	case "Node":
		node, ok := e.Object.(*v1.Node)
		if ok {
			object = node
		}
	case "Pod":
		pod, ok := e.Object.(*v1.Pod)
		if ok {
			object = pod
		}
	case "ReplicationController":
		rc, ok := e.Object.(*v1.ReplicationController)
		if ok {
			object = rc
		}
	case "Service":
		svc, ok := e.Object.(*v1.Service)
		if ok {
			object = svc
		}
	case "Deployment":
		dpl, ok := e.Object.(*appsv1.Deployment)
		if ok {
			object = dpl
		}
	case "ReplicaSet":
		rs, ok := e.Object.(*appsv1.ReplicaSet)
		if ok {
			object = rs
		}
	case "StatefulSet":
		sts, ok := e.Object.(*appsv1.StatefulSet)
		if ok {
			object = sts
		}
	case "DaemonSet":
		ds, ok := e.Object.(*appsv1.DaemonSet)
		if ok {
			object = ds
		}
	}

	event.Object = object
	return event
}
