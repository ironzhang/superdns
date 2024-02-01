package subscribe

import (
	"context"
	"time"

	"github.com/ironzhang/superdns/superdns-agent/internal/controller"
	"github.com/ironzhang/superdns/superdns-agent/internal/ready"
	"github.com/ironzhang/tlog"
)

// A Subscriber subscribes domains.
type Subscriber struct {
	notary     *notary
	controller *controller.Controller
	inspection *ready.Inspection
}

// NewSubscriber returns an instance of Subscriber.
func NewSubscriber(c *controller.Controller, inspection *ready.Inspection) *Subscriber {
	return &Subscriber{
		notary:     newNotary(),
		controller: c,
		inspection: inspection,
	}
}

// SubscribeDomain subscribes domain in ttl duration.
//
// if ttl <= 0, means forever.
func (p *Subscriber) SubscribeDomain(ctx context.Context, domain string, ttl time.Duration) (bool, error) {
	do := func(ctx context.Context) (context.CancelFunc, error) {
		logger := tlog.WithContext(ctx).WithArgs("domain", domain)
		ctx, cancel := context.WithCancel(ctx)

		logger.Infow("watch domain")
		err := p.controller.WatchDomain(ctx, domain)
		if err != nil {
			logger.Errorw("controller watch domain", "error", err)
			return func() {}, err
		}

		return func() {
			logger.Infow("cancel watch domain")
			cancel()
		}, nil
	}

	newSub, err := p.notary.NewLease(ctx, domain, ttl, do)
	if err != nil {
		return false, err
	}

	if !p.inspection.ServiceReady(domain) {
		p.controller.RefreshClusters(ctx, domain)
	}
	if !p.inspection.RouteReady(domain) {
		p.controller.RefreshRoutes(ctx, domain)
	}

	return newSub, nil
}

// ListSubscribeDomains returns the domains which are subscribing.
func (p *Subscriber) ListSubscribeDomains(ctx context.Context) []string {
	return p.notary.ListLeaseKeys(ctx)
}
