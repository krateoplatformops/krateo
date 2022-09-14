package claims

import (
	"context"

	"github.com/davecgh/go-spew/spew"
	"github.com/krateoplatformops/krateo/internal/core"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

type CreateCoreOpts struct {
	RESTConfig *rest.Config
	Data       map[string]interface{}
}

func Create(ctx context.Context, opts CreateCoreOpts) error {
	gvk := schema.GroupVersionKind{
		Group:   "modules.krateo.io",
		Version: "v1",
		Kind:    "Core",
	}

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

	spew.Dump(obj)

	return core.Create(ctx, core.CreateOpts{
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
			"hostname": "api",
		},
		"argo-cd": map[string]interface{}{
			"hostname": "argocd",
		},
		"socket-service": map[string]interface{}{
			"hostname": "socket",
		},
		"kong": map[string]interface{}{
			"postgresql": map[string]interface{}{
				"enabled": true,
			},
		},
	}
}
