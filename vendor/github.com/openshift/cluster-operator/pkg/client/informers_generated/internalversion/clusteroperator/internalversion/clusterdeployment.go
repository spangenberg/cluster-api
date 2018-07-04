/*
Copyright 2018 The Kubernetes Authors.

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

// This file was automatically generated by informer-gen

package internalversion

import (
	clusteroperator "github.com/openshift/cluster-operator/pkg/apis/clusteroperator"
	internalclientset "github.com/openshift/cluster-operator/pkg/client/clientset_generated/internalclientset"
	internalinterfaces "github.com/openshift/cluster-operator/pkg/client/informers_generated/internalversion/internalinterfaces"
	internalversion "github.com/openshift/cluster-operator/pkg/client/listers_generated/clusteroperator/internalversion"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
	time "time"
)

// ClusterDeploymentInformer provides access to a shared informer and lister for
// ClusterDeployments.
type ClusterDeploymentInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() internalversion.ClusterDeploymentLister
}

type clusterDeploymentInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewClusterDeploymentInformer constructs a new informer for ClusterDeployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewClusterDeploymentInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredClusterDeploymentInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredClusterDeploymentInformer constructs a new informer for ClusterDeployment type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredClusterDeploymentInformer(client internalclientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Clusteroperator().ClusterDeployments(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.Clusteroperator().ClusterDeployments(namespace).Watch(options)
			},
		},
		&clusteroperator.ClusterDeployment{},
		resyncPeriod,
		indexers,
	)
}

func (f *clusterDeploymentInformer) defaultInformer(client internalclientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredClusterDeploymentInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *clusterDeploymentInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&clusteroperator.ClusterDeployment{}, f.defaultInformer)
}

func (f *clusterDeploymentInformer) Lister() internalversion.ClusterDeploymentLister {
	return internalversion.NewClusterDeploymentLister(f.Informer().GetIndexer())
}
