package subscribe

import (
	"context"
	"sync"
	"time"

	"github.com/ironzhang/superlib/ctxutil"
)

var now = time.Now

type lease struct {
	key     string
	cancel  context.CancelFunc
	forever bool
	expire  time.Time
}

func newLease(key string, cancel context.CancelFunc, ttl time.Duration) *lease {
	les := &lease{
		key:    key,
		cancel: cancel,
	}
	les.KeepAlive(ttl)
	return les
}

func (p *lease) KeepAlive(ttl time.Duration) {
	if ttl <= 0 {
		p.forever = true
		return
	}

	p.forever = false
	p.expire = now().Add(ttl)
}

func (p *lease) Expired() bool {
	if p.forever {
		return false
	}
	return now().After(p.expire)
}

func (p *lease) Revoke() {
	p.cancel()
}

type doFunc func(ctx context.Context) (context.CancelFunc, error)

type notary struct {
	mu     sync.Mutex
	leases map[string]*lease
}

func newNotary() *notary {
	return new(notary).init()
}

func (p *notary) init() *notary {
	p.leases = make(map[string]*lease)
	go p.running()
	return p
}

func (p *notary) running() {
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-t.C:
			p.update()
		}
	}
}

func (p *notary) update() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for key, les := range p.leases {
		if les.Expired() {
			les.Revoke()
			delete(p.leases, key)
		}
	}
}

func (p *notary) NewLease(ctx context.Context, key string, ttl time.Duration, do doFunc) (bool, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	les, ok := p.leases[key]
	if ok {
		les.KeepAlive(ttl)
		return false, nil
	}

	cancel, err := do(ctxutil.CloneContext(ctx))
	if err != nil {
		return false, err
	}

	les = newLease(key, cancel, ttl)
	p.leases[key] = les
	return true, nil
}

func (p *notary) ListLeaseKeys(ctx context.Context) []string {
	p.mu.Lock()
	defer p.mu.Unlock()

	keys := make([]string, 0, len(p.leases))
	for key := range p.leases {
		keys = append(keys, key)
	}
	return keys
}
