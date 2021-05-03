/*
Copyright 2021.

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

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ConnectionSpec defines the desired state of Connection
type ConnectionSpec struct {
	//+kubebuilder:validation:Required
	//+kubebuilder:validation:MinLength=3
	// Type is the database type
	Type string `json:"type"`

	// Provider is the database cloud name
	Provider string `json:"provider"`

	//+kubebuilder:validation:Required
	// Database is the name of the database instance to import
	Database string `json:"database"`
}

// ConnectionStatus defines the observed state of Connection
type ConnectionStatus struct {
	//+kubebuilder:validation:Required
	// DBConfigMap is the name of the ConfigMap containing the connection info
	DBConfigMap string `json:"dbConfigMap"`

	//+kubebuilder:validation:Required
	// DBCredentials is the name of the Secret containing the database credentials
	DBCredentials string `json:"dbCredentials"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Connection is the Schema for the connections API
type Connection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ConnectionSpec   `json:"spec,omitempty"`
	Status ConnectionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ConnectionList contains a list of Connection
type ConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Connection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Connection{}, &ConnectionList{})
}
