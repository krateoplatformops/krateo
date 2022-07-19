package providers

import (
	"context"
	"fmt"
	"os"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	toolsWatch "k8s.io/client-go/tools/watch"
)

func DeleteInstalled(ctx context.Context, restConfig *rest.Config) error {
	sel, err := core.InstalledBySelector()
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	all, err := core.List(ctx, core.ListOpts{
		RESTConfig:    restConfig,
		GVR:           gvr,
		LabelSelector: sel.String(),
	})
	if err != nil {
		return err
	}

	dc, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	for _, el := range all {
		err := dc.Resource(gvr).
			Namespace(el.GetNamespace()).
			Delete(ctx, el.GetName(), metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}

		waitUntilDeleted(ctx, dc, el.GetName(), el.GetNamespace())
	}

	return nil
}

func waitUntilDeleted(ctx context.Context, dc dynamic.Interface, name, namespace string) error {
	req, err := labels.NewRequirement(core.PackageNameLabel, selection.Equals, []string{name})
	if err != nil {
		return err
	}

	sel := labels.NewSelector()
	sel = sel.Add(*req)

	watchFn := func(_ metav1.ListOptions) (watch.Interface, error) {
		timeoutSecs := int64(120)

		gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}

		return dc.Resource(gvr).
			Namespace(namespace).
			Watch(ctx, metav1.ListOptions{
				LabelSelector:  sel.String(),
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

		if et := event.Type; et == watch.Deleted {
			break
		}
	}

	return nil
}

type DeleteOpts struct {
	RESTConfig *rest.Config
	Object     *unstructured.Unstructured
}

func Delete(ctx context.Context, opts DeleteOpts) error {
	gvk := schema.FromAPIVersionAndKind(opts.Object.GetAPIVersion(), opts.Object.GetKind())
	mapping, err := core.FindGVR(opts.RESTConfig, &gvk)
	if err != nil {
		return err
	}

	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	// obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dc.Resource(mapping.Resource).Namespace(opts.Object.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dc.Resource(mapping.Resource)
	}

	err = dr.Delete(ctx, opts.Object.GetName(), metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
