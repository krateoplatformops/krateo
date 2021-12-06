package actions

import (
	"context"
	"os"
	"time"

	"github.com/krateoplatformops/krateo/pkg/clients/helmclient"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

type ChartInstallOption func(*ChartInstall)

func WaitTimeout(how time.Duration) ChartInstallOption {
	return func(act *ChartInstall) {
		act.waitTimeout = how
	}
}

func Debug(enable bool) ChartInstallOption {
	return func(act *ChartInstall) {
		act.debug = enable
	}
}

func ReleaseName(name string) ChartInstallOption {
	return func(act *ChartInstall) {
		act.releaseName = name
	}
}

type ChartInstall struct {
	waitTimeout time.Duration
	debug       bool
	releaseName string
}

func Install(opts ...ChartInstallOption) *ChartInstall {
	res := &ChartInstall{}

	for _, o := range opts {
		o(res)
	}

	return res
}

func (act *ChartInstall) Do(rc *rest.Config, cr *helmclient.ReleaseParameters) error {
	log := logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stderr)

	if act.debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	hc, err := helmclient.NewClient(log, rc, cr.Namespace)
	if err != nil {
		return err
	}

	// create the clientset
	kc, err := kclient.New(rc, kclient.Options{})
	if err != nil {
		return err
	}

	ctx := context.Background()

	cv, err := composeValuesFromSpec(ctx, kc, cr.ValuesSpec)
	if err != nil {
		return errors.Wrap(err, errFailedToComposeValues)
	}

	creds, err := repoCredsFromSecret(ctx, kc, cr.Chart.PullSecretRef)
	if err != nil {
		return errors.Wrap(err, errFailedToGetRepoCreds)
	}

	patch := newPatcher()
	p, err := patch.getFromSpec(ctx, kc, cr.PatchesFrom)
	if err != nil {
		return errors.Wrap(err, errFailedToLoadPatches)
	}

	chart, err := hc.PullAndLoadChart(&cr.Chart, creds)
	if err != nil {
		return err
	}
	if cr.Chart.Version == "" {
		cr.Chart.Version = chart.Metadata.Version
	}

	releaseName := act.releaseName
	if releaseName == "" {
		releaseName = chart.Name()
	}

	rel, err := hc.Install(releaseName, chart, cv, p, act.waitTimeout)
	if err != nil {
		return err
	}

	if rel == nil {
		return errors.New(errLastReleaseIsNil)
	}

	return nil
}
