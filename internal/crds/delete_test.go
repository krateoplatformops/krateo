//go:build integration
// +build integration

package crds

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDelete(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	items, err := Instances(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range items {
		err = patchAndDelete(context.TODO(), restConfig, &el)
		assert.Nil(t, err, "expecting nil error patching an deleting object: ", el.GetName())
	}

	items, err = List(context.TODO(), ListOpts{RESTConfig: restConfig})
	assert.Nil(t, err, "expecting nil error listing crds")

	for _, el := range items {
		err = patchAndDelete(context.TODO(), restConfig, &el)
		assert.Nil(t, err, "expecting nil error patching an deleting object: ", el.GetName())
	}

}

func patchAndDelete(ctx context.Context, restConfig *rest.Config, el *unstructured.Unstructured) error {
	gvk := schema.FromAPIVersionAndKind(el.GetAPIVersion(), el.GetKind())
	gvr, err := core.FindGVR(restConfig, &gvk)
	if err != nil {
		return err
	}

	err = core.Patch(ctx, core.PatchOpts{
		RESTConfig: restConfig,
		GVR:        gvr.Resource,
		PatchData:  []byte(`{"metadata":{"finalizers":[]}}`),
		Name:       el.GetName(),
	})
	if err != nil {
		return err
	}

	return core.Delete(ctx, core.DeleteOpts{
		RESTConfig: restConfig,
		Object:     el,
	})
}
