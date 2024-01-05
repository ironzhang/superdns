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

func newDomainAndLidcSelector(domain, lidc string) (labels.Selector, error) {
	r1, err := labels.NewRequirement("domain", selection.Equals, []string{domain})
	if err != nil {
		return nil, err
	}
	r2, err := labels.NewRequirement("lidc", selection.Equals, []string{lidc})
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().Add(*r1, *r2), nil
}
