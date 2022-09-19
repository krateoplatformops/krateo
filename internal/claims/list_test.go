//go:build integration
// +build integration

package claims

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	// kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	kubeconfig, err := ioutil.ReadFile("../../testdata/krateo-test.yml")
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range all {
		t.Logf("> %s\n", el.GetName())
	}
}
