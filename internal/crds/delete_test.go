//go:build integration
// +build integration

package crds

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

	items, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range items {
		err = PatchAndDelete(context.TODO(), restConfig, &el)
		assert.Nil(t, err, "expecting nil error patching an deleting object: ", el.GetName())
	}

}
