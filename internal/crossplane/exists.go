package crossplane

import (
	"context"
	"strings"

	"github.com/krateoplatformops/krateo/internal/core"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

func PODImageVersion(pod *corev1.Pod) (string, error) {
	if pod != nil && (len(pod.Spec.Containers) > 0) {
		img := pod.Spec.Containers[0].Image
		idx := strings.LastIndex(img, ":")
		if idx != -1 {
			return img[idx+1:], nil
		}
	}
	return "", nil
}

func InstalledPOD(ctx context.Context, restConfig *rest.Config) (*corev1.Pod, error) {
	items, err := core.List(ctx, core.ListOpts{
		RESTConfig:    restConfig,
		GVK:           schema.GroupVersionKind{Version: "v1", Kind: "Pod"},
		Namespace:     "",
		LabelSelector: "app=crossplane",
	})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(items) > 0 {
		pod := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(items[0].UnstructuredContent(), pod)
		if err != nil {
			return nil, err
		}

		return pod, nil
	}

	return nil, nil
}

type ExistOpts struct {
	RESTConfig *rest.Config
	Namespace  string
}

func Exists(ctx context.Context, opts ExistOpts) (bool, error) {
	sel, err := labels.Parse("app=crossplane")
	if err != nil {
		return false, err
	}

	list, err := core.List(ctx, core.ListOpts{
		RESTConfig: opts.RESTConfig,
		GVK: schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Pod",
		},
		Namespace:     opts.Namespace,
		LabelSelector: sel.String(),
	})
	if err != nil {
		return false, err
	}

	return len(list) > 0, nil
}
