package crossplane

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func GetPOD(ctx context.Context, restConfig *rest.Config) (*unstructured.Unstructured, error) {
	items, err := core.List(ctx, core.ListOpts{
		RESTConfig:    restConfig,
		GVK:           schema.GroupVersionKind{Version: "v1", Kind: "Pod"},
		Namespace:     "",
		LabelSelector: "app=crossplane",
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(items) > 0 {
		return &items[0], nil
	}

	return nil, nil
}

type ExistOpts struct {
	RESTConfig *rest.Config
	Namespace  string
}

func Exists(ctx context.Context, opts ExistOpts) (bool, error) {
	sel, err := labels.Parse("app=crossplane")
	if err != nil {
		return false, err
	}

	list, err := core.List(ctx, core.ListOpts{
		RESTConfig: opts.RESTConfig,
		GVK: schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Pod",
		},
		Namespace:     opts.Namespace,
		LabelSelector: sel.String(),
	})
	if err != nil {
		return false, err
	}

	return len(list) > 0, nil
}
