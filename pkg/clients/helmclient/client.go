package helmclient

import (
	"fmt"
	"io/ioutil"
	"time"

	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/rest"
	ktype "sigs.k8s.io/kustomize/api/types"
)

const (
	helmDriverSecret  = "secret"
	chartCache        = "/tmp/charts"
	releaseMaxHistory = 20
)

const (
	errFailedToCheckIfLocalChartExists = "failed to check if cached chart file exists"
	errFailedToPullChart               = "failed to pull chart"
	errFailedToLoadChart               = "failed to load chart"
	errUnexpectedDirContentTmpl        = "expected 1 .tgz chart file, got [%s]"
)

// HelmClient is the interface to interact with Helm
type HelmClient interface {
	GetLastRelease(release string) (*release.Release, error)
	Install(release string, chart *chart.Chart, vals map[string]interface{}, patches []ktype.Patch, timeout time.Duration) (*release.Release, error)
	Upgrade(release string, chart *chart.Chart, vals map[string]interface{}, patches []ktype.Patch) (*release.Release, error)
	Rollback(release string) error
	Uninstall(release string) error
	PullAndLoadChart(spec *ChartSpec, creds *RepoCreds) (*chart.Chart, error)
}

type helmClient struct {
	log             *logrus.Logger
	pullClient      *action.Pull
	getClient       *action.Get
	installClient   *action.Install
	upgradeClient   *action.Upgrade
	rollbackClient  *action.Rollback
	uninstallClient *action.Uninstall
}

// NewClient returns a new Helm Client with provided config
func NewClient(log *logrus.Logger, config *rest.Config, namespace string) (HelmClient, error) {
	rg := newRESTClientGetter(config, namespace)

	actionConfig := new(action.Configuration)
	// Always store helm state in the same cluster/namespace where chart is deployed
	if err := actionConfig.Init(rg, namespace, helmDriverSecret, func(format string, v ...interface{}) {
		log.Debug(fmt.Sprintf(format, v...))
	}); err != nil {
		return nil, err
	}

	pc := action.NewPull()

	if _, err := os.Stat(chartCache); os.IsNotExist(err) {
		err = os.Mkdir(chartCache, 0750)
		if err != nil {
			return nil, err
		}
	}

	pc.DestDir = chartCache
	pc.Settings = &cli.EnvSettings{}

	gc := action.NewGet(actionConfig)

	ic := action.NewInstall(actionConfig)
	ic.Namespace = namespace
	ic.CreateNamespace = true

	uc := action.NewUpgrade(actionConfig)
	uic := action.NewUninstall(actionConfig)

	rb := action.NewRollback(actionConfig)

	return &helmClient{
		log:             log,
		pullClient:      pc,
		getClient:       gc,
		installClient:   ic,
		upgradeClient:   uc,
		rollbackClient:  rb,
		uninstallClient: uic,
	}, nil
}

func getChartFileName(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	if len(files) != 1 {
		fileNames := make([]string, 0, len(files))
		for _, f := range files {
			fileNames = append(fileNames, f.Name())
		}
		return "", errors.Errorf(errUnexpectedDirContentTmpl, strings.Join(fileNames, ","))
	}
	return files[0].Name(), nil
}

// Pulls latest chart version. Returns absolute chartFilePath or error.
func (hc *helmClient) pullLatestChartVersion(spec *ChartSpec, creds *RepoCreds) (string, error) {
	tmpDir, err := ioutil.TempDir(chartCache, "")
	if err != nil {
		return "", err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			hc.log.WithFields(logrus.Fields{
				"tmpDir": tmpDir,
			}).Info("failed to remove temporary directory")
		}
	}()

	if err := hc.pullChart(spec, creds, tmpDir); err != nil {
		return "", nil
	}

	chartFileName, err := getChartFileName(tmpDir)
	if err != nil {
		return "", err
	}

	chartFilePath := filepath.Join(chartCache, chartFileName)
	if err := os.Rename(filepath.Join(tmpDir, chartFileName), chartFilePath); err != nil {
		return "", nil
	}
	return chartFilePath, nil
}

func (hc *helmClient) pullChart(spec *ChartSpec, creds *RepoCreds, chartDir string) error {
	pc := hc.pullClient

	pc.RepoURL = spec.Repository
	pc.Version = spec.Version
	pc.Username = creds.Username
	pc.Password = creds.Password

	pc.DestDir = chartDir

	o, err := pc.Run(spec.Name)
	hc.log.Debug(o)
	if err != nil {
		return errors.Wrap(err, errFailedToPullChart)
	}
	return nil
}

func (hc *helmClient) PullAndLoadChart(spec *ChartSpec, creds *RepoCreds) (*chart.Chart, error) {
	var chartFilePath string
	var err error
	if spec.Version == "" {
		chartFilePath, err = hc.pullLatestChartVersion(spec, creds)
		if err != nil {
			return nil, err
		}
	} else {
		chartFilePath = filepath.Join(chartCache, fmt.Sprintf("%s-%s.tgz", spec.Name, spec.Version))
		if _, err := os.Stat(chartFilePath); os.IsNotExist(err) {
			if err = hc.pullChart(spec, creds, chartCache); err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, errors.Wrap(err, errFailedToCheckIfLocalChartExists)
		}
	}

	chart, err := loader.Load(chartFilePath)
	if err != nil {
		return nil, errors.Wrap(err, errFailedToLoadChart)
	}
	return chart, nil
}

func (hc *helmClient) GetLastRelease(release string) (*release.Release, error) {
	return hc.getClient.Run(release)
}

func (hc *helmClient) Install(release string, chart *chart.Chart, vals map[string]interface{}, patches []ktype.Patch, timeout time.Duration) (*release.Release, error) {
	hc.installClient.ReleaseName = release

	if timeout > 0 {
		hc.installClient.Wait = true
		hc.installClient.Timeout = timeout
	}

	if len(patches) > 0 {
		hc.installClient.PostRenderer = &KustomizationRender{
			patches: patches,
			log:     hc.log,
		}
	}

	return hc.installClient.Run(chart, vals)
}

func (hc *helmClient) Upgrade(release string, chart *chart.Chart, vals map[string]interface{}, patches []ktype.Patch) (*release.Release, error) {
	// Reset values so that source of truth for desired state is always the CR itself
	hc.upgradeClient.ResetValues = true
	hc.upgradeClient.MaxHistory = releaseMaxHistory

	if len(patches) > 0 {
		hc.upgradeClient.PostRenderer = &KustomizationRender{
			patches: patches,
			log:     hc.log,
		}
	}

	return hc.upgradeClient.Run(release, chart, vals)
}

func (hc *helmClient) Rollback(release string) error {
	return hc.rollbackClient.Run(release)
}

func (hc *helmClient) Uninstall(release string) error {
	_, err := hc.uninstallClient.Run(release)
	if errors.Is(err, driver.ErrReleaseNotFound) {
		return nil
	}
	return err
}
