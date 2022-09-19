package crossplane

import (
	"bytes"
	"context"
	"fmt"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/krateoplatformops/krateo/internal/eventbus"
	"github.com/krateoplatformops/krateo/internal/events"
	"github.com/krateoplatformops/krateo/internal/helm"
	"github.com/krateoplatformops/krateo/internal/httputils"
	"github.com/krateoplatformops/krateo/internal/pods"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

const (
	chartReleaseName = "crossplane"
)

type InstallOpts struct {
	RESTConfig *rest.Config
	ChartURL   string
	Verbose    bool
	HttpProxy  string
	HttpsProxy string
	NoProxy    string
	Namespace  string
	EventBus   eventbus.Bus
}

func Install(ctx context.Context, opts InstallOpts) error {
	chartArchive := &bytes.Buffer{}
	err := httputils.Fetch(opts.ChartURL, chartArchive)
	if err != nil {
		return err
	}

	err = createNamespaceEventually(ctx, opts.RESTConfig, opts.Namespace)
	if err != nil {
		return fmt.Errorf("creating namespace '%s': %w", opts.Namespace, err)
	}

	helmOpts := helm.InstallOptions{
		RESTConfig:  opts.RESTConfig,
		Namespace:   opts.Namespace,
		ReleaseName: chartReleaseName,
		ChartSource: bytes.NewReader(chartArchive.Bytes()),
		ChartValues: map[string]interface{}{
			"args": []interface{}{
				"--debug",
			},
			"securityContextCrossplane": map[string]interface{}{
				"runAsUser":  nil,
				"runAsGroup": nil,
			},
			"securityContextRBACManager": map[string]interface{}{
				"runAsUser":  nil,
				"runAsGroup": nil,
			},
			"extraEnvVarsCrossplane": map[string]interface{}{},
		},
		LogFn: func(format string, v ...interface{}) {
			if opts.Verbose && opts.EventBus != nil {
				opts.EventBus.Publish(events.NewDebugEvent(format, v))
			}
		},
	}

	if opts.HttpProxy != "" {
		envVars := helmOpts.ChartValues["extraEnvVarsCrossplane"].(map[string]interface{})
		envVars["HTTP_PROXY"] = opts.HttpProxy
	}

	if opts.HttpsProxy != "" {
		envVars := helmOpts.ChartValues["extraEnvVarsCrossplane"].(map[string]interface{})
		envVars["HTTPS_PROXY"] = opts.HttpsProxy
	}

	if opts.NoProxy != "" {
		envVars := helmOpts.ChartValues["extraEnvVarsCrossplane"].(map[string]interface{})
		envVars["NO_PROXY"] = opts.NoProxy
	}

	err = helm.Install(helmOpts)
	if err != nil {
		return err
	}

	return waitUntilCrossplaneIdReady(opts.RESTConfig, opts.Namespace)
}

func createNamespaceEventually(ctx context.Context, restConfig *rest.Config, namespace string) error {
	obj := &unstructured.Unstructured{}
	obj.SetKind("Namespace")
	obj.SetName(namespace)

	return core.Create(ctx, core.CreateOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Namespace",
		},
		Object: obj,
	})
}

// waitUntilCrossplaneIdReady waits until Crossplane PODs are ready
func waitUntilCrossplaneIdReady(restConfig *rest.Config, namespace string) error {
	sel, err := labels.Parse("app=crossplane")
	if err != nil {
		return err
	}

	stopFn := func(cond corev1.PodCondition) bool {
		return cond.Type == corev1.PodReady &&
			cond.Status == corev1.ConditionTrue
	}

	dc, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	return pods.Watch(dc, pods.WatchOpts{
		Namespace: namespace,
		Selector:  sel,
		StopFunc:  stopFn,
	})
}
