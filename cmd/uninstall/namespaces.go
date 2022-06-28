package uninstall

import (
	"github.com/krateoplatformops/krateo/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func deleteNamespaceForcingFinalizers(rc *rest.Config, dryRun bool, name string) error {
	cli, err := kubernetes.Namespaces(rc)
	if err != nil {
		return err
	}

	do := metav1.DeleteOptions{}
	if dryRun {
		do.DryRun = []string{metav1.DryRunAll}
	}

	err = cli.Delete(name, do)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	src, err := cli.Get(name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	src.Spec.Finalizers = []corev1.FinalizerName{}

	uo := metav1.UpdateOptions{}
	if dryRun {
		uo.DryRun = []string{metav1.DryRunAll}
	}
	_, err = cli.Finalize(src, uo)
	if err != nil {
		return err
	}

	err = cli.Delete(name, do)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	return nil
}
