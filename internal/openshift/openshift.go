package openshift

import (
	"context"

	"github.com/krateoplatformops/krateo/internal/core"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/rest"
)

const (
	openshiftClusterRoleName        = "openshift-crossplane-full-finalizer-role"
	openshiftClusterRoleBindingName = "openshift-crossplane-full-finalizer-role-binding"
)

func CreateClusterRole(ctx context.Context, rc *rest.Config) error {
	cr := rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				core.InstalledByLabel: core.InstalledByValue,
			},
			Name: openshiftClusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"*"},
				Resources: []string{"*/finalizers"},
				Verbs:     []string{"*"},
			},
		},
	}

	dat, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&cr)
	if err != nil {
		return err
	}

	obj := unstructured.Unstructured{}
	obj.SetUnstructuredContent(dat)

	gvk := schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "ClusterRole",
	}

	dr, err := core.DynamicForGVR(rc, gvk, "")
	if err != nil {
		if core.IsNoKindMatchError(err) {
			return nil
		}
		return err
	}
	_, err = dr.Create(ctx, &obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

func CreateClusterRoleBinding(ctx context.Context, rc *rest.Config) error {
	crb := rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				core.InstalledByLabel: core.InstalledByValue,
			},
			Name: openshiftClusterRoleBindingName,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     openshiftClusterRoleName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      "rbac-manager",
				Namespace: "crossplane-system",
			},
			{
				Kind:      "ServiceAccount",
				Name:      "crossplane",
				Namespace: "crossplane-system",
			},
		},
	}

	dat, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&crb)
	if err != nil {
		return err
	}

	obj := unstructured.Unstructured{}
	obj.SetUnstructuredContent(dat)

	gvk := schema.GroupVersionKind{
		Group:   "rbac.authorization.k8s.io",
		Version: "v1",
		Kind:    "ClusterRoleBinding",
	}

	dr, err := core.DynamicForGVR(rc, gvk, "")
	if err != nil {
		if core.IsNoKindMatchError(err) {
			return nil
		}
		return err
	}
	_, err = dr.Create(ctx, &obj, metav1.CreateOptions{})
	if err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}
