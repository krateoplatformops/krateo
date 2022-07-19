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

type CreateOpts struct {
	RESTConfig *rest.Config
	GVR        schema.GroupVersionResource
	Object     *unstructured.Unstructured
	Namespace  string
}

// Create creates a resource if does not exists.
func Create(ctx context.Context, opts CreateOpts) error {
	// set 'installed-by' label
	labels := opts.Object.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	labels[InstalledByLabel] = InstalledByValue
	opts.Object.SetLabels(labels)

	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	namespace := opts.Namespace
	if namespace == "" {
		namespace = opts.Object.GetNamespace()
	}

	_, err = dc.Resource(opts.GVR).Create(ctx, opts.Object, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}
