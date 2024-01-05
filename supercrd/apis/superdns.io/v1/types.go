/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Cluster is a top-level type which represents cluster.
type Cluster struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec ClusterSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ClusterList ia a top-level list type which represents cluster list.
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Cluster `json:"items"`
}

type ClusterSpec struct {
	Cluster   string            `json:"cluster"`
	Features  map[string]string `json:"features"`
	Endpoints []Endpoint        `json:"endpoints"`
}

type Endpoint struct {
	Addr   string `json:"addr"`
	State  string `json:"state"`
	Weight int    `json:"weight"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Route is a top-level type which represents route strategy.
type Route struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`
	// +optional
	Spec RouteSpec `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RouteList ia a top-level list type which represents route strategy list.
type RouteList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Route `json:"items"`
}

type RouteSpec struct {
	EnableScript        bool          `json:"enableScript"`
	ScriptContent       string        `json:"scriptContent"`
	DefaultDestinations []Destination `json:"defaultDestinations"`
}

type Destination struct {
	Cluster string  `json:"cluster"`
	Percent float64 `json:"percent"`
}
