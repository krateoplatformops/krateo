package helmclient

import (
	"helm.sh/helm/v3/pkg/release"
)

// A ChartSpec defines the chart spec for a Release
type ChartSpec struct {
	Repository    string          `json:"repository"`
	Name          string          `json:"name"`
	Version       string          `json:"version"`
	PullSecretRef SecretReference `json:"pullSecretRef,omitempty"`
}

// NamespacedName represents a namespaced object name
type NamespacedName struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// DataKeySelector defines required spec to access a key of a configmap or secret
type DataKeySelector struct {
	NamespacedName `json:",inline,omitempty"`
	Key            string `json:"key,omitempty"`
	Optional       bool   `json:"optional,omitempty"`
}

// ValueFromSource represents source of a value
type ValueFromSource struct {
	ConfigMapKeyRef *DataKeySelector `json:"configMapKeyRef,omitempty"`
	SecretKeyRef    *DataKeySelector `json:"secretKeyRef,omitempty"`
}

// SetVal represents a "set" value override in a Release
type SetVal struct {
	Name      string           `json:"name"`
	Value     string           `json:"value,omitempty"`
	ValueFrom *ValueFromSource `json:"valueFrom,omitempty"`
}

// ValuesSpec defines the Helm value overrides spec for a Release
type ValuesSpec struct {
	// TODO: investigate using map[string]interface{} instead
	Values     string            `json:"values,omitempty"`
	ValuesFrom []ValueFromSource `json:"valuesFrom,omitempty"`
	Set        []SetVal          `json:"set,omitempty"`
}

// ReleaseParameters are the configurable fields of a Release.
type ReleaseParameters struct {
	Chart       ChartSpec         `json:"chart"`
	Namespace   string            `json:"namespace"`
	PatchesFrom []ValueFromSource `json:"patchesFrom,omitempty"`
	ValuesSpec  `json:",inline"`
}

// ReleaseObservation are the observable fields of a Release.
type ReleaseObservation struct {
	State              release.Status `json:"state,omitempty"`
	ReleaseDescription string         `json:"releaseDescription,omitempty"`
	Revision           int            `json:"revision,omitempty"`
}

// A SecretReference is a reference to a secret in an arbitrary namespace.
type SecretReference struct {
	// Name of the secret.
	Name string `json:"name"`

	// Namespace of the secret.
	Namespace string `json:"namespace"`
}
