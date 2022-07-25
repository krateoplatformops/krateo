package crds

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ListOpts struct {
	RESTConfig *rest.Config
	FilterFunc core.FilterFunc
}

func List(ctx context.Context, opts ListOpts) ([]unstructured.Unstructured, error) {
	res, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
		RESTConfig: opts.RESTConfig,
		Query:      "customresourcedefinitions",
	})
	if err != nil {
		if core.IsNoKindMatchError(err) {
			return nil, nil
		}
		return nil, err
	}

	all, err := core.ListByAPIResource(ctx, core.ListByAPIResourceOpts{
		RESTConfig:  opts.RESTConfig,
		APIResource: *res,
	})
	if err != nil {
		if core.IsNoKindMatchError(err) {
			return nil, nil
		}
		return nil, err
	}
	if opts.FilterFunc == nil {
		return all, nil
	}

	return core.Filter(all, opts.FilterFunc)
}

func Instances(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	gvk := schema.GroupVersionKind{
		Group:   "apiextensions.k8s.io",
		Version: "v1",
		Kind:    "CustomResourceDefinition",
	}

	crds, err := core.List(ctx, core.ListOpts{RESTConfig: restConfig, GVK: gvk})
	if err != nil {
		if core.IsNoKindMatchError(err) {
			return nil, nil
		}
		return nil, err
	}

	/*
		list, err := core.Filter(crds, func(el unstructured.Unstructured) bool {
			ok := strings.HasSuffix(el.GetAPIVersion(), "krateo.io")
			ok = ok || strings.HasSuffix(el.GetAPIVersion(), "crossplane.io")
			return ok
		})
		if err != nil {
			return nil, err
		}
	*/

	items := []unstructured.Unstructured{}

	for _, el := range crds {
		api, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
			RESTConfig: restConfig,
			Query:      el.GetName(),
		})
		if err != nil {
			if core.IsNoKindMatchError(err) {
				continue
			}
			return nil, err
		}
		fmt.Println("==> ", api.GroupVersionResource().Resource, api.Group)

		res, err := core.List(ctx, core.ListOpts{
			RESTConfig: restConfig,
			GVK: schema.GroupVersionKind{
				Group:   api.Group,
				Version: api.Version,
				Kind:    api.Kind,
			},
			Namespace: el.GetNamespace(),
		})
		if err != nil {
			continue
		}

		items = append(items, res...)
	}

	return items, nil
}
