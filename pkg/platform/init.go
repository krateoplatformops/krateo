package platform

import (
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/platform/crossplane"
	"github.com/krateoplatformops/krateo/pkg/platform/crossplane/providers"
	"github.com/krateoplatformops/krateo/pkg/platform/namespaces"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type InitOptions struct {
	Config     *rest.Config
	Bus        eventbus.Bus
	Verbose    bool
	HttpProxy  string
	HttpsProxy string
	NoProxy    string
}

func Init(opts InitOptions) error {
	dc, err := dynamic.NewForConfig(opts.Config)
	if err != nil {
		return err
	}

	// Create Crossplane namespace
	opts.Bus.Publish(events.NewStartWaitEvent("creating namespace '%s'", namespaces.CrossplaneSystem))
	err = namespaces.Create(dc, namespaces.CrossplaneSystem)
	if err != nil {
		return err
	}
	opts.Bus.Publish(events.NewDoneEvent("namespace '%s' successfully created", namespaces.CrossplaneSystem))

	// Install Crossplane
	err = opts.installCrossplaneEventually(dc)
	if err != nil {
		return err
	}

	for _, v := range providers.All() {
		// Install controller config
		opts.Bus.Publish(events.NewStartWaitEvent("creating provider %s controller configuration...", v.Name()))
		err = providers.CreateControllerConfig(dc, v)
		if err != nil {
			return err
		}
		opts.Bus.Publish(events.NewStartWaitEvent("provider %s controller configuration successfully created", v.Name()))

		opts.Bus.Publish(events.NewStartWaitEvent("installing provider %s (%s)...", v.Name(), v.Version()))
		if opts.Verbose {
			opts.Bus.Publish(events.NewDebugEvent("> image: %s", v.Image()))
		}
		err = providers.InstallEventually(dc, v)
		if err != nil {
			return err
		}
		opts.Bus.Publish(events.NewDoneEvent("provider %s (%s) successfully installed", v.Name(), v.Version()))
	}

	opts.Bus.Publish(events.NewStartWaitEvent("creating role bindings for crossplane providers..."))
	err = providers.CreateClusterRoleBindings(dc)
	if err != nil {
		return err
	}
	opts.Bus.Publish(events.NewStartWaitEvent("role bindings for crossplane providers successfully created"))

	/*
		opts.Bus.Publish(events.NewStartWaitEvent("creating namespace '%s'...", namespaces.KrateoSystem))
		err = namespaces.Create(dc, namespaces.KrateoSystem)
		if err != nil {
			return err
		}
		opts.Bus.Publish(events.NewDoneEvent("namespace '%s' successfully created", namespaces.KrateoSystem))
	*/
	return nil
}

func (o *InitOptions) installCrossplaneEventually(dc dynamic.Interface) error {
	o.Bus.Publish(events.NewStartWaitEvent("installing crossplane %s...", crossplane.ChartVersion))
	ok, err := crossplane.IsInstalled(dc)
	if err != nil {
		return err
	}

	if ok {
		o.Bus.Publish(events.NewDoneEvent("crossplane %s already installed", crossplane.ChartVersion))
		return nil
	}

	err = crossplane.InstallChart(o.Config, o.Bus, crossplane.ChartOpts{
		Verbose:    o.Verbose,
		HttpProxy:  o.HttpProxy,
		HttpsProxy: o.HttpsProxy,
		NoProxy:    o.NoProxy,
	})
	if err != nil {
		return err
	}
	o.Bus.Publish(events.NewDoneEvent("crossplane %s successfully installed", crossplane.ChartVersion))

	o.Bus.Publish(events.NewStartWaitEvent("waiting for crossplane pod ready..."))
	err = crossplane.WaitUntilReady(dc)
	if err != nil {
		return err
	}
	o.Bus.Publish(events.NewDoneEvent("crossplane pod up and running"))

	return nil
}
