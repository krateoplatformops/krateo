package actions

import (
	"context"
	"encoding/base64"

	"github.com/krateoplatformops/krateo/pkg/clients/kubeclient"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
)

const (
	kubeSystemNamespace = "kube-system"
)

func getClusterId(cfg *rest.Config) (string, error) {
	kc, err := kubeclient.NewKubeClient(cfg)
	if err != nil {
		return "", err
	}

	// Let's do the equivalent of this command:
	// kubectl get serviceaccount default \
	//    -o=jsonpath='{.secrets[0].name}' | xargs kubectl get secret \
	//    -ojsonpath='{.data.ca\.crt}' | base64 --decode
	sa := &v1.ServiceAccount{}
	err = kc.Get(context.TODO(), types.NamespacedName{
		Namespace: kubeSystemNamespace,
		Name:      "default",
	}, sa)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get service account")
	}
	if len(sa.Secrets) == 0 {
		return "", errors.New("no secrets found for service account")
	}

	s := &v1.Secret{}
	err = kc.Get(context.TODO(), types.NamespacedName{
		Namespace: kubeSystemNamespace,
		Name:      sa.Secrets[0].Name,
	}, s)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get %s secret", sa.Secrets[0].Name)
	}

	h := blake2b.Sum256(s.Data["ca.crt"])
	clusterId := base64.URLEncoding.EncodeToString(h[:])
	return clusterId, nil
}
