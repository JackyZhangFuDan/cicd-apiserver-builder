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
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource"
	"sigs.k8s.io/apiserver-runtime/pkg/builder/resource/resourcestrategy"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// JenkinsService
// +k8s:openapi-gen=true
type JenkinsService struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JenkinsServiceSpec   `json:"spec,omitempty"`
	Status JenkinsServiceStatus `json:"status,omitempty"`
}

// JenkinsServiceList
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type JenkinsServiceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []JenkinsService `json:"items"`
}

// JenkinsServiceSpec defines the desired state of JenkinsService
// 后面的json注解没有自动生成，需要手动加上，后部有注解的都是这种情况
type JenkinsServiceSpec struct {
	InstanceAmount int `json:"instanceamount,omitempty"`
	InstanceCpu    int `json:"instancecpu,omitempty"`
}

type JenkinsServiceInstance struct {
	Cpu     int  `json:"cpu,omitempty"`
	Running bool `json:"running"`
}

var _ resource.Object = &JenkinsService{}
var _ resourcestrategy.Validater = &JenkinsService{}

func (in *JenkinsService) GetObjectMeta() *metav1.ObjectMeta {
	return &in.ObjectMeta
}

func (in *JenkinsService) NamespaceScoped() bool {
	return true
}

func (in *JenkinsService) New() runtime.Object {
	return &JenkinsService{}
}

func (in *JenkinsService) NewList() runtime.Object {
	return &JenkinsServiceList{}
}

func (in *JenkinsService) GetGroupVersionResource() schema.GroupVersionResource {
	return schema.GroupVersionResource{
		Group:    "cicd.autobusi.com",
		Version:  "v1",
		Resource: "jenkinsservices",
	}
}

func (in *JenkinsService) IsStorageVersion() bool {
	return true
}

func (in *JenkinsService) Validate(ctx context.Context) field.ErrorList {
	// TODO(user): Modify it, adding your API validation here.
	return nil
}

var _ resource.ObjectList = &JenkinsServiceList{}

func (in *JenkinsServiceList) GetListMeta() *metav1.ListMeta {
	return &in.ListMeta
}

// JenkinsServiceStatus defines the observed state of JenkinsService
type JenkinsServiceStatus struct {
	ApprovalStatus string                   `json:"approvalstatus"`
	Instances      []JenkinsServiceInstance `json:"instances"`
}

func (in JenkinsServiceStatus) SubResourceName() string {
	return "status"
}

// JenkinsService implements ObjectWithStatusSubResource interface.
var _ resource.ObjectWithStatusSubResource = &JenkinsService{}

func (in *JenkinsService) GetStatus() resource.StatusSubResource {
	return in.Status
}

// JenkinsServiceStatus{} implements StatusSubResource interface.
var _ resource.StatusSubResource = &JenkinsServiceStatus{}

func (in JenkinsServiceStatus) CopyTo(parent resource.ObjectWithStatusSubResource) {
	parent.(*JenkinsService).Status = in
}

// // 为什么v1作为storage version还是需要实现 MultiVersionObject呢？？不实现runtime居然会报错？？
// // NewStorageVersionObject returns a new empty instance of storage version.
// func (in *JenkinsService) NewStorageVersionObject() runtime.Object {
// 	return &JenkinsService{}
// }

// // ConvertToStorageVersion receives an new instance of storage version object as the conversion target
// // and overwrites it to the equal form of the current resource version.
// func (in *JenkinsService) ConvertToStorageVersion(storageObj runtime.Object) error {
// 	storageObj.(*JenkinsService).ObjectMeta = in.ObjectMeta
// 	storageObj.(*JenkinsService).TypeMeta = in.TypeMeta
// 	storageObj.(*JenkinsService).Spec.InstanceAmount = in.Spec.InstanceAmount
// 	storageObj.(*JenkinsService).Spec.InstanceCpu = in.Spec.InstanceCpu
// 	return nil
// }

// // ConvertFromStorageVersion receives an instance of storage version as the conversion source and
// // in-place mutates the current object to the equal form of the storage version object.
// func (in *JenkinsService) ConvertFromStorageVersion(storageObj runtime.Object) error {
// 	in = storageObj.(*JenkinsService)
// 	return nil
// }

// var _ resource.MultiVersionObject = &JenkinsService{}
