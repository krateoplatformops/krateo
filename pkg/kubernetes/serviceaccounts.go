package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

type ServiceAccountsClient interface {
	List(namespace string) ([]string, error)
	Get(name, namespace string, opts metav1.GetOptions) (*corev1.ServiceAccount, error)
}

func ServiceAccounts(c *rest.Config) (ServiceAccountsClient, error) {
	config := *c
	config.APIPath = "/api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	rc, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &serviceAccountsClientImpl{
		client: rc,
	}, nil
}

type serviceAccountsClientImpl struct {
	client *rest.RESTClient
}

// kubectl get serviceaccounts -n NAMESPACE -o name
func (impl *serviceAccountsClientImpl) List(namespace string) ([]string, error) {
	list := &corev1.ServiceAccountList{}

	err := impl.client.Get().
		Namespace(namespace).
		Resource("serviceaccounts").
		Do(context.TODO()).
		Into(list)
	if err != nil {
		if errors.IsNotFound(err) {
			return []string{}, nil
		}
		return nil, err
	}

	var res []string
	for _, el := range list.Items {
		res = append(res, el.Name)
	}

	return res, nil
}

func (impl *serviceAccountsClientImpl) Get(name, namespace string, opts metav1.GetOptions) (*corev1.ServiceAccount, error) {
	res := &corev1.ServiceAccount{}

	err := impl.client.Get().
		Namespace(namespace).
		Resource("serviceaccounts").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(res)

	return res, err
}
