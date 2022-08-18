//go:build integration
// +build integration

package crossplane

import (
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestUninstall(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = Uninstall(UninstallOpts{
		RESTConfig: restConfig,
		Namespace:  "default",
	})
	assert.Nil(t, err, "expecting nil error uninstalling crossplane")
}
