package cmd

import (
	"github.com/krateoplatformops/krateoctl/cmd/flags"
	"github.com/krateoplatformops/krateoctl/pkg/actions"
	"github.com/krateoplatformops/krateoctl/pkg/eventbus"
	"github.com/krateoplatformops/krateoctl/pkg/events"
	"github.com/krateoplatformops/krateoctl/pkg/log"
	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "init",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: false,
		Args:                  cobra.NoArgs,
		Short:                 "Initialize Krateo Platform",
		RunE: func(cmd *cobra.Command, args []string) error {
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
			act := actions.NewInit(bus, cfg)
			act.SetVerboseEnabled(verbose)
			return act.Run()
		},
	}

	cmd.Flags().StringP(flags.Kubeconfig, "k", flags.DefaultKubeconfigValue(), "absolute path to the kubeconfig file")
	cmd.Flags().BoolP(flags.Verbose, "v", false, "dump verbose output")

	return cmd
}
