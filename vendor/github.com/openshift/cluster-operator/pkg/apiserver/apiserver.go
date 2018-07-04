/*
Copyright 2017 The Kubernetes Authors.

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

package apiserver

import (
	cov1alpha1 "github.com/openshift/cluster-operator/pkg/apis/clusteroperator/v1alpha1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	cav1alpha1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// ClusterOperatorAPIServer contains the base GenericAPIServer along with other
// configured runtime configuration
type ClusterOperatorAPIServer struct {
	GenericAPIServer *genericapiserver.GenericAPIServer
}

// PrepareRun prepares s to run. The returned value represents the runnable server
func (s ClusterOperatorAPIServer) PrepareRun() RunnableServer {
	return s.GenericAPIServer.PrepareRun()
}

// DefaultAPIResourceConfigSource returns a default API Resource config source
func DefaultAPIResourceConfigSource() *serverstorage.ResourceConfig {
	ret := serverstorage.NewResourceConfig()
	ret.EnableVersions(
		cov1alpha1.SchemeGroupVersion,
		cav1alpha1.SchemeGroupVersion,
	)

	return ret
}
