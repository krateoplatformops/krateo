package claims

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/krateoplatformops/krateo/internal/core"
	"github.com/krateoplatformops/krateo/internal/maps"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
)

func TestCore(t *testing.T) {
	kubeconfig, err := ioutil.ReadFile(clientcmd.RecommendedHomeFile)
	assert.Nil(t, err, "expecting nil error loading kubeconfig")

	restConfig, err := core.RESTConfigFromBytes(kubeconfig, "")
	assert.Nil(t, err, "expecting nil error creating rest.Config")

	data := maps.Map(map[string]interface{}{
		"platform":  "kubernetes",
		"namespace": "krateo-namespace",
		"domain":    "krateo.site",
	})

	dflts := maps.Map(CoreDefaultClaims())
	inp := data.MergeHere(dflts)

	err = Create(context.TODO(), CreateCoreOpts{
		RESTConfig: restConfig,
		Data:       inp,
	})
	assert.Nil(t, err, "expecting nil error creating claim")

	url := fmt.Sprintf("%s://%s.%s",
		data.GetString("protocol"),
		data.GetString("app.hostname"),
		data.GetString("domain"))

	t.Logf("%s\n", url)
}
