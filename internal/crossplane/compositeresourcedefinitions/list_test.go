//go:build integration
// +build integration

package compositions

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestList(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range all {
		t.Logf("> %s\n", el.GetName())
	}
}
