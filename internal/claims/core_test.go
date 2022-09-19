package claims

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/krateoplatformops/krateo/internal/strvals"
	"github.com/stretchr/testify/assert"
)

func TestCore(t *testing.T) {
	//kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	kubeconfig, err := ioutil.ReadFile("../../testdata/krateo-test.yml")
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	/*
		data := maps.Map(map[string]interface{}{
			"platform":  "kubernetes",
			"namespace": "krateo-namespace",
			"domain":    "krateo.site",
		})*/

	vals := []string{
		"platform=kubernetes",
		"namespace=krateo-namespace",
		"domain=krateo.site",
	}
	inp, err := strvals.ParseString(strings.Join(vals, ","))
	assert.Nil(t, err, "expecting nil error parsing string values")

	//dflts := maps.Map(CoreDefaultClaims())
	//inp := data.MergeHere(dflts)

	err = Create(context.TODO(), CreateCoreOpts{
		RESTConfig: restConfig,
		Data:       inp,
	})
	assert.Nil(t, err, "expecting nil error creating claim")

	url := fmt.Sprintf("https://app.%s", inp["domain"])

	t.Logf("%s\n", url)
}
