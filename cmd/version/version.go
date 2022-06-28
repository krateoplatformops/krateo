package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdVersion creates a command object for the "version" command
func NewCmdVersion(ver, build string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the current build information",
		Long:  ``,
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Printf("Krateo PlatformOps Installer v%s (build: %s)\n", ver, build)
		},
	}
}
