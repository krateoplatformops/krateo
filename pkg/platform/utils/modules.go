package utils

import (
	"fmt"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
)

const (
	moduleConfigurationGroupAndKind = "Configuration.pkg.crossplane.io"
	moduleClaimsGroup               = "modules.krateo.io"
)

// Decode Crossplane Package Configuration YAML into unstructured.Unstructured
func DecodeModuleConfiguration(dc dynamic.Interface, data []byte) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode(data, nil, obj)
	if err != nil {
		return nil, err
	}

	if gvk.GroupKind().String() != moduleConfigurationGroupAndKind {
		return nil, fmt.Errorf("kind: %s in apiGroup: %s is not allowed", gvk.Kind, gvk.Group)
	}

	return obj, nil
}

// DecodeModuleClaims Crossplane Composition Claims YAML into unstructured.Unstructured
func DecodeModuleClaims(dc dynamic.Interface, data []byte) (*unstructured.Unstructured, error) {
	obj := &unstructured.Unstructured{}

	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := dec.Decode(data, nil, obj)
	if err != nil {
		return nil, err
	}

	if grp := gvk.GroupKind().Group; grp != moduleClaimsGroup {
		return nil, fmt.Errorf("apiGroup: %s is not allowed", grp)
	}

	return obj, nil
}
