//go:build integration
// +build integration

package claims

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestWatchKO(t *testing.T) {
	//kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	kubeconfig, err := ioutil.ReadFile("../../testdata/krateo-test.yml")
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = WaitUntilReady(context.Background(), restConfig, "core")
	assert.ErrorIs(t, err, core.ErrWatcherTimeout)
}

func TestWatchOK(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	//kubeconfig, err := ioutil.ReadFile("../../testdata/krateo-test.yml")
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = WaitUntilReady(context.Background(), restConfig, "core")
	assert.Nil(t, err, "expecting  nil error watching")
}
