//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2020 The Crossplane Authors.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cache) DeepCopyInto(out *Cache) {
	*out = *in
	in.Response.DeepCopyInto(&out.Response)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cache.
func (in *Cache) DeepCopy() *Cache {
	if in == nil {
		return nil
	}
	out := new(Cache)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DependsOn) DeepCopyInto(out *DependsOn) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DependsOn.
func (in *DependsOn) DeepCopy() *DependsOn {
	if in == nil {
		return nil
	}
	out := new(DependsOn)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mapping) DeepCopyInto(out *Mapping) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Action != nil {
		in, out := &in.Action, &out.Action
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mapping.
func (in *Mapping) DeepCopy() *Mapping {
	if in == nil {
		return nil
	}
	out := new(Mapping)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Mappings) DeepCopyInto(out *Mappings) {
	*out = *in
	if in.Create != nil {
		in, out := &in.Create, &out.Create
		*out = new(Mapping)
		(*in).DeepCopyInto(*out)
	}
	if in.Get != nil {
		in, out := &in.Get, &out.Get
		*out = new(Mapping)
		(*in).DeepCopyInto(*out)
	}
	if in.Update != nil {
		in, out := &in.Update, &out.Update
		*out = new(Mapping)
		(*in).DeepCopyInto(*out)
	}
	if in.Delete != nil {
		in, out := &in.Delete, &out.Delete
		*out = new(Mapping)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Mappings.
func (in *Mappings) DeepCopy() *Mappings {
	if in == nil {
		return nil
	}
	out := new(Mappings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PatchesFrom) DeepCopyInto(out *PatchesFrom) {
	*out = *in
	out.DependsOn = in.DependsOn
	if in.FieldPath != nil {
		in, out := &in.FieldPath, &out.FieldPath
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PatchesFrom.
func (in *PatchesFrom) DeepCopy() *PatchesFrom {
	if in == nil {
		return nil
	}
	out := new(PatchesFrom)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Payload) DeepCopyInto(out *Payload) {
	*out = *in
	in.Body.DeepCopyInto(&out.Body)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Payload.
func (in *Payload) DeepCopy() *Payload {
	if in == nil {
		return nil
	}
	out := new(Payload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Readiness) DeepCopyInto(out *Readiness) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Readiness.
func (in *Readiness) DeepCopy() *Readiness {
	if in == nil {
		return nil
	}
	out := new(Readiness)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Reference) DeepCopyInto(out *Reference) {
	*out = *in
	if in.DependsOn != nil {
		in, out := &in.DependsOn, &out.DependsOn
		*out = new(DependsOn)
		**out = **in
	}
	if in.PatchesFrom != nil {
		in, out := &in.PatchesFrom, &out.PatchesFrom
		*out = new(PatchesFrom)
		(*in).DeepCopyInto(*out)
	}
	if in.ToFieldPath != nil {
		in, out := &in.ToFieldPath, &out.ToFieldPath
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Reference.
func (in *Reference) DeepCopy() *Reference {
	if in == nil {
		return nil
	}
	out := new(Reference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Request) DeepCopyInto(out *Request) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Request.
func (in *Request) DeepCopy() *Request {
	if in == nil {
		return nil
	}
	out := new(Request)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Request) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequestList) DeepCopyInto(out *RequestList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Request, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequestList.
func (in *RequestList) DeepCopy() *RequestList {
	if in == nil {
		return nil
	}
	out := new(RequestList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *RequestList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequestParameters) DeepCopyInto(out *RequestParameters) {
	*out = *in
	in.Mappings.DeepCopyInto(&out.Mappings)
	in.Payload.DeepCopyInto(&out.Payload)
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.WaitTimeout != nil {
		in, out := &in.WaitTimeout, &out.WaitTimeout
		*out = new(v1.Duration)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequestParameters.
func (in *RequestParameters) DeepCopy() *RequestParameters {
	if in == nil {
		return nil
	}
	out := new(RequestParameters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequestSpec) DeepCopyInto(out *RequestSpec) {
	*out = *in
	in.ResourceSpec.DeepCopyInto(&out.ResourceSpec)
	in.ForProvider.DeepCopyInto(&out.ForProvider)
	if in.References != nil {
		in, out := &in.References, &out.References
		*out = make([]Reference, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Readiness = in.Readiness
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequestSpec.
func (in *RequestSpec) DeepCopy() *RequestSpec {
	if in == nil {
		return nil
	}
	out := new(RequestSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RequestStatus) DeepCopyInto(out *RequestStatus) {
	*out = *in
	in.ResourceStatus.DeepCopyInto(&out.ResourceStatus)
	in.Response.DeepCopyInto(&out.Response)
	in.Cache.DeepCopyInto(&out.Cache)
	in.RequestDetails.DeepCopyInto(&out.RequestDetails)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RequestStatus.
func (in *RequestStatus) DeepCopy() *RequestStatus {
	if in == nil {
		return nil
	}
	out := new(RequestStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Response) DeepCopyInto(out *Response) {
	*out = *in
	if in.Headers != nil {
		in, out := &in.Headers, &out.Headers
		*out = make(map[string][]string, len(*in))
		for key, val := range *in {
			var outVal []string
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]string, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	in.BodyObject.DeepCopyInto(&out.BodyObject)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Response.
func (in *Response) DeepCopy() *Response {
	if in == nil {
		return nil
	}
	out := new(Response)
	in.DeepCopyInto(out)
	return out
}
