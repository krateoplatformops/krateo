package init

import (
	"os"

	"github.com/krateoplatformops/krateo/cmd/tools"
	"github.com/krateoplatformops/krateo/cmd/tools/flags"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/log"
	"github.com/krateoplatformops/krateo/pkg/platform"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func NewInitCmd() *cobra.Command {
	o := platform.InitOptions{
		Bus:     eventbus.New(),
		Verbose: false,
	}

	cmd := &cobra.Command{
		Use:                   "init",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: false,
		Args:                  cobra.NoArgs,
		Short:                 "Initialize Krateo Platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeconfig, err := cmd.Flags().GetString(clientcmd.RecommendedConfigPathFlag)
			if err != nil {
				return err
			}

			o.Verbose, _ = cmd.Flags().GetBool(flags.Verbose)
			o.HttpProxy, _ = cmd.Flags().GetString(flags.HttpProxy)
			o.HttpsProxy, _ = cmd.Flags().GetString(flags.HttpsProxy)
			o.NoProxy, _ = cmd.Flags().GetString(flags.NoProxy)

			o.Config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				return err
			}

			l := log.GetInstance()
			if o.Verbose {
				l.SetLevel(log.DebugLevel)
			}

			handler := tools.LogEventHandler(l)

			eids := []eventbus.Subscription{
				o.Bus.Subscribe(events.StartWaitEventID, handler),
				o.Bus.Subscribe(events.StopWaitEventID, handler),
				o.Bus.Subscribe(events.DoneEventID, handler),
				o.Bus.Subscribe(events.DebugEventID, handler),
			}
			defer func() {
				for _, e := range eids {
					o.Bus.Unsubscribe(e)
				}
			}()

			return platform.Init(o)
		},
	}

	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	// if KUBECONFIG is empty
	if len(defaultKubeconfig) == 0 {
		// look for file $HOME/.kube/config
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	cmd.Flags().BoolP(flags.Verbose, "v", false, "dump verbose output")
	cmd.Flags().String(clientcmd.RecommendedConfigPathFlag, defaultKubeconfig, "absolute path to the kubeconfig file")
	cmd.Flags().String(flags.HttpProxy, os.Getenv("HTTP_PROXY"), "use the specified HTTP proxy")
	cmd.Flags().String(flags.HttpsProxy, os.Getenv("HTTPS_PROXY"), "use the specified HTTPS proxy")
	cmd.Flags().String(flags.NoProxy, os.Getenv("NO_PROXY"), "comma-separated list of hosts and domains which do not use the proxy")

	return cmd
}
