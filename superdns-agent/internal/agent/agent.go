package agent

import (
	"context"
	"time"

	"github.com/ironzhang/superdns/superdns-agent/internal/subscribe"
)

// Agent represents an instance of an agent, which providers the superdns agent's functions.
type Agent struct {
	subscriber *subscribe.Subscriber
}

// New returns an instance of Agent.
func New(s *subscribe.Subscriber) *Agent {
	return &Agent{
		subscriber: s,
	}
}

// SubscribeDomains subscribe the given domains.
func (p *Agent) SubscribeDomains(ctx context.Context, domains []string, ttl time.Duration) error {
	for _, domain := range domains {
		_, err := p.subscriber.SubscribeDomain(ctx, domain, ttl)
		if err != nil {
			return err
		}
	}
	return nil
}

// ListSubscribeDomains returns the domains which are subscribing.
func (p *Agent) ListSubscribeDomains(ctx context.Context) []string {
	return p.subscriber.ListSubscribeDomains(ctx)
}
