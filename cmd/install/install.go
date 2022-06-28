package install

import (
	"errors"
	"fmt"
	"os"

	"github.com/krateoplatformops/krateo/cmd/tools"
	"github.com/krateoplatformops/krateo/cmd/tools/flags"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/log"
	"github.com/krateoplatformops/krateo/pkg/platform"
	"github.com/krateoplatformops/krateo/pkg/storage"
	"github.com/krateoplatformops/krateo/pkg/storage/git"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
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

			o := platform.InstallOptions{
				Bus:        eventbus.New(),
				Verbose:    false,
				ModuleName: args[0],
				Platform:   "kubernetes",
			}

			kubeconfig, err := cmd.Flags().GetString(clientcmd.RecommendedConfigPathFlag)
			if err != nil {
				return err
			}

			o.Verbose, _ = cmd.Flags().GetBool(flags.Verbose)
			o.Platform, _ = cmd.Flags().GetString(flags.Platform)

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

			o.ModulePackage, err = pullModulePackage(o.ModuleName, o.Bus, o.Verbose)
			if err != nil {
				return err
			}

			gitUrl, _ := cmd.Flags().GetString(flags.GitURL)
			gitToken, _ := cmd.Flags().GetString(flags.GitToken)

			o.ModuleClaims, err = pullModuleClaims(o.ModuleName, gitUrl, gitToken, o.Bus, o.Verbose)
			if err != nil {
				return err
			}

			return platform.Install(&o)
		},
	}

	cmd.Flags().StringP(flags.GitURL, "r", "", "git repository url for pushing module configuration")
	cmd.Flags().StringP(flags.GitToken, "t", "", "token for git repository authentication")
	cmd.Flags().StringP(flags.Platform, "p", "kubernetes", "platform selector; i.e. openshift, kubernetes")
	cmd.Flags().BoolP(flags.Verbose, "v", false, "dump verbose output")

	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	// if KUBECONFIG is empty
	if len(defaultKubeconfig) == 0 {
		// look for file $HOME/.kube/config
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}
	cmd.Flags().StringP(clientcmd.RecommendedConfigPathFlag, "k",
		defaultKubeconfig, "absolute path to the kubeconfig file")
	//nolint:errcheck
	cmd.MarkFlagRequired(flags.GitURL)
	//nolint:errcheck
	cmd.MarkFlagRequired(flags.GitToken)

	return cmd
}

func pullModulePackage(name string, bus eventbus.Bus, verbose bool) ([]byte, error) {
	entry, err := tools.GitPullModulePackage(name)
	if err != nil {
		return nil, err
	}

	if verbose {
		bus.Publish(events.NewDebugEvent("pulled '%s' (rev: %s)\n", entry.Meta.Name, entry.Meta.Version[0:8]))
		bus.Publish(events.NewDebugEvent(string(entry.Content)))
	}

	return entry.Content, nil
}

func pullModuleClaims(name string, gitUrl, gitToken string, bus eventbus.Bus, verbose bool) ([]byte, error) {
	entry, err := tools.GitPullModuleDefaultsFormUserRepo(name, gitUrl, git.WithGitToken(gitToken))
	if err != nil {
		if !errors.Is(err, storage.ErrEntryNotFound) {
			return nil, err
		}

		entry, err = tools.GitPullModuleDefaults(name)
		if err != nil {
			return nil, err
		}

		err = tools.GitPushEntry(gitUrl, entry, git.WithGitToken(gitToken))
		if err != nil {
			return nil, err
		}
	}

	if verbose {
		bus.Publish(events.NewDebugEvent("pulled '%s' (rev: %s)\n", entry.Meta.Name, entry.Meta.Version[0:8]))
	}

	return entry.Content, nil
}
