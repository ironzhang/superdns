package agentapi

import (
	"context"

	"github.com/ironzhang/superlib/httputils/echoutil/echorpc"
	"github.com/ironzhang/tlog"
	"github.com/labstack/echo"

	"github.com/ironzhang/superdns/pkg/protocol"
	"github.com/ironzhang/superdns/superdns-agent/internal/controller"
)

type Handler struct {
	controller *controller.Controller
}

func Register(e *echo.Echo, h *Handler) {
	e.POST("/superdns/agent/v1/api/watch/domains", echorpc.HandlerFunc(h.WatchDomains))
}

func NewHandler(c *controller.Controller) *Handler {
	return &Handler{controller: c}
}

func (p *Handler) WatchDomains(ctx context.Context, req *protocol.WatchDomainsReq, resp interface{}) error {
	tlog.WithContext(ctx).Debugw("start watching domains", "domains", req.Domains)
	err := p.controller.WatchDomains(ctx, req.Domains)
	if err != nil {
		tlog.WithContext(ctx).Errorw("watch domains", "domains", req.Domains, "error", err)
		return err
	}

	<-ctx.Done()
	tlog.WithContext(ctx).Debugw("canceled watching domains", "domains", req.Domains)
	return nil
}
