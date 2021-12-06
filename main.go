package main

import (
	"os"

	"github.com/krateoplatformops/krateoctl/cmd"
	"github.com/krateoplatformops/krateoctl/pkg/log"
)

func main() {
	ctl := cmd.KrateoCtl()
	if err := ctl.Execute(); err != nil {
		log.GetInstance().Error(err.Error())
		os.Exit(1)
	}
}
