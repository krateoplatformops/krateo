package providers

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo/internal/catalog"
	"github.com/krateoplatformops/krateo/internal/controllerconfigs"
	"github.com/krateoplatformops/krateo/internal/core"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

type InstallOpts struct {
	RESTConfig *rest.Config
	Info       *catalog.PackageInfo
	Namespace  string
}

func Install(ctx context.Context, opts InstallOpts) error {
	// fetch the manifest
	data, err := catalog.FetchManifest(opts.Info)
	if err != nil {
		return err
	}

	// decode the YAML
	obj, gvk, err := core.DecodeYAML(data)
	if err != nil {
		return err
	}

	if !isCrossplaneProvider(gvk) {
		return fmt.Errorf("%s is not a provider", obj.GetName())
	}

	// create a controller config
	ccf, err := controllerconfigs.Create(context.TODO(), controllerconfigs.CreateOpts{
		Info:       opts.Info,
		RESTConfig: opts.RESTConfig,
	})
	if err != nil {
		return err
	}

	// update the provider with the controller config reference
	controllerConfigRef := map[string]interface{}{
		"name": ccf.GetName(),
	}
	err = unstructured.SetNestedField(obj.Object, controllerConfigRef, "spec", "controllerConfigRef")
	if err != nil {
		return err
	}

	// install the provider
	err = core.Apply(ctx, core.ApplyOpts{RESTConfig: opts.RESTConfig, Object: obj, GVK: gvk})
	if err != nil {
		return err
	}

	// wait for it
	return waitUntilProviderIsReady(ctx, opts.RESTConfig, opts.Info.Name, opts.Namespace)
}

func waitUntilProviderIsReady(ctx context.Context, restConfig *rest.Config, name, namespace string) error {
	req, err := labels.NewRequirement(core.PackageNameLabel, selection.Equals, []string{name})
	if err != nil {
		return err
	}

	sel := labels.NewSelector()
	sel = sel.Add(*req)

	stopFn := func(et watch.EventType, obj *unstructured.Unstructured) (bool, error) {
		pod := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &pod)
		if err != nil {
			return false, err
		}

		for _, cond := range pod.Status.Conditions {
			if (cond.Type == corev1.PodReady) && (cond.Status == corev1.ConditionTrue) {
				return true, nil
			}
		}

		return false, nil
	}

	return core.Watch(ctx, core.WatchOpts{
		RESTConfig: restConfig,
		GVR:        schema.GroupVersionResource{Version: "v1", Resource: "pods"},
		Namespace:  namespace,
		Selector:   sel,
		StopFunc:   stopFn,
	})
}

func isCrossplaneProvider(gvk *schema.GroupVersionKind) bool {
	return gvk.Group == "pkg.crossplane.io"
}
