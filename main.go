package main

import (
	"os"

	"github.com/krateoplatformops/krateo/cmd"
	"github.com/krateoplatformops/krateo/pkg/log"
)

func main() {
	ctl := cmd.KrateoCtl()
	if err := ctl.Execute(); err != nil {
		log.GetInstance().Error(err.Error())
		os.Exit(1)
	}
}
