//go:build integration
// +build integration

package hacks

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestRemoveNamespaceFinalizers(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = RemoveNamespaceFinalizers(context.TODO(), RemoveNamespaceFinalizersOpts{
		RESTConfig: restConfig,
		Name:       "demo-system",
	})
	assert.Nil(t, err, "expecting nil error updating object")
}
