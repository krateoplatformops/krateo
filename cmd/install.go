package cmd

import (
	"fmt"

	"github.com/krateoplatformops/krateoctl/cmd/flags"
	"github.com/krateoplatformops/krateoctl/pkg/actions"
	"github.com/krateoplatformops/krateoctl/pkg/eventbus"
	"github.com/krateoplatformops/krateoctl/pkg/events"
	"github.com/krateoplatformops/krateoctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewInstallCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:                   "install <MODULE>",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ArbitraryArgs,
		Short:                 "Install Krateo modules",
		SilenceErrors:         true,
		Example:               "  krateo install core",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missed module name")
			}

			verbose, _ := cmd.Flags().GetBool(flags.Verbose)

			l := log.GetInstance()
			if verbose {
				l.SetLevel(log.DebugLevel)
			}

			handler := updateLog(l)
			bus := eventbus.New()
			eids := []eventbus.Subscription{
				bus.Subscribe(events.StartWaitEventID, handler),
				bus.Subscribe(events.StopWaitEventID, handler),
				bus.Subscribe(events.DoneEventID, handler),
				bus.Subscribe(events.DebugEventID, handler),
			}
			defer func() {
				for _, e := range eids {
					bus.Unsubscribe(e)
				}
			}()

			cfg, err := cmd.Flags().GetString(flags.Kubeconfig)
			if err != nil {
				return err
			}

			repoURL, err := cmd.Flags().GetString(flags.RepoURL)
			if err != nil {
				return err
			}

			repoToken, err := cmd.Flags().GetString(flags.RepoToken)
			if err != nil {
				return err
			}

			act := actions.NewInstall(bus, args[0], cfg)
			act.SetVerboseEnabled(verbose)
			act.SetModuleRepoURL(repoURL)
			act.SetModuleRepoToken(repoToken)
			return act.Run()
		},
	}

	cmd.Flags().StringP(flags.Kubeconfig, "k", flags.DefaultKubeconfigValue(), "absolute path to the kubeconfig file")
	cmd.Flags().StringP(flags.RepoURL, "r", "", "url of the git repository where the module resides")
	cmd.Flags().StringP(flags.RepoToken, "t", "", "token for git repository authentication")
	cmd.Flags().BoolP(flags.Verbose, "v", false, "dump verbose output")
	//nolint:errcheck
	cmd.MarkFlagRequired(flags.RepoURL)

	return cmd
}
