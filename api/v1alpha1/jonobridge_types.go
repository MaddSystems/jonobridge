package v1alpha1

import (
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JonobridgeSpec defines the desired state of Jonobridge
type JonobridgeSpec struct {
    LocalAddress  string `json:"localAddress"`
    RemoteAddress string `json:"remoteAddress"`
}

// JonobridgeStatus defines the observed state of Jonobridge
type JonobridgeStatus struct {
    ConnectedClients int `json:"connectedClients"`
}

// +kubebuilder:object:root=true
// Jonobridge is the Schema for the jonobridge API
type Jonobridge struct {
    metav1.TypeMeta   `json:",inline"`
    metav1.ObjectMeta `json:"metadata,omitempty"`

    Spec   JonobridgeSpec   `json:"spec,omitempty"`
    Status JonobridgeStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
type JonobridgeList struct {
    metav1.TypeMeta `json:",inline"`
    metav1.ListMeta `json:"metadata,omitempty"`
    Items           []Jonobridge `json:"items"`
}
