package uninstall

import (
	"strings"
	"time"

	"github.com/krateoplatformops/krateo/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func listCRDs(cli kubernetes.CrdsClient) ([]string, error) {
	res, err := cli.List(metav1.ListOptions{})
	if err != nil {
		return []string{}, err
	}

	names := []string{}

	for _, el := range res.Items {
		accept := strings.Contains(el.Name, "krateo")
		accept = accept || strings.Contains(el.Name, "crossplane")
		if accept {
			names = append(names, el.Name)
		}
	}

	return names, nil
}

func patchAndDeleteCRDs(cli kubernetes.CrdsClient, dryRun bool, names []string) error {
	patchData := []byte(`{"metadata":{"finalizers":[]}}`)

	for _, el := range names {
		_, err := cli.Get(el, metav1.GetOptions{})
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return err
		}

		po := metav1.PatchOptions{}
		if dryRun {
			po.DryRun = []string{metav1.DryRunAll}
		}

		_, err = cli.Patch(el, types.MergePatchType, patchData, po)
		if err != nil {
			return err
		}

		//var dat bytes.Buffer
		//marshallCrd(&dat, crd, "yaml")
		//fmt.Printf("\n ==> AFTER\n%s\n", dat.String())

		do := metav1.DeleteOptions{}
		if dryRun {
			do.DryRun = []string{metav1.DryRunAll}
		}
		err = cli.Delete(el, do)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return err
		}

		time.Sleep(time.Second * 1)
	}

	return nil
}
