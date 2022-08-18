//go:build integration
// +build integration

package crds

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
)

func TestList(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	items, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing compositions")

	for _, el := range items {
		t.Logf("  > %s\n", el.GetName())
	}
}

func TestCRDInstances(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	list, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing crds")

	t.Logf("found [%d] crds\n", len(list))

	items := []unstructured.Unstructured{}
	for _, el := range list {
		res := CRDInstances(context.TODO(), restConfig, el.GetName())
		if items != nil {
			items = append(items, res...)
		}
	}

	for _, el := range items {
		t.Logf("> %s - (%s)\n", el.GetName(), el.GetAPIVersion())
	}
}
