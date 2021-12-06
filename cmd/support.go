package cmd

import (
	"github.com/krateoplatformops/krateoctl/pkg/eventbus"
	"github.com/krateoplatformops/krateoctl/pkg/events"
	"github.com/krateoplatformops/krateoctl/pkg/log"
)

func updateLog(l log.Logger) eventbus.EventHandler {
	return func(e eventbus.Event) {
		switch e.EventID() {
		case events.DebugEventID:
			evt := e.(*events.DebugEvent)
			l.Debug(evt.Message())

		case events.StartWaitEventID:
			evt := e.(*events.StartWaitEvent)
			l.StartWait(evt.Message())

		case events.StopWaitEventID:
			l.StopWait()

		case events.DoneEventID:
			evt := e.(*events.DoneEvent)
			l.StopWait()
			l.Done(evt.Message())
		}
	}
}

/*
func createDockerSecretManifest(name, namespace, user, pass string) ([]byte, error) {
	creds, err := encodeDockerConfig(user, pass)
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"Name":        name,
		"Namespace":   namespace,
		"Credentials": creds,
	}

	src, err := tmpl.Execute("docker-secret.yaml", data)
	if err != nil {
		return nil, err
	}

	return src, err
}

func encodeDockerConfig(user, pass string) (string, error) {
	data := map[string]string{
		"User": user,
		"Pass": pass,
	}

	auth := fmt.Sprintf("%s:%s", data["User"], data["Pass"])
	data["Auth"] = base64.StdEncoding.EncodeToString([]byte(auth))

	src, err := tmpl.Execute("dockerconfig.json", data)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(src), nil
}

// listPodsWithLabels returns the list of currently
// scheduled or running pods in `namespace` with the given labels
func listPodsWithLabels(kc client.Client, namespace string, tag ...string) (*v1.PodList, error) {
	listOpts := []client.ListOption{
		client.HasLabels(tag),
		//"pkg.crossplane.io/revision",
	}

	list := &v1.PodList{}
	if err := kc.List(context.TODO(), list, listOpts...); err != nil {
		return nil, errors.Wrapf(err, "cannot get pod list with tags: %v", tag)
	}

	return list, nil
}
*/
