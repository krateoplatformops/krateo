//go:build integration
// +build integration

package crossplane

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestInstall(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = Install(context.TODO(), InstallOpts{
		RESTConfig: restConfig,
		Namespace:  core.DefaultNamespace,
	})
	assert.Nil(t, err, "expecting nil error installing crossplane")
}

func TestCrossplaneExists(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	ok, err := Exists(context.TODO(), ExistOpts{
		RESTConfig: restConfig,
		Namespace:  core.DefaultNamespace,
	})
	assert.Nil(t, err, "expecting nil error checking for crossplane existence")
	assert.True(t, ok, "expecting true checking for crossplane existence")
}
