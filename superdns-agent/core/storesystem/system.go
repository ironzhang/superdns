package storesystem

import "github.com/ironzhang/superdns/superdns-agent/core/storesystem/pod"

type System struct {
	PodStore *pod.Store
}

func NewSystem() *System {
	return &System{
		PodStore: pod.NewStore(),
	}
}
