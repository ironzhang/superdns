package app

import (
	"context"
	"net"

	"github.com/ironzhang/superlib/httputils/echoutil"
	"github.com/ironzhang/tlog"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/ironzhang/superdns/pkg/filewrite"
	"github.com/ironzhang/superdns/pkg/k8sclient"
	superdnsclient "github.com/ironzhang/superdns/supercrd/clients/clientset/versioned"
	"github.com/ironzhang/superdns/superdns-agent/internal/controller"
	"github.com/ironzhang/superdns/superdns-agent/internal/handlers/agentapi"
	"github.com/ironzhang/superdns/superdns-agent/internal/paths"
)

type Application struct {
	echosvr *echo.Echo
}

func newWatchClient() (*k8sclient.WatchClient, error) {
	cfg, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		tlog.Errorw("build k8s rest client config", "error", err)
		return nil, err
	}

	clientset, err := superdnsclient.NewForConfig(cfg)
	if err != nil {
		tlog.Errorw("new superdns client", "config", cfg, "error", err)
		return nil, err
	}

	return k8sclient.NewWatchClient(clientset.SuperdnsV1().RESTClient()), nil
}

func (p *Application) Init() error {
	ln, err := net.Listen("tcp", Conf.Listen.Addr)
	if err != nil {
		tlog.Errorw("net listen tcp", "addr", Conf.Listen.Addr, "error", err)
		return err
	}

	watchclient, err := newWatchClient()
	if err != nil {
		tlog.Errorw("new watch client", "error", err)
		return err
	}

	pathmgr := paths.NewPathManager(Conf.ResourcePath)
	fwriter := filewrite.NewFileWriter(pathmgr.TemporaryPath())
	control := controller.New(controller.Options{
		Namespace: Conf.Namespace,
	}, watchclient, pathmgr, fwriter)

	p.echosvr = echo.New()
	p.echosvr.Listener = ln
	p.echosvr.HTTPErrorHandler = echoutil.HTTPErrorHandler
	p.echosvr.Use(echoutil.Recover())
	p.echosvr.Use(echoutil.ServeMuxMiddleware(nil))
	p.echosvr.Use(middleware.RequestID())
	p.echosvr.Use(echoutil.TraceMiddleware())
	p.echosvr.Use(echoutil.AccessLogMiddleware())

	agentapi.Register(p.echosvr, agentapi.NewHandler(control))

	return nil
}

func (p *Application) Fini() error {
	return nil
}

func (p *Application) RunHTTPServer(ctx context.Context) error {
	tlog.Infof("serve http server on %s", Conf.Listen.Addr)
	go p.echosvr.StartServer(p.echosvr.Server)

	<-ctx.Done()
	p.echosvr.Close()
	return nil
}
