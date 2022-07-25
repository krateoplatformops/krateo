package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/kubectl/pkg/scheme"
)

func TestGetCRDs(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig)
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	res, err := Get(context.TODO(), GetOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Group:   "apiextensions.k8s.io",
			Version: "v1",
			Kind:    "CustomResourceDefinitions",
		},
		Name: "controllerconfigs.pkg.crossplane.io",
	})
	assert.Nil(t, err, "expecting nil error listing crds")

	// Serializer = Decoder + Encoder.
	serializer := jsonserializer.NewSerializerWithOptions(
		jsonserializer.DefaultMetaFactory, // jsonserializer.MetaFactory
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectCreater
		scheme.Scheme,                     // runtime.Scheme implements runtime.ObjectTyper
		jsonserializer.SerializerOptions{
			Yaml:   true,
			Pretty: true,
			Strict: false,
		},
	)

	// Typed -> YAML
	// Runtime.Encode() is just a helper function to invoke Encoder.Encode()
	yaml, err := runtime.Encode(serializer, res)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%s\n\n", yaml)

	pt := jsonpath.New("sample")
	pt.AllowMissingKeys(true)
	err = pt.Parse(`{range .items[*]}{.kind}/{.metadata.name}{"\n"}{end}`)
	assert.Nil(t, err, "expecting nil error creating json path")

	values, err := pt.FindResults(res.Object)
	assert.Nil(t, err, "expecting nil error creating json path")

	spew.Dump(values)
}
