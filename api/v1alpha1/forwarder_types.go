package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ForwarderSpec defines the desired state of Forwarder
type ForwarderSpec struct {
    Topic string `json:"topic"`
}

// ForwarderStatus defines the observed state of Forwarder
type ForwarderStatus struct {
    MessageCount int `json:"messageCount"`
}

// +kubebuilder:object:root=true
type Forwarder struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   ForwarderSpec   `json:"spec,omitempty"`
    Status ForwarderStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type ForwarderList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Forwarder `json:"items"`
}
