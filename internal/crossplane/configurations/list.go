package configurations

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

func List(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	res, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
		RESTConfig: restConfig,
		Query:      "configurations",
	})
	if err != nil {
		return nil, err
	}

	return core.ListByAPIResource(ctx, core.ListByAPIResourceOpts{
		RESTConfig:  restConfig,
		APIResource: *res,
	})
}
