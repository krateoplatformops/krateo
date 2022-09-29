//go:build integration
// +build integration

package compositeresourcedefinitions

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/krateoplatformops/krateo/internal/strvals"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"sigs.k8s.io/yaml"

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

func TestGetSpecFields_claim_data(t *testing.T) {
	obj, err := loadSampleCompositeResourceDefinition()
	assert.Nil(t, err, "expecting nil error decoding definition")
	assert.NotNil(t, obj, "expecting not nil unstructured object")

	xrd := obj.(*xpextv1.CompositeResourceDefinition)
	all, err := GetSpecFields(xrd)
	assert.Nil(t, err, "expecting nil error getting spec fields")

	values := []string{}

	for _, el := range all {
		if len(el.Default) == 0 {
			if !el.Required {
				continue
			}

			t.Logf("%s will be prompted\n", el.Name)
			continue
		}

		switch el.Type {
		case TypeBoolean:
			val, err := strconv.ParseBool(el.Default)
			if err == nil {
				values = append(values, fmt.Sprintf("%s=%t", el.Name, val))
			}

		case TypeInteger:
			val, err := strconv.Atoi(el.Default)
			if err == nil {
				values = append(values, fmt.Sprintf("%s=%d", el.Name, val))
			}
		case TypeNumber:
			val, err := strconv.ParseFloat(el.Default, 64)
			if err != nil {
				values = append(values, fmt.Sprintf("%s=%f", el.Name, val))
			}
		default:
			values = append(values, fmt.Sprintf("%s=%s", el.Name, el.Default))
		}
	}

	data := make(map[string]interface{})
	err = strvals.ParseInto(strings.Join(values, ","), data)
	assert.Nil(t, err, "expecting nil error creating claim data")

	gvk := schema.GroupVersionKind{
		Group:   "modules.krateo.io",
		Version: "v1",
		Kind:    "Core",
	}

	clm := &unstructured.Unstructured{}
	clm.SetKind(gvk.Kind)
	clm.SetAPIVersion(gvk.GroupVersion().String())
	clm.SetName("core")
	clm.SetLabels(map[string]string{
		core.InstalledByLabel: core.InstalledByValue,
	})
	err = unstructured.SetNestedField(clm.Object, data, "spec")
	assert.Nil(t, err, "expecting nil error creating claim unstructured object")

	b, err := yaml.Marshal(clm)
	assert.Nil(t, err, "expecting nil error marshalling claim unstructured object")

	t.Logf("%s\n", string(b))
}

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
