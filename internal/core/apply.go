package core

import (
	"context"
	"encoding/json"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

type ApplyOpts struct {
	RESTConfig *rest.Config
	Object     *unstructured.Unstructured
	GVK        *schema.GroupVersionKind
}

func Apply(ctx context.Context, opts ApplyOpts) error {
	// find GVR
	mapping, err := FindGVR(opts.RESTConfig, opts.GVK)
	if err != nil {
		return err
	}

	// prepare the dynamic client
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	// obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dc.Resource(mapping.Resource).Namespace(opts.Object.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dc.Resource(mapping.Resource)
	}

	// 6. Marshal object into JSON
	data, err := json.Marshal(opts.Object)
	if err != nil {
		return err
	}

	// create or Update the object with SSA (types.ApplyPatchType indicates SSA).
	_, err = dr.Patch(ctx, opts.Object.GetName(), types.ApplyPatchType, data, metav1.PatchOptions{
		FieldManager: InstalledByValue,
	})

	return err
}

// FindGVR find the corresponding GVR (available in *meta.RESTMapping) for gvk
func FindGVR(cfg *rest.Config, gvk *schema.GroupVersionKind) (*meta.RESTMapping, error) {
	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}
