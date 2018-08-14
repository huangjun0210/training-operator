// Copyright 2018 The Kubeflow Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package controller provides a Kubernetes controller for a TFJob resource.
package tfcontroller

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	tfv1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	"github.com/kubeflow/tf-operator/pkg/controller.v2/jobcontroller"
	tflogger "github.com/kubeflow/tf-operator/pkg/logger"
	train_util "github.com/kubeflow/tf-operator/pkg/util/train"
)

const (
	// tfConfig is the environment variable name of TensorFlow cluster spec.
	tfConfig = "TF_CONFIG"

	// podTemplateRestartPolicyReason is the warning reason when the restart
	// policy is setted in pod template.
	podTemplateRestartPolicyReason = "SettedPodTemplateRestartPolicy"
)

// reconcilePods checks and updates pods for each given TFReplicaSpec.
// It will requeue the tfjob in case of an error while creating/deleting pods.
func (tc *TFController) reconcilePods(
	tfjob *tfv1alpha2.TFJob,
	pods []*v1.Pod,
	rtype tfv1alpha2.TFReplicaType,
	spec *tfv1alpha2.TFReplicaSpec, rstatus map[string]v1.PodPhase) error {

	// Convert TFReplicaType to lower string.
	rt := strings.ToLower(string(rtype))
	logger := tflogger.LoggerForReplica(tfjob, rt)
	// Get all pods for the type rt.
	pods, err := filterPodsForTFReplicaType(pods, rt)
	if err != nil {
		return err
	}
	replicas := int(*spec.Replicas)
	restart := false

	initializeTFReplicaStatuses(tfjob, rtype)

	podSlices := getPodSlices(pods, replicas, logger)
	for index, podSlice := range podSlices {
		if len(podSlice) > 1 {
			logger.Warningf("We have too many pods for %s %d", rt, index)
			// TODO(gaocegege): Kill some pods.
		} else if len(podSlice) == 0 {
			logger.Infof("Need to create new pod: %s-%d", rt, index)
			err = tc.createNewPod(tfjob, rt, strconv.Itoa(index), spec)
			if err != nil {
				return err
			}
		} else {
			// Check the status of the current pod.
			pod := podSlice[0]
			// Check if the pod is retryable.
			if spec.RestartPolicy == tfv1alpha2.RestartPolicyExitCode {
				var exitCode int32
				for _, status := range pod.Status.ContainerStatuses {
					state := status.State
					// Get the exit code of the tensorflow container.
					if status.Name == tfv1alpha2.DefaultContainerName && state.Terminated != nil {
						exitCode = state.Terminated.ExitCode
					}
				}
				if pod.Status.Phase == v1.PodFailed && train_util.IsRetryableExitCode(exitCode) {
					logger.Infof("Need to restart the pod: %s-%d", rt, index)
					if err := tc.PodControl.DeletePod(pod.Namespace, pod.Name, tfjob); err != nil {
						return err
					}
					restart = true
				}
			}
			updateTFJobReplicaStatuses(tfjob, rtype, pod)
		}
	}

	return updateStatusSingle(tfjob, rtype, replicas, restart)
}

// getPodSlices returns a slice, which element is the slice of pod.
func getPodSlices(pods []*v1.Pod, replicas int, logger *log.Entry) [][]*v1.Pod {
	podSlices := make([][]*v1.Pod, replicas)
	for _, pod := range pods {
		if _, ok := pod.Labels[tfReplicaIndexLabel]; !ok {
			logger.Warning("The pod do not have the index label.")
			continue
		}
		index, err := strconv.Atoi(pod.Labels[tfReplicaIndexLabel])
		if err != nil {
			logger.Warningf("Error when strconv.Atoi: %v", err)
			continue
		}
		if index < 0 || index >= replicas {
			logger.Warningf("The label index is not expected: %d", index)
		} else {
			podSlices[index] = append(podSlices[index], pod)
		}
	}
	return podSlices
}

