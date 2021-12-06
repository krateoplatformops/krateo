package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Build information. Populated at build-time.
var (
	Version string
	Build   string
)

// NewCmdVersion creates a command object for the "version" command
func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			if Version == "" {
				Version = "0.0.0"
			}

			if Build == "" {
				Build = "SNAPSHOT"
			}

			fmt.Printf("Krateo Platform Installer v%s (build: %s)", Version, Build)
		},
	}
}
