package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodRestartSpec defines the desired state of PodRestart
type PodRestartSpec struct {
	// LabelSelector defines which pods should be restarted
	LabelSelector map[string]string `json:"labelSelector,omitempty"`

	// RestartInterval defines how often the pods should be restarted (e.g., "10m", "1h")
	RestartInterval string `json:"restartInterval,omitempty"`
}

// PodRestartStatus defines the observed state of PodRestart
type PodRestartStatus struct {
	LastRestart metav1.Time `json:"lastRestart,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodRestart is the Schema for the podrestart API
type PodRestart struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodRestartSpec   `json:"spec,omitempty"`
	Status PodRestartStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodRestartList contains a list of PodRestart
type PodRestartList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodRestart `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodRestart{}, &PodRestartList{})
}
