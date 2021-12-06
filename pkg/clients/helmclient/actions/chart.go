package actions

import (
	"context"

	"github.com/krateoplatformops/krateoctl/pkg/clients/helmclient"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	keyRepoUsername = "username"
	keyRepoPassword = "password"
)

const (
	errFailedToGetRepoPullSecret       = "failed to get repo pull secret"
	errChartPullSecretMissingNamespace = "namespace must be set in chart pull secret ref"
	errChartPullSecretMissingUsername  = "username missing in chart pull secret"
	errChartPullSecretMissingPassword  = "password missing in chart pull secret"
)

func repoCredsFromSecret(ctx context.Context, kube client.Client, secretRef helmclient.SecretReference) (*helmclient.RepoCreds, error) {
	repoUser := ""
	repoPass := ""
	if secretRef.Name != "" {
		if secretRef.Namespace == "" {
			return nil, errors.New(errChartPullSecretMissingNamespace)
		}
		d, err := getSecretData(ctx, kube, types.NamespacedName{Name: secretRef.Name, Namespace: secretRef.Namespace})
		if err != nil {
			return nil, errors.Wrap(err, errFailedToGetRepoPullSecret)
		}
		repoUser = string(d[keyRepoUsername])
		if repoUser == "" {
			return nil, errors.New(errChartPullSecretMissingUsername)
		}
		repoPass = string(d[keyRepoPassword])
		if repoPass == "" {
			return nil, errors.New(errChartPullSecretMissingPassword)
		}
	}

	return &helmclient.RepoCreds{
		Username: repoUser,
		Password: repoPass,
	}, nil
}
