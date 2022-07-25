//go:build integration
// +build integration

package pods

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
)

func TestList(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	items, err := core.List(context.TODO(), core.ListOpts{
		RESTConfig:    restConfig,
		GVK:           schema.GroupVersionKind{Version: "v1", Kind: "Pod"},
		Namespace:     "",
		LabelSelector: "app=crossplane",
	})
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range items {
		t.Logf("  > %s\n", el.GetName())
	}
}
