package compositeresourcedefinitions

import (
	"context"

	xpextv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func Get(ctx context.Context, rc *rest.Config, name string) (*xpextv1.CompositeResourceDefinition, error) {
	obj, err := core.Get(context.TODO(), core.GetOpts{
		RESTConfig: rc,
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.crossplane.io",
			Version: "v1",
			Kind:    "CompositeResourceDefinition",
		},
		Name: name,
	})
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, nil
	}

	var xrd *xpextv1.CompositeResourceDefinition
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &xrd)
	return xrd, err
}
