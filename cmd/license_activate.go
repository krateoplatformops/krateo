package cmd

import (
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo/cmd/flags"
	"github.com/krateoplatformops/krateo/pkg/actions"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/log"
	"github.com/spf13/cobra"
)

func NewLicenseActivateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                   "activate <ORDER NUMBER>",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ArbitraryArgs,
		Short:                 "Activate the Krateo license",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("order number of the license purchase is mandatory")
			}

			cfg, err := cmd.Flags().GetString(flags.Kubeconfig)
			if err != nil {
				return err
			}

			l := log.GetInstance()

			bus := eventbus.New()
			eids := []eventbus.Subscription{
				bus.Subscribe(events.StartWaitEventID, updateLog(l)),
				bus.Subscribe(events.StopWaitEventID, updateLog(l)),
				bus.Subscribe(events.DoneEventID, updateLog(l)),
			}
			defer func() {
				for _, e := range eids {
					bus.Unsubscribe(e)
				}
			}()

			orderNr := strings.TrimSpace(args[0])

			act := actions.NewLicenseActivate(bus, cfg, orderNr)
			return act.Run()
		},
	}

	cmd.Flags().StringP(flags.Kubeconfig, "k", flags.DefaultKubeconfigValue(), "absolute path to the kubeconfig file")

	return cmd
}
