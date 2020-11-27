/*


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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// LabelsSpec defines the desired state of Labels
type LabelsSpec struct {
	Labels map[string]string `json:"labels,omitempty"`
}

// LabelsStatus defines the observed state of Labels
type LabelsStatus struct {
	ManagedDeployments int    `json:"managedDeployments,omitempty"`
	Labels             string `json:"labels,omitempty"`
}

// +kubebuilder:object:root=true
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Deployments",type="integer",JSONPath=".status.managedDeployments",description="Number of deployments managed by controller"
// +kubebuilder:printcolumn:name="Labels",type="string",JSONPath=".status.labels",description=""

// Labels is the Schema for the labels API
type Labels struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LabelsSpec   `json:"spec,omitempty"`
	Status LabelsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// LabelsList contains a list of Labels
type LabelsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Labels `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Labels{}, &LabelsList{})
}
