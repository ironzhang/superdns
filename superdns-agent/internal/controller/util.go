package controller

import (
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func newDomainSelector(domain string) (labels.Selector, error) {
	r, err := labels.NewRequirement("domain", selection.Equals, []string{domain})
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().Add(*r), nil
}
