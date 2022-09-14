//go:build integration
// +build integration

package configurations

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestWatchKO(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = waitUntilHealtyAndInstalled(context.Background(), restConfig, "krateo-module-core-xxx")
	assert.ErrorIs(t, err, core.ErrWatcherTimeout)
}

func TestWatchOK(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = waitUntilHealtyAndInstalled(context.Background(), restConfig, "krateo-module-core-configuration")
	assert.Nil(t, err, "expecting not nil error watching")
}
