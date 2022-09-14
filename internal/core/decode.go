package core

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

func DecodeYAML(src []byte) (*unstructured.Unstructured, *schema.GroupVersionKind, error) {
	obj := &unstructured.Unstructured{}
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode(src, nil, obj)
	if err != nil {
		return nil, nil, err
	}

	labels := obj.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}
	labels[InstalledByLabel] = InstalledByValue
	unstructured.SetNestedStringMap(obj.Object, labels, "metadata", "labels")

	return obj, gvk, nil
}

func FromUnstructuredViaJSON(u map[string]interface{}, obj interface{}) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, obj)
}
