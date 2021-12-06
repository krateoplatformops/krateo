package kubeclient

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewKubeClient returns a kubernetes client given a secret with connection
// information.
func NewKubeClient(config *rest.Config) (client.Client, error) {
	kc, err := client.New(config, client.Options{})
	if err != nil {
		return nil, errors.Wrap(err, "cannot create Kubernetes client")
	}

	return kc, nil
}

// NewKubeClientSet returns a kubernetes clientset given a secret with connection
// information.
func NewKubeClientSet(config *rest.Config) (*kubernetes.Clientset, error) {
	kc, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create Kubernetes clientset")
	}

	return kc, nil
}
