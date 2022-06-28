package kubernetes

import (
	"context"
	"time"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type ClusterRoleBindingsClient interface {
	Get(name string) (*rbacv1.ClusterRoleBinding, error)
	Create(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error)
	Delete(name string) error
	DeleteCollection(opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
}

func ClusterRoleBindings(c *rest.Config) (ClusterRoleBindingsClient, error) {
	config := *c
	config.APIPath = "/apis"
	config.GroupVersion = &rbacv1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	rc, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &clusterRoleBindingsClientImpl{
		client: rc,
	}, nil
}

type clusterRoleBindingsClientImpl struct {
	client *rest.RESTClient
}

// Get takes name of the clusterRoleBinding, and returns the corresponding object.
// If the clusterRoleBinding with this name does not exists return nil.
func (impl *clusterRoleBindingsClientImpl) Get(name string) (*rbacv1.ClusterRoleBinding, error) {
	opts := metav1.GetOptions{}

	res := &rbacv1.ClusterRoleBinding{}

	err := impl.client.Get().
		Resource("clusterrolebindings").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(res)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// Create takes the representation of a clusterRoleBinding and creates it.
// Returns the server's representation of the clusterRoleBinding, and an error, if there is any.
//
// Equivalent to:
// $
func (impl *clusterRoleBindingsClientImpl) Create(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
	res := &rbacv1.ClusterRoleBinding{}

	opts := metav1.CreateOptions{
		FieldManager: DefaultFieldManager,
	}

	err := impl.client.Post().
		Resource("clusterrolebindings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(crb).
		Do(context.TODO()).
		Into(res)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// Delete takes name of the clusterRoleBinding and deletes it.
// Returns an error if one occurs.
func (impl *clusterRoleBindingsClientImpl) Delete(name string) error {
	opts := metav1.DeleteOptions{}

	return impl.client.Delete().
		Resource("clusterrolebindings").
		Name(name).
		Body(&opts).
		Do(context.TODO()).
		Error()
}

func (impl *clusterRoleBindingsClientImpl) DeleteCollection(opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return impl.client.Delete().
		Resource("clusterrolebindings").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(context.TODO()).
		Error()
}
