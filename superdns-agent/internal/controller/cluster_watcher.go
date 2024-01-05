package controller

import (
	"encoding/json"

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
	c, ok := event.Object.(*superdnsv1.Cluster)
	if !ok {
		tlog.Errorw("object is not a cluster", "object", event.Object)
		return nil
	}

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
