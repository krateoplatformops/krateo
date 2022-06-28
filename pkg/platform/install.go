package platform

import (
	"fmt"
	"strings"

	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/platform/crds"
	"github.com/krateoplatformops/krateo/pkg/platform/utils"
	"github.com/krateoplatformops/krateo/pkg/text"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type InstallOptions struct {
	Config        *rest.Config
	Bus           eventbus.Bus
	Verbose       bool
	Platform      string
	ModuleName    string
	ModulePackage []byte
	ModuleClaims  []byte
}

func Install(opts *InstallOptions) error {
	dc, err := dynamic.NewForConfig(opts.Config)
	if err != nil {
		return err
	}

	err = installModulePackage(dc, opts)
	if err != nil {
		return err
	}

	err = installModuleClaims(dc, opts)
	if err != nil {
		return err
	}

	return nil
}

func installModulePackage(dc dynamic.Interface, o *InstallOptions) error {
	obj, err := utils.DecodeModuleConfiguration(dc, o.ModulePackage)
	if err != nil {
		return err
	}

	o.Bus.Publish(events.NewStartWaitEvent("installing module package '%s'...", o.ModuleName))

	err = utils.CreateOrUpdateResourceFromUnstructured(o.Config, dc, obj)
	if err != nil {
		return err
	}

	err = waitForModuleCRDs(o.Config, o.ModuleName)
	if err != nil {
		return err
	}
	o.Bus.Publish(events.NewDoneEvent("module package '%s' successfully installed", o.ModuleName))

	return nil
}

func installModuleClaims(dc dynamic.Interface, o *InstallOptions) error {
	obj, err := utils.DecodeModuleClaims(dc, o.ModuleClaims)
	if err != nil {
		return err
	}

	obj, err = updateCompositionSelector(obj, map[string]interface{}{
		"platform": o.Platform,
	})
	if err != nil {
		return err
	}

	if o.Verbose {
		o.Bus.Publish(events.NewDebugEvent("patched manifest using platform '%s'\n", o.Platform))
	}

	o.Bus.Publish(events.NewStartWaitEvent("installing claims form module '%s' ...", o.ModuleName))
	err = utils.CreateOrUpdateResourceFromUnstructured(o.Config, dc, obj)
	if err != nil {
		return err
	}

	o.Bus.Publish(events.NewDoneEvent("claims for module '%s' successfully installed", o.ModuleName))

	return nil
}

func waitForModuleCRDs(rc *rest.Config, module string) error {
	return crds.WaitFor(rc, []*apiextensionsv1.CustomResourceDefinition{
		{
			Spec: apiextensionsv1.CustomResourceDefinitionSpec{
				Group: "modules.krateo.io",
				Versions: []apiextensionsv1.CustomResourceDefinitionVersion{
					{
						Name:    "v1alpha1",
						Served:  false,
						Storage: false,
					},
				},
				Names: apiextensionsv1.CustomResourceDefinitionNames{
					Kind:   text.Title(module),
					Plural: strings.ToLower(module),
				},
			},
		},
	},
	)
}

func updateCompositionSelector(el *unstructured.Unstructured, kvs map[string]interface{}) (*unstructured.Unstructured, error) {
	spec, ok := el.Object["spec"].(map[string]interface{})
	if !ok {
		return el, fmt.Errorf("missed 'spec' in manifest")
	}

	sel, ok := spec["compositionSelector"].(map[string]interface{})
	if !ok {
		return el, fmt.Errorf("missed 'spec.compositionSelector' in manifest")
	}

	mal, ok := sel["matchLabels"].(map[string]interface{})
	if !ok {
		return el, fmt.Errorf("missed 'spec.compositionSelector.matchLabels' in manifest")
	}

	for k, v := range kvs {
		mal[k] = v
	}

	return el, nil
}
