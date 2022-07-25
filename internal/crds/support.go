package crds

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

func PatchAndDelete(ctx context.Context, restConfig *rest.Config, el *unstructured.Unstructured) error {
	err := core.Patch(ctx, core.PatchOpts{
		RESTConfig: restConfig,
		GVK:        el.GroupVersionKind(),
		PatchData:  []byte(`{"metadata":{"finalizers":[]}}`),
		Name:       el.GetName(),
		Namespace:  el.GetNamespace(),
	})
	if err != nil {
		return err
	}

	return core.Delete(ctx, core.DeleteOpts{
		RESTConfig: restConfig,
		Object:     el,
	})
}
