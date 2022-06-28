package namespaces

import (
	"context"

	"github.com/krateoplatformops/krateo/pkg/platform/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

const (
	KubeSystem       = "kube-system"
	CrossplaneSystem = "crossplane-system"
	KrateoSystem     = "krateo-system"
)

// Create creates a namespace if not exists.
func Create(dc dynamic.Interface, name string) error {
	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	obj := &unstructured.Unstructured{}
	obj.SetKind("Namespace")
	obj.SetName(name)
	obj.SetLabels(map[string]string{
		utils.LabelManagedBy: utils.DefaultFieldManager,
	})

	_, err := dc.Resource(gvr).
		Create(context.Background(), obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}
