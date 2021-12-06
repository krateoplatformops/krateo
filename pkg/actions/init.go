package actions

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/krateoplatformops/krateo/pkg/clients/helmclient/actions"
	"github.com/krateoplatformops/krateo/pkg/clients/kubeclient"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/retrier"
	"github.com/krateoplatformops/krateo/pkg/tmpl"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	kerrors "k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func NewInit(bus eventbus.Bus, kubecfg string) *initAction {
	return &initAction{
		bus:        bus,
		wait:       time.Minute * 5,
		kubeconfig: kubecfg,
	}
}

var _ Action = (*initAction)(nil)

type initAction struct {
	wait       time.Duration
	verbose    bool
	kubeconfig string

	bus eventbus.Bus
	rc  *rest.Config
}

func (o *initAction) SetVerboseEnabled(verbose bool) {
	o.verbose = verbose
}

func (o *initAction) Run() (err error) {
	o.rc, err = clientcmd.BuildConfigFromFlags("", o.kubeconfig)
	if err != nil {
		return err
	}

	// Step 1 - Install crossplane
	o.bus.Publish(events.NewStartWaitEvent("Runtime installation preparation"))

	err = o.installChart("chart-crossplane.yaml", nil,
		actions.WaitTimeout(o.wait), actions.Debug(o.verbose))
	if err != nil {
		// horrible hack to skip error if chart is already installed
		if !strings.HasPrefix(err.Error(), "cannot re-use a name that is still in use") {
			return err
		}
	}
	o.bus.Publish(events.NewDoneEvent("Runtime preparation finished"))

	o.bus.Publish(events.NewStartWaitEvent("Preparing charts provider"))
	err = o.waitUntilProviderHelmIsReady()
	if err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent("Chart provider installed"))

	o.bus.Publish(events.NewStartWaitEvent("Creating role bindings"))
	err = o.createClusterRoleBindingForProviderHelm()
	if err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent("Role bindings created"))

	return nil
}

// installChart use the helm API to execute the command equivalent to:
//
// $ helm repo add [NAME] [URL]
// $ helm repo update
// $ helm install [NAME] [CHART] [flags]
func (o *initAction) installChart(tpl string, data interface{}, opts ...actions.ChartInstallOption) error {
	chart, err := tmpl.Chart(tpl, data)
	if err != nil {
		return err
	}

	return actions.Install(opts...).Do(o.rc, chart)
}

// createClusterRoleBindingForProviderHelm implements these commands:
//
// SA=$(kubectl -n crossplane-system get sa -o name | grep provider-helm | sed -e 's|serviceaccount\/|crossplane-system:|g')
// kubectl create clusterrolebinding provider-helm-admin-binding --clusterrole cluster-admin --serviceaccount="${SA}"
//
// see: https://github.com/crossplane-contrib/provider-helm#testing-in-local-cluster
func (o *initAction) createClusterRoleBindingForProviderHelm() error {
	sas, err := o.getServiceAccounts("crossplane-system")
	if err != nil {
		return err
	}

	var providerHelmSA string
	for _, sa := range sas {
		if strings.Contains(sa, "provider-helm-") {
			providerHelmSA = sa
			break
		}
	}

	src, err := tmpl.Execute(
		"provider-helm-clusterrolebinding.yaml",
		map[string]string{
			"ServiceAccountName": providerHelmSA,
		})
	if err != nil {
		return err
	}

	err = kubeclient.Apply(context.Background(), o.rc, src)
	return ignoreStatusErrorEventually(err, http.StatusConflict)
}

func (o *initAction) getServiceAccounts(ns string) ([]string, error) {
	kc, err := kubeclient.NewKubeClient(o.rc)
	if err != nil {
		return nil, err
	}

	// Let's do the equivalent of this command:
	// kubectl -n crossplane-system get sa -o name
	list := &v1.ServiceAccountList{}

	if err := kc.List(context.TODO(), list, &client.ListOptions{Namespace: ns}); err != nil {
		return nil, errors.Wrapf(err, "cannot get service account list")
	}
	if len(list.Items) == 0 {
		return nil, errors.Errorf("no service account found in namespace: %s", ns)
	}

	var res []string
	for _, el := range list.Items {
		res = append(res, el.Name)
	}

	return res, nil
}

func (o *initAction) waitUntilProviderHelmIsReady() error {
	kc, err := kubeclient.NewKubeClient(o.rc)
	if err != nil {
		return err
	}

	errPodNotReady := errors.New("POD not Ready")

	// create a retrier with constant backoff, retries number of attempts (20) with a 5s sleep between retries.
	r := retrier.New(retrier.ConstantBackoff(20, 5*time.Second), retrier.WhitelistClassifier{errPodNotReady})

	// this counter is just for getting some logging for showcasing, remove in production code.
	attempt := 0

	// retrier works similar to hystrix, we pass the actual work (doing the http request) in a func.
	return r.Run(func() error {
		attempt++

		ready, err := isProviderHelmReady(kc)
		if err == nil && !ready {
			err = errPodNotReady
		}

		return err
	})
}

func isProviderHelmReady(kc client.Client) (bool, error) {
	listOpts := []client.ListOption{
		client.HasLabels{
			"pkg.crossplane.io/revision",
		},
	}
	list := &v1.PodList{}
	if err := kc.List(context.TODO(), list, listOpts...); err != nil {
		return false, errors.Wrapf(err, "cannot get service account list")
	}

	var res bool
	for _, el := range list.Items {
		if strings.Contains(el.GetName(), "provider-helm") {
			if el.Status.Phase == v1.PodRunning {
				res = true
				break
			}
		}
	}

	return res, nil
}

func ignoreStatusErrorEventually(err error, code int32) error {
	if err == nil {
		return nil
	}

	statusErr := k8sStatusErrorLookup(err)
	if statusErr != nil && statusErr.Status().Code != code {
		return err
	}

	return nil
}

func k8sStatusErrorLookup(err error) *kerrors.StatusError {
	if merr, ok := err.(*multierror.Error); ok {
		for _, x := range merr.Errors {
			se, ok := x.(*kerrors.StatusError)
			if ok {
				return se
			}
		}
	}

	res, ok := err.(*kerrors.StatusError)
	if ok {
		return res
	}

	return nil
}
