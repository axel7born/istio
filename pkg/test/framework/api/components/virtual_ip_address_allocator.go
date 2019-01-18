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

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/ids"
)

// Ingress represents a deployed Ingress Gateway instance.
type VirtualIPAddressAllocator interface {
	component.Instance
	// Allocate a new IP address
	AllocateIPAddress(port int, name string) (string, error)
	// Allocate a new IP address
	AllocateIPAddressOrFail(port int, name string, t testing.TB) string
}


// GetIngress from the repository
func GetVirtualIPAddressAllocator(e component.Repository, t testing.TB) VirtualIPAddressAllocator {
	return e.GetComponentOrFail(ids.VirtualIPAddressAllocator, t).(VirtualIPAddressAllocator)
}
