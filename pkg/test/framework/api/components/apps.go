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

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/application/echo"
	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/ids"
)

// AppProtocol enumerates the protocol options for calling an DeployedAppEndpoint endpoint.
type AppProtocol string

const (
	// AppProtocolHTTP calls the app with HTTP
	AppProtocolHTTP = "http"
	// AppProtocolGRPC calls the app with GRPC
	AppProtocolGRPC = "grpc"
	// AppProtocolWebSocket calls the app with WebSocket
	AppProtocolWebSocket = "ws"
)

// Apps is a component that provides access to all deployed test services.
type Apps interface {
	component.Instance
	GetApp(name string) (App, error)
	GetAppOrFail(name string, t testing.TB) App
}

// Service represents a deployed service within the mesh.
type Service interface {
	ClusterIP() string
}

// App represents a deployed fake App within the mesh.
type App interface {
	Name() string
	Endpoints() []AppEndpoint
	EndpointsForProtocol(protocol model.Protocol) []AppEndpoint
	Call(e AppEndpoint, opts AppCallOptions) ([]*echo.ParsedResponse, error)
	CallOrFail(e AppEndpoint, opts AppCallOptions, t testing.TB) []*echo.ParsedResponse
	Service() Service
	CallURL(url *url.URL, dst App, opts AppCallOptions) ([]*echo.ParsedResponse, error)
	CallURLOrFail(url *url.URL, dst App, opts AppCallOptions, t testing.TB) []*echo.ParsedResponse
}

// AppCallOptions defines options for calling a DeployedAppEndpoint.
type AppCallOptions struct {
	// Count indicates the number of exchanges that should be made with the service endpoint. If not set (i.e. 0), defaults to 1.
	Count int
}

// AppEndpoint represents a single endpoint in a DeployedApp.
type AppEndpoint interface {
	URL() *url.URL
	Name() string
	Owner() App
	Protocol() model.Protocol
}

// GetApps from the repository
func GetApps(e component.Repository, t testing.TB) Apps {
	return e.GetComponentOrFail("", ids.Apps, t).(Apps)
}
