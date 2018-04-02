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
package fake

import (
	v1alpha2 "github.com/kubeflow/tf-operator/pkg/apis/tensorflow/v1alpha2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTFJobs implements TFJobInterface
type FakeTFJobs struct {
	Fake *FakeKubeflowV1alpha2
	ns   string
}

var tfjobsResource = schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1alpha2", Resource: "tfjobs"}

var tfjobsKind = schema.GroupVersionKind{Group: "kubeflow.org", Version: "v1alpha2", Kind: "TFJob"}

// Get takes name of the tFJob, and returns the corresponding tFJob object, and an error if there is any.
func (c *FakeTFJobs) Get(name string, options v1.GetOptions) (result *v1alpha2.TFJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(tfjobsResource, c.ns, name), &v1alpha2.TFJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.TFJob), err
}

// List takes label and field selectors, and returns the list of TFJobs that match those selectors.
func (c *FakeTFJobs) List(opts v1.ListOptions) (result *v1alpha2.TFJobList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(tfjobsResource, tfjobsKind, c.ns, opts), &v1alpha2.TFJobList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha2.TFJobList{}
	for _, item := range obj.(*v1alpha2.TFJobList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tFJobs.
func (c *FakeTFJobs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(tfjobsResource, c.ns, opts))

}

// Create takes the representation of a tFJob and creates it.  Returns the server's representation of the tFJob, and an error, if there is any.
func (c *FakeTFJobs) Create(tFJob *v1alpha2.TFJob) (result *v1alpha2.TFJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(tfjobsResource, c.ns, tFJob), &v1alpha2.TFJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.TFJob), err
}

// Update takes the representation of a tFJob and updates it. Returns the server's representation of the tFJob, and an error, if there is any.
func (c *FakeTFJobs) Update(tFJob *v1alpha2.TFJob) (result *v1alpha2.TFJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(tfjobsResource, c.ns, tFJob), &v1alpha2.TFJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.TFJob), err
}

// Delete takes name of the tFJob and deletes it. Returns an error if one occurs.
func (c *FakeTFJobs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(tfjobsResource, c.ns, name), &v1alpha2.TFJob{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTFJobs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(tfjobsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha2.TFJobList{})
	return err
}

// Patch applies the patch and returns the patched tFJob.
func (c *FakeTFJobs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha2.TFJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(tfjobsResource, c.ns, name, data, subresources...), &v1alpha2.TFJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha2.TFJob), err
}
