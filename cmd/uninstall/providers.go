package uninstall

import (
	"context"

	"github.com/krateoplatformops/krateo/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func deleteCrossplaneProviders(dc dynamic.Interface, dryRun bool) error {
	opts := metav1.DeleteOptions{}
	if dryRun {
		opts.DryRun = []string{metav1.DryRunAll}
	}

	// delete controller config
	err := dc.Resource(schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1alpha1",
		Resource: "controllerconfigs",
	}).
		Namespace(kubernetes.CrossplaneSystemNamespace).
		Delete(context.Background(), "krateo-controllerconfig", opts)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	// delete crossplane provider-helm
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	err = dc.Resource(gvr).
		Namespace(kubernetes.CrossplaneSystemNamespace).
		Delete(context.Background(), "crossplane-provider-helm", opts)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	// delete crossplane provider-kubernetes
	err = dc.Resource(gvr).
		Namespace(kubernetes.CrossplaneSystemNamespace).
		Delete(context.Background(), "crossplane-provider-kubernetes", opts)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
