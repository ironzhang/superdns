package ready

import (
	"github.com/ironzhang/superlib/fileutil"

	"github.com/ironzhang/superdns/superdns-agent/internal/paths"
)

type Inspection struct {
	pathmgr *paths.PathManager
}

func NewInspection(pm *paths.PathManager) *Inspection {
	return &Inspection{pathmgr: pm}
}

func (p *Inspection) ServiceReady(domain string) bool {
	return fileutil.FileExist(p.pathmgr.ServiceModelPath(domain))
}

func (p *Inspection) RouteReady(domain string) bool {
	if !fileutil.FileExist(p.pathmgr.RouteModelPath(domain)) {
		return false
	}
	if !fileutil.FileExist(p.pathmgr.RouteScriptPath(domain)) {
		return false
	}
	return true
}
