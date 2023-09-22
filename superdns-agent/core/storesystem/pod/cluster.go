package pod

import "github.com/ironzhang/superdns/superlib/model"

type Cluster struct {
	app       string
	group     string
	cluster   string
	endpoints []model.Endpoint
}

func newCluster(app, group, cluster string) *Cluster {
	return &Cluster{
		app:     app,
		group:   group,
		cluster: cluster,
	}
}

func (p *Cluster) addEndpoint(endpoint model.Endpoint) {
	p.endpoints = append(p.endpoints, endpoint)
}

func (p *Cluster) Name() string {
	return p.cluster
}

func (p *Cluster) ListEndpoints() []model.Endpoint {
	return p.endpoints
}
