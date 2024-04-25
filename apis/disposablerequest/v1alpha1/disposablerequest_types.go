/*
Copyright 2022 The Crossplane Authors.

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

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// DisposableRequestParameters are the configurable fields of a DisposableRequest.
type DisposableRequestParameters struct {
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.url' is immutable"
	URL string `json:"url"`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.method' is immutable"
	Method string `json:"method"`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.headers' is immutable"
	Headers map[string][]string `json:"headers,omitempty"`
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Field 'forProvider.body' is immutable"
	Body string `json:"body,omitempty"`

	WaitTimeout *metav1.Duration `json:"waitTimeout,omitempty"`

	// RollbackRetriesLimit is max number of attempts to retry HTTP request by sending again the request.
	RollbackRetriesLimit *int32 `json:"rollbackRetriesLimit,omitempty"`

	// InsecureSkipTLSVerify, when set to true, skips TLS certificate checks for the HTTP request
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`

	// ExpectedResponse is a jq filter expression used to evaluate the HTTP response and determine if it matches the expected criteria.
	// The expression should return a boolean; if true, the response is considered expected.
	// Example: '.Body.job_status == "success"'
	ExpectedResponse string `json:"expectedResponse,omitempty"`
}

// DependsOn refers to an object by Name, Kind, APIVersion, etc. It is used to
// reference other Object or arbitrary Kubernetes resource which is either
// cluster or namespace scoped.
type DependsOn struct {
	// APIVersion of the referenced object.
	// +kubebuilder:default=kubernetes.crossplane.io/v1alpha1
	// +optional
	APIVersion string `json:"apiVersion,omitempty"`
	// Kind of the referenced object.
	// +kubebuilder:default=Object
	// +optional
	Kind string `json:"kind,omitempty"`
	// Name of the referenced object.
	Name string `json:"name"`
	// Namespace of the referenced object.
	// +optional
	Namespace string `json:"namespace,omitempty"`

	// FieldPath defines the object source to assess value (optional)
	FieldPath *string `json:"fieldPath"`

	// ExpectedValue defines value to assert on resolved reference field (optional)
	ExpectedValue *string `json:"expectedValue"`
}

// Reference refers to an Object or arbitrary Kubernetes resource and optionally
// patch values from that resource to the current Object.
type Reference struct {
	// DependsOn is used to declare dependency on other Object or arbitrary
	// Kubernetes resource.
	// +optional
	*DependsOn `json:"dependsOn,omitempty"`
}

// A DisposableRequestSpec defines the desired state of a DisposableRequest.
type DisposableRequestSpec struct {
	xpv1.ResourceSpec `json:",inline"`

	ForProvider DisposableRequestParameters `json:"forProvider"`
	References  []Reference                 `json:"references,omitempty"`
}

type Response struct {
	StatusCode int                 `json:"statusCode,omitempty"`
	Body       string              `json:"body,omitempty"`
	Headers    map[string][]string `json:"headers,omitempty"`
}

type Mapping struct {
	Method  string              `json:"method"`
	Body    string              `json:"body,omitempty"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// A DisposableRequestStatus represents the observed state of a DisposableRequest.
type DisposableRequestStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	Response            Response `json:"response,omitempty"`
	Failed              int32    `json:"failed,omitempty"`
	Error               string   `json:"error,omitempty"`
	Synced              bool     `json:"synced,omitempty"`
	RequestDetails      Mapping  `json:"requestDetails,omitempty"`
}

// +kubebuilder:object:root=true

// A DisposableRequest is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,http}
type DisposableRequest struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DisposableRequestSpec   `json:"spec"`
	Status DisposableRequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DisposableRequestList contains a list of DisposableRequest
type DisposableRequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DisposableRequest `json:"items"`
}

// DisposableRequest type metadata.
var (
	DisposableRequestKind             = reflect.TypeOf(DisposableRequest{}).Name()
	DisposableRequestGroupKind        = schema.GroupKind{Group: Group, Kind: DisposableRequestKind}.String()
	DisposableRequestKindAPIVersion   = DisposableRequestKind + "." + SchemeGroupVersion.String()
	DisposableRequestGroupVersionKind = SchemeGroupVersion.WithKind(DisposableRequestKind)
)

func init() {
	SchemeBuilder.Register(&DisposableRequest{}, &DisposableRequestList{})
}
