package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
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
	ls, err := newDomainSelector(domain)
	if err != nil {
		return err
	}

	p.wc.Watch(ctx, p.namespace, "clusters", &superdnsv1.Cluster{}, ls, fields.Everything(), cache.Indexers{}, &p.cw)

	return nil
}
