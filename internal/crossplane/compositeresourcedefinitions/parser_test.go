//go:build integration
// +build integration

package compositeresourcedefinitions

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestGetFields(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	xrd, err := Get(context.TODO(), restConfig, "core.modules.krateo.io")
	assert.Nil(t, err, "expecting nil error getting composite resource definition")
	assert.NotNil(t, xrd, "expecting not nil getting composite resource definition")

	fields, err := GetFields(xrd, false)
	assert.Nil(t, err, "expecting nil error parsing composite resource definition")
	assert.NotNil(t, fields, "expecting not nil getting composite resource definition fields")

	fmt.Printf("%v\n", fields)
}
