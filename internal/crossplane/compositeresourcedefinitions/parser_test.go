//go:build integration
// +build integration

package compositeresourcedefinitions

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"

	xpextv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
)

func TestGetSpecFields(t *testing.T) {
	obj, err := loadSampleCompositeResourceDefinition()
	assert.Nil(t, err, "expecting nil error decoding definition")
	assert.NotNil(t, obj, "expecting not nil unstructured object")

	xrd := obj.(*xpextv1.CompositeResourceDefinition)
	res, err := GetSpecFields(xrd)
	assert.Nil(t, err, "expecting nil error getting spec fields")

	for _, el := range res {
		if el.Required {
			fmt.Printf("%+v\n", el)
		}
	}
}

/*
// GetPropFields returns the fields from a map of schema properties
func GetPropFields(props map[string]extv1.JSONSchemaProps) []string {
	propFields := make([]string, len(props))
	i := 0
	for k := range props {
		propFields[i] = k
		i++
	}
	return propFields
}
*/

func loadSampleCompositeResourceDefinition() (runtime.Object, error) {
	const (
		sampleFile = "../../../testdata/definition.yaml"
	)

	scheme := runtime.NewScheme()
	_ = xpextv1.AddToScheme(scheme)

	decode := serializer.NewCodecFactory(scheme).UniversalDeserializer().Decode
	stream, err := ioutil.ReadFile(sampleFile)
	if err != nil {
		return nil, err
	}

	obj, _, err := decode(stream, nil, nil)
	return obj, err
}
