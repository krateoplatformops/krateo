package providers

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func ListInstalled(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	sel, err := core.InstalledBySelector()
	if err != nil {
		return nil, err
	}

	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	return core.List(ctx, core.ListOpts{
		RESTConfig:    restConfig,
		GVR:           gvr,
		LabelSelector: sel.String(),
	})
}

func ListAll(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	return core.List(ctx, core.ListOpts{
		RESTConfig: restConfig,
		GVR:        gvr,
	})
}
