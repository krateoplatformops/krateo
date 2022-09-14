package configurations

import (
	"context"
	"time"

	"github.com/krateoplatformops/krateo/internal/core"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

// Condition types.
const (
	// A TypeInstalled indicates whether a package has been installed.
	TypeInstalled ConditionType = "Installed"

	// A TypeHealthy indicates whether a package is healthy.
	TypeHealthy ConditionType = "Healthy"
)

// A ConditionType represents a condition a resource could be in.
type ConditionType string

// A ConditionReason represents the reason a resource is in a condition.
type ConditionReason string

// A Condition that may apply to a resource.
type Condition struct {
	// Type of this condition. At most one of each condition type may apply to
	// a resource at any point in time.
	Type ConditionType `json:"type"`

	// Status of this condition; is it currently True, False, or Unknown?
	Status corev1.ConditionStatus `json:"status"`

	// LastTransitionTime is the last time this condition transitioned from one
	// status to another.
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`

	// A Reason for this condition's last transition from one status to another.
	Reason ConditionReason `json:"reason"`

	// A Message containing details about this condition's last transition from
	// one status to another, if any.
	Message string `json:"message,omitempty"`
}

// A ConditionedStatus reflects the observed status of a resource. Only
// one condition of each type may exist.
type ConditionedStatus struct {
	// Conditions of the resource.
	Conditions []Condition `json:"conditions,omitempty"`
}

func waitUntilHealtyAndInstalled(ctx context.Context, restConfig *rest.Config, name string) error {
	stopFn := func(et watch.EventType, obj *unstructured.Unstructured) (bool, error) {
		if obj.GetName() != name {
			return false, nil
		}

		val, ok, err := unstructured.NestedFieldCopy(obj.Object, "status")
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
		//fmt.Printf("%T  ==>>  %+v\n", val, val)
		var status ConditionedStatus
		err = core.FromUnstructuredViaJSON(val.(map[string]interface{}), &status)
		if err != nil {
			return false, err
		}

		healty := false
		installed := false
		for _, cond := range status.Conditions {
			if cond.Type == TypeHealthy {
				healty = (cond.Status == corev1.ConditionTrue)
			}

			if cond.Type == TypeInstalled {
				installed = (cond.Status == corev1.ConditionTrue)
			}
		}

		return (healty && installed), nil
	}

	return core.Watch(ctx, core.WatchOpts{
		RESTConfig: restConfig,
		GVR: schema.GroupVersionResource{
			Group:    "pkg.crossplane.io",
			Version:  "v1",
			Resource: "configurations",
		},
		StopFn:  stopFn,
		Timeout: time.Minute * 2,
	})
}
