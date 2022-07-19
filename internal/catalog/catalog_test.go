//go:build integration
// +build integration

package catalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetch(t *testing.T) {
	all, err := Fetch()
	if err != nil {
		t.Fatal(err)
	}

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, all, "expecting non-nil result")
	assert.Greater(t, len(all.Items), 0)
}

func TestPackagesToInstall(t *testing.T) {
	all, err := Fetch()

	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, all, "expecting non-nil result")

	toInstall, err := FilterBy(ForCLI())
	assert.Nil(t, err, "expecting nil error")
	assert.NotNil(t, toInstall, "expecting non-nil result")

	assert.Greater(t, len(all.Items), len(toInstall.Items))
}
