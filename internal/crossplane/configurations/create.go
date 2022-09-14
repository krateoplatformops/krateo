package configurations

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
		Version:  "v1",
		Resource: "configurations",
	}

	obj := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"spec": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						core.InstalledByLabel: core.InstalledByValue,
						core.PackageNameLabel: opts.Info.Name,
					},
				},
				"package":                  fmt.Sprintf("%s:%s", opts.Info.Image, opts.Info.Version),
				"packagePullPolicy":        "IfNotPresent",
				"revisionActivationPolicy": "Automatic",
				"revisionHistoryLimit":     1,
			},
		},
	}

	obj.SetKind("Configuration")
	obj.SetAPIVersion("pkg.crossplane.io/v1")
	obj.SetName(fmt.Sprintf("%s-configuration", opts.Info.Name))
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
