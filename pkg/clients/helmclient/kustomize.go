package helmclient

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

const (
	kustomizationFileName  = "kustomization.yaml"
	helmOutputFileName     = "helm-output.yaml"
	helmTempDirNamePattern = "helm-post-render"
)

// KustomizationRender Implements helm PostRenderer interface
type KustomizationRender struct {
	patches []types.Patch
	log     *logrus.Logger
}

// Run runs a set of Kustomize patches against yaml input and returns the patched content.
func (kr KustomizationRender) Run(renderedManifests *bytes.Buffer) (modifiedManifests *bytes.Buffer, err error) {
	d, err := ioutil.TempDir("", helmTempDirNamePattern)
	if err != nil {
		return nil, err
	}

	fsys := filesys.MakeFsOnDisk()
	defer func() {
		if err := fsys.RemoveAll(d); err != nil {
			log.Printf("Failed to cleanup tmp data (path=%s, err=%s)", d, err)
			//kr.logger.Info().Msgf("Failed to cleanup tmp data", "path", d, "err", err)
		}
	}()

	k := types.Kustomization{
		Resources: []string{helmOutputFileName},
		Patches:   kr.patches,
	}

	kdata, err := json.Marshal(k)
	if err != nil {
		return nil, err
	}

	err = fsys.WriteFile(filepath.Join(d, kustomizationFileName), kdata)
	if err != nil {
		return nil, err
	}

	err = fsys.WriteFile(filepath.Join(d, helmOutputFileName), renderedManifests.Bytes())
	if err != nil {
		return nil, err
	}

	opts := &krusty.Options{
		DoLegacyResourceSort: false,
		LoadRestrictions:     types.LoadRestrictionsRootOnly,
		DoPrune:              false,
		PluginConfig:         disabledPluginConfig(),
	}

	kust := krusty.MakeKustomizer(opts)
	m, err := kust.Run(fsys, d)
	if err != nil {
		return nil, err
	}

	yml, err := m.AsYaml()
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(yml), nil
}

func disabledPluginConfig() *types.PluginConfig {
	return types.MakePluginConfig(
		types.PluginRestrictionsBuiltinsOnly,
		types.BploUseStaticallyLinked)
}
