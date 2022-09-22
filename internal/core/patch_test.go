//go:build integration
// +build integration

package core

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/clientcmd"
)

func TestPatchNamespace(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = Patch(context.TODO(), PatchOpts{
		RESTConfig: restConfig,
		Name:       "krateo-system",
		GVK: schema.GroupVersionKind{
			Group:   "",
			Version: "v1",
			Kind:    "Namespace",
		},
		PatchData: []byte(`[
			{ "op": "remove", "path": "/spec/finalizers/1" }
		  ]`),
		PatchType: types.JSONPatchType,
	})
	assert.Nil(t, err, "expecting nil error patching resource")
	//if status, ok := err.(kerrors.APIStatus); ok || errors.As(err, &status) {
	//	fmt.Println(status.Status().Reason)
	//}
}
