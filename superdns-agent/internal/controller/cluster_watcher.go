package controller

import (
	"github.com/ironzhang/tlog"

	"k8s.io/client-go/tools/cache"

	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
)

type clusterWatcher struct {
}

func (p *clusterWatcher) OnWatch(store cache.Store, key string) error {
	obj, ok, err := store.GetByKey(key)
	if err != nil {
		tlog.Errorw("store get", "key", key, "error", err)
		return err
	}
	if !ok {
		tlog.Infow("cluster does not exist", "key", key)
		return nil
	}

	cluster, ok := obj.(*superdnsv1.Cluster)
	if !ok {
		tlog.Errorw("object is not a cluster", "key", key)
		return nil
	}

	tlog.Infow("on watch cluster", "name", cluster.GetName(), "spec", cluster.Spec)

	return nil
}
