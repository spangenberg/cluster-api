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

package validation

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/cluster-operator/pkg/apis/clusteroperator"
)

// getValidClusterDeployment gets a cluster deployment that passes all validity checks.
func getValidClusterDeployment() *clusteroperator.ClusterDeployment {
	return &clusteroperator.ClusterDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-cluster",
		},
		Spec: getValidClusterDeploymentSpec(),
	}
}

func getValidClusterDeploymentSpec() clusteroperator.ClusterDeploymentSpec {
	return clusteroperator.ClusterDeploymentSpec{
		ClusterID: "cluster-id",
		MachineSets: []clusteroperator.ClusterMachineSet{
			{
				MachineSetConfig: clusteroperator.MachineSetConfig{
					NodeType: clusteroperator.NodeTypeMaster,
					Infra:    true,
					Size:     1,
				},
			},
		},
		ClusterVersionRef: clusteroperator.ClusterVersionReference{
			Name: "v3-9",
		},
	}
}

func getClusterVersionReference() corev1.ObjectReference {
	return corev1.ObjectReference{
		Namespace: "openshift-cluster-operator",
		Name:      "v3-9",
		UID:       "fakeuid",
	}
}

// getTestMachineSet gets a ClusterMachineSet initialized with either compute or master node type
func getTestMachineSet(size int, shortName string, master bool, infra bool) clusteroperator.ClusterMachineSet {
	nodeType := clusteroperator.NodeTypeCompute
	if master {
		nodeType = clusteroperator.NodeTypeMaster
	}
	return clusteroperator.ClusterMachineSet{
		ShortName: shortName,
		MachineSetConfig: clusteroperator.MachineSetConfig{
			NodeType: nodeType,
			Size:     size,
			Infra:    infra,
		},
	}
}

