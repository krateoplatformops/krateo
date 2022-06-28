package providers

import (
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo/pkg/platform/pods"
	"github.com/krateoplatformops/krateo/pkg/platform/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/dynamic"
)

const (
	controllerConfigNamePattern = "provider-%s-controllerconfig"
	providerNamePattern         = "provider-%s"
)

func InstallEventually(dc dynamic.Interface, info *providerInfo) error {
	ok, err := isProviderAlreadyInstalled(dc, info)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	err = createProvider(dc, info)
	if err != nil {
		return err
	}

	return waitForProvider(dc, info)
}

func CreateControllerConfig(dc dynamic.Interface, info *providerInfo) error {
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1alpha1",
		Resource: "controllerconfigs",
	}

	obj := unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "ControllerConfig",
			"apiVersion": "pkg.crossplane.io/v1alpha1",
			"metadata": map[string]interface{}{
				"name": fmt.Sprintf(controllerConfigNamePattern, info.Name()),
			},
			"labels": map[string]string{
				utils.LabelManagedBy: utils.DefaultFieldManager,
			},
			"spec": map[string]interface{}{
				"securityContext":    map[string]interface{}{},
				"podSecurityContext": map[string]interface{}{},
				/*
					"resources": map[string]interface{}{
						"limits": map[string]interface{}{
							"cpu":    "100m",
							"memory": "128Mi",
						},
						"requests": map[string]interface{}{
							"cpu":    "50m",
							"memory": "64Mi",
						},
					},
				*/
			},
		},
	}

	_, err := dc.Resource(gvr).Create(context.TODO(), &obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func isProviderAlreadyInstalled(dc dynamic.Interface, info *providerInfo) (bool, error) {
	req, err := labels.NewRequirement("pkg.crossplane.io/provider",
		selection.Equals, []string{fmt.Sprintf(providerNamePattern, info.Name())})
	if err != nil {
		return false, err
	}

	sel := labels.NewSelector()
	sel = sel.Add(*req)

	return pods.Exists(dc, sel)
}

func createProvider(dc dynamic.Interface, info *providerInfo) error {
	gvr := schema.GroupVersionResource{
		Group:    "pkg.crossplane.io",
		Version:  "v1",
		Resource: "providers",
	}

	obj := unstructured.Unstructured{
		Object: map[string]interface{}{
			"kind":       "Provider",
			"apiVersion": "pkg.crossplane.io/v1",
			"metadata": map[string]interface{}{
				"name": fmt.Sprintf(providerNamePattern, info.Name()),
				"annotations": map[string]interface{}{
					"metaUrl": info.MetaUrl(),
				},
			},
			"labels": map[string]string{
				utils.LabelManagedBy: utils.DefaultFieldManager,
			},
			"spec": map[string]interface{}{
				"package": info.Image(),
				"controllerConfigRef": map[string]interface{}{
					"name": fmt.Sprintf(controllerConfigNamePattern, info.Name()),
				},
			},
		},
	}

	_, err := dc.Resource(gvr).Create(context.Background(), &obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// waitForProvider waits until the specified crossplane provider is ready
func waitForProvider(dc dynamic.Interface, info *providerInfo) error {
	req, err := labels.NewRequirement("pkg.crossplane.io/provider",
		selection.Equals, []string{fmt.Sprintf(providerNamePattern, info.Name())})
	if err != nil {
		return err
	}

	sel := labels.NewSelector()
	sel = sel.Add(*req)

	stopFn := func(cond corev1.PodCondition) bool {
		return cond.Type == corev1.PodReady &&
			cond.Status == corev1.ConditionTrue
	}

	return pods.Watch(dc, sel, stopFn)
}
