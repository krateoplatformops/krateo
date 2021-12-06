package crds

import (
	"github.com/krateoplatformops/krateoctl/pkg/strvals"
	"sigs.k8s.io/yaml"
)

//nolint:errcheck
func ToYAML(src []string) ([]byte, error) {
	res := make(map[string]interface{})

	for _, s := range src {
		strvals.ParseInto(s, res)
	}

	return yaml.Marshal(res)
}
