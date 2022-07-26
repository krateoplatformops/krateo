package crds

import (
	"context"
	"strings"

	"github.com/krateoplatformops/krateo/internal/core"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ListOpts struct {
	RESTConfig *rest.Config
	FilterFunc core.FilterFunc
}

func CRDInstances(ctx context.Context, restConfig *rest.Config, crd string) []unstructured.Unstructured {
	api, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
		RESTConfig: restConfig,
		Query:      crd,
	})
	if err != nil || api == nil {
		return nil
	}

	items, err := core.List(ctx, core.ListOpts{
		RESTConfig: restConfig,
		GVK:        api.GroupVersionKind(),
	})
	if err != nil {
		return nil
	}

	return items
}

func List(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	all, err := core.List(ctx, core.ListOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return core.Filter(all, func(el unstructured.Unstructured) bool {
		ok := strings.HasSuffix(el.GetName(), "krateo.io")
		ok = ok || strings.HasSuffix(el.GetName(), "crossplane.io")
		return ok
	})
}
