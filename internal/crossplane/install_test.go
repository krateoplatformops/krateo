//go:build integration
// +build integration

package crossplane

import (
	"context"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestInstall(t *testing.T) {
	os.Setenv("KUBECONFIG", "/Users/lucasepe/Projects/Kiratech/krateo-integr-v11.yaml")

	kubeconfig, err := ioutil.ReadFile(os.Getenv(clientcmd.RecommendedConfigPathEnvVar))
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	err = Install(context.TODO(), InstallOpts{
		RESTConfig: restConfig,
		Namespace:  "krateo-system",
	})
	assert.Nil(t, err, "expecting nil error installing crossplane")
}

func TestGetPOD(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	pod, err := InstalledPOD(context.TODO(), restConfig)
	assert.Nil(t, err, "expecting nil error getting crossplane pod")

	if pod != nil && len(pod.Spec.Containers) > 0 {
		img := pod.Spec.Containers[0].Image
		t.Logf("%s\n", img)
		idx := strings.LastIndex(img, ":")
		if idx != -1 {
			t.Logf("%s\n", img[idx+1:])
		}
	}
}
