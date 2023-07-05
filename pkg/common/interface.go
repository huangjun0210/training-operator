/*
Copyright 2023 The Kubeflow Authors.

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

package common

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	apiv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
)

// ControllerInterface defines the Interface to be implemented by custom operators. e.g. tf-operator needs to implement this interface
type ControllerInterface interface {
	// Returns the Controller name
	ControllerName() string

	// Returns the GroupVersionKind of the API
	GetAPIGroupVersionKind() schema.GroupVersionKind

	// Returns the GroupVersion of the API
	GetAPIGroupVersion() schema.GroupVersion

	// Returns the Group Name(value) in the labels of the job
	GetGroupNameLabelValue() string

	// Returns the Job from Informer Cache
	GetJobFromInformerCache(namespace, name string) (metav1.Object, error)

	// Returns the Job from API server
	GetJobFromAPIClient(namespace, name string) (metav1.Object, error)

	// GetPodsForJob returns the pods managed by the job. This can be achieved by selecting pods using label key "job-name"
	// i.e. all pods created by the job will come with label "job-name" = <this_job_name>
	GetPodsForJob(job interface{}) ([]*v1.Pod, error)

	// GetServicesForJob returns the services managed by the job. This can be achieved by selecting services using label key "job-name"
	// i.e. all services created by the job will come with label "job-name" = <this_job_name>
	GetServicesForJob(job interface{}) ([]*v1.Service, error)

	// DeleteJob deletes the job
	DeleteJob(job interface{}) error

	// UpdateJobStatus updates the job status and job conditions
	UpdateJobStatus(job interface{}, replicas map[apiv1.ReplicaType]*apiv1.ReplicaSpec, jobStatus *apiv1.JobStatus) error

	// UpdateJobStatusInApiServer updates the job status in API server
	UpdateJobStatusInApiServer(job interface{}, jobStatus *apiv1.JobStatus) error

	// SetClusterSpec sets the cluster spec for the pod
	SetClusterSpec(job interface{}, podTemplate *v1.PodTemplateSpec, rtype, index string) error

	// Returns the default container name in pod
	GetDefaultContainerName() string

	// Get the default container port name
	GetDefaultContainerPortName() string

	// Returns if this replica type with index specified is a master role.
	// MasterRole pod will have "job-role=master" set in its label
	IsMasterRole(replicas map[apiv1.ReplicaType]*apiv1.ReplicaSpec, rtype apiv1.ReplicaType, index int) bool

	// ReconcileJobs checks and updates replicas for each given ReplicaSpec of a job.
	// Common implementation will be provided and User can still override this to implement their own reconcile logic
	ReconcileJobs(job interface{}, replicas map[apiv1.ReplicaType]*apiv1.ReplicaSpec, jobStatus apiv1.JobStatus, runPolicy *apiv1.RunPolicy) error

	// ReconcilePods checks and updates pods for each given ReplicaSpec.
	// It will requeue the job in case of an error while creating/deleting pods.
	// Common implementation will be provided and User can still override this to implement their own reconcile logic
	ReconcilePods(job interface{}, jobStatus *apiv1.JobStatus, pods []*v1.Pod, rtype apiv1.ReplicaType, spec *apiv1.ReplicaSpec,
		replicas map[apiv1.ReplicaType]*apiv1.ReplicaSpec) error

	// ReconcileServices checks and updates services for each given ReplicaSpec.
	// It will requeue the job in case of an error while creating/deleting services.
	// Common implementation will be provided and User can still override this to implement their own reconcile logic
	ReconcileServices(job metav1.Object, services []*v1.Service, rtype apiv1.ReplicaType, spec *apiv1.ReplicaSpec) error

	// GetFrameworkName returns framework name (e.g., tensorflow).
	GetFrameworkName() string
}
