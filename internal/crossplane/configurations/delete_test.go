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

func TestDelete(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing configurations")

	t.Logf("found [%d] configurations\n", len(all))
	for _, el := range all {
		t.Logf("deleting: %s\n", el.GetName())
		err := core.Delete(context.TODO(), core.DeleteOpts{
			RESTConfig: restConfig,
			Object:     &el,
		})
		assert.Nil(t, err, "expecting nil error deleting configuration: %s", el.GetName())
	}
}
