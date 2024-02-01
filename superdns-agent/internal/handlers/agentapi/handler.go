package agentapi

import (
	"context"
	"time"

	"github.com/ironzhang/superlib/httputils/echoutil/echorpc"
	"github.com/ironzhang/tlog"
	"github.com/labstack/echo"

	"github.com/ironzhang/superdns/pkg/protocol"
	"github.com/ironzhang/superdns/superdns-agent/internal/agent"
)

// Handler is an echo rpc handler to handle agent api.
type Handler struct {
	agent *agent.Agent
}

// Register registers agent api.
func Register(e *echo.Echo, h *Handler) {
	e.POST("/superdns/agent/v1/api/subscribe/domains", echorpc.HandlerFunc(h.SubscribeDomains))
	e.GET("/superdns/agent/v1/api/list/subscribe/domains", echorpc.HandlerFunc(h.ListSubscribeDomains))
}

// NewHandler returns an instance of Handler.
func NewHandler(a *agent.Agent) *Handler {
	return &Handler{agent: a}
}

// SubscribeDomains handles the subscribe domains request.
func (p *Handler) SubscribeDomains(ctx context.Context, req *protocol.SubscribeDomainsReq, resp interface{}) error {
	fn := func() error {
		err := p.agent.SubscribeDomains(ctx, req.Domains, time.Duration(req.TTL))
		if err != nil {
			tlog.WithContext(ctx).Errorw("subscribe domains", "domains", req.Domains, "error", err)
			return err
		}
		return nil
	}

	if req.Asynchronous {
		go fn()
		return nil
	}
	return fn()
}

// ListSubscribeDomains handles the list subscribe domains request.
func (p *Handler) ListSubscribeDomains(ctx context.Context, req interface{}, resp *protocol.ListSubscribeDomainsResp) error {
	resp.Domains = p.agent.ListSubscribeDomains(ctx)
	return nil
}
