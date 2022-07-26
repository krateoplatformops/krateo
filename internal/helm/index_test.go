//go:build integration
// +build integration

package helm

import (
	"sort"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/assert"
)

const (
	crossplaneHelmIndexURL = "https://charts.crossplane.io/stable/index.yaml"
)

func TestIndexFromURL(t *testing.T) {
	idx, err := IndexFromURL(crossplaneHelmIndexURL)
	assert.Nil(t, err, "expecting nil error fetching index")

	vs := map[*semver.Version]string{}
	keys := []*semver.Version{}

	for _, cvs := range idx.Entries {
		for i := len(cvs) - 1; i >= 0; i-- {
			if len(cvs[i].URLs) > 0 {
				v, err := semver.NewVersion(cvs[i].AppVersion)
				assert.Nil(t, err, "expecting nil error parsing semver: ", cvs[i].AppVersion)
				vs[v] = cvs[i].URLs[0]
				keys = append(keys, v)
			}
		}
	}

	sort.Sort(sort.Reverse(semver.Collection(keys)))

	latest := vs[keys[0]]
	t.Logf("latest: %s\n", latest)
}
