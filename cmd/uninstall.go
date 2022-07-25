package cmd

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/krateoplatformops/krateo/internal/clusterrolebindings"
	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/krateoplatformops/krateo/internal/crds"
	"github.com/krateoplatformops/krateo/internal/crossplane"
	"github.com/krateoplatformops/krateo/internal/crossplane/compositions"
	"github.com/krateoplatformops/krateo/internal/crossplane/configurations"
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
	cmd.Flags().BoolVar(&o.dryRun, "dry-run", false, "preview the object that would be deleted, without really deleting it")
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
	dryRun     bool
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

	if err := o.deleteModules(ctx); err != nil {
		return err
	}

	if err := o.deletePackages(ctx); err != nil {
		return err
	}

	if err := o.deleteControllerConfigs(ctx); err != nil {
		return err
	}

	if err := o.deleteCrossplane(ctx); err != nil {
		return err
	}

	o.bus.Publish(events.NewStartWaitEvent("finishing cleaning..."))
	o.deletCompositions(ctx)
	o.deleteCRDsQuietly(ctx)
	o.deleteClusterRoleBindingsQuietly(ctx)
	o.bus.Publish(events.NewStartWaitEvent("cleaning done"))

	return nil
}

func (o *uninstallOpts) deleteCrossplane(ctx context.Context) error {
	pod, err := crossplane.GetPOD(ctx, o.restConfig)
	if err != nil {
		return err
	}
	if pod == nil {
		if o.verbose {
			o.bus.Publish(events.NewDebugEvent("crossplane not found"))
		}
		return nil
	}

	if o.dryRun {
		o.bus.Publish(events.NewDebugEvent(
			"found crossplane pod: %s in namespace: %s",
			pod.GetName(), pod.GetNamespace()))
		return nil
	}

	o.bus.Publish(events.NewStartWaitEvent("uninstalling crossplane %s...", crossplane.ChartVersion))

	err = crossplane.Uninstall(crossplane.UninstallOpts{
		RESTConfig: o.restConfig,
		EventBus:   o.bus,
		Namespace:  pod.GetNamespace(),
		Verbose:    o.verbose,
	})
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewDoneEvent("crossplane %s uninstalled", crossplane.ChartVersion))

	return nil
}

func (o *uninstallOpts) deletePackages(ctx context.Context) error {
	all, err := providers.List(ctx, o.restConfig)
	if err != nil {
		return err
	}

	if len(all) == 0 {
		return nil
	}

	if o.dryRun {
		o.bus.Publish(events.NewDebugEvent("found [%d] packages", len(all)))
	}

	for _, el := range all {
		if o.dryRun {
			o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
			continue
		}

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

func (o *uninstallOpts) deleteControllerConfigs(ctx context.Context) error {
	all, err := controllerconfigs.ListAll(ctx, o.restConfig)
	if err != nil {
		return err
	}

	if len(all) == 0 {
		return nil
	}

	if o.dryRun {
		o.bus.Publish(events.NewDebugEvent("found [%d] controller configs", len(all)))
	}

	for _, el := range all {
		if o.dryRun {
			o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
			continue
		}

		err := controllerconfigs.Delete(ctx, controllerconfigs.DeleteOpts{
			RESTConfig: o.restConfig,
			Name:       el.GetName(),
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (o *uninstallOpts) deleteModules(ctx context.Context) error {
	all, err := configurations.List(ctx, o.restConfig)
	if err != nil {
		return err
	}
	if len(all) == 0 {
		return nil
	}

	if o.dryRun {
		o.bus.Publish(events.NewDebugEvent("found [%d] modules", len(all)))
	}

	for _, el := range all {
		if o.dryRun {
			o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
			continue
		}

		o.bus.Publish(events.NewStartWaitEvent("uninstalling module %s...", el.GetName()))
		err := core.Delete(ctx, core.DeleteOpts{
			RESTConfig: o.restConfig,
			Object:     &el,
		})
		if err != nil {
			return err
		}
		o.bus.Publish(events.NewDoneEvent("module %s uninstalled", el.GetName()))
	}

	return nil
}

func (o *uninstallOpts) deletCompositions(ctx context.Context) {
	all, err := compositions.List(ctx, o.restConfig)
	if err != nil {
		return
	}

	if len(all) == 0 {
		return
	}

	if o.dryRun {
		o.bus.Publish(events.NewDebugEvent("found [%d] compositions", len(all)))
	}

	for _, el := range all {
		if o.dryRun {
			o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
			continue
		}
		_ = core.Delete(ctx, core.DeleteOpts{
			RESTConfig: o.restConfig,
			Object:     &el,
		})
	}
}

func (o *uninstallOpts) deleteCRDsQuietly(ctx context.Context) {
	/*
		items, err := crds.Instances(ctx, o.restConfig)
		if err == nil {
			if tot := len(items); tot > 0 && o.dryRun {
				o.bus.Publish(events.NewDebugEvent("found [%d] crds", tot))
			}

			for _, el := range items {
				if o.dryRun {
					o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
					continue
				}

				crds.PatchAndDelete(ctx, o.restConfig, &el)
			}
		}
	*/

	items, err := crds.List(ctx, crds.ListOpts{RESTConfig: o.restConfig})
	if err == nil {
		if tot := len(items); tot > 0 && o.dryRun {
			o.bus.Publish(events.NewDebugEvent("found [%d] crds", tot))
		}

		for _, el := range items {
			if o.dryRun {
				o.bus.Publish(events.NewDebugEvent(" > %s", el.GetName()))
				continue
			}

			crds.PatchAndDelete(ctx, o.restConfig, &el)
		}
	}
}

func (o *uninstallOpts) deleteClusterRoleBindingsQuietly(ctx context.Context) {
	all, err := clusterrolebindings.List(ctx, o.restConfig)
	if err != nil {
		return
	}

	res, err := core.Filter(all, func(obj unstructured.Unstructured) bool {
		accept := (obj.GetName() == "provider-helm-admin-binding")
		accept = accept || (obj.GetName() == "provider-kubernetes-admin-binding")
		return accept
	})

	if len(res) == 0 {
		return
	}

	if o.dryRun {
		o.bus.Publish(events.NewDebugEvent("found [%d] cluster role bindings", len(res)))
	}

	for _, el := range res {
		if o.dryRun {
			o.bus.Publish(events.NewDebugEvent("> %s", el.GetName()))
			continue
		}
		_ = clusterrolebindings.Delete(ctx, clusterrolebindings.DeleteOpts{
			RESTConfig: o.restConfig,
			Name:       el.GetName(),
		})
	}
}
