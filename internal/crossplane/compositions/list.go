package compositions

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func List(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	all, err := core.List(ctx, core.ListOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.crossplane.io",
			Version: "v1",
			Kind:    "Composition",
		},
	})
	if err != nil {
		if core.IsNoKindMatchError(err) || errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return all, nil
}
