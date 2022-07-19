package core

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func InstalledBySelector() (labels.Selector, error) {
	req, err := labels.NewRequirement(InstalledByLabel,
		selection.Equals, []string{InstalledByValue})
	if err != nil {
		return nil, err
	}

	sel := labels.NewSelector()
	sel = sel.Add(*req)

	return sel, nil
}
