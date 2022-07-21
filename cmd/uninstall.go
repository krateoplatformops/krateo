package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/krateoplatformops/krateo/internal/crossplane"
	"github.com/krateoplatformops/krateo/internal/crossplane/compositions"
	"github.com/krateoplatformops/krateo/internal/crossplane/controllerconfigs"
	"github.com/krateoplatformops/krateo/internal/crossplane/providers"
	"github.com/krateoplatformops/krateo/internal/eventbus"
	"github.com/krateoplatformops/krateo/internal/events"
	"github.com/krateoplatformops/krateo/internal/log"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func newUninstallCmd() *cobra.Command {
	o := uninstallOpts{
		bus:     eventbus.New(),
		verbose: false,
	}

	cmd := &cobra.Command{
		Use:                   "uninstall",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.NoArgs,
		Short:                 "Uninstall Krateo",
		SilenceErrors:         true,
		Example:               "  krateo uninstall",
		RunE: func(cmd *cobra.Command, args []string) error {
			l := log.GetInstance()
			if o.verbose {
				l.SetLevel(log.DebugLevel)
			}

			handler := events.LogHandler(l)
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

			if err := o.complete(); err != nil {
				return err
			}

			return o.run()
		},
	}

	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if len(defaultKubeconfig) == 0 {
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	cmd.Flags().BoolVarP(&o.verbose, "verbose", "v", false, "dump verbose output")
	cmd.Flags().StringVar(&o.kubeconfig, clientcmd.RecommendedConfigPathFlag, defaultKubeconfig, "absolute path to the kubeconfig file")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "n", "default", "namespace where to install krateo runtime")

	return cmd
}

type uninstallOpts struct {
	kubeconfig string
	bus        eventbus.Bus
	restConfig *rest.Config
	namespace  string
	verbose    bool
}

func (o *uninstallOpts) complete() (err error) {
	yml, err := ioutil.ReadFile(o.kubeconfig)
	if err != nil {
		return err
	}

	o.restConfig, err = core.RESTConfigFromBytes(yml)
	if err != nil {
		return err
	}

	return nil
}

func (o *uninstallOpts) run() error {
	ctx := context.TODO()

	if err := o.uninstallCompositions(ctx); err != nil {
		return err
	}

	if err := o.uninstallPackages(ctx); err != nil {
		return err
	}

	if err := o.uninstallControllerConfigs(ctx); err != nil {
		return err
	}

	if err := o.uninstallCrossplane(ctx); err != nil {
		return err
	}

	return nil
}

func (o *uninstallOpts) uninstallCrossplane(ctx context.Context) error {
	ok, err := crossplane.Exists(ctx, crossplane.ExistOpts{
		RESTConfig: o.restConfig,
		Namespace:  o.namespace,
	})
	if err != nil {
		return err
	}
	if !ok {
		if o.verbose {
			o.bus.Publish(events.NewDebugEvent("crossplane not found in namespace '%s'", o.namespace))
		}
		return nil
	}

	o.bus.Publish(events.NewStartWaitEvent("uninstalling crossplane %s...", crossplane.ChartVersion))

	err = crossplane.Uninstall(crossplane.UninstallOpts{
		RESTConfig: o.restConfig,
		EventBus:   o.bus,
		Namespace:  o.namespace,
		Verbose:    o.verbose,
	})
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewDoneEvent("crossplane %s uninstalled", crossplane.ChartVersion))

	return nil
}

func (o *uninstallOpts) uninstallPackages(ctx context.Context) error {
	all, err := providers.All(ctx, o.restConfig)
	if err != nil {
		return err
	}

	if o.verbose {
		o.bus.Publish(events.NewDebugEvent("found [%d] packages", len(all)))
	}

	if len(all) == 0 {
		return nil
	}

	for _, el := range all {
		o.bus.Publish(events.NewStartWaitEvent("uninstalling package %s...", el.GetName()))
		err := core.Delete(ctx, core.DeleteOpts{
			RESTConfig: o.restConfig,
			Object:     &el,
		})
		if err != nil {
			return err
		}

		// Start Watching
		req, err := labels.NewRequirement(core.PackageNameLabel, selection.Equals, []string{el.GetName()})
		if err != nil {
			return err
		}

		sel := labels.NewSelector()
		sel = sel.Add(*req)

		err = core.Watch(ctx, core.WatchOpts{
			RESTConfig: o.restConfig,
			Selector:   sel,
			GVR:        schema.GroupVersionResource{Version: "v1", Resource: "pods"},
			Namespace:  el.GetNamespace(),
			StopFunc: func(et watch.EventType, obj *unstructured.Unstructured) (bool, error) {
				return (et == watch.Deleted), nil
			},
		})
		if err != nil {
			return err
		}
		o.bus.Publish(events.NewDoneEvent("package %s uninstalled", el.GetName()))
	}

	return nil
}

func (o *uninstallOpts) uninstallControllerConfigs(ctx context.Context) error {
	all, err := controllerconfigs.ListAll(ctx, o.restConfig)
	if err != nil {
		return err
	}

	if o.verbose {
		o.bus.Publish(events.NewDebugEvent("found [%d] controller configs", len(all)))
	}

	if len(all) == 0 {
		return nil
	}

	o.bus.Publish(events.NewStartWaitEvent("deleting controller configs"))
	for _, el := range all {
		if o.verbose {
			o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
		}
		err := controllerconfigs.Delete(ctx, controllerconfigs.DeleteOpts{
			RESTConfig: o.restConfig,
			Name:       el.GetName(),
		})
		if err != nil {
			return err
		}
	}
	o.bus.Publish(events.NewDoneEvent("controller configs deleted"))

	return nil
}

func (o *uninstallOpts) uninstallCompositions(ctx context.Context) error {
	all, err := compositions.List(ctx, o.restConfig)
	if err != nil {
		return err
	}

	if o.verbose {
		o.bus.Publish(events.NewDebugEvent("found [%d] compositions", len(all)))
	}

	if len(all) == 0 {
		return nil
	}

	for _, el := range all {
		o.bus.Publish(events.NewStartWaitEvent("uninstalling composition %s...", el.GetName()))
		err := core.Delete(ctx, core.DeleteOpts{
			RESTConfig: o.restConfig,
			Object:     &el,
		})
		if err != nil {
			return err
		}

		// Start Watching
		/*
			err = core.Watch(ctx, core.WatchOpts{
				RESTConfig: o.restConfig,
				GVR: schema.GroupVersionResource{
					Group:    "apiextensions.crossplane.io",
					Version:  "v1",
					Resource: "compositions",
				},
				Namespace: el.GetNamespace(),
				StopFunc: func(et watch.EventType, obj *unstructured.Unstructured) (bool, error) {
					return (obj.GetName() == el.GetName() && et == watch.Deleted), nil
				},
			})
			if err != nil {
				return err
			}
		*/
		o.bus.Publish(events.NewDoneEvent("composition %s uninstalled", el.GetName()))
	}

	return nil
}
