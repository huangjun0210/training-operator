// Copyright 2021 The Kubeflow Authors
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

package common_test

import (
	"context"
	"reflect"
	"strings"
	"testing"

	kubeflowv1 "github.com/kubeflow/training-operator/pkg/apis/kubeflow.org/v1"
	testjobv1 "github.com/kubeflow/training-operator/test_job/apis/test_job/v1"
	"github.com/kubeflow/training-operator/test_job/reconciler.v1/test_job"
	testutilv1 "github.com/kubeflow/training-operator/test_job/test_util/v1"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestCreateNewService(t *testing.T) {
	type tc struct {
		testJob         *testjobv1.TestJob
		testRType       kubeflowv1.ReplicaType
		testSpec        *kubeflowv1.ReplicaSpec
		testIndex       string
		expectedService *corev1.Service
	}
	testCase := []tc{
		func() tc {
			tj := testutilv1.NewTestJob(3)
			jobName := "testjob1"
			tj.SetName(jobName)
			idx := "0"
			svc := &corev1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Name:      jobName + "-worker-" + idx,
					Namespace: corev1.NamespaceDefault,
				},
				Spec: corev1.ServiceSpec{
					Ports: []corev1.ServicePort{
						corev1.ServicePort{
							Name: testjobv1.DefaultPortName,
							Port: testjobv1.DefaultPort,
						},
					},
					ClusterIP: corev1.ClusterIPNone,
					Selector: map[string]string{
						kubeflowv1.OperatorNameLabel: "Test Reconciler",
						kubeflowv1.JobNameLabel:      jobName,
						kubeflowv1.ReplicaTypeLabel:  strings.ToLower(string(testjobv1.TestReplicaTypeWorker)),
						kubeflowv1.ReplicaIndexLabel: idx,
					},
				},
			}
			return tc{
				testJob:         tj,
				testRType:       kubeflowv1.ReplicaType(testjobv1.TestReplicaTypeWorker),
				testSpec:        tj.Spec.TestReplicaSpecs[testjobv1.TestReplicaTypeWorker],
				testIndex:       idx,
				expectedService: svc,
			}
		}(),
	}
	testReconciler := test_job.NewTestReconciler()

	for _, c := range testCase {
		err := testReconciler.CreateNewService(c.testJob, c.testRType, c.testSpec, c.testIndex)
		if err != nil {
			t.Errorf("Got error when CreateNewService: %v", err)
			continue
		}
		var got corev1.Service
		if err = testReconciler.FC.Get(context.Background(), client.ObjectKeyFromObject(c.expectedService), &got, &client.GetOptions{}); err != nil {
			if apierrors.IsNotFound(err) {
				t.Errorf("Cannot find Service %s/%s created", c.expectedService.GetNamespace(), c.expectedService.GetName())
			}
			t.Errorf("Got error when Get service: %v", err)
		}
		if !reflect.DeepEqual(c.expectedService.Spec, got.Spec) {
			t.Errorf("Spec mismatch for service %s/%s", c.expectedService.GetNamespace(), c.expectedService.GetName())
		}
	}
}
