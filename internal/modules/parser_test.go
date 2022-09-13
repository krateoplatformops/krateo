package modules

import (
	"fmt"
	"testing"

	v1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

func TestParse(t *testing.T) {
	src := "ghcr.io/krateoplatformops/krateo-module-core:latest"

	xrd, err := PullXRD(src, PullOpts{ModuleName: "core"})
	if err != nil {
		t.Fatal(err)
	}
	if xrd == nil {
		t.Log("XRD not found\n")
		return
	}

	spec, required, err := XRDSpecs(xrd)
	if err != nil {
		t.Fatal(err)
	}

	fields := []Field{}

	for _, k := range required {
		v := spec[k]
		flatten(k, v, &fields)
	}

	for _, el := range fields {
		fmt.Printf("%+v\n\n", el)
	}
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
