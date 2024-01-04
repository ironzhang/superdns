package controller

import (
	"context"
	"os"
	"testing"

	"k8s.io/client-go/tools/clientcmd"

	"github.com/ironzhang/superdns/pkg/filewrite"
	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsclient "github.com/ironzhang/superdns/supercrd/clients/clientset/versioned"
	"github.com/ironzhang/superdns/superdns-agent/internal/paths"
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

	// new path manager
	pm := paths.NewPathManager("./testdata")

	// new file writer
	fw := filewrite.NewFileWriter(pm.TemporaryPath())

	// new controller
	controller := New("superdns", wc, pm, fw)

	// watch clusters
	err = controller.WatchClusters(context.TODO(), "example.app.com")
	if err != nil {
		t.Fatalf("watch clusters: %v", err)
	}
	err = controller.WatchClusters(context.TODO(), "example1.app.com")
	if err != nil {
		t.Fatalf("watch clusters: %v", err)
	}

	<-context.TODO().Done()
}

func TestMain(m *testing.M) {
	os.RemoveAll("./testdata")
	m.Run()
	os.RemoveAll("./testdata")
}
