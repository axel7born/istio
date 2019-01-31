//  Copyright 2018 Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package components

import (
	"testing"

	"k8s.io/api/core/v1"

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/ids"
)

// Ingress represents a deployed Ingress Gateway instance.
type Egress interface {
	component.Instance
	// Configure a secret and wait for the existence
	ConfigureSecretAndWaitForExistence(secret *v1.Secret) (*v1.Secret, error)
}

// GetIngress from the repository
func GetEgress(e component.Repository, t testing.TB) Egress {
	return e.GetComponentOrFail(ids.Egress, t).(Egress)
}
