package helm

import (
	"errors"
	"io"
	"strings"
	"time"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	helmDriver = "secret"
)

type InstallOptions struct {
	RESTConfig  *rest.Config
	Namespace   string
	ReleaseName string
	ChartSource io.Reader
	ChartValues map[string]interface{}
	LogFn       func(format string, v ...interface{})
}

func Install(opts InstallOptions) error {
	rg := newRESTClientGetter(opts.RESTConfig, opts.Namespace)

	actionConfig := new(action.Configuration)
	err := actionConfig.Init(rg, opts.Namespace, helmDriver, opts.LogFn)
	if err != nil {
		return err
	}

	chart, err := loader.LoadArchive(opts.ChartSource)
	if err != nil {
		return err
	}

	iCli := action.NewInstall(actionConfig)
	iCli.Namespace = opts.Namespace
	iCli.ReleaseName = opts.ReleaseName
	iCli.Wait = false
	iCli.Timeout = 10 * time.Second

	_, err = iCli.Run(chart, opts.ChartValues)
	if err != nil {
		if !strings.HasPrefix(err.Error(), "cannot re-use a name that is still in use") {
			return err
		}
	}

	return nil
}

type UninstallOptions struct {
	RESTConfig  *rest.Config
	Namespace   string
	ReleaseName string
	LogFn       func(format string, v ...interface{})
}

func Uninstall(opts UninstallOptions) error {
	rg := newRESTClientGetter(opts.RESTConfig, opts.Namespace)

	actionConfig := new(action.Configuration)
	err := actionConfig.Init(rg, opts.Namespace, helmDriver, opts.LogFn)
	if err != nil {
		return err
	}

	act := action.NewUninstall(actionConfig)
	_, err = act.Run(opts.ReleaseName)
	if !errors.Is(err, driver.ErrReleaseNotFound) {
		return err
	}

	return nil
}

type restClientGetter struct {
	Namespace string
	config    *rest.Config
}

func newRESTClientGetter(config *rest.Config, namespace string) *restClientGetter {
	return &restClientGetter{
		Namespace: namespace,
		config:    config,
	}
}

func (c *restClientGetter) ToRESTConfig() (*rest.Config, error) {
	return c.config, nil
}

func (c *restClientGetter) ToDiscoveryClient() (discovery.CachedDiscoveryInterface, error) {
	config, err := c.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	// The more groups you have, the more discovery requests you need to make.
	// given 25 groups (our groups + a few custom conf) with one-ish version each, discovery needs to make 50 requests
	// double it just so we don't end up here again for a while.  This config is only used for discovery.
	// config.Burst = 100

	discoveryClient, _ := discovery.NewDiscoveryClientForConfig(config)
	return memory.NewMemCacheClient(discoveryClient), nil
}

func (c *restClientGetter) ToRESTMapper() (meta.RESTMapper, error) {
	discoveryClient, err := c.ToDiscoveryClient()
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(discoveryClient)
	expander := restmapper.NewShortcutExpander(mapper, discoveryClient)
	return expander, nil
}

func (c *restClientGetter) ToRawKubeConfigLoader() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// use the standard defaults for this client command
	// DEPRECATED: remove and replace with something more accurate
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig

	overrides := &clientcmd.ConfigOverrides{ClusterDefaults: clientcmd.ClusterDefaults}
	overrides.Context.Namespace = c.Namespace

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, overrides)
}
