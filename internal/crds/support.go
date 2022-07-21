package crds

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func PatchAndDelete(ctx context.Context, restConfig *rest.Config, el *unstructured.Unstructured) error {
	gvk := schema.FromAPIVersionAndKind(el.GetAPIVersion(), el.GetKind())
	gvr, err := core.FindGVR(restConfig, &gvk)
	if err != nil {
		return err
	}

	err = core.Patch(ctx, core.PatchOpts{
		RESTConfig: restConfig,
		GVR:        gvr.Resource,
		PatchData:  []byte(`{"metadata":{"finalizers":[]}}`),
		Name:       el.GetName(),
	})
	if err != nil {
		return err
	}

	return core.Delete(ctx, core.DeleteOpts{
		RESTConfig: restConfig,
		Object:     el,
	})
}
