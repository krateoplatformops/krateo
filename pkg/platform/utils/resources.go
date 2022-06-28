package utils

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

/*
func createResourceFromYAML(rc *rest.Config, dc dynamic.Interface, src []byte) error {
	obj := &unstructured.Unstructured{}

	// decode YAML into unstructured.Unstructured
	dec := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, _, err := dec.Decode(src, nil, obj)
	if err != nil {
		return err
	}

	return createOrUpdateResourceFromUnstructured(rc, dc, obj)
}
*/

func CreateOrUpdateResourceFromUnstructured(rc *rest.Config, dc dynamic.Interface, obj *unstructured.Unstructured) error {
	gvk := obj.GroupVersionKind()

	mapping, err := findGVR(rc, &gvk)
	if err != nil {
		return err
	}

	cli := dc.Resource(mapping.Resource)

	res, err := cli.Get(context.Background(), obj.GetName(), metav1.GetOptions{})
	if err == nil {
		obj.SetResourceVersion(res.GetResourceVersion())
		_, err = cli.Update(context.Background(), obj, metav1.UpdateOptions{})
		return err
	} else {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	_, err = cli.Create(context.Background(), obj, metav1.CreateOptions{})
	return err

}

// find the corresponding GVR (available in *meta.RESTMapping) for gvk
func findGVR(cfg *rest.Config, gvk *schema.GroupVersionKind) (*meta.RESTMapping, error) {
	// DiscoveryClient queries API server about the resources
	dc, err := discovery.NewDiscoveryClientForConfig(cfg)
	if err != nil {
		return nil, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	return mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
}
