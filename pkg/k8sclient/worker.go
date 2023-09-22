package k8sclient

import (
	"context"
	"time"

	"github.com/ironzhang/tlog"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type worker struct {
	watcher    Watcher
	queue      workqueue.RateLimitingInterface
	store      cache.Store
	controller cache.Controller
}

func (p *worker) HandleErr(err error, key interface{}) {
	if err == nil {
		p.queue.Forget(key)
		return
	}

	// retries 3 times if something goes wrong
	retries := p.queue.NumRequeues(key)
	if retries < 3 {
		tlog.Infow("retrying", "key", key, "retries", retries, "error", err)
		p.queue.AddRateLimited(key)
		return
	}

	// stop retrying
	p.queue.Forget(key)
	runtime.HandleError(err)
	tlog.Infow("dropping", "key", key, "error", err)
}

func (p *worker) Process() bool {
	key, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(key)

	err := p.watcher.OnWatch(p.store, key.(string))
	p.HandleErr(err, key)
	return true
}

func (p *worker) RunWorker() {
	for p.Process() {
	}
}

func (p *worker) Run(ctx context.Context) {
	defer p.queue.ShutDown()

	go p.controller.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), p.controller.HasSynced) {
		tlog.WithContext(ctx).Errorw("timed out waiting for caches to sync")
		return
	}

	go wait.Until(p.RunWorker, 500*time.Millisecond, ctx.Done())

	<-ctx.Done()
}
