package actions

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/krateoplatformops/krateo/pkg/clients/kubeclient"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/gitutils"
	"github.com/krateoplatformops/krateo/pkg/osutils"
	"helm.sh/helm/v3/pkg/strvals"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

// TODO type InstallOption func(*installAction)

func NewInstall(bus eventbus.Bus, module, kubecofig string) *installAction {
	return &installAction{
		bus:        bus,
		module:     module,
		kubeconfig: kubecofig,
	}
}

const (
	moduleNamePattern    = "krateo-module-%s"
	moduleConfigPattern  = "krateo-module-%s.cfg"
	modulePackagePattern = "examples/krateo-package-module-%s.yaml"
)

var _ Action = (*installAction)(nil)

type installAction struct {
	module       string
	kubeconfig   string
	gitRepoToken string
	gitRepoURL   string
	verbose      bool
	bus          eventbus.Bus
}

func (o *installAction) SetVerboseEnabled(verbose bool) {
	o.verbose = verbose
}

func (o *installAction) SetModuleRepoURL(s string) {
	o.gitRepoURL = s
}

func (o *installAction) SetModuleRepoToken(s string) {
	o.gitRepoToken = s
}

func (o *installAction) Run() (err error) {
	dir, err := osutils.GetAppDir("krateo")
	if err != nil {
		return err
	}

	moduleCfgFileName := filepath.Join(dir, fmt.Sprintf(moduleConfigPattern, o.module))
	ok, err := osutils.FileExists(moduleCfgFileName)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("config file: %s not found, please use 'krateo config %s' first", moduleCfgFileName, o.module)
	}

	if o.verbose {
		o.bus.Publish(events.NewDebugEvent(fmt.Sprintf("Found module config @ %s", moduleCfgFileName)))
	}

	rc, err := clientcmd.BuildConfigFromFlags("", o.kubeconfig)
	if err != nil {
		return err
	}

	// Fetch module composition package
	moduleRepoURL := joinURL(o.gitRepoURL, fmt.Sprintf(moduleNamePattern, o.module))
	if o.verbose {
		o.bus.Publish(events.NewDebugEvent(fmt.Sprintf("Contacting git repository @ %s", moduleRepoURL)))
	}

	o.bus.Publish(events.NewStartWaitEvent("Fetching package config"))

	fs, err := gitutils.Clone(moduleRepoURL, o.gitRepoToken)
	if err != nil {
		return fmt.Errorf("%w: %s", err, moduleRepoURL)
	}

	pkg, err := gitutils.ReadFile(fs, fmt.Sprintf(modulePackagePattern, o.module))
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewDoneEvent("Package config successfully fetched"))

	// Generate module specs
	o.bus.Publish(events.NewStartWaitEvent("Generating module specs"))
	specs, err := o.makeModuleSpecs(moduleCfgFileName)
	if err != nil {
		return err
	}

	if o.verbose {
		o.bus.Publish(events.NewDebugEvent(string(specs)))
	}

	o.bus.Publish(events.NewDoneEvent("Module specs generated"))

	// Apply module specs
	o.bus.Publish(events.NewStartWaitEvent("Installing module specs"))
	err = kubeclient.Apply(context.Background(), rc, pkg)
	if err != nil {
		return err
	}

	err = kubeclient.Apply(context.Background(), rc, specs)
	if err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent("Module specs successfully installed"))

	return nil
}

func (o *installAction) makeModuleSpecs(fn string) ([]byte, error) {
	values, err := readLines(fn, "spec")
	if err != nil {
		return nil, err
	}

	mod := map[string]interface{}{
		"apiVersion": "modules.krateo.io/v1alpha1",
		"kind":       strings.Title(o.module),
		"metadata": map[string]string{
			"name": fmt.Sprintf(moduleNamePattern, o.module),
		},
	}

	for _, line := range values {
		if err := strvals.ParseInto(line, mod); err != nil {
			return nil, err
		}
	}

	return yaml.Marshal(mod)
}
