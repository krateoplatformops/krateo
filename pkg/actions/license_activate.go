package actions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/krateoplatformops/krateo/pkg/clients/kubeclient"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/tmpl"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	activationURL = "http://localhost:8090/activate"
)

func NewLicenseActivate(bus eventbus.Bus, kubecfg, orderNr string) *licenseActivateAction {
	return &licenseActivateAction{
		bus:        bus,
		debug:      false,
		kubeconfig: kubecfg,
		orderNr:    orderNr,
	}
}

var _ Action = (*licenseActivateAction)(nil)

type licenseActivateAction struct {
	orderNr    string
	kubeconfig string
	debug      bool

	bus eventbus.Bus
	rc  *rest.Config
}

func (o *licenseActivateAction) Run() (err error) {
	// Use the current context in kubeconfig
	o.rc, err = clientcmd.BuildConfigFromFlags("", o.kubeconfig)
	if err != nil {
		return err
	}

	// Derive the Cluster Identifier
	o.bus.Publish(events.NewStartWaitEvent("Computing cluster identifier"))
	time.Sleep(time.Second * 2)

	clusterId, err := getClusterId(o.rc)
	if err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent("Cluster identifier successfully computed"))

	o.bus.Publish(events.NewStartWaitEvent("Synching information"))
	time.Sleep(time.Second * 2)

	key, err := o.sendInfo(o.orderNr, clusterId)
	if err != nil {
		return err
	}
	o.bus.Publish(events.NewDoneEvent("Synch completed"))

	o.bus.Publish(events.NewStartWaitEvent("Storing license data"))
	time.Sleep(time.Second * 2)

	src, err := tmpl.Execute("license-secret.yaml", map[string]string{
		"License": key,
	})
	if err != nil {
		return err
	}

	if err := kubeclient.Apply(context.TODO(), o.rc, src); err != nil {
		if merr, ok := err.(*multierror.Error); ok {
			return multierror.Flatten(merr)
		}
		return err
	}
	o.bus.Publish(events.NewDoneEvent("License data successfully stored"))

	return nil
}

func (o *licenseActivateAction) sendInfo(orderNr, clusterId string) (string, error) {
	data, err := json.Marshal(map[string]string{
		"orderNumber": orderNr,
		"clusterId":   clusterId,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, activationURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	// debug(httputil.DumpRequestOut(req, true))

	client := &http.Client{
		Timeout: time.Second * 40,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// debug(httputil.DumpResponse(res, true))

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	msg := struct {
		Status int    `json:"status"`
		Data   string `json:"data,omitempty"`
	}{}
	if err := json.Unmarshal(body, &msg); err != nil {
		return "", err
	}

	if msg.Status != 201 {
		return "", fmt.Errorf(msg.Data)
	}
	return msg.Data, nil
}

/*
func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n\n", data)
	} else {
		log.Fatalf("%s\n\n", err)
	}
}
*/
