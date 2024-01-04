package controller

import (
	"context"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsclient "github.com/ironzhang/superdns/supercrd/clients/clientset/versioned"
)

func TestControllerWatchClusters(t *testing.T) {
	// build config
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		t.Errorf("build config from flags: %v", err)
		return
	}

	// new superdns client
	clientset, err := superdnsclient.NewForConfig(cfg)
	if err != nil {
		t.Errorf("new client set: %v", err)
		return
	}

	// new watch client
	wc := k8sclient.NewWatchClient(clientset.SuperdnsV1().RESTClient())

	// new controller
	controller := New("superdns", wc)

	// watch clusters
	err = controller.WatchClusters(context.TODO(), "example.app.com")
	if err != nil {
		t.Fatalf("watch clusters: %v", err)
	}

	<-context.TODO().Done()
}
