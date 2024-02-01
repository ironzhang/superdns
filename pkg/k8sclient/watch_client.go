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

// Action event action
type Action string

const (
	Add    Action = "Add"
	Update Action = "Update"
	Delete Action = "Delete"
)

// Event watch event
type Event struct {
	Action Action
	Key    string
	Object interface{}
}

// Watcher watcher interface
type Watcher interface {
	OnWatch(indexer cache.Indexer, event Event) error
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
	lselector labels.Selector, fselector fields.Selector, indexers cache.Indexers, watcher Watcher) cache.Indexer {
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
			queue.Add(&Event{
				Action: Add,
				Key:    key,
				Object: obj,
			})
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			if err != nil {
				tlog.Errorw("get meta namespace key", "error", err)
				return
			}
			queue.Add(&Event{
				Action: Update,
				Key:    key,
				Object: newObj,
			})
		},
		DeleteFunc: func(obj interface{}) {
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err != nil {
				tlog.Errorw("get deletion handling meta namespace key", "error", err)
				return
			}
			queue.Add(&Event{
				Action: Delete,
				Key:    key,
				Object: obj,
			})
		},
	}

	// new informer
	lw := cache.NewFilteredListWatchFromClient(p.rest, resource, namespace, func(options *metav1.ListOptions) {
		options.LabelSelector = lselector.String()
		options.FieldSelector = fselector.String()
	})
	indexer, controller := cache.NewIndexerInformer(lw, object, 0, h, indexers)

	// run worker
	w := worker{
		watcher:    watcher,
		queue:      queue,
		indexer:    indexer,
		controller: controller,
	}
	w.Run(ctx)

	return indexer
}
