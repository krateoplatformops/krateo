package main

import (
	"os"

	"github.com/krateoplatformops/krateo/cmd"
	"github.com/krateoplatformops/krateo/pkg/log"
)

// Build information. Populated at build-time.
var (
	Version string
	Build   string
)

func main() {
	ctl := cmd.KrateoCtl(Version, Build)
	if err := ctl.Execute(); err != nil {
		log.GetInstance().Error(err.Error())
		os.Exit(1)
	}
}
