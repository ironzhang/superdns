package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superdns/pkg/filewrite"
	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
	"github.com/ironzhang/superdns/superdns-agent/internal/paths"
)

type Options struct {
	Namespace string
	LIDC      string
}

type Controller struct {
	opts Options
	wc   *k8sclient.WatchClient
	cw   clusterWatcher
	rw   routeWatcher
}

func New(opts Options, wc *k8sclient.WatchClient, pm *paths.PathManager, fw *filewrite.FileWriter) *Controller {
	return &Controller{
		opts: opts,
		wc:   wc,
		cw: clusterWatcher{
			pathmgr: pm,
			fwriter: fw,
		},
		rw: routeWatcher{
			pathmgr: pm,
			fwriter: fw,
		},
	}
}

func (p *Controller) WatchClusters(ctx context.Context, domain string) error {
	ls, err := newDomainLabelSelector(domain)
	if err != nil {
		return err
	}

	p.wc.Watch(ctx, p.opts.Namespace, "clusters", &superdnsv1.Cluster{}, ls, fields.Everything(), cache.Indexers{}, &p.cw)

	return nil
}

func (p *Controller) WatchRoutes(ctx context.Context, domain string) error {
	fs, err := newDomainFieldSelector(domain)
	if err != nil {
		return err
	}

	p.wc.Watch(ctx, p.opts.Namespace, "routes", &superdnsv1.Route{}, labels.Everything(), fs, cache.Indexers{}, &p.rw)

	return nil
}

func (p *Controller) WatchDomain(ctx context.Context, domain string) (err error) {
	if err = p.WatchClusters(ctx, domain); err != nil {
		return err
	}
	if err = p.WatchRoutes(ctx, domain); err != nil {
		return err
	}
	return nil
}

func (p *Controller) WatchDomains(ctx context.Context, domains []string) (err error) {
	for _, domain := range domains {
		err = p.WatchDomain(ctx, domain)
		if err != nil {
			return err
		}
	}
	return nil
}
