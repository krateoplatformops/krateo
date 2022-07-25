//go:build integration
// +build integration

package clusterrolebindings

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
)

func TestListCRDs(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing clusterrolebindings")

	res, err := core.Filter(all, func(obj unstructured.Unstructured) bool {
		accept := (obj.GetName() == "provider-helm-admin-binding")
		accept = accept || (obj.GetName() == "provider-kubernetes-admin-binding")
		return accept
	})
	assert.Nil(t, err, "expecting nil error listing crds")

	for _, el := range res {
		t.Logf("%s\n", el.GetName())
	}
}
