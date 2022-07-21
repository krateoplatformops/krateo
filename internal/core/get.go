package core

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type GetOpts struct {
	RESTConfig *rest.Config
	GVR        schema.GroupVersionResource
	Name       string
	Namespace  string
}

func Get(ctx context.Context, opts GetOpts) (*unstructured.Unstructured, error) {
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	var res *unstructured.Unstructured
	if opts.Namespace == "" {
		res, err = dc.Resource(opts.GVR).Get(ctx, opts.Name, metav1.GetOptions{})
	} else {
		res, err = dc.Resource(opts.GVR).Namespace(opts.Namespace).Get(ctx, opts.Name, metav1.GetOptions{})
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}
