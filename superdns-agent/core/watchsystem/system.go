package watchsystem

import (
	"context"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	"github.com/ironzhang/superdns/superdns-agent/core/filesystem"
	"github.com/ironzhang/superdns/superdns-agent/core/storesystem"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
)

type System struct {
	watchclient *k8sclient.WatchClient
	storesystem *storesystem.System
	filesystem  *filesystem.System
}

func NewSystem(wc *k8sclient.WatchClient, ss *storesystem.System, fs *filesystem.System) *System {
	return &System{
		watchclient: wc,
		storesystem: ss,
		filesystem:  fs,
	}
}

func (p *System) Watch(ctx context.Context, app string) error {
	r, err := labels.NewRequirement("app", selection.Equals, []string{app})
	if err != nil {
		return err
	}
	s := labels.NewSelector()
	s.Add(*r)

	pw := podwatcher{
		app:        app,
		store:      p.storesystem.PodStore,
		filesystem: p.filesystem,
	}
	p.watchclient.Watch(ctx, "dev", "pods", &v1.Pod{}, s, fields.Everything(), &pw)
	return nil
}
