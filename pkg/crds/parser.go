package crds

import (
	"encoding/json"
	"fmt"
	"sort"
)

/*
func GetKeys(xrd *CompositeResourceDefinition) map[string]string {
	res := make(map[string]string)

	spec := getRootProperties(xrd)
	if spec == nil {
		return res
	}

	var keys []string
	for k := range spec.Properties {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		res[k] = spec.Properties[k].Type
	}
	return res
}
*/

type Field struct {
	Name        string
	Description string
	Type        string
}

// Flatten flatten the composite resource
// definition openapi schema properties
func Parse(src CompositeResourceDefinition) []Field {
	res := []Field{}

	if len(src.Spec.Versions) == 0 {
		return res
	}

	defs := src.Spec.Versions[0].Schema.OpenAPIV3Schema
	if defs == nil {
		return res
	}

	spec, ok := defs.Properties["spec"]
	if !ok {
		return res
	}

	sections := getPropFields(spec.Properties)
	for _, el := range sections {
		out := parseProps(spec.Properties[el], el)
		res = append(res, out...)
	}

	return res
}

func parseProps(el JSONSchemaProps, prefix string) []Field {
	res := []Field{}
	parseObject(&res, el, prefix)
	return res
}

//nolint:errcheck
func parseArray(dst *[]Field, el JSONSchemaProps, root string) {
	var props JSONSchemaProps
	json.Unmarshal(el.Items.RawMessage, &props)
	parseObject(dst, props, fmt.Sprintf("%s[0]", root))
}

func parseObject(dst *[]Field, el JSONSchemaProps, root string) {
	for key, val := range el.Properties {
		name := key
		if len(root) > 0 {
			name = fmt.Sprintf("%s.%s", root, key)
		}

		switch val.Type {
		case TypeObject:
			parseObject(dst, val, name)
		case TypeArray:
			parseArray(dst, val, name)
		default:
			item := Field{
				Name:        name,
				Description: val.Description,
				Type:        val.Type,
			}
			*dst = append(*dst, item)
		}
	}
}

// Flatten flatten the composite resource
// definition openapi schema properties
func Flatten(src CompositeResourceDefinition) []string {
	res := []string{}

	if len(src.Spec.Versions) == 0 {
		return res
	}

	defs := src.Spec.Versions[0].Schema.OpenAPIV3Schema
	if defs == nil {
		return res
	}

	spec, ok := defs.Properties["spec"]
	if !ok {
		return res
	}

	sections := getPropFields(spec.Properties)
	for _, el := range sections {
		out := flattenProps(spec.Properties[el], el)
		res = append(res, out...)
	}

	return res
}

// getPropFields returns the fields from a map of schema properties
func getPropFields(props map[string]JSONSchemaProps) []string {
	propFields := make([]string, len(props))

	i := 0
	for k := range props {
		propFields[i] = k
		i++
	}
	sort.Strings(propFields)

	return propFields
}

func flattenProps(el JSONSchemaProps, prefix string) []string {
	res := []string{}
	flattenObject(&res, el, prefix)
	return res
}

//nolint:errcheck
func flattenArray(dst *[]string, el JSONSchemaProps, root string) {
	var props JSONSchemaProps
	json.Unmarshal(el.Items.RawMessage, &props)
	flattenObject(dst, props, fmt.Sprintf("%s[0]", root))
}

func flattenObject(dst *[]string, el JSONSchemaProps, root string) {
	for key, val := range el.Properties {
		name := key
		if len(root) > 0 {
			name = fmt.Sprintf("%s.%s", root, key)
		}

		switch val.Type {
		case TypeObject:
			flattenObject(dst, val, name)
		case TypeArray:
			flattenArray(dst, val, name)
		case TypeBoolean:
			vals := fmt.Sprintf("%s=<BOOL>", name)
			*dst = append(*dst, vals)
		case TypeInteger:
			vals := fmt.Sprintf("%s=<INT>", name)
			*dst = append(*dst, vals)
		case TypeNumber:
			vals := fmt.Sprintf("%s=<FLOAT>", name)
			*dst = append(*dst, vals)
		default:
			vals := fmt.Sprintf("%s=<STRING>", name)
			*dst = append(*dst, vals)
		}
	}
}
