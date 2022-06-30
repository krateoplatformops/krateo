package providers

import "fmt"

type providerInfo struct {
	name            string
	version         string
	imageUriPattern string
	metaUrl         string
}

func (pi *providerInfo) Name() string {
	return pi.name
}

func (pi *providerInfo) Version() string {
	return pi.version
}

func (pi *providerInfo) Image() string {
	return fmt.Sprintf(pi.imageUriPattern, pi.name, pi.version)
}

func (pi *providerInfo) MetaUrl() string {
	return pi.metaUrl
}

func Helm() *providerInfo {
	return &providerInfo{
		name:            "helm",
		version:         "v0.10.0",
		imageUriPattern: "registry.upbound.io/crossplane/provider-%s:%s",
		metaUrl:         "https://raw.githubusercontent.com/crossplane-contrib/provider-helm/master/package/crossplane.yaml",
	}
}

func Kubernetes() *providerInfo {
	return &providerInfo{
		name:            "kubernetes",
		version:         "v0.3.0",
		imageUriPattern: "registry.upbound.io/crossplane/provider-%s:%s",
		metaUrl:         "https://raw.githubusercontent.com/crossplane-contrib/provider-kubernetes/master/package/crossplane.yaml",
	}
}

func Git() *providerInfo {
	return &providerInfo{
		name:            "git",
		version:         "v1.0.0",
		imageUriPattern: "ghcr.io/krateoplatformops/provider-%s:%s",
		metaUrl:         "https://raw.githubusercontent.com/krateoplatformops/provider-git/main/package/crossplane.yaml",
	}
}

func GitHub() *providerInfo {
	return &providerInfo{
		name:            "github",
		version:         "v1.0.0",
		imageUriPattern: "ghcr.io/krateoplatformops/provider-%s:%s",
		metaUrl:         "https://raw.githubusercontent.com/krateoplatformops/provider-github/main/package/crossplane.yaml",
	}
}

func ArgoCdToken() *providerInfo {
	return &providerInfo{
		name:            "argocd-token",
		version:         "v1.0.0",
		imageUriPattern: "ghcr.io/krateoplatformops/provider-%s:%s",
		metaUrl:         "https://raw.githubusercontent.com/krateoplatformops/provider-argocd-token/main/package/crossplane.yaml",
	}
}

func All() []*providerInfo {
	return []*providerInfo{
		Helm(),
		Kubernetes(),
		ArgoCdToken(),
		Git(),
		GitHub(),
	}
}
