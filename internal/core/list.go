package core

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type FilterFunc func(unstructured.Unstructured) bool

type ListOpts struct {
	RESTConfig    *rest.Config
	GVK           schema.GroupVersionKind
	Namespace     string
	LabelSelector string
}

func List(ctx context.Context, opts ListOpts) ([]unstructured.Unstructured, error) {
	listOpts := metav1.ListOptions{}
	if len(opts.LabelSelector) > 0 {
		listOpts.LabelSelector = opts.LabelSelector
	}

	dr, err := DynamicForGVR(opts.RESTConfig, opts.GVK, opts.Namespace)
	if err != nil {
		if IsNoKindMatchError(err) {
			return nil, nil
		}
		return nil, err
	}

	list, err := dr.List(ctx, listOpts)
	if err != nil {
		if apierrors.IsNotFound(err) {
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
