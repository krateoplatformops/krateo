package actions

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/krateoplatformops/krateoctl/pkg/clients/helmclient"
	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	ktypes "sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

const (
	keyDefaultPatchFrom        = "patch.yaml"
	errFailedToUnmarshallPatch = "failed to unmarshal patch"
)

// Patcher interface for managing Kustomize patches and detecting updates
type Patcher interface {
	patchGetter
	patchHasher
}

type patchGetter interface {
	getFromSpec(ctx context.Context, kube client.Client, vals []helmclient.ValueFromSource) ([]ktypes.Patch, error)
}

type patchHasher interface {
	shaOf(patches []ktypes.Patch) (string, error)
}

func newPatcher() Patcher {
	return patch{
		patchHasher: patchSha{},
		patchGetter: patchGet{},
	}
}

type patch struct {
	patchHasher
	patchGetter
}

type patchSha struct{}

func (patchSha) shaOf(patches []ktypes.Patch) (string, error) {
	if len(patches) == 0 {
		return "", nil
	}

	jb, err := json.Marshal(patches)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", sha256.Sum256(jb)), nil
}

type patchGet struct{}

func (patchGet) getFromSpec(ctx context.Context, kube client.Client, vals []helmclient.ValueFromSource) ([]ktypes.Patch, error) {
	var base []ktypes.Patch // nolint:prealloc

	for _, vf := range vals {
		s, err := getDataValueFromSource(ctx, kube, vf, keyDefaultPatchFrom)
		if err != nil {
			return nil, errors.Wrap(err, errFailedToGetValueFromSource)
		}

		if s == "" {
			continue
		}

		var p struct {
			Patches []ktypes.Patch `json:"patches"`
		}
		if err = yaml.Unmarshal([]byte(s), &p); err != nil {
			return nil, errors.Wrap(err, errFailedToUnmarshallPatch)
		}
		base = append(base, p.Patches...)
	}

	return base, nil
}
