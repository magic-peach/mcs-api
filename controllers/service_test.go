/*
Copyright 2020 The Kubernetes Authors.

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

package controllers

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/mcs-api/pkg/apis/v1beta1"
)

func TestServiceImportOwner(t *testing.T) {
	cases := map[string]struct {
		refs []metav1.OwnerReference
		want string
	}{
		"no owner references": {},
		"unrelated owner": {
			refs: []metav1.OwnerReference{{APIVersion: "v1", Kind: "ReplicaSet", Name: "other"}},
		},
		"wrong api version": {
			refs: []metav1.OwnerReference{{APIVersion: "v1", Kind: serviceImportKind, Name: "import"}},
		},
		"serviceimport owner": {
			refs: []metav1.OwnerReference{{APIVersion: v1beta1.GroupVersion.String(), Kind: serviceImportKind, Name: "import"}},
			want: "import",
		},
		"serviceimport owner among others": {
			refs: []metav1.OwnerReference{
				{APIVersion: "v1", Kind: "ReplicaSet", Name: "other"},
				{APIVersion: v1beta1.GroupVersion.String(), Kind: serviceImportKind, Name: "import"},
			},
			want: "import",
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := serviceImportOwner(c.refs); got != c.want {
				t.Errorf("serviceImportOwner() = %q, want %q", got, c.want)
			}
		})
	}
}
