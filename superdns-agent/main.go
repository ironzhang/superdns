package main

import (
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ironzhang/superdns/pkg/k8sclient"
	"github.com/ironzhang/superdns/superdns-agent/core/filesystem"
	"github.com/ironzhang/superdns/superdns-agent/core/storesystem"
	"github.com/ironzhang/superdns/superdns-agent/core/watchsystem"
	"github.com/ironzhang/superdns/superlib/parameter"
	"github.com/ironzhang/tlog"
)

func main() {
	tlog.Info("%v", parameter.Param)

	// build config
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		tlog.Errorf("build config from flags: %v", err)
		return
	}

	// new client
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		tlog.Errorf("new client set: %v", err)
		return
	}

	wc := k8sclient.NewWatchClient(clientset)
	ss := storesystem.NewSystem()
	fs := filesystem.NewSystem(parameter.Param.ResourcePath, ss)
	ws := watchsystem.NewSystem(wc, ss, fs)

	ctx := context.TODO()
	ws.Watch(ctx, "nginx")

	<-ctx.Done()
}
