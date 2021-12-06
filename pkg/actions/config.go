package actions

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/krateoplatformops/krateo/pkg/crds"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/gitutils"
	"github.com/krateoplatformops/krateo/pkg/osutils"
	"sigs.k8s.io/yaml"
)

func NewConfig(bus eventbus.Bus, module string) *configAction {
	return &configAction{
		bus:    bus,
		module: module,
	}
}

const (
	moduleDefinitionPattern = "cluster/definition.yaml"
)

var _ Action = (*configAction)(nil)

type configAction struct {
	module       string
	gitRepoToken string
	gitRepoURL   string
	verbose      bool
	bus          eventbus.Bus
}

func (o *configAction) SetVerboseEnabled(verbose bool) {
	o.verbose = verbose
}

func (o *configAction) SetModuleRepoURL(s string) {
	o.gitRepoURL = s
}

func (o *configAction) SetModuleRepoToken(s string) {
	o.gitRepoToken = s
}

func (o *configAction) Run() (err error) {
	dir, err := osutils.GetAppDir("krateo")
	if err != nil {
		return err
	}

	// Fetch module definition from Github
	o.bus.Publish(events.NewStartWaitEvent("Fetching module definition"))
	time.Sleep(time.Second * 2)

	moduleRepoURL := joinURL(o.gitRepoURL, fmt.Sprintf(moduleNamePattern, o.module))
	fs, err := gitutils.Clone(moduleRepoURL, o.gitRepoToken)
	if err != nil {
		return err
	}

	def, err := gitutils.ReadFile(fs, moduleDefinitionPattern)
	if err != nil {
		return err
	}

	xrd := crds.CompositeResourceDefinition{}
	if err := yaml.Unmarshal(def, &xrd); err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent(xrd.Spec.Versions[0].Schema.OpenAPIV3Schema.Description))

	props := crds.Parse(xrd)

	values := []string{}

	for _, el := range props {
		switch el.Type {
		case crds.TypeBoolean:
			res := yesNoPrompt(el.Description, true)
			values = append(values, fmt.Sprintf("%s=%t", el.Name, res))
		default:
			res := stringPrompt(el.Description)
			values = append(values, fmt.Sprintf("%s=%s", el.Name, res))
		}
	}

	// Save configuration
	o.bus.Publish(events.NewStartWaitEvent("Saving module configuration"))
	time.Sleep(time.Second * 2)

	fn := filepath.Join(dir, fmt.Sprintf("module-%s.cfg", o.module))
	if err := printLines(fn, values); err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent("Module configuration saved"))

	return nil
}

// StringPrompt asks for a string value using the label
func stringPrompt(label string) string {
	var s string
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprintf(os.Stderr, ">> %s: ", label)
		s, _ = r.ReadString('\n')
		if s != "" {
			break
		}
	}
	return strings.TrimSpace(s)
}

// YesNoPrompt asks yes/no questions using the label.
func yesNoPrompt(label string, def bool) bool {
	choices := "Y/n"
	if !def {
		choices = "y/N"
	}

	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, ">> %s (%s): ", label, choices)
		s, _ = r.ReadString('\n')
		s = strings.TrimSpace(s)
		if s == "" {
			return def
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
