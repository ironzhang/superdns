package controller

import (
	"context"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
	"k8s.io/apimachinery/pkg/fields"
)

type Controller struct {
	namespace string
	wc        *k8sclient.WatchClient
	cw        clusterWatcher
}

func New(ns string, wc *k8sclient.WatchClient) *Controller {
	return &Controller{
		namespace: ns,
		wc:        wc,
	}
}

func (p *Controller) WatchClusters(ctx context.Context, domain string) error {
	//ls := labels.Everything()
	ls, err := newDomainSelector(domain)
	if err != nil {
		return err
	}

	p.wc.Watch(ctx, p.namespace, "clusters", &superdnsv1.Cluster{}, ls, fields.Everything(), &p.cw)

	return nil
}
