package core

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type PatchOpts struct {
	RESTConfig *rest.Config
	PatchData  []byte
	GVR        schema.GroupVersionResource
	Name       string
}

func Patch(ctx context.Context, opts PatchOpts) error {
	dc, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return err
	}

	_, err = dc.Resource(opts.GVR).
		Patch(ctx, opts.Name, types.MergePatchType, opts.PatchData, metav1.PatchOptions{
			FieldManager: InstalledByValue,
		})
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	return nil
}
