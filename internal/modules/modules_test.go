//go:build integration
// +build integration

package modules

import (
	"fmt"
	"testing"
)

func TestPullSchema(t *testing.T) {
	src := "ghcr.io/krateoplatformops/krateo-module-core:latest"

	spec, required, err := PullSchema(src, PullOpts{ModuleName: "core"})
	if err != nil {
		t.Fatal(err)
	}

	if spec == nil {
		return
	}

	fmt.Printf("REQUIRED:\n\t%v\n\n", required)

	for _, el := range required {
		vals := spec[el]
		fmt.Println(vals.Title)
		fmt.Printf("%+v\n\n", vals)
	}
}
