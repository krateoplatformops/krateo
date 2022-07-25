package core

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type GetOpts struct {
	RESTConfig *rest.Config
	GVK        schema.GroupVersionKind
	Name       string
	Namespace  string
}

func Get(ctx context.Context, opts GetOpts) (*unstructured.Unstructured, error) {
	dr, err := DynamicForGVR(opts.RESTConfig, opts.GVK, opts.Namespace)
	if err != nil {
		if IsNoKindMatchError(err) {
			return nil, nil
		}
		return nil, err
	}

	res := &unstructured.Unstructured{}
	res, err = dr.Get(ctx, opts.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}
