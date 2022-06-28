package kubernetes

import (
	"context"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

type CrdsClient interface {
	Get(name string, opts metav1.GetOptions) (*apiextensionsv1.CustomResourceDefinition, error)
	List(opts metav1.ListOptions) (*apiextensionsv1.CustomResourceDefinitionList, error)
	Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*apiextensionsv1.CustomResourceDefinition, error)
	Delete(name string, opts metav1.DeleteOptions) error
	DeleteCollection(opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
}

func Crds(c *rest.Config) (CrdsClient, error) {
	config := *c
	config.APIPath = "/apis"
	config.GroupVersion = &apiextensionsv1.SchemeGroupVersion
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	rc, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &crdsClientImpl{
		client: rc,
	}, nil
}

type crdsClientImpl struct {
	client *rest.RESTClient
}

// Get takes name of the customResourceDefinition, and returns the corresponding customResourceDefinition object, and an error if there is any.
func (impl *crdsClientImpl) Get(name string, opts metav1.GetOptions) (*apiextensionsv1.CustomResourceDefinition, error) {
	res := &apiextensionsv1.CustomResourceDefinition{}

	err := impl.client.Get().
		Resource("customresourcedefinitions").
		Name(name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(context.TODO()).
		Into(res)

	return res, err
}

func (impl *crdsClientImpl) List(opts metav1.ListOptions) (*apiextensionsv1.CustomResourceDefinitionList, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}

	res := &apiextensionsv1.CustomResourceDefinitionList{}

	err := impl.client.Get().
		Resource("customresourcedefinitions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(context.TODO()).
		Into(res)

	return res, err
}

// Patch applies the patch and returns the patched customResourceDefinition.
func (impl *crdsClientImpl) Patch(name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (*apiextensionsv1.CustomResourceDefinition, error) {
	res := &apiextensionsv1.CustomResourceDefinition{}

	err := impl.client.Patch(pt).
		Resource("customresourcedefinitions").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(context.TODO()).
		Into(res)

	return res, err
}

// Delete takes name of the customResourceDefinition and deletes it. Returns an error if one occurs.
func (impl *crdsClientImpl) Delete(name string, opts metav1.DeleteOptions) error {
	return impl.client.Delete().
		Resource("customresourcedefinitions").
		Name(name).
		Body(&opts).
		Do(context.TODO()).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (impl *crdsClientImpl) DeleteCollection(opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return impl.client.Delete().
		Resource("customresourcedefinitions").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(context.TODO()).
		Error()
}
