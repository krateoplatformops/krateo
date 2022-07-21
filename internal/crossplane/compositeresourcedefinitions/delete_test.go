package compositions

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestDelete(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = Delete(context.TODO(), DeleteOpts{
		RESTConfig:      restConfig,
		PatchFinalizers: true,
		Name:            "core.modules.krateo.io",
	})
	assert.Nil(t, err, "expecting nil error deleting composition")
}
