package uninstall

import (
	"context"
	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func listXRs(dc dynamic.Interface) (map[string]schema.GroupVersionResource, error) {
	gvr := schema.GroupVersionResource{
		Group:    "apiextensions.crossplane.io",
		Version:  "v1",
		Resource: "compositeresourcedefinitions",
	}

	list, err := dc.Resource(gvr).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	res := make(map[string]schema.GroupVersionResource)

	for _, el := range list.Items {
		if !strings.HasSuffix(el.GetName(), "modules.krateo.io") {
			continue
		}

		spec, ok := el.Object["spec"].(map[string]interface{})
		if !ok {
			continue
		}

		versions, ok := spec["versions"].([]interface{})
		if !ok {
			continue
		}

		names, ok := spec["names"].(map[string]interface{})
		if !ok {
			continue
		}

		g := spec["group"].(string)
		v := versions[0].(map[string]interface{})["name"].(string)
		r := names["plural"].(string)

		res[el.GetName()] = schema.GroupVersionResource{
			Group:    g,
			Version:  v,
			Resource: r,
		}
	}

	return res, nil
}

func deleteXRs(dc dynamic.Interface, gvrs map[string]schema.GroupVersionResource) error {
	for name, el := range gvrs {
		err := dc.Resource(el).Delete(context.Background(), name, metav1.DeleteOptions{})
		if err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}

	return nil
}
