package claims

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ModuleOpts struct {
	RESTConfig *rest.Config
	Data       map[string]interface{}
}

func ApplyCoreModule(ctx context.Context, opts ModuleOpts) error {
	gvk := getGroupVersionKind()

	obj := &unstructured.Unstructured{}
	obj.SetKind(gvk.Kind)
	obj.SetAPIVersion(gvk.GroupVersion().String())
	obj.SetName("core")
	obj.SetLabels(map[string]string{
		core.InstalledByLabel: core.InstalledByValue,
	})
	err := unstructured.SetNestedField(obj.Object, opts.Data, "spec")
	if err != nil {
		return err
	}

	return core.Apply(ctx, core.ApplyOpts{
		RESTConfig: opts.RESTConfig,
		GVK:        gvk,
		Object:     obj,
	})
}

func getGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "modules.krateo.io",
		Version: "v1",
		Kind:    "Core",
	}
}

func getGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "modules.krateo.io",
		Version:  "v1",
		Resource: "core",
	}
}
