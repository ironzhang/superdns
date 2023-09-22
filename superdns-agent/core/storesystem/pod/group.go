package pod

import "fmt"

type Group struct {
	app      string
	group    string
	clusters map[string]*Cluster
}

func newGroup(app, group string) *Group {
	return &Group{
		app:      app,
		group:    group,
		clusters: make(map[string]*Cluster),
	}
}

func (p *Group) getOrNewCluster(cluster string) *Cluster {
	c, ok := p.clusters[cluster]
	if !ok {
		c = newCluster(p.app, p.group, cluster)
		p.clusters[cluster] = c
	}
	return c
}

func (p *Group) Domain() string {
	return fmt.Sprintf("%s.%s", p.group, p.app)
}

func (p *Group) ListClusters() []*Cluster {
	results := make([]*Cluster, 0, len(p.clusters))
	for _, c := range p.clusters {
		results = append(results, c)
	}
	return results
}
