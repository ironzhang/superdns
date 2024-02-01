package controller

import (
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

func newDomainLabelSelector(domain string) (labels.Selector, error) {
	r, err := labels.NewRequirement("domain", selection.Equals, []string{domain})
	if err != nil {
		return nil, err
	}
	return labels.NewSelector().Add(*r), nil
}

func newDomainFieldSelector(domain string) (fields.Selector, error) {
	set := fields.Set{
		"metadata.name": domain,
	}
	return set.AsSelector(), nil
}

func clusterIndexerKey(domain string) string {
	return fmt.Sprintf("cluster:%s", domain)
}

func routeIndexerKey(domain string) string {
	return fmt.Sprintf("route:%s", domain)
}
