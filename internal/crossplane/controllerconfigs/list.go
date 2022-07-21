package controllerconfigs

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func ListAll(ctx context.Context, restConfig *rest.Config) ([]unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1alpha1",
		Resource: "controllerconfigs",
	}

	return core.List(ctx, core.ListOpts{
		RESTConfig: restConfig,
		GVR:        gvr,
	})
}
