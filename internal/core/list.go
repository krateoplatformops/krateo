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

type FilterFunc func(unstructured.Unstructured) bool

type ListOpts struct {
	RESTConfig    *rest.Config
	GVR           schema.GroupVersionResource
	Namespace     string
	LabelSelector string
}

func List(ctx context.Context, opts ListOpts) ([]unstructured.Unstructured, error) {
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	listOpts := metav1.ListOptions{}
	if len(opts.LabelSelector) > 0 {
		listOpts.LabelSelector = opts.LabelSelector
	}

	var list *unstructured.UnstructuredList
	if opts.Namespace == "" {
		list, err = dc.Resource(opts.GVR).List(ctx, listOpts)
	} else {
		list, err = dc.Resource(opts.GVR).Namespace(opts.Namespace).List(ctx, listOpts)
	}

	if err != nil {
		if errors.IsNotFound(err) {
			return []unstructured.Unstructured{}, nil
		}
		return nil, err
	}

	return list.Items, nil
}

func Filter(list []unstructured.Unstructured, accept FilterFunc) ([]unstructured.Unstructured, error) {
	if accept == nil {
		return list, nil
	}

	res := []unstructured.Unstructured{}
	for _, el := range list {
		if accept(el) {
			res = append(res, el)
		}
	}

	return res, nil
}
