//go:build integration
// +build integration

package configurations

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/catalog"
	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreate(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	obj, err := Create(context.TODO(), CreateOpts{
		RESTConfig: restConfig,
		Info: &catalog.PackageInfo{
			Name:    "krateo-module-core",
			Image:   "ghcr.io/krateoplatformops/krateo-module-core",
			Version: "latest",
		},
	})
	assert.Nil(t, err, "expecting nil error creating Configuration")
	assert.NotNil(t, obj, "expecting not nil Configuration")
}
