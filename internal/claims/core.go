package claims

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type ModuleOpts struct {
	RESTConfig *rest.Config
	Data       map[string]interface{}
}

func ApplyCoreModule(ctx context.Context, opts ModuleOpts) error {
	gvk := getGroupVersionKind()

	obj := &unstructured.Unstructured{}
	obj.SetKind(gvk.Kind)
	obj.SetAPIVersion(gvk.GroupVersion().String())
	obj.SetName("core")
	obj.SetLabels(map[string]string{
		core.InstalledByLabel: core.InstalledByValue,
	})
	err := unstructured.SetNestedField(obj.Object, opts.Data, "spec")
	if err != nil {
		return err
	}

	return core.Apply(ctx, core.ApplyOpts{
		RESTConfig: opts.RESTConfig,
		GVK:        gvk,
		Object:     obj,
	})
}

func CoreDefaultClaims() map[string]interface{} {
	return map[string]interface{}{
		//"platform":   "kubernetes",
		//"namespace":  "krateo-namespace",
		//"domain":     "krateo.site",
		"protocol":   "https",
		"domainPort": int64(443),
		"app": map[string]interface{}{
			"hostname": "app",
		},
		"api": map[string]interface{}{
			"version":  "1.0.1",
			"hostname": "api",
		},
		"argo-cd": map[string]interface{}{
			"hostname": "argocd",
		},
		"socket-service": map[string]interface{}{
			"hostname": "socket",
		},
		"deployment-service": map[string]interface{}{
			"version":  "1.0.18",
			"hostname": "deployment",
		},
		"kongapigw": map[string]interface{}{
			"postgresql": map[string]interface{}{
				"enabled": true,
			},
		},
	}
}

func getGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "modules.krateo.io",
		Version: "v1",
		Kind:    "Core",
	}
}

func getGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "modules.krateo.io",
		Version:  "v1",
		Resource: "core",
	}
}
