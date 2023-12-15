package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Endpoint struct {
	Addr     string
	State    string
	PodState string
	Weight   string
}

type ClusterSpec struct {
	Domain    string
	Cluster   string
	Endpoints []Endpoint
}

type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec ClusterSpec `json:"spec,omitempty"`
}
