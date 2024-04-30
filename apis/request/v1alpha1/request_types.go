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

	"github.com/crossplane/crossplane-runtime/pkg/fieldpath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// RequestParameters are the configurable fields of a Request.
type RequestParameters struct {
	Mappings []Mapping           `json:"mappings"`
	Payload  Payload             `json:"payload"`
	Headers  map[string][]string `json:"headers,omitempty"`

	WaitTimeout *metav1.Duration `json:"waitTimeout,omitempty"`

	// InsecureSkipTLSVerify, when set to true, skips TLS certificate checks for the HTTP request
	InsecureSkipTLSVerify bool `json:"insecureSkipTLSVerify,omitempty"`
}

type Mapping struct {
	// +kubebuilder:validation:Enum=POST;GET;PUT;DELETE
	Method  string              `json:"method"`
	Body    string              `json:"body,omitempty"`
	URL     string              `json:"url"`
	Headers map[string][]string `json:"headers,omitempty"`

	// +kubebuilder:validation:Enum=CREATE;GET;UPDATE;DELETE
	Action string `json:"action"`
}

type Payload struct {
	BaseUrl string `json:"baseUrl,omitempty"`
	Body    string `json:"body,omitempty"`
	// Raw JSON representation of the kubernetes object to be created.

	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	BodyObject runtime.RawExtension `json:"body-object,omitempty"`
}

// A RequestSpec defines the desired state of a Request.
type RequestSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       RequestParameters `json:"forProvider"`
	References        []Reference       `json:"references,omitempty"`
	Readiness         Readiness         `json:"readiness,omitempty"`
}

// RequestObservation are the observable fields of a Request.
type Response struct {
	Headers map[string][]string `json:"headers,omitempty"`
	Body    string              `json:"body,omitempty"`

	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	BodyObject runtime.RawExtension `json:"body-object,omitempty"`

	StatusCode int `json:"statusCode,omitempty"`
}

// A RequestStatus represents the observed state of a Request.
type RequestStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	Response            Response `json:"response,omitempty"`
	Cache               Cache    `json:"cache,omitempty"`
	Failed              int32    `json:"failed,omitempty"`
	Error               string   `json:"error,omitempty"`
	RequestDetails      Mapping  `json:"requestDetails,omitempty"`
}

type Cache struct {
	LastUpdated string   `json:"lastUpdated,omitempty"`
	Response    Response `json:"response,omitempty"`
}

// +kubebuilder:object:root=true

// A Request is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,http}
type Request struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RequestSpec   `json:"spec"`
	Status RequestStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RequestList contains a list of Request
type RequestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Request `json:"items"`
}

// Request type metadata.
var (
	RequestKind             = reflect.TypeOf(Request{}).Name()
	RequestGroupKind        = schema.GroupKind{Group: Group, Kind: RequestKind}.String()
	RequestKindAPIVersion   = RequestKind + "." + SchemeGroupVersion.String()
	RequestGroupVersionKind = SchemeGroupVersion.WithKind(RequestKind)
)

func init() {
	SchemeBuilder.Register(&Request{}, &RequestList{})
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
}

// PatchesFrom refers to an object by Name, Kind, APIVersion, etc., and patch
// fields from this object.
type PatchesFrom struct {
	DependsOn `json:",inline"`
	// FieldPath is the path of the field on the resource whose value is to be
	// used as input.
	FieldPath *string `json:"fieldPath"`
}

// Reference refers to an Object or arbitrary Kubernetes resource and optionally
// patch values from that resource to the current Object.
type Reference struct {
	// DependsOn is used to declare dependency on other Object or arbitrary
	// Kubernetes resource.
	// +optional
	*DependsOn `json:"dependsOn,omitempty"`
	// PatchesFrom is used to declare dependency on other Object or arbitrary
	// Kubernetes resource, and also patch fields from this object.
	// +optional
	*PatchesFrom `json:"patchesFrom,omitempty"`
	// ToFieldPath is the path of the field on the resource whose value will
	// be changed with the result of transforms. Leave empty if you'd like to
	// propagate to the same path as patchesFrom.fieldPath.
	// +optional
	ToFieldPath *string `json:"toFieldPath,omitempty"`
}

// ReadinessPolicy defines how the Object's readiness condition should be computed.
type ReadinessPolicy string

const (
	// ReadinessPolicySuccessfulCreate means the object is marked as ready when the
	// underlying external resource is successfully created.
	ReadinessPolicySuccessfulCreate ReadinessPolicy = "SuccessfulCreate"
	// ReadinessPolicyDeriveFromObject means the object is marked as ready if and only if the underlying
	// external resource is considered ready.
	ReadinessPolicyDeriveFromObject ReadinessPolicy = "DeriveFromObject"
	// ReadinessPolicyAllTrue means that all conditions have status true on the object.
	// There must be at least one condition.
	ReadinessPolicyAllTrue ReadinessPolicy = "AllTrue"
)

// Readiness defines how the object's readiness condition should be computed,
// if not specified it will be considered ready as soon as the underlying external
// resource is considered up-to-date.
type Readiness struct {
	// Policy defines how the Object's readiness condition should be computed.
	// +optional
	// +kubebuilder:validation:Enum=SuccessfulCreate;DeriveFromObject;AllTrue
	// +kubebuilder:default=SuccessfulCreate
	Policy ReadinessPolicy `json:"policy,omitempty"`
}

// ApplyFromFieldPathPatch patches the "to" resource, using a source field
// on the "from" resource.
func (r *Reference) ApplyFromFieldPathPatch(from, to runtime.Object) error {
	// Default to patch the same field on the "to" resource.
	if r.ToFieldPath == nil {
		r.ToFieldPath = r.PatchesFrom.FieldPath
	}

	paved, err := fieldpath.PaveObject(from)
	if err != nil {
		return err
	}

	out, err := paved.GetValue(*r.PatchesFrom.FieldPath)
	if err != nil {
		return err
	}

	return PatchFieldValueToObject(*r.ToFieldPath, out, to)
}

// patchFieldValueToObject, given a path, value and "to" object, will
// apply the value to the "to" object at the given path, returning
// any errors as they occur.
func PatchFieldValueToObject(path string, value interface{}, to runtime.Object) error {
	paved, err := fieldpath.PaveObject(to)
	if err != nil {
		return err
	}

	err = paved.SetValue(path, value)
	if err != nil {
		return err
	}

	return runtime.DefaultUnstructuredConverter.FromUnstructured(paved.UnstructuredContent(), to)
}
