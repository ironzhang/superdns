package controller

import (
	"github.com/ironzhang/tlog"

	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
)

type clusterWatcher struct {
}

func (p *clusterWatcher) OnWatch(indexer cache.Indexer, event k8sclient.Event) error {
	clusters := make(map[string]superdnsv1.ClusterSpec, 0)

	objects := indexer.List()
	for _, obj := range objects {
		cluster, ok := obj.(*superdnsv1.Cluster)
		if !ok {
			tlog.Errorw("object is not a cluster", "obj", obj)
			return nil
		}
		clusters[cluster.Spec.Cluster] = cluster.Spec
	}
	tlog.Infow("on watch cluster", "clusters", clusters)

	return nil
}
