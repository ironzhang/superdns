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

func TestControllerWatchDomains(t *testing.T) {
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
	opts := Options{
		Namespace: "superdns",
	}
	controller := New(opts, wc, pm, fw)

	// watch domain
	err = controller.WatchDomain(context.TODO(), "example.app.com")
	if err != nil {
		t.Fatalf("watch domains: %v", err)
	}

	// refresh domain
	controller.RefreshClusters(context.TODO(), "example.app.com")
	controller.RefreshRoutes(context.TODO(), "example.app.com")
}

func TestMain(m *testing.M) {
	os.RemoveAll("./testdata")
	m.Run()
	os.RemoveAll("./testdata")
}
