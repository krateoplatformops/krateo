//go:build integration
// +build integration

package providers

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/catalog"
	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestInstall(t *testing.T) {
	list, err := catalog.FilterBy(catalog.ForCLI())
	assert.Nil(t, err, "expecting nil error loading catalog")

	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	for _, el := range list.Items {
		t.Logf("> installing: %s\n", el.Name)
		err := Install(context.TODO(), InstallOpts{
			RESTConfig: restConfig,
			Info:       &el,
		})
		assert.Nil(t, err, "expecting nil error installing package ", el.Name)
	}
}
