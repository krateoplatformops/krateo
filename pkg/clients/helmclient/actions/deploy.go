package actions

import (
	"context"
	"os"

	"github.com/krateoplatformops/krateo/pkg/clients/helmclient"
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

const (
	errLastReleaseIsNil = "last helm release is nil"
	// errFailedToInstall            = "failed to install release"
	// errFailedToUninstall          = "failed to uninstall release"
	errFailedToGetRepoCreds  = "failed to get user name and password from secret reference"
	errFailedToComposeValues = "failed to compose values"
	// errFailedToCreateRestConfig   = "cannot create new rest config using provider secret"
	errFailedToLoadPatches = "failed to load patches"
)

/*
func Install(ctx context.Context, rc *rest.Config, cr *helmclient.ReleaseParameters, releaseName string, timeout time.Duration, debug bool) error {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	if debug {
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

	rel, err := hc.Install(releaseName, chart, cv, p, timeout)
	if err != nil {
		return err
	}

	if rel == nil {
		return errors.New(errLastReleaseIsNil)
	}

	return nil
}
*/
func Uninstall(ctx context.Context, rc *rest.Config, namespace, releaseName string, debug bool) error {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	if debug {
		log.SetLevel(logrus.DebugLevel)
	} else {
		log.SetLevel(logrus.InfoLevel)
	}

	hc, err := helmclient.NewClient(log, rc, namespace)
	if err != nil {
		return err
	}

	return hc.Uninstall(releaseName)
}
