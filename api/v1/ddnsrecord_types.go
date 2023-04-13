/*
Copyright 2023.

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

// DdnsRecordSpec defines the desired state of DdnsRecord
type DdnsRecordSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Record to be updated for the given Zone.
	Record string `json:"record"`

	// Zone (DNS domain) of the record to updated.
	Zone string `json:"zone"`

	// TTL (time-to-live) of the DNS record (and update interval)/
	Ttl *metav1.Duration `json:"ttl,omitempty"`

	// DNS provider (configured via secret "ddns-config")
	Provider string `json:"provider"`
}

// DdnsRecordStatus defines the observed state of DdnsRecord
type DdnsRecordStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DdnsRecord is the Schema for the ddnsrecords API
type DdnsRecord struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DdnsRecordSpec   `json:"spec,omitempty"`
	Status DdnsRecordStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DdnsRecordList contains a list of DdnsRecord
type DdnsRecordList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DdnsRecord `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DdnsRecord{}, &DdnsRecordList{})
}
