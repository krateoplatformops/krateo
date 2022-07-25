package core

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type CreateOpts struct {
	RESTConfig *rest.Config
	GVK        schema.GroupVersionKind
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
	dr, err := DynamicForGVR(opts.RESTConfig, opts.GVK, opts.Object.GetNamespace())
	if err != nil {
		if IsNoKindMatchError(err) {
			return nil
		}
		return err
	}

	_, err = dr.Create(ctx, opts.Object, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}
