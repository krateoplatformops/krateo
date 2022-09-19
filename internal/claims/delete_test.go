package claims

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	//kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	kubeconfig, err := ioutil.ReadFile("../../testdata/krateo-test.yml")
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	obj, err := core.Get(context.TODO(), core.GetOpts{
		RESTConfig: restConfig,
		Name:       "core",
		GVK:        getGroupVersionKind(),
	})
	if err != nil {
		t.Fatal(err)
	}
	err = core.Delete(context.TODO(), core.DeleteOpts{
		RESTConfig: restConfig,
		Object:     obj,
	})
	assert.Nil(t, err, "expecting nil error deleting claim: %s", obj.GetName())

}
