package config

import (
	"fmt"
	"strconv"

	"github.com/krateoplatformops/krateo/cmd/tools"
	"github.com/krateoplatformops/krateo/cmd/tools/flags"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/log"
	"github.com/krateoplatformops/krateo/pkg/storage"
	"github.com/krateoplatformops/krateo/pkg/storage/git"
	"github.com/krateoplatformops/krateo/pkg/strvals"
	"github.com/krateoplatformops/krateo/pkg/text"
	"github.com/krateoplatformops/krateo/pkg/xrds"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

type configOptions struct {
	bus eventbus.Bus
	//kubeconfig string
	verbose bool

	gitUrl   string
	gitToken string
	module   string
}

func NewConfigCmd() *cobra.Command {
	o := &configOptions{}

	cmd := &cobra.Command{
		Use:                   "config <MODULE>",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ArbitraryArgs,
		Short:                 "Manage Krateo modules configuration",
		SilenceErrors:         true,
		Example:               "  krateo config core",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("missed module name")
			}

			o.module = args[0]

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

	cmd.Flags().StringVarP(&o.gitUrl, flags.GitURL, "r", "", "git repository url for pushing module configuration")
	cmd.Flags().StringVarP(&o.gitToken, flags.GitToken, "t", "", "token for git repository authentication")
	cmd.Flags().BoolVarP(&o.verbose, flags.Verbose, "v", false, "dump verbose output")
	//nolint:errcheck
	cmd.MarkFlagRequired(flags.GitURL)
	//nolint:errcheck
	cmd.MarkFlagRequired(flags.GitToken)

	return cmd
}

func (o *configOptions) Run() error {
	o.bus.Publish(events.NewStartWaitEvent(fmt.Sprintf("Fetching definition of module '%s'", o.module)))
	entry, err := tools.GitPullModuleDefinition(o.module)
	if err != nil {
		return err
	}

	xrd := &xrds.CompositeResourceDefinition{}
	if err := yaml.Unmarshal(entry.Content, xrd); err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent(fmt.Sprintf("%s fetched", xrd.Spec.Versions[0].Schema.OpenAPIV3Schema.Description)))

	values, err := o.promptForCompositionValues(xrd)
	if err != nil {
		return err
	}

	return o.pushModuleConfig(values)
}

func (o *configOptions) promptForCompositionValues(xrd *xrds.CompositeResourceDefinition) ([]string, error) {
	props, err := xrds.ParseCompositeResourceDefinition(xrd)
	if err != nil {
		return nil, err
	}

	valuesMap := map[string]string{}

	for _, el := range props {
		if len(el.Default) > 0 {
			valuesMap[el.Name] = el.Default
		}

		if el.Required {
			switch el.Type {
			case xrds.TypeBoolean:
				//nolint:errcheck
				def, _ := strconv.ParseBool(el.Default)
				res := yesNoPrompt(el.Description, def)
				valuesMap[el.Name] = fmt.Sprintf("%t", res)
			default:
				res := stringPrompt(el.Description, el.Default, el.Required)
				valuesMap[el.Name] = res
			}
		}
	}

	values := []string{}
	for key, val := range valuesMap {
		values = append(values, fmt.Sprintf("%s=%s", key, val))
	}
	return values, nil
}

func (o *configOptions) pushModuleConfig(values []string) error {
	dat, err := o.buildModuleConfig(values)
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewStartWaitEvent(fmt.Sprintf("Pushing module '%s' configuration to '%s'", o.module, o.gitUrl)))

	entry := storage.Entry{
		Path: fmt.Sprintf("defaults/krateo-module-%s.yaml", o.module),
		Meta: storage.Metadata{
			Name: fmt.Sprintf("krateo-module-%s.yaml", o.module),
		},
		Content: dat,
	}

	err = tools.GitPushEntry(o.gitUrl, entry, git.WithGitToken(o.gitToken))
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewDoneEvent(fmt.Sprintf("Module '%s' configuration pushed to '%s", o.module, o.gitUrl)))

	return nil
}

func (o *configOptions) buildModuleConfig(values []string) ([]byte, error) {
	mod := map[string]interface{}{
		"apiVersion": "modules.krateo.io/v1alpha1",
		"kind":       text.Title(o.module),
		"metadata": map[string]string{
			"name": fmt.Sprintf("krateo-module-%s", o.module),
		},
	}

	for _, line := range values {
		line = fmt.Sprintf("spec.%s", line)
		err := strvals.ParseInto(line, mod)
		if err != nil {
			return nil, err
		}
	}

	return yaml.Marshal(mod)
}
