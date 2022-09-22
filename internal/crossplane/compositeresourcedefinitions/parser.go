package compositeresourcedefinitions

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	xpextv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	"github.com/pkg/errors"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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

func GetSpecFields(xrd *xpextv1.CompositeResourceDefinition) ([]Field, error) {
	vr := xrd.Spec.Versions[0]
	spec, required, err := getProps("spec", vr.Schema)
	if err != nil {
		return nil, err
	}

	res := []Field{}
	for key, el := range spec {
		flattenProps(key, el, required, &res)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res, nil
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

func flattenProps(prefix string, val extv1.JSONSchemaProps, req []string, fields *[]Field) {
	switch val.Type {
	case TypeObject:
		for key, el := range val.Properties {
			flattenProps(prefix+"."+key, el, val.Required, fields)
		}
	case TypeArray:
	default:
		//root := parent(prefix)
		//fmt.Printf("root: %s ==> %s =]> %+v\n", root, prefix, req)
		*fields = append(*fields, Field{
			Name:        prefix,
			Description: val.Description,
			Type:        val.Type,
			Default:     strval(val.Default),
			Required:    contains(req, leaf(prefix)),
		})
	}
}

func strval(inp *extv1.JSON) string {
	if inp == nil || inp.Raw == nil {
		return ""
	}

	var v interface{}
	if err := json.Unmarshal(inp.Raw, &v); err != nil {
		return err.Error()
	}

	switch v := v.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case error:
		return v.Error()
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", v)
	}
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func parent(prefix string) string {
	parts := strings.Split(prefix, ".")
	if t := len(parts); t >= 2 {
		return parts[t-2]
	}

	return prefix
}

func leaf(prefix string) string {
	for i := len(prefix) - 1; i >= 0; i-- {
		if prefix[i] == '.' {
			return prefix[i+1:]
		}
	}
	return prefix
}
