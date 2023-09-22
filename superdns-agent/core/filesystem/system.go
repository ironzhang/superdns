package filesystem

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ironzhang/superdns/superlib/model"
	"github.com/ironzhang/superdns/superlib/superutil"

	"github.com/ironzhang/superdns/superdns-agent/core/storesystem"
	"github.com/ironzhang/superdns/superdns-agent/core/storesystem/pod"
)

type System struct {
	resourcePath string
	store        *storesystem.System
}

func NewSystem(resourcePath string, store *storesystem.System) *System {
	return &System{resourcePath: resourcePath, store: store}
}

func (p *System) UpdateServiceModel(app string) error {
	a, ok := p.store.PodStore.GetApplication(app)
	if !ok {
		return nil
	}

	for _, g := range a.ListGroups() {
		sm := makeServiceModel(g)
		p.writeServiceModel(sm)

		rm := makeRouteModel(g)
		p.writeRouteModel(rm)
	}

	return nil
}

func (p *System) writeServiceModel(m model.ServiceModel) error {
	filename := p.serviceFilePath(m.Domain)
	dir := filepath.Dir(filename)
	os.MkdirAll(dir, os.ModePerm)
	return superutil.WriteJSON(filename, m)
}

func (p *System) writeRouteModel(m model.RouteModel) error {
	filename := p.routeFilePath(m.Domain)
	dir := filepath.Dir(filename)
	os.MkdirAll(dir, os.ModePerm)
	return superutil.WriteJSON(filename, m)
}

func (p *System) serviceFilePath(domain string) string {
	return fmt.Sprintf("%s/services/%s.json", p.resourcePath, domain)
}

func (p *System) routeFilePath(domain string) string {
	return fmt.Sprintf("%s/routes/%s.json", p.resourcePath, domain)
}

func makeServiceModel(g *pod.Group) model.ServiceModel {
	m := model.ServiceModel{
		Domain:   g.Domain(),
		Clusters: make(map[string]model.Cluster),
	}
	for _, c := range g.ListClusters() {
		m.Clusters[c.Name()] = makeClusterModel(c)
	}
	return m
}

func makeClusterModel(c *pod.Cluster) model.Cluster {
	return model.Cluster{
		Name:      c.Name(),
		Endpoints: c.ListEndpoints(),
	}
}

func makeRouteModel(g *pod.Group) model.RouteModel {
	m := model.RouteModel{
		Domain: g.Domain(),
		Strategy: model.RouteStrategy{
			DefaultDestinations: []model.Destination{
				{
					Cluster: "default",
					Percent: 1,
				},
			},
		},
	}
	return m
}
