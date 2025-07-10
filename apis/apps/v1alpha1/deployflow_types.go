/*
Copyright 2021 The Triton Authors.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/intstr"
)

// 部署阶段
type DeployPhase string

const (
	Pending       DeployPhase = "Pending"
	Initializing  DeployPhase = "Initializing"
	BatchStarted  DeployPhase = "BatchStarted"
	BatchFinished DeployPhase = "BatchFinished"
	Success       DeployPhase = "Success"
	Failed        DeployPhase = "Failed"
	Aborted       DeployPhase = "Aborted"
	Canceled      DeployPhase = "Canceled"
)

// 部署类型
type DeployMode string

const (
	Auto   DeployMode = "auto"
	Manual DeployMode = "manual"
)

// 批次阶段
type BatchPhase string

const (
	BatchPending     BatchPhase = "Pending"
	BatchSmoking     BatchPhase = "Smoking"
	BatchSmoked      BatchPhase = "Smoked"
	BatchBaking      BatchPhase = "Baking"
	BatchBaked       BatchPhase = "Baked"
	BatchSmokeFailed BatchPhase = "SmokeFailed"
	BatchBakeFailed  BatchPhase = "BakeFailed"
)

/**
主要功能模块：

1、身份标识系统
AppID + GroupID：构成全局唯一标识
AppName：人类可读的语义化名称
CloneSetName：关联底层 Kubernetes 资源
2、部署规格控制
Replicas：控制应用副本数量
Template：定义 Pod 的完整配置（容器镜像、环境变量、资源限制等）
Selector：确保 Pod 与 Controller 的关联关系
3、高级特性
VolumeClaimTemplates：支持有状态应用的持久化存储
ApplicationType：兼容不同 workload 类型（如 cloneset/advanced statefulset）
ApplicationLabel：实现应用维度的监控/日志采集
*/
// ApplicationSpec describes the new application state which will be created or updated by a Deploy
type ApplicationSpec struct {
	AppID   int    `json:"appID"`
	GroupID int    `json:"groupID"`
	AppName string `json:"appName"`
	// CloneSetName，cloneset层面唯一标识。
	// 命名建议格式：<appID>-<appName>-<groupID>-<env>，示例：12122-my-web-app-10010-prod
	CloneSetName string `json:"clonesetName"`

	// Selector is a label query over pods that should match the replica count.
	// It must match the pod template's labels.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	// 选择器，Controller 和 Pod 之间的桥梁。
	// 这个字段的作用是告诉控制器哪些POD属于这个应用。确保在更新、扩缩容时操作的是正确的Pod
	Selector *metav1.LabelSelector `json:"selector"`

	// Template describes the pods that will be created.
	Template corev1.PodTemplateSpec `json:"template"`

	// Replicas defines the replicas num of app
	Replicas *int32 `json:"replicas,omitempty"`

	// ApplicationType defines the type of app, ex: cloneset, advanced statefulset
	ApplicationType string `json:"applicationType,omitempty"`

	// VolumeClaimTemplates is a list of claims that pods are allowed to reference.
	// Note that PVC will be deleted when its pod has been deleted.
	VolumeClaimTemplates []corev1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`

	// ApplicationLabel defines the label of app
	ApplicationLabel map[string]string `json:"applicationLabel,omitempty"`
}

// DeployFlowSpec defines the desired state of DeployFlow
// 发布策略的定义
type DeployFlowSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Action defines the action of deployflow, ex: create, update, restart, scale...
	// 部署类型：create/update/restart/scale
	Action string `json:"action"`

	// Application defines somethings about app to deploy
	// 声明应用终态。
	Application *ApplicationSpec `json:"application"`

	// 更新策略，金丝雀发布，普通批次发布。声明过程约束
	UpdateStrategy *DeployUpdateStrategy `json:"updateStrategy,omitempty"`

	// 非更新操作策略，重启/扩缩容
	NonUpdateStrategy *DeployNonUpdateStrategy `json:"nonUpdateStrategy,omitempty"`
}

type DeployUpdateStrategy struct {
	BaseStrategy `json:",inline"`

	// NoPullIn indicates that the pullIn step in batch Baking phase will be skipped, which
	// means that as long as the pod is ready, traffic from outside will come in.
	// Default value is false
	NoPullIn bool `json:"noPullIn,omitempty"`

	// +kubebuilder:validation:Optional
	Canary int `json:"canary,omitempty"`

	// +kubebuilder:validation:Optional

	// Stage describes the desired stage you want to go to.
	Stage BatchPhase `json:"stage,omitempty"`
}

type DeployNonUpdateStrategy struct {
	BaseStrategy `json:",inline"`

	// +kubebuilder:validation:Optional
	// +nullable

	// PodsToDelete is the names of Pod should be deleted.
	PodsToDelete []string `json:"podsToDelete,omitempty"`
}

type BaseStrategy struct {
	// Paused indicates that the Deploy should be paused or resumed.
	// Set true to pause the deploy, false to resume the deploy.
	Paused *bool `json:"paused,omitempty"`

	// Canceled indicates that the Deploy should be canceled.
	// Default value is false
	Canceled bool `json:"canceled,omitempty"`

	// +kubebuilder:validation:Optional

	// number of pods that can be scheduled at a time. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// Absolute number is calculated from percentage by rounding up. Defaults to the same value with Replicas
	// Value can be changed during a deploy. If it is changed, .status.batches needs to be calculated again.
	BatchSize *intstr.IntOrString `json:"batchSize,omitempty"`

	// Minimum time interval to wait between two batches
	BatchIntervalSeconds int32 `json:"batchIntervalSeconds,omitempty"`

	// Deploy mode, candidates are "auto" and "manual", if not set, default to "manual".
	// "manual" indicates that the DeployFlow is controlled by user, he can make progress by updating "Batches",
	// "auto" indicates that the DeployFlow will always move forward no matter what "Batches" is.
	Mode DeployMode `json:"mode,omitempty"`

	// +kubebuilder:default=1

	// Batches is the number of batch you want to finish
	Batches int `json:"batches,omitempty"`
}

// DeployFlowStatus defines the observed state of DeployFlow
type DeployFlowStatus struct {
	// Important: Run "make" to regenerate code after modifying this file

	// Replicas is the number of Pods created by the CloneSet controller.
	Replicas int32 `json:"replicas"`

	// ReplicasToProcess is the number of Pods which will be created/restarted/deleted in this Deploy.
	ReplicasToProcess int32 `json:"replicasToProcess"`

	// AvailableReplicas is the number of Pods created by the CloneSet controller that have a Ready Condition for at least minReadySeconds.
	AvailableReplicas int32 `json:"availableReplicas"`

	// UpdatedReplicas is the number of Pods created by the CloneSet controller from the CloneSet version
	// indicated by updateRevision.
	UpdatedReplicas int32 `json:"updatedReplicas"`

	// UpdatedReadyReplicas is the number of Pods created by the CloneSet controller from the CloneSet version
	// indicated by updateRevision and have a Ready Condition.
	UpdatedReadyReplicas int32 `json:"updatedReadyReplicas"`

	// UpdateRevision, if not empty, indicates the latest revision of the CloneSet.
	UpdateRevision string `json:"updateRevision,omitempty"`

	// 批次状态明细（子状态）
	Conditions []BatchCondition `json:"conditions"`

	// +kubebuilder:validation:Optional
	// +nullable
	Pods []string `json:"pods"`

	// +kubebuilder:validation:Optional
	Paused bool `json:"paused"`
	// 全局阶段（Pending/Processing/Success）
	Phase DeployPhase `json:"phase"`
	// +kubebuilder:validation:Optional
	Finished bool `json:"finished"`
	// 总批次数
	Batches int `json:"batches"`
	// +kubebuilder:validation:Optional
	FinishedBatches int `json:"finishedBatches"`
	// +kubebuilder:validation:Optional
	FinishedReplicas int `json:"finishedReplicas"`
	FailedReplicas   int `json:"failedReplicas"`

	// +nullable
	StartedAt metav1.Time `json:"startedAt,omitempty"`

	// +nullable
	UpdatedAt metav1.Time `json:"updatedAt,omitempty"`

	// +nullable
	FinishedAt metav1.Time `json:"finishedAt,omitempty"`
}

type BatchCondition struct {
	// 当前批次
	Batch int `json:"batch"`
	// 批次大小
	BatchSize int  `json:"batchSize"`
	Canary    bool `json:"canary"`
	// 批次阶段（Smoked/Baking/Baked）
	Phase          BatchPhase `json:"phase"`
	FailedReplicas int        `json:"failedReplicas"`

	// +kubebuilder:validation:Optional
	// +nullable
	Pods []PodInfo `json:"pods"`

	// +nullable
	StartedAt metav1.Time `json:"startedAt,omitempty"`

	// +nullable
	PulledInAt metav1.Time `json:"pulledInAt,omitempty"`

	// +nullable
	FinishedAt metav1.Time `json:"finishedAt,omitempty"`
}

type PodInfo struct {
	Name         string `json:"name"`
	IP           string `json:"ip"`
	Port         int32  `json:"port"`
	Phase        string `json:"phase"`
	PullInStatus string `json:"pullInStatus"`
}

// +kubebuilder:subresource:status
// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=df
// +kubebuilder:printcolumn:name="REPLICAS",type="integer",JSONPath=".status.replicas",description="Replicas of this Deploy"
// +kubebuilder:printcolumn:name="UPDATED_READY_REPLICAS",type="integer",JSONPath=".status.updatedReadyReplicas",description="Updated and ready replicas"
// +kubebuilder:printcolumn:name="FINISHED_REPLICAS",type="integer",JSONPath=".status.finishedReplicas",description="Total replicas in all finished batches"
// +kubebuilder:printcolumn:name="PHASE",type="string",JSONPath=".status.phase",description="Phase of this Deploy"
// +kubebuilder:printcolumn:name="BATCHES",type="integer",JSONPath=".status.batches",description="Total batches of this Deploy"
// +kubebuilder:printcolumn:name="CURRENT_BATCH_SIZE",type="integer",JSONPath=".status.conditions[-1].batchSize",description="Size of current batch"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC."
// +kubebuilder:printcolumn:name="UPDATE_AT",type="date",priority = 1,JSONPath=".status.updatedAt",description="The last update time of deployflow. It is represented in RFC3339 form and is in UTC."
// +kubebuilder:printcolumn:name="PODS",type="string",priority = 1,JSONPath=".status.pods",description="The pods deployed by this deployflow."

// DeployFlow is the Schema for the deployflows API
type DeployFlow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeployFlowSpec   `json:"spec,omitempty"`
	Status DeployFlowStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeployFlowList contains a list of DeployFlow
type DeployFlowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeployFlow `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeployFlow{}, &DeployFlowList{})
}
