package controller

import (
	"encoding/json"

	"github.com/ironzhang/superlib/superutil/supermodel"
	"github.com/ironzhang/tlog"

	"k8s.io/client-go/tools/cache"

	"github.com/ironzhang/superdns/pkg/filewrite"
	"github.com/ironzhang/superdns/pkg/k8sclient"
	"github.com/ironzhang/superdns/pkg/superconv"
	superdnsv1 "github.com/ironzhang/superdns/supercrd/apis/superdns.io/v1"
	"github.com/ironzhang/superdns/superdns-agent/internal/paths"
)

type routeWatcher struct {
	pathmgr *paths.PathManager
	fwriter *filewrite.FileWriter
}

func (p *routeWatcher) OnWatch(indexer cache.Indexer, event k8sclient.Event) error {
	tlog.Debugw("on watch", "event", event)

	r, ok := event.Object.(*superdnsv1.Route)
	if !ok {
		tlog.Errorw("object is not a route", "object", event.Object)
		return nil
	}

	model := supermodel.RouteModel{
		Domain:   r.ObjectMeta.Name,
		Strategy: superconv.ToSupermodelRoute(*r),
	}
	err := p.writeModel(model)
	if err != nil {
		tlog.Errorw("write route model", "model", model, "error", err)
		return err
	}
	err = p.writeScript(model.Domain, r.Spec.ScriptContent)
	if err != nil {
		tlog.Errorw("write route script", "domain", model.Domain, "error", err)
		return err
	}

	return nil
}

func (p *routeWatcher) writeModel(m supermodel.RouteModel) error {
	data, err := json.MarshalIndent(m, "", "\t")
	if err != nil {
		return err
	}

	path := p.pathmgr.RouteModelPath(m.Domain)
	if err = p.fwriter.WriteFile(path, data); err != nil {
		return err
	}
	return nil
}

func (p *routeWatcher) writeScript(domain, content string) error {
	path := p.pathmgr.RouteScriptPath(domain)
	if err := p.fwriter.WriteFile(path, []byte(content)); err != nil {
		return err
	}
	return nil
}
