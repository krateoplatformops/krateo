package kubeclient

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// A Secret is a secret in an arbitrary namespace.
type Secret struct {
	// Name of the secret.
	Name string

	// Namespace of the secret.
	Namespace string

	// Data allows specifying non-binary secret data in string form.
	Data map[string]string

	// Map of string keys and values that can be used to organize and categorize
	// (scope and select) secrets.
	Labels map[string]string
}

func CreateSecret(ctx context.Context, kc *kubernetes.Clientset, ref *Secret) error {
	metadata := metav1.ObjectMeta{
		Name:   ref.Name,
		Labels: ref.Labels,
	}
	if ref.Namespace != "" {
		metadata.Namespace = ref.Namespace
	}

	secret := corev1.Secret{
		ObjectMeta: metadata,
		StringData: ref.Data,
	}

	_, err := kc.CoreV1().Secrets(metadata.Namespace).Create(ctx, &secret, metav1.CreateOptions{})
	return err
}

func ExistsSecret(ctx context.Context, kc *kubernetes.Clientset, name, namespace string) (bool, error) {
	_, err := kc.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			return false, nil
		}
		//log.Printf("[DEBUG] Received error: %#v", err)
	}

	return true, err
}

func GetSecret(ctx context.Context, kc *kubernetes.Clientset, name, namespace string) (*Secret, error) {
	exists, err := ExistsSecret(ctx, kc, name, namespace)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, nil
	}

	secret, err := kc.CoreV1().Secrets(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	res := &Secret{
		Name:      secret.Name,
		Namespace: secret.Namespace,
		Labels:    map[string]string{},
		Data:      map[string]string{},
	}

	for k, v := range secret.Labels {
		res.Labels[k] = v
	}

	for k, v := range secret.Data {
		res.Data[k] = string(v)
	}

	return res, nil
}

func DeleteSecret(ctx context.Context, kc *kubernetes.Clientset, name, namespace string) error {
	return kc.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}
