package core

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/rest"
)

type DeleteOpts struct {
	RESTConfig *rest.Config
	Object     *unstructured.Unstructured
}

func Delete(ctx context.Context, opts DeleteOpts) error {
	dr, err := DynamicForGVR(opts.RESTConfig, opts.Object.GroupVersionKind(), opts.Object.GetNamespace())
	if err != nil {
		if IsNoKindMatchError(err) {
			return nil
		}
		return err
	}

	err = dr.Delete(ctx, opts.Object.GetName(), metav1.DeleteOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
