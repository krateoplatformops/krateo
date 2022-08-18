//go:build integration
// +build integration

package controllerconfigs

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
)

func TestList(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	//sel, err := core.InstalledBySelector()
	//assert.Nil(t, err, "expecting nil error creating labels selector")

	opts := core.ListOpts{
		GVK: schema.GroupVersionKind{
			Group:   "pkg.crossplane.io",
			Version: "v1alpha1",
			Kind:    "ControllerConfig",
		},
		RESTConfig: restConfig,
		//LabelSelector: sel.String(),
	}

	list, err := core.List(context.TODO(), opts)
	assert.Nil(t, err, "expecting nil error listing controller configs")

	t.Logf("found [%d] controller configs\n", len(list))
	for _, el := range list {
		t.Logf("> %s\n", el.GetName())
	}
}
