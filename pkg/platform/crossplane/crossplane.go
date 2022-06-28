package crossplane

import (
	"embed"
	"fmt"

	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/helm"
	"github.com/krateoplatformops/krateo/pkg/kubernetes"

	"github.com/krateoplatformops/krateo/pkg/platform/pods"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

// https://charts.crossplane.io/stable
const (
	ChartVersion     = "1.8.1"
	ChartReleaseName = "crossplane"
)

//go:embed assets/*
var assetsFS embed.FS

type ChartOpts struct {
	Verbose    bool
	HttpProxy  string
	HttpsProxy string
	NoProxy    string
}

func InstallChart(rc *rest.Config, bus eventbus.Bus, co ChartOpts) error {
	chartArchive := fmt.Sprintf("assets/crossplane-%s.tgz", ChartVersion)
	fp, err := assetsFS.Open(chartArchive)
	if err != nil {
		return err
	}
	defer fp.Close()

	opts := &helm.InstallOptions{
		Namespace:   kubernetes.CrossplaneSystemNamespace,
		ReleaseName: ChartReleaseName,
		ChartSource: fp,
		ChartValues: map[string]interface{}{
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
			if co.Verbose && bus != nil {
				bus.Publish(events.NewDebugEvent(format, v))
			}
		},
	}

	if co.HttpProxy != "" {
		envVars := opts.ChartValues["extraEnvVarsCrossplane"].(map[string]interface{})
		envVars["HTTP_PROXY"] = co.HttpProxy
	}

	if co.HttpsProxy != "" {
		envVars := opts.ChartValues["extraEnvVarsCrossplane"].(map[string]interface{})
		envVars["HTTPS_PROXY"] = co.HttpsProxy
	}

	if co.NoProxy != "" {
		envVars := opts.ChartValues["extraEnvVarsCrossplane"].(map[string]interface{})
		envVars["NO_PROXY"] = co.NoProxy
	}

	return helm.Install(rc, opts)
}

func IsInstalled(dc dynamic.Interface) (bool, error) {
	sel, err := labels.Parse("app=crossplane")
	if err != nil {
		return false, err
	}

	return pods.Exists(dc, sel)
}

// WaitUntilReady waits until Crossplane PODs are ready
func WaitUntilReady(dc dynamic.Interface) error {
	sel, err := labels.Parse("app=crossplane")
	if err != nil {
		return err
	}

	stopFn := func(cond corev1.PodCondition) bool {
		return cond.Type == corev1.PodReady &&
			cond.Status == corev1.ConditionTrue
	}

	return pods.Watch(dc, sel, stopFn)
}
