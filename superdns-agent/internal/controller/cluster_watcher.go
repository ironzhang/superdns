package controller

import (
	"encoding/json"
	"sort"

	"github.com/ironzhang/tlog"

	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superlib/superutil/supermodel"

	"github.com/ironzhang/superdns/pkg/filewrite"
	"github.com/ironzhang/superdns/pkg/k8sclient"
	"github.com/ironzhang/superdns/pkg/superconv"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
	"github.com/ironzhang/superdns/superdns-agent/internal/paths"
)

type clusterWatcher struct {
	pathmgr *paths.PathManager
	fwriter *filewrite.FileWriter
}

func (p *clusterWatcher) OnWatch(indexer cache.Indexer, event k8sclient.Event) error {
	tlog.Debugw("on watch", "event", event)

	c, ok := event.Object.(*superdnsv1.Cluster)
	if !ok {
		tlog.Errorw("object is not a cluster", "object", event.Object)
		return nil
	}
	return p.refresh(indexer, c)
}

func (p *clusterWatcher) OnRefresh(indexer cache.Indexer) {
	for _, obj := range indexer.List() {
		c, ok := obj.(*superdnsv1.Cluster)
		if !ok {
			tlog.Errorw("object is not a cluster", "object", obj)
			continue
		}

		tlog.Debugw("on refresh", "domain", c.ObjectMeta.Labels["domain"])
		p.refresh(indexer, c)
		return
	}
}

func (p *clusterWatcher) refresh(indexer cache.Indexer, c *superdnsv1.Cluster) error {
	model := supermodel.ServiceModel{
		Domain:   c.ObjectMeta.Labels["domain"],
		Clusters: objectsToClusters(indexer.List()),
	}
	err := p.writeModel(model)
	if err != nil {
		tlog.Errorw("write service model", "model", model, "error", err)
		return err
	}
	return nil
}

func (p *clusterWatcher) writeModel(m supermodel.ServiceModel) error {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}

	path := p.pathmgr.ServiceModelPath(m.Domain)
	if err = p.fwriter.WriteFile(path, data); err != nil {
		return err
	}
	return nil
}

func objectsToClusters(objects []interface{}) []supermodel.Cluster {
	clusters := make([]supermodel.Cluster, 0, len(objects))
	for _, obj := range objects {
		c, ok := obj.(*superdnsv1.Cluster)
		if !ok {
			tlog.Errorw("object is not a cluster", "obj", obj)
			continue
		}
		clusters = append(clusters, superconv.ToSupermodelCluster(*c))
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name < clusters[j].Name
	})

	return clusters
}
