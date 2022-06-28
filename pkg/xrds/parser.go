package xrds

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"sigs.k8s.io/yaml"
)

// Field is a property in the composite resource definition
type Field struct {
	Name        string
	Description string
	Type        string
	Default     string
	Required    bool
}

func ParseCompositeResourceDefinition(xrd *CompositeResourceDefinition) ([]Field, error) {
	p := &parser{
		fields: make([]Field, 0),
	}

	return p.parse(xrd)
}

func ParseBytes(dat []byte) ([]Field, error) {
	xrd := &CompositeResourceDefinition{}
	err := yaml.Unmarshal(dat, xrd)
	if err != nil {
		return nil, err
	}

	return ParseCompositeResourceDefinition(xrd)
}

type parser struct {
	fields []Field
}

func (p *parser) parse(src *CompositeResourceDefinition) ([]Field, error) {
	if len(src.Spec.Versions) == 0 {
		return nil, fmt.Errorf("no spec version found")
	}

	defs := src.Spec.Versions[0].Schema.OpenAPIV3Schema
	if defs == nil {
		return nil, fmt.Errorf("openapi schema not found")
	}

	spec, ok := defs.Properties["spec"]
	if !ok {
		return nil, fmt.Errorf("missed 'spec' in openapi schema")
	}

	err := p.parseObject(spec, "")
	if err != nil {
		return nil, err
	}

	for i, el := range p.fields {
		if contains(defs.Required, el.Name) {
			p.fields[i].Required = true
		}
	}

	sort.Slice(p.fields, func(i, j int) bool {
		return p.fields[i].Name < p.fields[j].Name
	})

	return p.fields, nil
}

func (p *parser) parseArray(el JSONSchemaProps, prefix string) error {
	var props JSONSchemaProps
	err := json.Unmarshal(el.Items.RawMessage, &props)
	if err != nil {
		return err
	}
	return p.parseObject(props, fmt.Sprintf("%s[0]", prefix))
}

func (p *parser) parseObject(el JSONSchemaProps, prefix string) error {
	for key, val := range el.Properties {
		name := strings.ReplaceAll(key, ".", "\\.")
		if len(prefix) > 0 {
			name = fmt.Sprintf("%s.%s", prefix, name)
		}

		switch val.Type {
		case TypeObject:
			err := p.parseObject(val, name)
			if err != nil {
				return err
			}
		case TypeArray:
			err := p.parseArray(val, name)
			if err != nil {
				return err
			}
		default:
			item := Field{
				Name:        name,
				Description: val.Description,
				Type:        val.Type,
				Required:    contains(el.Required, key),
			}

			if val.Default != nil {
				item.Default = val.Default.String()
			}

			p.fields = append(p.fields, item)
		}
	}

	return nil
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}
