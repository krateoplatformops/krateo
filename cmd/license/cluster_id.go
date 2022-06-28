package license

import (
	"encoding/base64"
	"fmt"

	"github.com/krateoplatformops/krateo/pkg/kubernetes"
	"golang.org/x/crypto/blake2b"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func getClusterId(rc *rest.Config) (string, error) {
	sacl, err := kubernetes.ServiceAccounts(rc)
	if err != nil {
		return "", err
	}

	// Let's do the equivalent of this command:
	// kubectl get serviceaccount default \
	//    -o=jsonpath='{.secrets[0].name}' | xargs kubectl get secret \
	//    -ojsonpath='{.data.ca\.crt}' | base64 --decode
	sa, err := sacl.Get("default", kubernetes.KubeSystemNamespace, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("cannot get service account: %w", err)
	}

	if len(sa.Secrets) == 0 {
		return "", fmt.Errorf("no secrets found for service account")
	}

	secrets, err := kubernetes.Secrets(rc)
	if err != nil {
		return "", err
	}

	res, err := secrets.Get(sa.Secrets[0].Name, kubernetes.KubeSystemNamespace, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("cannot get %s secret: %w", sa.Secrets[0].Name, err)
	}

	h := blake2b.Sum256(res.Data["ca.crt"])
	clusterId := base64.URLEncoding.EncodeToString(h[:])
	return clusterId, nil
}

func createLicenseSecret(rc *rest.Config, key string) error {
	cli, err := kubernetes.Secrets(rc)
	if err != nil {
		return err
	}

	sec := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "krateo-license",
			Namespace: kubernetes.KrateoSystemNamespace,
		},
		Type: "Opaque",
		StringData: map[string]string{
			"payload": key,
		},
	}

	_, err = cli.Create(kubernetes.KrateoSystemNamespace, &sec, metav1.CreateOptions{
		FieldManager: kubernetes.DefaultFieldManager,
	})

	return err
}
