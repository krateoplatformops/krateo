package compositeresourcedefinitions

import (
	"encoding/json"
	"fmt"

	xpextv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	"github.com/pkg/errors"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

const (
	//moduleApiGroup     = "modules.krateo.io"
	errParseValidation = "cannot parse validation schema"
)

const (
	TypeArray   = "array"
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeObject  = "object"
	TypeString  = "string"
)

// Field is a property in the composite resource definition
type Field struct {
	Name        string
	Description string
	Type        string
	Default     string
	Required    bool
}

func GetFields(xrd *xpextv1.CompositeResourceDefinition, requiredOnly bool) ([]Field, error) {
	spec, required, err := getOpenAPISpecs(xrd)
	if err != nil {
		return nil, err
	}

	fields := []Field{}

	if requiredOnly {
		for _, k := range required {
			v := spec[k]
			flatten(k, v, &fields)
		}
		return fields, nil
	}

	for k, v := range spec {
		flatten(k, v, &fields)
	}

	return fields, nil
}

func flatten(prefix string, src v1.JSONSchemaProps, fields *[]Field) {
	switch src.Type {
	case TypeObject:
		flattenObject(prefix, src, fields)
	case TypeArray:
		flattenArray(prefix, src, fields)
	default:
		*fields = append(*fields, Field{
			Name:        prefix,
			Description: src.Description,
			Type:        src.Type,
		})
	}

}

func flattenArray(prefix string, src v1.JSONSchemaProps, fields *[]Field) {
	fmt.Println("NOT IMPLEMENTED")
}

func flattenObject(prefix string, src v1.JSONSchemaProps, fields *[]Field) {
	for k, v := range src.Properties {
		flatten(prefix+"."+k, v, fields)
	}
}

func getOpenAPISpecs(xrd *xpextv1.CompositeResourceDefinition) (map[string]extv1.JSONSchemaProps, []string, error) {
	vr := xrd.Spec.Versions[0]
	return getProps("spec", vr.Schema)
}

func getProps(field string, v *xpextv1.CompositeResourceValidation) (map[string]extv1.JSONSchemaProps, []string, error) {
	if v == nil {
		return nil, nil, nil
	}

	s := &extv1.JSONSchemaProps{}
	if err := json.Unmarshal(v.OpenAPIV3Schema.Raw, s); err != nil {
		return nil, nil, errors.Wrap(err, errParseValidation)
	}

	spec, ok := s.Properties[field]
	if !ok {
		return nil, nil, nil
	}

	return spec.Properties, spec.Required, nil
}
