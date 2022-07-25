package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"
)

func TestListCRDs(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), ListOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
	})

	//res, err := Filter(all, func(obj unstructured.Unstructured) bool {
	//	return strings.HasSuffix(obj.GetName(), "modules.krateo.io")
	//})
	assert.Nil(t, err, "expecting nil error listing crds")

	for _, el := range all {
		t.Logf("%s\n", el.GetName())
	}
}

func TestJSONPath(t *testing.T) {
	pt := jsonpath.New("sample")
	pt.AllowMissingKeys(true)
	err := pt.Parse(`{range .items[*]}{.kind}/{.metadata.name}{"\n"}{end}`)
	assert.Nil(t, err, "expecting nil error creating json path")

	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := List(context.TODO(), ListOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinition",
		},
	})

	for _, el := range all {
		t.Logf("%s\n", el.GetName())

		spew.Dump(el.Object)
		values, _ := pt.FindResults(el.Object)

		spew.Dump(values)
		fmt.Printf("\n\n")
	}

}
