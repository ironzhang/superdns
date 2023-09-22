package pod

import (
	"net"
	"strconv"

	v1 "k8s.io/api/core/v1"

	"github.com/ironzhang/superdns/superlib/model"
)

type Application struct {
	name   string
	groups map[string]*Group
}

func NewApplication(name string) *Application {
	return &Application{
		name:   name,
		groups: make(map[string]*Group),
	}
}

func (p *Application) getOrNewGroup(group string) *Group {
	g, ok := p.groups[group]
	if !ok {
		g = newGroup(p.name, group)
		p.groups[group] = g
	}
	return g
}

func (p *Application) AddPod(pod *v1.Pod) {
	cluster := pod.ObjectMeta.Labels["cluster"]
	for _, c := range pod.Spec.Containers {
		for _, cp := range c.Ports {
			c := p.getOrNewGroup(cp.Name).getOrNewCluster(cluster)
			c.addEndpoint(model.Endpoint{
				Addr:   net.JoinHostPort(pod.Status.PodIP, strconv.Itoa(int(cp.ContainerPort))),
				State:  podPhaseToState(pod.Status.Phase),
				Weight: 100,
			})
		}
	}
}

func (p *Application) ListGroups() []*Group {
	results := make([]*Group, 0, len(p.groups))
	for _, g := range p.groups {
		results = append(results, g)
	}
	return results
}

func podPhaseToState(phase v1.PodPhase) model.State {
	if phase == v1.PodRunning {
		return model.Enabled
	}
	return model.Disabled
}
