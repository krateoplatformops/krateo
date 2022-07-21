//go:build integration
// +build integration

package providers

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGetProviders(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := All(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing providers")

	t.Logf("found [%d] providers\n", len(all))
	for _, el := range all {
		t.Logf("> %s\n", el.GetName())
	}
}

func TestGetConfigurations(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := GetConfigurations(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing configurations")

	t.Logf("found [%d] providers\n", len(all))
	for _, el := range all {
		t.Logf("> %s\n", el.GetName())
	}
}
