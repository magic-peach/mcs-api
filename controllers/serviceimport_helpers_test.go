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

	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/mcs-api/pkg/apis/v1beta1"
)

func TestServicePorts(t *testing.T) {
	if got := servicePorts(&v1beta1.ServiceImport{}); len(got) != 0 {
		t.Errorf("no ports: got %v, want an empty slice", got)
	}

	appProtocol := "kubernetes.io/h2c"
	svcImport := &v1beta1.ServiceImport{
		Spec: v1beta1.ServiceImportSpec{
			Ports: []v1beta1.ServicePort{
				{Name: "http", Protocol: v1.ProtocolTCP, Port: 80},
				{Name: "grpc", Protocol: v1.ProtocolTCP, Port: 443, AppProtocol: &appProtocol},
			},
		},
	}

	got := servicePorts(svcImport)
	want := []v1.ServicePort{
		{Name: "http", Protocol: v1.ProtocolTCP, Port: 80},
		{Name: "grpc", Protocol: v1.ProtocolTCP, Port: 443, AppProtocol: &appProtocol},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d ports, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i].Name != want[i].Name || got[i].Protocol != want[i].Protocol || got[i].Port != want[i].Port {
			t.Errorf("port %d = %+v, want %+v", i, got[i], want[i])
		}
		if (got[i].AppProtocol == nil) != (want[i].AppProtocol == nil) {
			t.Errorf("port %d AppProtocol = %v, want %v", i, got[i].AppProtocol, want[i].AppProtocol)
		} else if got[i].AppProtocol != nil && *got[i].AppProtocol != *want[i].AppProtocol {
			t.Errorf("port %d AppProtocol = %q, want %q", i, *got[i].AppProtocol, *want[i].AppProtocol)
		}
	}
}