// createNewPod creates a new pod for the given index and type.
func (tc *TFController) createNewPod(tfjob *tfv1alpha2.TFJob, rt, index string, spec *tfv1alpha2.TFReplicaSpec) error {
	tfjobKey, err := KeyFunc(tfjob)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("Couldn't get key for tfjob object %#v: %v", tfjob, err))
		return err
	}
	expectationPodsKey := jobcontroller.GenExpectationPodsKey(tfjobKey, rt)
	err = tc.Expectations.ExpectCreations(expectationPodsKey, 1)
	if err != nil {
		return err
	}
	logger := tflogger.LoggerForReplica(tfjob, rt)
	// Create OwnerReference.
	controllerRef := tc.GenOwnerReference(tfjob)

	// Set type and index for the worker.
	labels := tc.GenLabels(tfjob.Name)
	labels[tfReplicaTypeLabel] = rt
	labels[tfReplicaIndexLabel] = index

	podTemplate := spec.Template.DeepCopy()

	// Set name for the template.
	podTemplate.Name = jobcontroller.GenGeneralName(tfjob.Name, rt, index)

	if podTemplate.Labels == nil {
		podTemplate.Labels = make(map[string]string)
	}

	for key, value := range labels {
		podTemplate.Labels[key] = value
	}

	if err := setClusterSpec(podTemplate, tfjob, rt, index); err != nil {
		return err
	}

	// Submit a warning event if the user specifies restart policy for
	// the pod template. We recommend to set it from the replica level.
	if podTemplate.Spec.RestartPolicy != v1.RestartPolicy("") {
		errMsg := "Restart policy in pod template will be overwritten by restart policy in replica spec"
		logger.Warning(errMsg)
		tc.Recorder.Event(tfjob, v1.EventTypeWarning, podTemplateRestartPolicyReason, errMsg)
	}
	setRestartPolicy(podTemplate, spec)

	err = tc.PodControl.CreatePodsWithControllerRef(tfjob.Namespace, podTemplate, tfjob, controllerRef)
	if err != nil && errors.IsTimeout(err) {
		// Pod is created but its initialization has timed out.
		// If the initialization is successful eventually, the
		// controller will observe the creation via the informer.
		// If the initialization fails, or if the pod keeps
		// uninitialized for a long time, the informer will not
		// receive any update, and the controller will create a new
		// pod when the expectation expires.
		return nil
	} else if err != nil {
		return err
	}
	return nil
}

func setClusterSpec(podTemplateSpec *v1.PodTemplateSpec, tfjob *tfv1alpha2.TFJob, rt, index string) error {
	// Generate TF_CONFIG JSON string.
	tfConfigStr, err := genTFConfigJSONStr(tfjob, rt, index)
	if err != nil {
		return err
	}

	if tfConfigStr == "" {
		return nil
	}
	// Add TF_CONFIG environment variable.
	for i := range podTemplateSpec.Spec.Containers {
		if len(podTemplateSpec.Spec.Containers[i].Env) == 0 {
			podTemplateSpec.Spec.Containers[i].Env = make([]v1.EnvVar, 0)
		}
		podTemplateSpec.Spec.Containers[i].Env = append(podTemplateSpec.Spec.Containers[i].Env, v1.EnvVar{
			Name:  tfConfig,
			Value: tfConfigStr,
		})
	}
	return nil
}

func setRestartPolicy(podTemplateSpec *v1.PodTemplateSpec, spec *tfv1alpha2.TFReplicaSpec) {
	if spec.RestartPolicy == tfv1alpha2.RestartPolicyExitCode {
		podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicyNever
	} else {
		podTemplateSpec.Spec.RestartPolicy = v1.RestartPolicy(spec.RestartPolicy)
	}
}

// filterPodsForTFReplicaType returns pods belong to a TFReplicaType.
func filterPodsForTFReplicaType(pods []*v1.Pod, tfReplicaType string) ([]*v1.Pod, error) {
	var result []*v1.Pod

	tfReplicaSelector := &metav1.LabelSelector{
		MatchLabels: make(map[string]string),
	}

	tfReplicaSelector.MatchLabels[tfReplicaTypeLabel] = tfReplicaType

	for _, pod := range pods {
		selector, err := metav1.LabelSelectorAsSelector(tfReplicaSelector)
		if err != nil {
			return nil, err
		}
		if !selector.Matches(labels.Set(pod.Labels)) {
			continue
		}
		result = append(result, pod)
	}
	return result, nil
}
