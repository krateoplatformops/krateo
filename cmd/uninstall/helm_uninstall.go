package uninstall

import (
	"fmt"

	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/helm"
	"github.com/krateoplatformops/krateo/pkg/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	chartReleaseName = "crossplane"
)

func uninstallCrossplaneChart(rc *rest.Config, bus eventbus.Bus, verbose bool) error {
	opts := &helm.UninstallOptions{
		Namespace:   kubernetes.CrossplaneSystemNamespace,
		ReleaseName: chartReleaseName,
		LogFn: func(format string, v ...interface{}) {
			if verbose && bus != nil {
				msg := fmt.Sprintf(format, v...)
				bus.Publish(events.NewDebugEvent(msg))
			}
		},
	}

	return helm.Uninstall(rc, opts)
}