// TestValidateClusterDeployment tests the ValidateCluster function.
func TestValidateClusterDeployment(t *testing.T) {
	cases := []struct {
		name              string
		clusterDeployment *clusteroperator.ClusterDeployment
		valid             bool
	}{
		{
			name:              "valid",
			clusterDeployment: getValidClusterDeployment(),
			valid:             true,
		},
		{
			name: "invalid name",
			clusterDeployment: func() *clusteroperator.ClusterDeployment {
				c := getValidClusterDeployment()
				c.Name = "###"
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid spec",
			clusterDeployment: func() *clusteroperator.ClusterDeployment {
				c := getValidClusterDeployment()
				c.Spec.MachineSets[0].Size = 0
				return c
			}(),
			valid: false,
		},
		{
			name: "invalid status",
			clusterDeployment: func() *clusteroperator.ClusterDeployment {
				c := getValidClusterDeployment()
				c.Status.MachineSetCount = -1
				return c
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		errs := ValidateClusterDeployment(tc.clusterDeployment)
		if len(errs) != 0 && tc.valid {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
			continue
		} else if len(errs) == 0 && !tc.valid {
			t.Errorf("%v: unexpected success", tc.name)
		}
	}
}

// TestValidateClusterDeploymentUpdate tests the ValidateClusterDeploymentUpdate function.
func TestValidateClusterDeploymentUpdate(t *testing.T) {
	cases := []struct {
		name  string
		old   *clusteroperator.ClusterDeployment
		new   *clusteroperator.ClusterDeployment
		valid bool
	}{
		{
			name:  "valid",
			old:   getValidClusterDeployment(),
			new:   getValidClusterDeployment(),
			valid: true,
		},
		{
			name: "invalid spec",
			old:  getValidClusterDeployment(),
			new: func() *clusteroperator.ClusterDeployment {
				c := getValidClusterDeployment()
				c.Spec.MachineSets[0].Size = 0
				return c
			}(),
			valid: false,
		},
		{
			name: "mutated clusterID",
			old:  getValidClusterDeployment(),
			new: func() *clusteroperator.ClusterDeployment {
				c := getValidClusterDeployment()
				c.Spec.ClusterID = "mutated-cluster-id"
				return c
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		errs := ValidateClusterDeploymentUpdate(tc.new, tc.old)
		if len(errs) != 0 && tc.valid {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
			continue
		} else if len(errs) == 0 && !tc.valid {
			t.Errorf("%v: unexpected success", tc.name)
		}
	}
}

// TestValidateClusterDeploymentStatusUpdate tests the ValidateClusterDeploymentStatusUpdate function.
func TestValidateClusterDeploymentStatusUpdate(t *testing.T) {
	cases := []struct {
		name  string
		old   *clusteroperator.ClusterDeployment
		new   *clusteroperator.ClusterDeployment
		valid bool
	}{
		{
			name:  "valid",
			old:   getValidClusterDeployment(),
			new:   getValidClusterDeployment(),
			valid: true,
		},
		{
			name: "invalid status",
			old:  getValidClusterDeployment(),
			new: func() *clusteroperator.ClusterDeployment {
				c := getValidClusterDeployment()
				c.Status.MachineSetCount = -1
				return c
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		errs := ValidateClusterDeploymentStatusUpdate(tc.new, tc.old)
		if len(errs) != 0 && tc.valid {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
			continue
		} else if len(errs) == 0 && !tc.valid {
			t.Errorf("%v: unexpected success", tc.name)
		}
	}
}

// TestValidateClusterDeploymentSpec tests the validateClusterDeploymentSpec function.
func TestValidateClusterDeploymentSpec(t *testing.T) {
	cases := []struct {
		name  string
		spec  *clusteroperator.ClusterDeploymentSpec
		valid bool
	}{
		{
			name: "missing clusterID",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.ClusterID = ""
				return &cs
			}(),
			valid: false,
		},
		{
			name: "valid master only",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				return &cs
			}(),
			valid: true,
		},
		{
			name: "invalid master size",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(0, "", true, true),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "valid single compute",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, false),
					getTestMachineSet(1, "one", false, true),
				}
				return &cs
			}(),
			valid: true,
		},
		{
			name: "valid multiple computes",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, true),
					getTestMachineSet(1, "one", false, false),
					getTestMachineSet(5, "two", false, false),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: true,
		},
		{
			name: "invalid compute name",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, true),
					getTestMachineSet(1, "one", false, false),
					getTestMachineSet(5, "", false, false),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "invalid compute size",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, true),
					getTestMachineSet(1, "one", false, false),
					getTestMachineSet(0, "two", false, false),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "invalid duplicate compute name",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, true),
					getTestMachineSet(1, "one", false, false),
					getTestMachineSet(5, "one", false, false),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "no master machineset",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "one", false, true),
					getTestMachineSet(5, "two", false, false),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "no infra machineset",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, false),
					getTestMachineSet(1, "one", false, false),
					getTestMachineSet(5, "one", false, false),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "more than one master",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, true),
					getTestMachineSet(1, "", true, false),
					getTestMachineSet(5, "one", false, false),
					getTestMachineSet(2, "two", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "more than one infra",
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.MachineSets = []clusteroperator.ClusterMachineSet{
					getTestMachineSet(1, "", true, false),
					getTestMachineSet(1, "one", false, true),
					getTestMachineSet(5, "two", false, true),
					getTestMachineSet(2, "three", false, false),
				}
				return &cs
			}(),
			valid: false,
		},
		{
			name: "missing cluster version name", // namespace is optional
			spec: func() *clusteroperator.ClusterDeploymentSpec {
				cs := getValidClusterDeploymentSpec()
				cs.ClusterVersionRef.Name = ""
				return &cs
			}(),
			valid: false,
		},
	}

	for _, tc := range cases {
		errs := validateClusterDeploymentSpec(tc.spec, field.NewPath("spec"))
		if len(errs) != 0 && tc.valid {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
			continue
		} else if len(errs) == 0 && !tc.valid {
			t.Errorf("%v: unexpected success", tc.name)
		}
	}
}

// TestValidateClusterDeploymentStatus tests the validateClusterDeploymentStatus function.
func TestValidateClusterDeploymentStatus(t *testing.T) {
	cases := []struct {
		name   string
		status *clusteroperator.ClusterDeploymentStatus
		valid  bool
	}{
		{
			name:   "empty",
			status: &clusteroperator.ClusterDeploymentStatus{},
			valid:  true,
		},
		{
			name: "positive machinesets",
			status: &clusteroperator.ClusterDeploymentStatus{
				MachineSetCount: 1,
			},
			valid: true,
		},
		{
			name: "negative machinesets",
			status: &clusteroperator.ClusterDeploymentStatus{
				MachineSetCount: -1,
			},
			valid: false,
		},
	}

	for _, tc := range cases {
		errs := validateClusterDeploymentStatus(tc.status, field.NewPath("status"))
		if len(errs) != 0 && tc.valid {
			t.Errorf("%v: unexpected error: %v", tc.name, errs)
			continue
		} else if len(errs) == 0 && !tc.valid {
			t.Errorf("%v: unexpected success", tc.name)
		}
	}
}
