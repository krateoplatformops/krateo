package core

import (
	"context"
	"encoding/json"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

type ApplyOpts struct {
	RESTConfig *rest.Config
	Object     *unstructured.Unstructured
	GVK        *schema.GroupVersionKind
}

func Apply(ctx context.Context, opts ApplyOpts) error {
	dr, err := DynamicForGVR(opts.RESTConfig, *opts.GVK, opts.Object.GetNamespace())
	if err != nil {
		if IsNoKindMatchError(err) {
			return nil
		}
		return err
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
