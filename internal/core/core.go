package core

import (
	"errors"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	InstalledByLabel = "app.kubernetes.io/installed-by"
	InstalledByValue = "krateo"
	PackageNameLabel = "krateo.io/package-name"
	DefaultNamespace = "crossplane-system"
)

var (
	NoKindMatchError = errors.New("RESTMapper can't find any match for kind")
)

func RESTConfigFromBytes(data []byte, withContext string) (*rest.Config, error) {
	config, err := clientcmd.Load(data)
	if err != nil {
		return nil, err
	}

	currentContext := config.CurrentContext
	if len(withContext) > 0 {
		currentContext = withContext
	}

	restConfig, err := clientcmd.NewNonInteractiveClientConfig(*config,
		currentContext, &clientcmd.ConfigOverrides{}, nil).ClientConfig()
	if err != nil {
		return nil, err
	}
	// Set QPS and Burst to a threshold that ensures the client doesn't generate throttling log messages
	restConfig.QPS = 20
	restConfig.Burst = 100

	return restConfig, nil
}

// FindGVR find the corresponding GVR (available in *meta.RESTMapping) for gvk
func FindGVR(cfg *rest.Config, gvk schema.GroupVersionKind) (*meta.RESTMapping, error) {
	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}

func DynamicForGVR(restConfig *rest.Config, gvk schema.GroupVersionKind, namespace string) (dynamic.ResourceInterface, error) {
	mapping, err := FindGVR(restConfig, gvk)
	if err != nil {
		return nil, err
	}

	dc, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	// obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dc.Resource(mapping.Resource).Namespace(namespace)
	} else {
		// for cluster-wide resources
		dr = dc.Resource(mapping.Resource)
	}

	return dr, nil
}

func IsNoKindMatchError(err error) bool {
	var noKindMatchError *meta.NoKindMatchError
	return errors.As(err, &noKindMatchError)
}
