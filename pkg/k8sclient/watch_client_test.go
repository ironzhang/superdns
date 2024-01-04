package k8sclient

import (
	"context"
	"errors"
	"testing"

	"github.com/ironzhang/tlog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func PrintObjects(objects []interface{}) {
	tlog.Infof("---------- start ----------")
	for _, obj := range objects {
		pod := obj.(*v1.Pod)
		tlog.Infof("pod %s %s", pod.GetName(), pod.Status.Phase)
	}
	tlog.Infof("---------- stop -----------")
}

type testPodWatcher struct {
}

func (p *testPodWatcher) OnWatch(indexer cache.Indexer, event Event) error {
	objects, err := indexer.Index("app_index", event.Object)
	if err != nil {
		tlog.Errorw("index", "obj", event.Object, "error", err)
		return err
	}
	PrintObjects(objects)
	return nil
}

func TestWatchClient(t *testing.T) {
	// build config
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		t.Errorf("build config from flags: %v", err)
		return
	}

	// new client
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		t.Errorf("new client set: %v", err)
		return
	}

	// build indexers
	indexers := cache.Indexers{
		"app_index": func(obj interface{}) ([]string, error) {
			pod, ok := obj.(*v1.Pod)
			if !ok {
				return nil, errors.New("object is not a pod")
			}
			return []string{pod.ObjectMeta.Labels["app"]}, nil
		},
	}

	wc := NewWatchClient(clientset.CoreV1().RESTClient())
	wc.Watch(context.TODO(), "dev", "pods", &v1.Pod{}, labels.Everything(), fields.Everything(), indexers, &testPodWatcher{})
	<-context.TODO().Done()
}
