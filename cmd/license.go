package cmd

import (
	"github.com/spf13/cobra"
)

func NewLicenseCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "license",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: false,
		Args:                  cobra.NoArgs,
		Short:                 "Manage Krateo License",
		SilenceErrors:         true,
	}

	cmd.AddCommand(NewLicenseActivateCmd())
	//cmd.AddCommand(NewLicenseCheckCmd())

	return cmd
}
