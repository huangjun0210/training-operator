package helper

import (
	"fmt"

	tfv1 "github.com/tensorflow/k8s/pkg/apis/tensorflow/v1alpha1"
	"github.com/tensorflow/k8s/pkg/util"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AsOwner make OwnerReference according to the parameter
func AsOwner(tfJob *tfv1.TfJob) metav1.OwnerReference {
	trueVar := true
	// Both api.OwnerReference and metatypes.OwnerReference are combined into that.
	return metav1.OwnerReference{
		APIVersion:         tfJob.APIVersion,
		Kind:               tfJob.Kind,
		Name:               tfJob.ObjectMeta.Name,
		UID:                tfJob.ObjectMeta.UID,
		Controller:         &trueVar,
		BlockOwnerDeletion: &trueVar,
	}
}

// ConfigureAcceleratorsForTfJobSpec adds any accelerator specific configuration to the pods.
func ConfigureAcceleratorsForTfJobSpec(c *tfv1.TfJobSpec, accelerators map[string]tfv1.AcceleratorConfig) error {
	for _, r := range c.ReplicaSpecs {
		if r.Template == nil {
			return fmt.Errorf("Replica is missing Template; %v", util.Pformat(r))
		}
		for i, c := range r.Template.Spec.Containers {
			if c.Name == string(tfv1.TENSORFLOW) {
				// Identify the accelerators attached to this container.
				a := map[string]tfv1.AcceleratorConfig{}

				lists := []v1.ResourceList{c.Resources.Limits, c.Resources.Requests}
				for _, resources := range lists {
					for name, _ := range resources {

						if _, ok := accelerators[string(name)]; !ok {
							continue
						}

						// Add the expected mounts to the pods.
						a[string(name)] = accelerators[string(name)]
					}
				}

				// Add accelerator information to the pod.
				for _, config := range a {
					for _, v := range config.Volumes {
						r.Template.Spec.Volumes = append(r.Template.Spec.Volumes,
							v1.Volume{
								Name: v.Name,
								VolumeSource: v1.VolumeSource{
									HostPath: &v1.HostPathVolumeSource{
										Path: v.HostPath,
									},
								},
							})
						c.VolumeMounts = append(c.VolumeMounts, v1.VolumeMount{
							Name:      v.Name,
							MountPath: v.MountPath,
						})
					}

					for _, envVar := range config.EnvVars {
						c.Env = append(c.Env, v1.EnvVar{
							Name:  envVar.Name,
							Value: envVar.Value,
						})
					}
				}
				r.Template.Spec.Containers[i] = c
				break
			}
		}
	}
	return nil
}

// Cleanup cleans up user passed spec, e.g. defaulting, transforming fields.
// TODO: move this to admission controller
func Cleanup(c *tfv1.TfJobSpec) {
	// TODO(jlewi): Add logic to cleanup user provided spec; e.g. by filling in defaults.
	// We should have default container images so user doesn't have to provide these.
}

func CRDName() string {
	return fmt.Sprintf("%s.%s", tfv1.CRDKindPlural, tfv1.CRDGroup)
}

func scalingReason(from, to int) string {
	return fmt.Sprintf("Current cluster size: %d, desired cluster size: %d", from, to)
}
