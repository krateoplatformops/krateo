package kubeclient

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"github.com/pkg/errors"

	"github.com/hashicorp/go-multierror"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/apimachinery/pkg/types"
	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/discovery"
	memory "k8s.io/client-go/discovery/cached"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
)

func Apply(ctx context.Context, config *rest.Config, data []byte) error {
	var result *multierror.Error

	chanMes, chanErr := readYaml(data)
	for {
		select {
		case dataBytes, ok := <-chanMes:
			{
				if !ok {
					return result.ErrorOrNil()
				}

				// Get obj and dr
				obj, dr, err := buildDynamicResourceClient(config, dataBytes)
				if err != nil {
					result = multierror.Append(result, err)
					continue
				}

				// Create or Update
				_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, dataBytes, metav1.PatchOptions{
					FieldManager: "kubectl-golang",
				})
				if err != nil {
					result = multierror.Append(result, err)
				}
			}
		case err, ok := <-chanErr:
			if !ok {
				return result.ErrorOrNil()
			}
			if err == nil {
				continue
			}
			result = multierror.Append(result, err)
		}
	}
}

/*
func Apply(ctx context.Context, config *rest.Config, data []byte) (result []string, err error) {
	chanMes, chanErr := readYaml(data)
	for {
		select {
		case dataBytes, ok := <-chanMes:
			{
				if !ok {
					return result, nil
				}

				// Get obj and dr
				obj, dr, err := buildDynamicResourceClient(config, dataBytes)
				if err != nil {
					result = append(result, err.Error())
					continue
				}

				// Create or Update
				_, err = dr.Patch(ctx, obj.GetName(), types.ApplyPatchType, dataBytes, metav1.PatchOptions{
					FieldManager: "kubectl-golang",
				})
				if err != nil {
					result = append(result, err.Error())
				} else {
					result = append(result, obj.GetName()+" patched.")
				}
			}
		case err, ok := <-chanErr:
			if !ok {
				return result, nil
			}
			if err == nil {
				continue
			}
			result = append(result, err.Error())
		}
	}
}
*/
func Delete(ctx context.Context, config *rest.Config, data []byte) (result []string, err error) {
	chanMes, chanErr := readYaml(data)
	for {
		select {
		case dataBytes, ok := <-chanMes:
			{
				if !ok {
					return result, nil
				}

				// Get obj and dr
				obj, dr, err := buildDynamicResourceClient(config, dataBytes)
				if err != nil {
					result = append(result, err.Error())
				}

				// Delete
				deletePolicy := metav1.DeletePropagationBackground
				err = dr.Delete(ctx, obj.GetName(), metav1.DeleteOptions{
					PropagationPolicy: &deletePolicy,
				})
				if err != nil {
					result = append(result, err.Error())
				} else {
					result = append(result, obj.GetName()+" patched.")
				}
			}
		case err, ok := <-chanErr:
			if !ok {
				return result, nil
			}
			if err == nil {
				continue
			}
			result = append(result, err.Error())
		}
	}
}

func readYaml(data []byte) (<-chan []byte, <-chan error) {
	var (
		chanErr        = make(chan error)
		chanBytes      = make(chan []byte)
		multidocReader = utilyaml.NewYAMLReader(bufio.NewReader(bytes.NewReader(data)))
	)

	go func() {
		defer close(chanErr)
		defer close(chanBytes)

		for {
			buf, err := multidocReader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}
				chanErr <- errors.Wrap(err, "failed to read yaml data")
				return
			}
			chanBytes <- buf
		}
	}()
	return chanBytes, chanErr
}

func buildDynamicResourceClient(config *rest.Config, data []byte) (obj *unstructured.Unstructured, dr dynamic.ResourceInterface, err error) {
	// Decode YAML manifest into unstructured.Unstructured
	obj = &unstructured.Unstructured{}
	_, gvk, err := decUnstructured.Decode(data, nil, obj)
	if err != nil {
		return obj, dr, errors.Wrap(err, "Decode yaml failed. ")
	}

	// Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return obj, dr, errors.Wrap(err, "Prepare discovery mapper failed")
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return obj, dr, errors.Wrap(err, "Mapping kind with version failed")
	}

	// Prepare dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return obj, dr, errors.Wrap(err, "Prepare dynamic client failed.")
	}

	// Obtain REST interface for the GVR
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
	} else {
		// for cluster-wide resources
		dr = dynamicClient.Resource(mapping.Resource)
	}

	return obj, dr, nil
}
