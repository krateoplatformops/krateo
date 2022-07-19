package controllerconfigs

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo/internal/catalog"
	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type CreateOpts struct {
	RESTConfig *rest.Config
	Info       *catalog.PackageInfo
}

func Create(ctx context.Context, opts CreateOpts) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1alpha1",
		Resource: "controllerconfigs",
	}

	obj := &unstructured.Unstructured{}
	obj.SetKind("ControllerConfig")
	obj.SetAPIVersion("pkg.crossplane.io/v1alpha1")
	obj.SetName(fmt.Sprintf("%s-controllerconfig", opts.Info.Name))
	obj.SetLabels(map[string]string{
		core.InstalledByLabel: core.InstalledByValue,
	})
	unstructured.SetNestedField(obj.Object, map[string]interface{}{}, "spec", "securityContext")
	unstructured.SetNestedField(obj.Object, map[string]interface{}{}, "spec", "podSecurityContext")

	labelsForController := map[string]interface{}{
		core.InstalledByLabel: core.InstalledByValue,
		core.PackageNameLabel: opts.Info.Name,
	}
	unstructured.SetNestedField(obj.Object, labelsForController, "spec", "metadata", "labels")

	// prepare the dynamic client
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	_, err = dc.Resource(gvr).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	return obj, nil
}
