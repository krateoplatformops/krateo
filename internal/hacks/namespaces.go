package hacks

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type RemoveNamespaceFinalizersOpts struct {
	RESTConfig *rest.Config
	Name       string
}

func TestCreateNamespace(t *testing.T) {
	namespace := "demo-system"

	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	obj := &unstructured.Unstructured{}
	obj.SetKind("Namespace")
	obj.SetName(namespace)

	err = core.Create(context.TODO(), core.CreateOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Namespace",
		},
		Object: obj,
	})
	assert.Nil(t, err, "expecting nil error creating namespace")
}

func RemoveNamespaceFinalizers(ctx context.Context, opts RemoveNamespaceFinalizersOpts) error {
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	gvr := schema.GroupVersionResource{
		Version:  "v1",
		Resource: "namespaces",
	}

	obj, err := dc.Resource(gvr).Get(ctx, opts.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	err = unstructured.SetNestedSlice(obj.Object, nil, "spec", "finalizers")
	if err != nil {
		return err
	}

	_, err = dc.Resource(gvr).Update(ctx, obj, metav1.UpdateOptions{
		FieldManager: core.InstalledByValue,
	}, "finalize")
	return err
}
