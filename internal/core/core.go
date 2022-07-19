package core

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	InstalledByLabel = "app.kubernetes.io/installed-by"
	InstalledByValue = "krateo"
	PackageNameLabel = "krateo.io/package-name"
	DefaultNamespace = "crossplane-system"
)

func RESTConfigFromBytes(data []byte) (*rest.Config, error) {
	config, err := clientcmd.Load(data)
	if err != nil {
		return nil, err
	}
	//currentContext := config.CurrentContext
	//t.Logf("current context: %s", currentContext)

	restConfig, err := clientcmd.NewDefaultClientConfig(*config, nil).ClientConfig()
	if err != nil {
		return nil, err
	}
	// Set QPS and Burst to a threshold that ensures the client doesn't generate throttling log messages
	restConfig.QPS = 20
	restConfig.Burst = 100

	return restConfig, nil
}
