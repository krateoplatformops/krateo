//go:build integration
// +build integration

package clusterrolebindings

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCreate(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	all, err := core.List(context.TODO(), core.ListOpts{
		RESTConfig: restConfig,
		GVK: schema.GroupVersionKind{
			Version: "v1",
			Kind:    "Serviceaccount",
		},
		Namespace: core.DefaultNamespace,
	})
	assert.Nil(t, err, "expecting nil error listing service accounts")

	acceptFn := func(el unstructured.Unstructured) bool {
		keep := strings.HasPrefix(el.GetName(), "provider-helm")
		keep = keep || strings.HasPrefix(el.GetName(), "provider-kubernetes")
		return keep
	}

	res, err := core.Filter(all, acceptFn)
	assert.Nil(t, err, "expecting nil error filtering service accounts")

	t.Logf("found [%d] service accounts\n", len(res))
	for _, el := range res {
		idx := strings.LastIndex(el.GetName(), "-")
		name := fmt.Sprintf("%s-admin-binding", el.GetName()[0:idx])
		t.Logf("\n > creating clusterrolebinding: %s\n", name)

		err := Create(context.TODO(), CreateOptions{
			RESTConfig:       restConfig,
			Name:             name,
			SubjectName:      el.GetName(),
			SubjectNamespace: core.DefaultNamespace,
		})
		assert.Nil(t, err, "expecting nil error creating clusterrolebindings")
	}
}
