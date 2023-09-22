package watchsystem

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superdns/superdns-agent/core/filesystem"
	"github.com/ironzhang/superdns/superdns-agent/core/storesystem/pod"
)

type podwatcher struct {
	app        string
	store      *pod.Store
	filesystem *filesystem.System
}

func (p *podwatcher) OnWatch(store cache.Store, key string) error {
	app := pod.NewApplication(p.app)

	objects := store.List()
	for _, obj := range objects {
		pod := obj.(*v1.Pod)
		if pod.Status.PodIP == "" {
			continue
		}
		app.AddPod(pod)
	}

	p.store.SetApplication(app)
	p.filesystem.UpdateServiceModel(p.app)

	return nil
}
