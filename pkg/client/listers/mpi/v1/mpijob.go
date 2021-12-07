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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/kubeflow/training-operator/pkg/apis/mpi/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MPIJobLister helps list MPIJobs.
// All objects returned here must be treated as read-only.
type MPIJobLister interface {
	// List lists all MPIJobs in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.MPIJob, err error)
	// MPIJobs returns an object that can list and get MPIJobs.
	MPIJobs(namespace string) MPIJobNamespaceLister
	MPIJobListerExpansion
}

// mPIJobLister implements the MPIJobLister interface.
type mPIJobLister struct {
	indexer cache.Indexer
}

// NewMPIJobLister returns a new MPIJobLister.
func NewMPIJobLister(indexer cache.Indexer) MPIJobLister {
	return &mPIJobLister{indexer: indexer}
}

// List lists all MPIJobs in the indexer.
func (s *mPIJobLister) List(selector labels.Selector) (ret []*v1.MPIJob, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.MPIJob))
	})
	return ret, err
}

// MPIJobs returns an object that can list and get MPIJobs.
func (s *mPIJobLister) MPIJobs(namespace string) MPIJobNamespaceLister {
	return mPIJobNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MPIJobNamespaceLister helps list and get MPIJobs.
// All objects returned here must be treated as read-only.
type MPIJobNamespaceLister interface {
	// List lists all MPIJobs in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.MPIJob, err error)
	// Get retrieves the MPIJob from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.MPIJob, error)
	MPIJobNamespaceListerExpansion
}

// mPIJobNamespaceLister implements the MPIJobNamespaceLister
// interface.
type mPIJobNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all MPIJobs in the indexer for a given namespace.
func (s mPIJobNamespaceLister) List(selector labels.Selector) (ret []*v1.MPIJob, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.MPIJob))
	})
	return ret, err
}

// Get retrieves the MPIJob from the indexer for a given namespace and name.
func (s mPIJobNamespaceLister) Get(name string) (*v1.MPIJob, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("mpijob"), name)
	}
	return obj.(*v1.MPIJob), nil
}
