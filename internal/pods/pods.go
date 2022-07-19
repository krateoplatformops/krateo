package pods

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

type WatchOpts struct {
	Selector  labels.Selector
	Namespace string
	StopFunc  func(corev1.PodCondition) bool
}

func Watch(dc dynamic.Interface, opts WatchOpts) error {
	watchFn := func(_ metav1.ListOptions) (watch.Interface, error) {
		timeoutSecs := int64(120)

		gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}

		return dc.Resource(gvr).
			Namespace(opts.Namespace).
			Watch(context.Background(), metav1.ListOptions{
				LabelSelector:  opts.Selector.String(),
				TimeoutSeconds: &timeoutSecs,
			})
	}

	// create a `RetryWatcher` using initial version "1" and our specialized watcher
	rw, err := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFn})
	if err != nil {
		return err
	}
	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintf(os.Stderr, "run time panic: %v", x)
		}
		rw.Stop()
	}()

	// process incoming event notifications
	for {
		// grab the event object
		event, ok := <-rw.ResultChan()
		if !ok {
			return fmt.Errorf("closed channel")
		}

		if et := event.Type; et != watch.Added && et != watch.Modified {
			continue
		}

		obj, ok := event.Object.(*unstructured.Unstructured)
		if !ok {
			return fmt.Errorf("invalid type '%T'", event.Object)
		}
		pod := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.
			FromUnstructured(obj.UnstructuredContent(), &pod)
		if err != nil {
			return err
		}

		for _, cond := range pod.Status.Conditions {
			if opts.StopFunc(cond) {
				return nil
			}
		}
	}
}
