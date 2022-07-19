package crossplane

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

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
		GVR: schema.GroupVersionResource{
			Version:  "v1",
			Resource: "pods",
		},
		Namespace:     opts.Namespace,
		LabelSelector: sel.String(),
	})
	if err != nil {
		return false, err
	}

	return len(list) > 0, nil
}
