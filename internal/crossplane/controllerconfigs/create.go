package controllerconfigs

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo/internal/catalog"
	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// DeepCopy is a deepcopy function, copying the receiver, creating a new EnvVar.
func (in *EnvVar) DeepCopy() *EnvVar {
	if in == nil {
		return nil
	}
	out := new(EnvVar)
	*out = *in
	return out
}

type CreateOpts struct {
	RESTConfig *rest.Config
	Info       *catalog.PackageInfo
	HttpProxy  string
	HttpsProxy string
	NoProxy    string
}

func Create(ctx context.Context, opts CreateOpts) (*unstructured.Unstructured, error) {
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1alpha1",
		Resource: "controllerconfigs",
	}

	envVars := []interface{}{}
	if opts.HttpProxy != "" {
		envVars = append(envVars, map[string]string{
			"name":  "HTTP_PROXY",
			"value": opts.HttpProxy,
		})
	}

	if opts.HttpsProxy != "" {
		envVars = append(envVars, map[string]string{
			"name":  "HTTPS_PROXY",
			"value": opts.HttpsProxy,
		})
	}

	if opts.NoProxy != "" {
		envVars = append(envVars, map[string]string{
			"name":  "NO_PROXY",
			"value": opts.NoProxy,
		})
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": map[string]interface{}{
				"securityContext":    map[string]interface{}{},
				"podSecurityContext": map[string]interface{}{},
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						core.InstalledByLabel: core.InstalledByValue,
						core.PackageNameLabel: opts.Info.Name,
					},
				},
				"env": envVars,
			},
		},
	}

	obj.SetKind("ControllerConfig")
	obj.SetAPIVersion("pkg.crossplane.io/v1alpha1")
	obj.SetName(fmt.Sprintf("%s-controllerconfig", opts.Info.Name))
	obj.SetLabels(map[string]string{
		core.InstalledByLabel: core.InstalledByValue,
	})

	// prepare the dynamic client
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	_, err = dc.Resource(gvr).Create(ctx, obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return nil, err
		}
	}

	return obj, nil
}
