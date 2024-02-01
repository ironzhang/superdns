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

// Options Controller options.
type Options struct {
	Namespace string
}

// A Controller watches domains and dumps the domains' data to files.
type Controller struct {
	opts     Options
	wc       *k8sclient.WatchClient
	cw       clusterWatcher
	rw       routeWatcher
	indexers map[string]cache.Indexer
}

// New returns an instance of Controller.
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
		indexers: make(map[string]cache.Indexer),
	}
}

func (p *Controller) watchClusters(ctx context.Context, domain string) error {
	ls, err := newDomainLabelSelector(domain)
	if err != nil {
		return err
	}

	key := clusterIndexerKey(domain)
	p.indexers[key] = p.wc.Watch(ctx, p.opts.Namespace, "clusters", &superdnsv1.Cluster{}, ls, fields.Everything(), cache.Indexers{}, &p.cw)

	return nil
}

func (p *Controller) watchRoutes(ctx context.Context, domain string) error {
	fs, err := newDomainFieldSelector(domain)
	if err != nil {
		return err
	}

	key := routeIndexerKey(domain)
	p.indexers[key] = p.wc.Watch(ctx, p.opts.Namespace, "routes", &superdnsv1.Route{}, labels.Everything(), fs, cache.Indexers{}, &p.rw)

	return nil
}

// WatchDomain watches the given domain.
func (p *Controller) WatchDomain(ctx context.Context, domain string) (err error) {
	if err = p.watchClusters(ctx, domain); err != nil {
		return err
	}
	if err = p.watchRoutes(ctx, domain); err != nil {
		return err
	}
	return nil
}

// RefreshClusters refresh the given domain's cluster file.
func (p *Controller) RefreshClusters(ctx context.Context, domain string) {
	key := clusterIndexerKey(domain)
	indexer, ok := p.indexers[key]
	if ok {
		p.cw.OnRefresh(indexer)
	}
}

// RefreshRoutes refresh the given domain's route file.
func (p *Controller) RefreshRoutes(ctx context.Context, domain string) {
	key := routeIndexerKey(domain)
	indexer, ok := p.indexers[key]
	if ok {
		p.rw.OnRefresh(indexer)
	}
}
