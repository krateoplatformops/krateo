package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/clientcmd"
)

func TestWatch(t *testing.T) {
	const namespace = "krateo-system"

	kubeconfig, err := ioutil.ReadFile(os.Getenv(clientcmd.RecommendedConfigPathEnvVar))
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	sel, err := labels.Parse("app=crossplane")
	assert.Nil(t, err, "expecting nil error parsing label")

	stopFn := func(et watch.EventType, obj *unstructured.Unstructured) (bool, error) {
		pod := &corev1.Pod{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(obj.UnstructuredContent(), &pod)
		if err != nil {
			return false, err
		}
		fmt.Printf("%+v\n", pod)
		for _, cond := range pod.Status.Conditions {
			if (cond.Type == corev1.PodReady) && (cond.Status == corev1.ConditionTrue) {
				return true, nil
			}
		}

		return false, nil
	}

	err = Watch(context.TODO(), WatchOpts{
		RESTConfig: restConfig,
		GVR:        schema.GroupVersionResource{Version: "v1", Resource: "pods"},
		Namespace:  namespace,
		Selector:   sel,
		StopFn:     stopFn,
	})
	assert.Nil(t, err, "expecting nil error watching pod")
}
