package helm

import (
	"bytes"
	"sort"

	"github.com/Masterminds/semver"
	"github.com/krateoplatformops/krateo/internal/httputils"
	"helm.sh/helm/v3/pkg/repo"
	"sigs.k8s.io/yaml"
)

func LatestVersionAndURL(idx *repo.IndexFile) (string, string, error) {
	vs := map[*semver.Version]string{}
	keys := []*semver.Version{}

	for _, cvs := range idx.Entries {
		for i := len(cvs) - 1; i >= 0; i-- {
			if len(cvs[i].URLs) > 0 {
				v, err := semver.NewVersion(cvs[i].AppVersion)
				if err != nil {
					return "", "", err
				}
				vs[v] = cvs[i].URLs[0]
				keys = append(keys, v)
			}
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(keys)))

	return keys[0].String(), vs[keys[0]], nil
}

// IndexFromURL loads an index file from an URL.
func IndexFromURL(url string) (*repo.IndexFile, error) {
	buf := &bytes.Buffer{}
	if err := httputils.Fetch(url, buf); err != nil {
		return nil, err
	}
	return IndexFromBytes(buf.Bytes())
}

// IndexFromBytes loads an index file and does minimal validity checking.
// This will fail if API Version is not set (ErrNoAPIVersion) or if the unmarshal fails.
func IndexFromBytes(data []byte) (*repo.IndexFile, error) {
	i := &repo.IndexFile{}

	if len(data) == 0 {
		return i, repo.ErrEmptyIndexYaml
	}

	if err := yaml.Unmarshal(data, i); err != nil {
		return i, err
	}

	i.SortEntries()
	//if i.APIVersion == "" {
	//	return i, repo.ErrNoAPIVersion
	//}
	return i, nil
}
