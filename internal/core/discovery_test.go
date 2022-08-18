package core

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestResolveAPIResource(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	res, err := ResolveAPIResource(ResolveAPIResourceOpts{
		RESTConfig: restConfig,
		Query:      "providers",
	})
	assert.Nil(t, err, "expecting nil error resolving API resource")

	t.Logf("%+v\n", res)
}

func TestListByAPIResource(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	res, err := ResolveAPIResource(ResolveAPIResourceOpts{
		RESTConfig: restConfig,
		Query:      "customresourcedefinitions",
	})
	assert.Nil(t, err, "expecting nil error resolving API resource")

	items, err := ListByAPIResource(context.TODO(), ListByAPIResourceOpts{
		RESTConfig:  restConfig,
		APIResource: *res,
	})
	assert.Nil(t, err, "expecting nil error listing by API resource")

	for _, el := range items {
		t.Logf("%s\n", el.GetName())
	}

}
