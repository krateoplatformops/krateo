//go:build integration
// +build integration

package crds

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestList(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	items, err := Instances(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range items {
		t.Logf("  > %s\n", el.GetName())
	}
}

func TestGet(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	res, err := core.ResolveAPIResource(core.ResolveAPIResourceOpts{
		RESTConfig: restConfig,
		Query:      "mongodb.core.modules.krateo.io",
	})

	items, err := core.ListByAPIResource(context.TODO(), core.ListByAPIResourceOpts{
		RESTConfig:  restConfig,
		APIResource: *res,
	})

	assert.Nil(t, err, "expecting nil error listing crds")
	for _, el := range items {
		t.Logf("> %s\n", el.GetName())
	}

}
