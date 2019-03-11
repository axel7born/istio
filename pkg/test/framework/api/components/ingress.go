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
	"net/url"
	"testing"

	v1 "k8s.io/api/core/v1"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/ids"
)

// Ingress represents a deployed Ingress Gateway instance.
type Ingress interface {
	component.Instance
	// URL returns the external URL of the ingress gateway (or the NodePort address,
	// when running under Minikube) for the given protocol
	URL(protocol model.Protocol) (*url.URL, error)

	//  Call makes an HTTP call through ingress, where the URL has the given path.
	Call(path string) (IngressCallResponse, error)

	// Configure a secret and wait for the existence
	ConfigureSecretAndWaitForExistence(secret *v1.Secret) (*v1.Secret, error)

	// Add addition secret mountpoint
	AddSecretMountPoint(path string) error
}

// IngressCallResponse is the result of a call made through Istio Ingress.
type IngressCallResponse struct {
	// Response status code
	Code int

	// Response body
	Body string
}

// GetIngress from the repository
func GetIngress(e component.Repository, t testing.TB) Ingress {
	return e.GetComponentOrFail(ids.Ingress, t).(Ingress)
}
