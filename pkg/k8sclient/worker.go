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
	indexer    cache.Indexer
	controller cache.Controller
}

func (p *worker) HandleErr(err error, item interface{}) {
	if err == nil {
		p.queue.Forget(item)
		return
	}

	// retries 3 times if something goes wrong
	retries := p.queue.NumRequeues(item)
	if retries < 3 {
		tlog.Infow("retrying", "item", item, "retries", retries, "error", err)
		p.queue.AddRateLimited(item)
		return
	}

	// stop retrying
	p.queue.Forget(item)
	runtime.HandleError(err)
	tlog.Infow("dropping", "item", item, "error", err)
}

func (p *worker) Process() bool {
	item, quit := p.queue.Get()
	if quit {
		return false
	}
	defer p.queue.Done(item)

	err := p.watcher.OnWatch(p.indexer, *(item.(*Event)))
	p.HandleErr(err, item)
	return true
}

func (p *worker) RunWorker() {
	for p.Process() {
	}
}

func (p *worker) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		p.queue.ShutDown()
	}()

	go p.controller.Run(ctx.Done())

	if !cache.WaitForCacheSync(ctx.Done(), p.controller.HasSynced) {
		tlog.WithContext(ctx).Errorw("timed out waiting for caches to sync")
		return
	}

	go wait.Until(p.RunWorker, 500*time.Millisecond, ctx.Done())
}
