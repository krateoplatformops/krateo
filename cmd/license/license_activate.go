package license

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/krateoplatformops/krateo/cmd/tools"
	"github.com/krateoplatformops/krateo/cmd/tools/flags"
	"github.com/krateoplatformops/krateo/pkg/eventbus"
	"github.com/krateoplatformops/krateo/pkg/events"
	"github.com/krateoplatformops/krateo/pkg/log"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultLicenseServerURL = "https://license.krateo.io/activate"
	//activationURL = "http://localhost:8090/activate"
)

type licenseActivateOptions struct {
	bus        eventbus.Bus
	kubeconfig string
	verbose    bool

	orderNr   string
	serverUrl string
}

func NewLicenseActivateCmd() *cobra.Command {
	o := &licenseActivateOptions{}

	cmd := &cobra.Command{
		Use:                   "activate <ORDER NUMBER>",
		DisableSuggestions:    true,
		DisableFlagsInUseLine: true,
		Args:                  cobra.ArbitraryArgs,
		Short:                 "Activate the Krateo license",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("order number of the license purchase is mandatory")
			}

			l := log.GetInstance()
			handler := tools.LogEventHandler(l)

			o.bus = eventbus.New()
			eids := []eventbus.Subscription{
				o.bus.Subscribe(events.StartWaitEventID, handler),
				o.bus.Subscribe(events.StopWaitEventID, handler),
				o.bus.Subscribe(events.DoneEventID, handler),
			}
			defer func() {
				for _, e := range eids {
					o.bus.Unsubscribe(e)
				}
			}()

			o.orderNr = strings.TrimSpace(args[0])

			return o.Run()
		},
	}

	defaultKubeconfig := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	// if KUBECONFIG is empty
	if len(defaultKubeconfig) == 0 {
		// look for file $HOME/.kube/config
		defaultKubeconfig = clientcmd.RecommendedHomeFile
	}

	cmd.Flags().BoolVarP(&o.verbose, flags.Verbose, "v", false, "dump verbose output")
	cmd.Flags().StringVarP(&o.kubeconfig, clientcmd.RecommendedConfigPathFlag, "k",
		defaultKubeconfig, "absolute path to the kubeconfig file")
	cmd.Flags().StringVarP(&o.serverUrl, flags.LicenseServerURL, "u",
		defaultLicenseServerURL, "Krateo license verification server URL")
	//nolint:errcheck
	cmd.Flags().MarkHidden(flags.LicenseServerURL)

	return cmd
}

func (o *licenseActivateOptions) Run() error {
	rc, err := clientcmd.BuildConfigFromFlags("", o.kubeconfig)
	if err != nil {
		return err
	}

	// Derive the Cluster Identifier
	o.bus.Publish(events.NewStartWaitEvent("Computing cluster identifier"))
	time.Sleep(time.Second * 2)

	clusterId, err := getClusterId(rc)
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

	err = createLicenseSecret(rc, key)
	if err != nil {
		return err
	}

	o.bus.Publish(events.NewDoneEvent("License data successfully stored"))

	return nil
}

func (o *licenseActivateOptions) sendInfo(orderNr, clusterId string) (string, error) {
	data, err := json.Marshal(map[string]string{
		"orderNumber": orderNr,
		"clusterId":   clusterId,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, o.serverUrl, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	client := &http.Client{
		Timeout: time.Second * 40,
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

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
