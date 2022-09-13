//go:build integration
// +build integration

package modules

import (
	"fmt"
	"testing"
)

func TestPullXRD(t *testing.T) {
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
	if spec == nil {
		return
	}

	fmt.Printf("REQUIRED:\n\t%v\n\n", required)

	for _, el := range required {
		vals := spec[el]
		fmt.Println("Descr.  ", vals.Description)
		fmt.Println("Type    ", vals.Type)
		fmt.Println("Default ", vals.Default)
	}

	fmt.Print("\n\n")
	/*
		for k, v := range spec {
			fmt.Printf("Key  %s\n", k)
			fmt.Printf("Val  %+v\n\n", v)
		}*/

	k := "kongapigw"
	v := spec[k]
	fmt.Printf("Key  %s\n", k)
	//fmt.Printf("Val  %+v\n\n", v)

	for kk, vv := range v.Properties {
		fmt.Printf("Key  %s\n", kk)
		fmt.Printf("Val  %+v\n\n", vv)
		fmt.Print("\n\n")
	}

}
