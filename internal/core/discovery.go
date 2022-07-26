package core

import (
	"context"
	"errors"
	"fmt"
	"strings"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
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

var (
	ErrCannotResolveResourceType = errors.New("cannot resolve resource type")
	ErrCannotListByAPIResource   = errors.New("cannot list objects by API resource")
)

// APIResource represents a Kubernetes API resource.
type APIResource struct {
	// name is the plural name of the resource.
	Name string `json:"name"`
	// namespaced indicates if a resource is namespaced or not.
	Namespaced bool `json:"namespaced"`
	// group is the preferred group of the resource.  Empty implies the group of the containing resource list.
	// For subresources, this may have a different value, for example: Scale".
	Group string `json:"group,omitempty"`
	// version is the preferred version of the resource.  Empty implies the version of the containing resource list
	// For subresources, this may have a different value, for example: v1 (while inside a v1beta1 version of the core resource's group)".
	Version string `json:"version,omitempty"`
	// kind is the kind for the resource (e.g. 'Foo' is the kind for a resource 'foo')
	Kind string `json:"kind"`
}

func (r *APIResource) GroupKind() schema.GroupKind {
	return schema.GroupKind{
		Group: r.Group,
		Kind:  r.Kind,
	}
}

func (r *APIResource) GroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   r.Group,
		Version: r.Version,
		Kind:    r.Kind,
	}
}

func (r *APIResource) GroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    r.Group,
		Version:  r.Version,
		Resource: r.Name,
	}
}

func (r *APIResource) String() string {
	if len(r.Group) == 0 {
		return fmt.Sprintf("%s.%s", r.Name, r.Version)
	}
	return fmt.Sprintf("%s.%s.%s", r.Name, r.Version, r.Group)
}

func (r *APIResource) WithGroupString() string {
	if len(r.Group) == 0 {
		return r.Name
	}
	return r.Name + "." + r.Group
}

type ResolveAPIResourceOpts struct {
	RESTConfig *rest.Config
	Query      string
}

func ResolveAPIResource(opts ResolveAPIResourceOpts) (*APIResource, error) {
	dc, err := discovery.NewDiscoveryClientForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	var gvr schema.GroupVersionResource
	var gvk schema.GroupVersionKind

	// Resolve type string into GVR
	fullySpecifiedGVR, gr := schema.ParseResourceArg(strings.ToLower(opts.Query))
	if fullySpecifiedGVR != nil {
		gvr, _ = mapper.ResourceFor(*fullySpecifiedGVR)
	}
	if gvr.Empty() {
		gvr, err = mapper.ResourceFor(gr.WithVersion(""))
		if err != nil {
			return nil, ErrCannotResolveResourceType
		}
	}
	// Obtain Kind from GVR
	gvk, err = mapper.KindFor(gvr)
	if gvk.Empty() {
		if err != nil {
			return nil, ErrCannotResolveResourceType
		}
	}
	// Determine scope of resource
	mapping, err := mapper.RESTMapping(gvk.GroupKind())
	if err != nil {
		return nil, ErrCannotResolveResourceType
	}
	// NOTE: This is a rather incomplete APIResource object, but it has enough
	//       information inside for our use case, which is to fetch API objects
	res := &APIResource{
		Name:       gvr.Resource,
		Namespaced: mapping.Scope.Name() == meta.RESTScopeNameNamespace,
		Group:      gvk.Group,
		Version:    gvk.Version,
		Kind:       gvk.Kind,
	}

	return res, nil
}

type ListByAPIResourceOpts struct {
	RESTConfig  *rest.Config
	APIResource APIResource
	Namespace   string
}

// listByAPI list all objects of the provided API & namespace. If listing the
// API at the cluster scope, set the namespace argument as an empty string.
func ListByAPIResource(ctx context.Context, opts ListByAPIResourceOpts) ([]unstructured.Unstructured, error) {
	dynamicClient, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	var ri dynamic.ResourceInterface
	var items []unstructured.Unstructured
	var next string

	isClusterScopeRequest := !opts.APIResource.Namespaced || opts.Namespace == ""
	if isClusterScopeRequest {
		ri = dynamicClient.Resource(opts.APIResource.GroupVersionResource())
	} else {
		ri = dynamicClient.Resource(opts.APIResource.GroupVersionResource()).Namespace(opts.Namespace)
	}

	for {
		objectList, err := ri.List(ctx, metav1.ListOptions{
			Limit:    250,
			Continue: next,
		})
		if err != nil {
			switch {
			case apierrors.IsForbidden(err):
				return nil, err
			case apierrors.IsNotFound(err):
				break
			default:
				return nil, ErrCannotListByAPIResource
			}
		}

		if objectList == nil {
			break
		}

		items = append(items, objectList.Items...)
		next = objectList.GetContinue()
		if len(next) == 0 {
			break
		}
	}

	return items, nil
}

type GetByAPIResourceOpts struct {
	RESTConfig  *rest.Config
	APIResource APIResource
	Name        string
	Namespace   string
}

// GetByAPIResource returns an object that matches the provided name & options on the server.
func GetByAPIResource(ctx context.Context, opts GetByAPIResourceOpts) (*unstructured.Unstructured, error) {
	dynamicClient, err := dynamic.NewForConfig(opts.RESTConfig)
	if err != nil {
		return nil, err
	}

	gvr := opts.APIResource.GroupVersionResource()

	var ri dynamic.ResourceInterface
	if opts.APIResource.Namespaced {
		ri = dynamicClient.Resource(gvr).Namespace(opts.Namespace)
	} else {
		ri = dynamicClient.Resource(gvr)
	}

	return ri.Get(ctx, opts.Name, metav1.GetOptions{})
}
