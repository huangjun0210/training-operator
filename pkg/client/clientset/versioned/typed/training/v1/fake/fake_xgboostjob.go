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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	trainingv1 "github.com/kubeflow/training-operator/pkg/apis/training/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeXGBoostJobs implements XGBoostJobInterface
type FakeXGBoostJobs struct {
	Fake *FakeKubeflowV1
	ns   string
}

var xgboostjobsResource = schema.GroupVersionResource{Group: "kubeflow.org", Version: "v1", Resource: "xgboostjobs"}

var xgboostjobsKind = schema.GroupVersionKind{Group: "kubeflow.org", Version: "v1", Kind: "XGBoostJob"}

// Get takes name of the xGBoostJob, and returns the corresponding xGBoostJob object, and an error if there is any.
func (c *FakeXGBoostJobs) Get(ctx context.Context, name string, options v1.GetOptions) (result *trainingv1.XGBoostJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(xgboostjobsResource, c.ns, name), &trainingv1.XGBoostJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*trainingv1.XGBoostJob), err
}

// List takes label and field selectors, and returns the list of XGBoostJobs that match those selectors.
func (c *FakeXGBoostJobs) List(ctx context.Context, opts v1.ListOptions) (result *trainingv1.XGBoostJobList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(xgboostjobsResource, xgboostjobsKind, c.ns, opts), &trainingv1.XGBoostJobList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &trainingv1.XGBoostJobList{ListMeta: obj.(*trainingv1.XGBoostJobList).ListMeta}
	for _, item := range obj.(*trainingv1.XGBoostJobList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested xGBoostJobs.
func (c *FakeXGBoostJobs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(xgboostjobsResource, c.ns, opts))

}

// Create takes the representation of a xGBoostJob and creates it.  Returns the server's representation of the xGBoostJob, and an error, if there is any.
func (c *FakeXGBoostJobs) Create(ctx context.Context, xGBoostJob *trainingv1.XGBoostJob, opts v1.CreateOptions) (result *trainingv1.XGBoostJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(xgboostjobsResource, c.ns, xGBoostJob), &trainingv1.XGBoostJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*trainingv1.XGBoostJob), err
}

// Update takes the representation of a xGBoostJob and updates it. Returns the server's representation of the xGBoostJob, and an error, if there is any.
func (c *FakeXGBoostJobs) Update(ctx context.Context, xGBoostJob *trainingv1.XGBoostJob, opts v1.UpdateOptions) (result *trainingv1.XGBoostJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(xgboostjobsResource, c.ns, xGBoostJob), &trainingv1.XGBoostJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*trainingv1.XGBoostJob), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeXGBoostJobs) UpdateStatus(ctx context.Context, xGBoostJob *trainingv1.XGBoostJob, opts v1.UpdateOptions) (*trainingv1.XGBoostJob, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(xgboostjobsResource, "status", c.ns, xGBoostJob), &trainingv1.XGBoostJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*trainingv1.XGBoostJob), err
}

// Delete takes name of the xGBoostJob and deletes it. Returns an error if one occurs.
func (c *FakeXGBoostJobs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(xgboostjobsResource, c.ns, name, opts), &trainingv1.XGBoostJob{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeXGBoostJobs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(xgboostjobsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &trainingv1.XGBoostJobList{})
	return err
}

// Patch applies the patch and returns the patched xGBoostJob.
func (c *FakeXGBoostJobs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *trainingv1.XGBoostJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(xgboostjobsResource, c.ns, name, pt, data, subresources...), &trainingv1.XGBoostJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*trainingv1.XGBoostJob), err
}
