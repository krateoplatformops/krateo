package compositions

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type DeleteOpts struct {
	RESTConfig      *rest.Config
	Name            string
	PatchFinalizers bool
}

func Delete(ctx context.Context, opts DeleteOpts) (err error) {
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.crossplane.io",
		Version:  "v1",
		Resource: "compositeresourcedefinitions",
	}

	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	if opts.PatchFinalizers {
		err := core.Patch(ctx, core.PatchOpts{
			RESTConfig: opts.RESTConfig,
			PatchData:  []byte(`{"metadata":{"finalizers":[]}}`),
			GVR:        gvr,
			Name:       opts.Name,
		})
		if err != nil {
			return err
		}
	}

	err = dc.Resource(gvr).Delete(ctx, opts.Name, metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return err
}