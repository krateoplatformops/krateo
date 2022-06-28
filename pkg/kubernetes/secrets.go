package kubernetes

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

type SecretsClient interface {
	Get(name, namespace string, opts metav1.GetOptions) (*corev1.Secret, error)
	Create(namespace string, secret *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error)
	Delete(name, namespace string, opts metav1.DeleteOptions) error
}

func Secrets(c *rest.Config) (SecretsClient, error) {
	config := *c
	config.APIPath = "/api"
	config.GroupVersion = &corev1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs

	rc, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &secretsClientImpl{
		client: rc,
	}, nil
}

type secretsClientImpl struct {
	client *rest.RESTClient
}

func (impl *secretsClientImpl) Get(name, namespace string, opts metav1.GetOptions) (*corev1.Secret, error) {
	res := &corev1.Secret{}

	err := impl.client.Get().
		Namespace(namespace).
		Resource("secrets").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(res)

	return res, err
}

func (impl *secretsClientImpl) Delete(name, namespace string, opts metav1.DeleteOptions) error {
	err := impl.client.Delete().
		Namespace(namespace).
		Resource("secrets").
		Name(name).
		Body(&opts).
		Do(context.TODO()).
		Error()
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
	}

	return err
}

func (impl *secretsClientImpl) Create(namespace string, secret *corev1.Secret, opts metav1.CreateOptions) (*corev1.Secret, error) {
	res := &corev1.Secret{}

	err := impl.client.Post().
		Namespace(namespace).
		Resource("secrets").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(secret).
		Do(context.TODO()).
		Into(res)
	if err != nil {
		if errors.IsAlreadyExists(err) {
			return impl.Get(secret.Name, namespace, metav1.GetOptions{})
		}
	}

	return res, err
}
