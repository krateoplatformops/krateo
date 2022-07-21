package crds

import (
	"context"
	"strings"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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
		return nil, err
	}

	all, err := core.ListByAPIResource(ctx, core.ListByAPIResourceOpts{
		RESTConfig:  opts.RESTConfig,
		APIResource: *res,
	})
	if err != nil {
		return nil, err
	}
	if opts.FilterFunc == nil {
		return all, nil
	}

	return core.Filter(all, opts.FilterFunc)
}

func Instances(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	crds, err := List(ctx, ListOpts{
		RESTConfig: restConfig,
		FilterFunc: func(el unstructured.Unstructured) bool {
			ok := strings.HasSuffix(el.GetName(), "krateo.io")
			ok = ok || strings.HasSuffix(el.GetName(), "crossplane.io")
			return ok
		},
	})
	if err != nil {
		return nil, err
	}

	items := []unstructured.Unstructured{}

	for _, el := range crds {
		api, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
			RESTConfig: restConfig,
			Query:      el.GetName(),
		})
		if err != nil {
			return nil, err
		}

		instances, err := core.ListByAPIResource(ctx, core.ListByAPIResourceOpts{
			RESTConfig:  restConfig,
			APIResource: *api,
		})
		if err != nil {
			return nil, err
		}

		items = append(items, instances...)
	}

	return items, nil
}
