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
	"regexp"
	"testing"

	"k8s.io/apimachinery/pkg/types"
)

func TestDerivedName(t *testing.T) {
	a := types.NamespacedName{Namespace: "ns1", Name: "svc-a"}
	b := types.NamespacedName{Namespace: "ns1", Name: "svc-b"}
	c := types.NamespacedName{Namespace: "ns2", Name: "svc-a"}

	if derivedName(a) != derivedName(a) {
		t.Error("derivedName is not deterministic for the same input")
	}
	if derivedName(a) == derivedName(b) {
		t.Error("expected different names to produce different derived names")
	}
	if derivedName(a) == derivedName(c) {
		t.Error("expected different namespaces to produce different derived names")
	}

	want := regexp.MustCompile(`^derived-[0-9a-v]{10}$`)
	if got := derivedName(a); !want.MatchString(got) {
		t.Errorf("derivedName() = %q, want a match for %s", got, want)
	}
}
