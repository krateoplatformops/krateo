package modules

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	xpextv1 "github.com/crossplane/crossplane/apis/apiextensions/v1"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
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

type parser struct {
	fields []Field
}

func (p *parser) parse(src *xpextv1.CompositeResourceDefinition) ([]Field, error) {
	if len(src.Spec.Versions) == 0 {
		return nil, fmt.Errorf("no spec version found")
	}

	v := src.Spec.Versions[0].Schema //.(*xpextv1.CompositeResourceValidation)
	if v == nil {
		return nil, fmt.Errorf("no OpenAPI schema found")
	}

	s := &extv1.JSONSchemaProps{}
	if err := json.Unmarshal(v.OpenAPIV3Schema.Raw, s); err != nil {
		return nil, err
	}

	err := p.parseObject(s, "")
	if err != nil {
		return nil, err
	}

	for i, el := range p.fields {
		if contains(s.Required, el.Name) {
			p.fields[i].Required = true
		}
	}

	sort.Slice(p.fields, func(i, j int) bool {
		return p.fields[i].Name < p.fields[j].Name
	})

	return p.fields, nil
}

func (p *parser) parseArray(el *extv1.JSONSchemaProps, prefix string) error {
	return nil
	//var props extv1.JSONSchemaProps
	//err := json.Unmarshal(el.Items.Schema, &props)
	//if err != nil {
	//	return err
	//}
	//return p.parseObject(props, fmt.Sprintf("%s[0]", prefix))
}

func (p *parser) parseObject(el *extv1.JSONSchemaProps, prefix string) error {
	for key, val := range el.Properties {
		name := strings.ReplaceAll(key, ".", "\\.")
		if len(prefix) > 0 {
			name = fmt.Sprintf("%s.%s", prefix, name)
		}

		switch val.Type {
		case TypeObject:
			err := p.parseObject(&val, name)
			if err != nil {
				return err
			}
		case TypeArray:
			err := p.parseArray(&val, name)
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
