//go:build integration
// +build integration

package apiextensions

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/tools/clientcmd"
)

func TestList(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error listing crds")

	for _, el := range all {
		spew.Dump(el)
		fmt.Println()
		//t.Logf("%s (%s)\n", el.GetName(), el.GetAPIVersion())
		group, ok, err := unstructured.NestedString(el.Object, "spec", "group")
		assert.Nil(t, err, "expecting nil error getting crd group")
		assert.True(t, ok, "expecting true getting crd group")

		names, ok, err := unstructured.NestedMap(el.Object, "spec", "names")
		assert.Nil(t, err, "expecting nil error getting crd group")
		//if ok {
		//	resource = names["plural"].(string)
		//}

		versions, ok, err := unstructured.NestedSlice(el.Object, "spec", "versions")

		v := versions[0].(map[string]interface{})["name"].(string)
		t.Logf("g: %s, v: %s, r: %s\n", group, v, names["plural"].(string))
	}

}
