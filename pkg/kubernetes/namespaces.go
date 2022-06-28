package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/kubectl/pkg/scheme"
)

type NamespacesClient interface {
	Get(name string, opts metav1.GetOptions) (*corev1.Namespace, error)
	Create(name string) (*corev1.Namespace, error)
	Update(item *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error)
	Delete(name string, opts metav1.DeleteOptions) error
	Finalize(item *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error)
}

func Namespaces(c *rest.Config) (NamespacesClient, error) {
	config := *c
	config.APIPath = "/api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	rc, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &namespacesClientImpl{
		client: rc,
	}, nil
}

type namespacesClientImpl struct {
	client *rest.RESTClient
}

func (impl *namespacesClientImpl) Get(name string, opts metav1.GetOptions) (*corev1.Namespace, error) {
	res := &corev1.Namespace{}

	err := impl.client.Get().
		Resource("namespaces").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(res)

	return res, err
}

func (impl *namespacesClientImpl) Update(item *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error) {
	res := &corev1.Namespace{}

	err := impl.client.Put().
		Resource("namespaces").
		Name(item.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(item).
		Do(context.TODO()).
		Into(res)

	return res, err
}

// Create takes the representation of a namespace and creates it.
// Returns the server's representation of the namespace, and an error, if there is any.
//
// Equivalent to:
// $ kubectl create namespace <NAME>
func (impl *namespacesClientImpl) Create(name string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				LabelManagedBy: DefaultFieldManager,
			},
		},
	}

	res := &corev1.Namespace{}

	opts := metav1.CreateOptions{
		FieldManager: DefaultFieldManager,
	}

	err := impl.client.Post().
		Resource("namespaces").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ns).
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

// Delete takes name of the namespace and deletes it.
// Returns an error if one occurs.
//
// Equivalent to:
// $ kubectl delete namespaces <NAME>
func (impl *namespacesClientImpl) Delete(name string, opts metav1.DeleteOptions) error {
	return impl.client.Delete().
		Resource("namespaces").
		Name(name).
		Body(&opts).
		Do(context.TODO()).
		Error()
}

func (impl *namespacesClientImpl) Finalize(item *corev1.Namespace, opts metav1.UpdateOptions) (*corev1.Namespace, error) {
	res := &corev1.Namespace{}

	err := impl.client.Put().
		Resource("namespaces").
		Name(item.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		SubResource("finalize").
		Body(item).
		Do(context.TODO()).
		Into(res)

	return res, err
}
