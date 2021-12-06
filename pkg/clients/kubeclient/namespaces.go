package kubeclient

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateNamespace(ctx context.Context, kc *kubernetes.Clientset, name string) error {
	ns := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Labels: map[string]string{
				"name": name,
			},
		},
	}
	_, err := kc.CoreV1().Namespaces().Create(ctx, &ns, metav1.CreateOptions{})
	return err
}

func ExistsNamespace(ctx context.Context, kc *kubernetes.Clientset, name string) (bool, error) {
	_, err := kc.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && errors.IsNotFound(statusErr) {
			return false, nil
		}
	}

	return true, err
}

func DeleteNamespace(ctx context.Context, kc *kubernetes.Clientset, name string) error {
	return kc.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{})
}
