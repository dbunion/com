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

	if cf, ok := e.Object.(*v1.ConfigMap); ok {
		object = convertToConfig(cf)
	}

	if ns, ok := e.Object.(*v1.Namespace); ok {
		object = convertToNamespace(ns)
	}

	if node, ok := e.Object.(*v1.Node); ok {
		object = convertToNode(node)
	}

	if pod, ok := e.Object.(*v1.Pod); ok {
		object = convertToPod(pod)
	}

	if rc, ok := e.Object.(*v1.ReplicationController); ok {
		object = convertToRC(rc)
	}

	if svc, ok := e.Object.(*v1.Service); ok {
		object = convertToService(svc)
	}

	if dpl, ok := e.Object.(*appsv1.Deployment); ok {
		object = convertToDeployment(dpl)
	}

	if rs, ok := e.Object.(*appsv1.ReplicaSet); ok {
		object = convertToReplicaSet(rs)
	}

	if sts, ok := e.Object.(*appsv1.StatefulSet); ok {
		object = convertToSTS(sts)
	}

	if ds, ok := e.Object.(*appsv1.DaemonSet); ok {
		object = convertToDaemonSet(ds)
	}

	event.Object = object
	return event
}
