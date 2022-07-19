package core

import (
	"context"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

type WatchOpts struct {
	RESTConfig *rest.Config
	GVR        schema.GroupVersionResource
	Selector   labels.Selector
	Namespace  string
	StopFunc   func(et watch.EventType, obj *unstructured.Unstructured) (bool, error)
}

func Watch(ctx context.Context, opts WatchOpts) error {
	watchFn := func(_ metav1.ListOptions) (watch.Interface, error) {
		timeoutSecs := int64(120)

		dc, err := dynamic.NewForConfig(opts.RESTConfig)
		if err != nil {
			return nil, err
		}

		listOpts := metav1.ListOptions{
			TimeoutSeconds: &timeoutSecs,
		}
		if opts.Selector != nil {
			listOpts.LabelSelector = opts.Selector.String()
		}

		return dc.Resource(opts.GVR).Namespace(opts.Namespace).Watch(ctx, listOpts)
	}

	// create a `RetryWatcher` using initial version "1" and our specialized watcher
	rw, err := toolsWatch.NewRetryWatcher("1", &cache.ListWatch{WatchFunc: watchFn})
	if err != nil {
		return err
	}

	defer func() {
		if x := recover(); x != nil {
			fmt.Fprintf(os.Stderr, "recoverd from run time panic: %v", x)
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

		if opts.StopFunc == nil {
			break
		}

		obj, ok := event.Object.(*unstructured.Unstructured)
		if !ok {
			return fmt.Errorf("invalid type '%T'", event.Object)
		}

		exit, err := opts.StopFunc(event.Type, obj)
		if err != nil {
			return err
		}
		if exit {
			break
		}
	}

	return nil
}
