package k8sclient

import (
	"context"
	"testing"

	"github.com/ironzhang/tlog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"

	superv1 "github.com/ironzhang/superdns/v1"
)

type testPodWatcher struct {
}

func (p *testPodWatcher) OnWatch(store cache.Store, key string) error {
	obj, exists, err := store.GetByKey(key)
	if err != nil {
		tlog.Errorw("store get", "key", key, "error", err)
		return err
	}
	if !exists {
		tlog.Infow("pod does not exist", "key", key)
		return nil
	}

	pod := obj.(*v1.Pod)
	tlog.Infof("pod %s %s", pod.GetName(), pod.Status.Phase)
	return nil
}

type testClusterWatcher struct {
}

func (p *testClusterWatcher) OnWatch(store cache.Store, key string) error {
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

	wc := NewWatchClient(clientset)
	wc.Watch(context.TODO(), "dev", "pods", &v1.Pod{}, labels.Everything(), fields.Everything(), &testPodWatcher{})
	<-context.TODO().Done()
}

func TestWatchClientWatchCluster(t *testing.T) {
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

	wc := NewWatchClient(clientset)
	wc.Watch(context.TODO(), "superdns", "clusters", &superv1.Cluster{}, labels.Everything(), fields.Everything(), &testClusterWatcher{})
	<-context.TODO().Done()
}
