package core

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type DeleteOpts struct {
	RESTConfig *rest.Config
	Object     *unstructured.Unstructured
}

func Delete(ctx context.Context, opts DeleteOpts) error {
	gvk := schema.FromAPIVersionAndKind(opts.Object.GetAPIVersion(), opts.Object.GetKind())
	mapping, err := FindGVR(opts.RESTConfig, &gvk)
	if err != nil {
		return err
	}

	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	// obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dc.Resource(mapping.Resource).Namespace(opts.Object.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dc.Resource(mapping.Resource)
	}

	err = dr.Delete(ctx, opts.Object.GetName(), metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
