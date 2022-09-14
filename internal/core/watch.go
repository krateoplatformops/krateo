package core

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

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

var (
	ErrWatcherTimeout = errors.New("watcher timed out")
)

type StopFunc func(et watch.EventType, obj *unstructured.Unstructured) (bool, error)

type WatchOpts struct {
	RESTConfig *rest.Config
	GVR        schema.GroupVersionResource
	Selector   labels.Selector
	Namespace  string
	Timeout    time.Duration
	StopFn     StopFunc
}

func Watch(ctx context.Context, opts WatchOpts) error {
	watchFn := func(_ metav1.ListOptions) (watch.Interface, error) {
		timeoutSecs := int64(180)

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

		if len(opts.Namespace) == 0 {
			return dc.Resource(opts.GVR).Watch(ctx, listOpts)
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
	if opts.Timeout <= 0 {
		opts.Timeout = 3 * time.Minute
	}
	timer := time.NewTimer(opts.Timeout)
	return watchOnce(rw, opts.StopFn, timer)
}

func watchOnce(w watch.Interface, stopFn StopFunc, timer *time.Timer) error {
	var err error

loop:
	for {
		select {
		// grab the event object
		case event, ok := <-w.ResultChan():
			if !ok {
				err = fmt.Errorf("closed channel")
				break loop
			}

			if stopFn == nil {
				break loop
			}

			obj, ok := event.Object.(*unstructured.Unstructured)
			if !ok {
				err = fmt.Errorf("invalid type '%T'", event.Object)
				break loop
			}

			exit, err := stopFn(event.Type, obj)
			if err != nil {
				return err
			}
			if exit {
				return nil
			}
		case <-timer.C:
			return ErrWatcherTimeout
		}
	}

	return err
}
