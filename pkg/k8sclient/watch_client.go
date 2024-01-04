package k8sclient

import (
	"context"

	"github.com/ironzhang/tlog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

// Watcher watcher interface
type Watcher interface {
	OnWatch(store cache.Store, key string) error
}

// WatchClient k8s watch client
type WatchClient struct {
	rest rest.Interface
}

// NewWatchClient new watch client
func NewWatchClient(rest rest.Interface) *WatchClient {
	return &WatchClient{rest: rest}
}

// Watch watch k8s resource
func (p *WatchClient) Watch(ctx context.Context, namespace, resource string, object runtime.Object,
	lselector labels.Selector, fselector fields.Selector, watcher Watcher) {
	// new workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// new resource event handler
	h := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err != nil {
				tlog.Errorw("get meta namespace key", "error", err)
				return
			}
			queue.Add(key)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				tlog.Errorw("get meta namespace key", "error", err)
				return
			}
			queue.Add(key)
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				tlog.Errorw("get deletion handling meta namespace key", "error", err)
				return
			}
			queue.Add(key)
		},
	}

	// new informer
	lw := cache.NewFilteredListWatchFromClient(p.rest, resource, namespace, func(options *metav1.ListOptions) {
		options.LabelSelector = lselector.String()
		options.FieldSelector = fselector.String()
	})
	store, controller := cache.NewInformer(lw, object, 0, h)

	// run worker
	w := worker{
		watcher:    watcher,
		queue:      queue,
		store:      store,
		controller: controller,
	}
	go w.Run(ctx)
}
