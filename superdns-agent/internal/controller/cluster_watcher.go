package controller

import (
	"github.com/ironzhang/tlog"

	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superlib/superutil/supermodel"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	"github.com/ironzhang/superdns/pkg/superconv"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
)

type clusterWatcher struct {
}

func (p *clusterWatcher) OnWatch(indexer cache.Indexer, event k8sclient.Event) error {
	c, ok := event.Object.(*superdnsv1.Cluster)
	if !ok {
		tlog.Errorw("object is not a cluster", "object", event.Object)
		return nil
	}

	model := supermodel.ServiceModel{
		Domain:   c.Spec.Domain,
		Clusters: objectsToClusters(indexer.List()),
	}
	tlog.Infow("on watch", "model", model)

	return nil
}

func objectsToClusters(objects []interface{}) map[string]supermodel.Cluster {
	clusters := make(map[string]supermodel.Cluster, len(objects))
	for _, obj := range objects {
		c, ok := obj.(*superdnsv1.Cluster)
		if !ok {
			tlog.Errorw("object is not a cluster", "obj", obj)
			continue
		}
		clusters[c.Spec.Cluster] = superconv.ToSupermodelCluster(*c)
	}
	return clusters
}
