package composite

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

func List(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	res, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
		RESTConfig: restConfig,
		Query:      "compositeresourcedefinitions",
	})
	if err != nil {
		return nil, err
	}

	all, err := core.ListByAPIResource(ctx, core.ListByAPIResourceOpts{
		RESTConfig:  restConfig,
		APIResource: *res,
	})
	if err != nil {
		return nil, err
	}

	items := []unstructured.Unstructured{}

	for _, el := range all {
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
