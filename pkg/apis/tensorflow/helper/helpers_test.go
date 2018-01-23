package helper

import (
	"reflect"
	"testing"

	"github.com/gogo/protobuf/proto"
	tfv1 "github.com/tensorflow/k8s/pkg/apis/tensorflow/v1alpha1"
	"github.com/tensorflow/k8s/pkg/util"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

func TestAddAccelertor(t *testing.T) {
	type testCase struct {
		in       *tfv1.TFJobSpec
		expected *tfv1.TFJobSpec
		config   map[string]tfv1.AcceleratorConfig
	}

	testCases := []testCase{
		// Case 1 checks that we look at requests.
		{
			in: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Replicas: proto.Int32(2),
						TFPort:   proto.Int32(10),
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
										Resources: v1.ResourceRequirements{
											Requests: map[v1.ResourceName]resource.Quantity{
												"nvidia-gpu": resource.MustParse("1"),
											},
										},
									},
								},
							},
						},
						TFReplicaType: tfv1.PS,
					},
				},
			},
			expected: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Replicas: proto.Int32(2),
						TFPort:   proto.Int32(10),
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
										Resources: v1.ResourceRequirements{
											Requests: map[v1.ResourceName]resource.Quantity{
												"nvidia-gpu": resource.MustParse("1"),
											},
										},
										VolumeMounts: []v1.VolumeMount{
											{
												Name:      "cuda-lib",
												MountPath: "/usr/local/cuda",
											},
										},
									},
								},
								Volumes: []v1.Volume{
									{
										Name: "cuda-lib",
										VolumeSource: v1.VolumeSource{
											HostPath: &v1.HostPathVolumeSource{
												Path: "/home/cuda",
											},
										},
									},
								},
							},
						},
						TFReplicaType: tfv1.PS,
					},
				},
			},
			config: map[string]tfv1.AcceleratorConfig{
				"nvidia-gpu": tfv1.AcceleratorConfig{
					Volumes: []tfv1.AcceleratorVolume{
						{
							Name:      "cuda-lib",
							HostPath:  "/home/cuda",
							MountPath: "/usr/local/cuda",
						},
					},
				},
			},
		},
		// Case 2 checks that we look at limit.
		{
			in: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Replicas: proto.Int32(2),
						TFPort:   proto.Int32(10),
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
										Resources: v1.ResourceRequirements{
											Limits: map[v1.ResourceName]resource.Quantity{
												"nvidia-gpu": resource.MustParse("1"),
											},
										},
									},
								},
							},
						},
						TFReplicaType: tfv1.PS,
					},
				},
			},
			expected: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Replicas: proto.Int32(2),
						TFPort:   proto.Int32(10),
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
										Resources: v1.ResourceRequirements{
											Limits: map[v1.ResourceName]resource.Quantity{
												"nvidia-gpu": resource.MustParse("1"),
											},
										},
										VolumeMounts: []v1.VolumeMount{
											{
												Name:      "cuda-lib",
												MountPath: "/usr/local/cuda",
											},
										},
									},
								},
								Volumes: []v1.Volume{
									{
										Name: "cuda-lib",
										VolumeSource: v1.VolumeSource{
											HostPath: &v1.HostPathVolumeSource{
												Path: "/home/cuda",
											},
										},
									},
								},
							},
						},
						TFReplicaType: tfv1.PS,
					},
				},
			},
			config: map[string]tfv1.AcceleratorConfig{
				"nvidia-gpu": tfv1.AcceleratorConfig{
					Volumes: []tfv1.AcceleratorVolume{
						{
							Name:      "cuda-lib",
							HostPath:  "/home/cuda",
							MountPath: "/usr/local/cuda",
						},
					},
				},
			},
		},
		// Case 3 no GPUs
		{
			in: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Replicas: proto.Int32(2),
						TFPort:   proto.Int32(10),
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
									},
								},
							},
						},
						TFReplicaType: tfv1.PS,
					},
				},
			},
			expected: &tfv1.TFJobSpec{
				ReplicaSpecs: []*tfv1.TFReplicaSpec{
					{
						Replicas: proto.Int32(2),
						TFPort:   proto.Int32(10),
						Template: &v1.PodTemplateSpec{
							Spec: v1.PodSpec{
								Containers: []v1.Container{
									{
										Name: "tensorflow",
									},
								},
							},
						},
						TFReplicaType: tfv1.PS,
					},
				},
			},
			config: map[string]tfv1.AcceleratorConfig{
				"nvidia-gpu": tfv1.AcceleratorConfig{
					Volumes: []tfv1.AcceleratorVolume{
						{
							Name:      "cuda-lib",
							HostPath:  "/home/cuda",
							MountPath: "/usr/local/cuda",
						},
					},
				},
			},
		},
	}

	for _, c := range testCases {
		if err := ConfigureAcceleratorsForTFJobSpec(c.in, c.config); err != nil {
			t.Errorf("ConfigureAccelerators error; %v", err)
		}
		if !reflect.DeepEqual(c.in, c.expected) {
			t.Errorf("Want\n%v; Got\n %v", util.Pformat(c.expected), util.Pformat(c.in))
		}
	}
}
