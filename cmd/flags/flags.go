package flags

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	Kubeconfig = "kubeconfig"
	RepoToken  = "repo-token"
	RepoURL    = "repo-url"
	Verbose    = "verbose"
)

func DefaultKubeconfigValue() string {
	value := os.Getenv("KUBECONFIG")

	if value == "" {
		home := homedir.HomeDir()
		if home != "" {
			value = filepath.Join(home, ".kube", "config")
		}
	}

	fp, err := os.Open(value)
	if errors.Is(err, os.ErrNotExist) {
		value = ""
	}
	defer fp.Close()

	return value
}

func GetKubeconfig(cmd *cobra.Command) (*rest.Config, error) {
	// Grab the kubeconfig path
	cfg, err := cmd.Flags().GetString(Kubeconfig)
	if err != nil {
		return nil, err
	}

	// Use the current context in kubeconfig
	return clientcmd.BuildConfigFromFlags("", cfg)
}
