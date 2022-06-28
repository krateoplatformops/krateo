package uninstall

import (
	"fmt"

	"github.com/krateoplatformops/krateo/cmd/tools"
	"github.com/krateoplatformops/krateo/cmd/tools/flags"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/kubernetes"
	"github.com/krateoplatformops/krateo/pkg/log"
	"github.com/krateoplatformops/krateo/pkg/platform/crossplane/providers"
	"github.com/spf13/cobra"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type uninstallOptions struct {
	bus        eventbus.Bus
	kubeconfig string
	verbose    bool
}

func NewUninstallCmd() *cobra.Command {
	o := &uninstallOptions{}

	cmd := &cobra.Command{
		Use:                   "uninstall",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Short:                 "Uninstall Krateo",
		SilenceErrors:         true,
		Example:               "  krateo unnstall",
		RunE: func(cmd *cobra.Command, args []string) error {
			verbose, _ := cmd.Flags().GetBool(flags.Verbose)

			l := log.GetInstance()
			if verbose {
				l.SetLevel(log.DebugLevel)
			}

			handler := tools.LogEventHandler(l)
			o.bus = eventbus.New()
			eids := []eventbus.Subscription{
				o.bus.Subscribe(events.StartWaitEventID, handler),
				o.bus.Subscribe(events.StopWaitEventID, handler),
				o.bus.Subscribe(events.DoneEventID, handler),
				o.bus.Subscribe(events.DebugEventID, handler),
			}
			defer func() {
				for _, e := range eids {
					o.bus.Unsubscribe(e)
				}
			}()

			return o.Run()
		},
	}

	cmd.Flags().BoolVarP(&o.verbose, flags.Verbose, "v", false, "dump verbose output")
	cmd.Flags().StringVarP(&o.kubeconfig, clientcmd.RecommendedConfigPathFlag, "k",
		clientcmd.RecommendedHomeFile, "absolute path to the kubeconfig file")
	//cmd.Flags().DurationVarP(flags.Timeout, "t", time.Second*30, "time to wait for operations to complete; e.g. 1s, 2m, 3h")
	//cmd.Flags().Bool(flags.DryRun, false, "do not send requests to the API server")
	//nolint:errcheck
	//cmd.Flags().MarkHidden(flags.DryRun)

	return cmd
}

func (o *uninstallOptions) Run() error {
	rc, err := clientcmd.BuildConfigFromFlags("", o.kubeconfig)
	if err != nil {
		return err
	}

	dc, err := dynamic.NewForConfig(rc)
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewStartWaitEvent("Uninstalling Krateo modules"))
	err = o.uninstallModule(dc)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}
	o.bus.Publish(events.NewDoneEvent("Krateo module uninstalled"))

	o.bus.Publish(events.NewStartWaitEvent("Uninstalling CRDs"))
	err = o.deleteCRDs(rc)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}
	o.bus.Publish(events.NewDoneEvent("CRDs uninstalled"))

	o.bus.Publish(events.NewStartWaitEvent("Uninstalling Crossplane providers"))
	err = deleteCrossplaneProviders(dc, false)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}
	o.bus.Publish(events.NewDoneEvent("Crossplane providers uninstalled"))

	o.bus.Publish(events.NewStartWaitEvent("Uninstalling providers role bindings"))
	err = providers.DeleteClusterRoleBindings(dc)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}
	o.bus.Publish(events.NewDoneEvent("Crossplane providers role binding uninstalled"))

	o.bus.Publish(events.NewStartWaitEvent("Uninstalling Crossplane"))
	err = uninstallCrossplaneChart(rc, o.bus, o.verbose)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}
	o.bus.Publish(events.NewDoneEvent("Crossplane uninstalled"))

	/*
		o.bus.Publish(events.NewStartWaitEvent(fmt.Sprintf("Deleting namespace '%s'", kubernetes.CrossplaneSystemNamespace)))
		err = deleteNamespaceForcingFinalizers(rc, false, kubernetes.CrossplaneSystemNamespace)
		if err != nil && o.verbose {
			o.bus.Publish(events.NewWarningEvent(err.Error()))
		}
		o.bus.Publish(events.NewDoneEvent(fmt.Sprintf("Namespace '%s' deleted", kubernetes.CrossplaneSystemNamespace)))
	*/

	err = deleteCrossplaneProviders(dc, false)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}

	o.bus.Publish(events.NewStartWaitEvent(fmt.Sprintf("Deleting namespace '%s'", kubernetes.KrateoSystemNamespace)))
	err = deleteNamespaceForcingFinalizers(rc, false, kubernetes.KrateoSystemNamespace)
	if err != nil && o.verbose {
		o.bus.Publish(events.NewWarningEvent(err.Error()))
	}
	o.bus.Publish(events.NewDoneEvent(fmt.Sprintf("Namespace '%s' deleted", kubernetes.KrateoSystemNamespace)))

	return nil
}

func (o *uninstallOptions) uninstallModule(dc dynamic.Interface) error {
	gvrs, err := listXRs(dc)
	if err != nil {
		return err
	}

	if o.verbose {
		msg := fmt.Sprintf("Found [%d] krateo composite resources", len(gvrs))
		o.bus.Publish(events.NewDebugEvent(msg))
		for name, el := range gvrs {
			msg = fmt.Sprintf("%s (%s/%s - %s)", name, el.Group, el.Version, el.Resource)
			o.bus.Publish(events.NewDebugEvent(msg))
		}
	}

	return deleteXRs(dc, gvrs)
}

func (o *uninstallOptions) deleteCRDs(rc *rest.Config) error {
	cli, err := kubernetes.Crds(rc)
	if err != nil {
		return err
	}

	crds, err := listCRDs(cli)
	if err != nil {
		return err
	}

	if o.verbose {
		if tot := len(crds); tot > 0 {
			o.bus.Publish(events.NewDebugEvent(fmt.Sprintf("Found [%d] CRDs to remove...", tot)))
			for _, x := range crds {
				o.bus.Publish(events.NewDebugEvent(x))
			}
		}
	}

	return patchAndDeleteCRDs(cli, false, crds)
}
