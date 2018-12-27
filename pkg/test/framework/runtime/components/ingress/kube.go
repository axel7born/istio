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

package ingress

import (
	"errors"
	"fmt"
	"io/ioutil"
	"istio.io/istio/pilot/pkg/model"
	"k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/url"
	"strings"
	"time"

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/context"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/api"
	"istio.io/istio/pkg/test/framework/runtime/components/environment/kube"
	kube2 "istio.io/istio/pkg/test/kube"
	"istio.io/istio/pkg/test/scopes"
	"istio.io/istio/pkg/test/util/retry"
)

const (
	serviceName = "istio-ingressgateway"
	istioLabel  = "ingressgateway"
	// Specifies how long we wait before a secret becomes existent.
	secretWaitTime = 20 * time.Second
	// Name of secret used by egress
	secretName = "istio-ingressgateway-certs"
)

var (
	retryTimeout = retry.Timeout(1 * time.Minute)
	retryDelay   = retry.Delay(5 * time.Second)

	_ components.Ingress = &kubeComponent{}
	_ api.Component      = &kubeComponent{}
)

type kubeComponent struct {
	scope                lifecycle.Scope
	url                  func(model.Protocol) (*url.URL, error)
	accessor             *kube2.Accessor
	istioSystemNamespace string
}

// NewKubeComponent factory function for the component
func NewKubeComponent() (api.Component, error) {
	return &kubeComponent{}, nil
}

func (c *kubeComponent) Descriptor() component.Descriptor {
	return descriptors.Ingress
}

func (c *kubeComponent) Scope() lifecycle.Scope {
	return c.scope
}

func (c *kubeComponent) Start(ctx context.Instance, scope lifecycle.Scope) (err error) {
	c.scope = scope

	env, err := kube.GetEnvironment(ctx)
	if err != nil {
		return err
	}

	c.accessor = env.Accessor
	c.istioSystemNamespace = env.SystemNamespace()
	address, err := retry.Do(func() (interface{}, bool, error) {

		// In Minikube, we don't have the ingress gateway. Instead we do a little bit of trickery to to get the Node
		// port.
		n := env.SystemNamespace()
		if env.MinikubeIngress() {
			pods, err := env.GetPods(n, fmt.Sprintf("istio=%s", istioLabel))
			if err != nil {
				return nil, false, err
			}
			if len(pods) == 0 {
				return nil, false, errors.New("no ingress pod found")
			}
			ip := pods[0].Status.HostIP
			if ip == "" {
				return nil, false, errors.New("no Host IP availale on the ingress node yet")
			}

			svc, err := env.Accessor.GetService(n, serviceName)
			if err != nil {
				return nil, false, err
			}

			if len(svc.Spec.Ports) == 0 {
				return nil, false, fmt.Errorf("no ports found in service: %s/%s", n, "istio-ingressgateway")
			}

			return func(protocol model.Protocol ) (*url.URL, error) {
				for _, p := range svc.Spec.Ports {
					if supportsProtocol(p.Name,protocol) {
						return &url.URL{ Scheme: strings.ToLower(string(protocol)), Host: fmt.Sprintf("%s:%d",  ip, p.NodePort)}, nil
					}
				}
				return nil, errors.New("No port found")
			}, true, nil
		}

		// Otherwise, get the load balancer IP.
		svc, err := env.Accessor.GetService(n, serviceName)
		if err != nil {
			return nil, false, err
		}

		if len(svc.Status.LoadBalancer.Ingress) == 0 || svc.Status.LoadBalancer.Ingress[0].IP == "" {
			return nil, false, fmt.Errorf("service ingress is not available yet: %s/%s", svc.Name, svc.Namespace)
		}

		ip := svc.Status.LoadBalancer.Ingress[0].IP
		return func(protocol model.Protocol ) (*url.URL, error) {
			for _, p := range svc.Spec.Ports {
				if supportsProtocol(p.Name,protocol)  {
					return &url.URL{ Scheme: strings.ToLower(string(protocol)), Host: fmt.Sprintf("%s:%d",  ip, p.Port)}, nil
				}
			}
			return nil, errors.New("No port found")
		}, true, nil
	}, retryTimeout, retryDelay)

	if err != nil {
		return err
	}

	c.url = address.(func(protocol model.Protocol) (*url.URL, error))
	return nil
}
func supportsProtocol(name string, protocol model.Protocol) bool {
	return  name == "http" && (protocol == model.ProtocolHTTP ||  protocol == model.ProtocolHTTP2 )||
		name == "http2" && ( protocol == model.ProtocolHTTP || protocol == model.ProtocolHTTP2 ) ||
		name == "https" && protocol == model.ProtocolHTTPS
}

// URL implements environment.DeployedIngress
func (c *kubeComponent) URL(protocol model.Protocol) (*url.URL, error) {
	return c.url(protocol)
}

func (c *kubeComponent) Call(path string) (components.IngressCallResponse, error) {
	client := &http.Client{
		Timeout: 1 * time.Minute,
	}

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url, err := c.url(model.ProtocolHTTP)
	if err != nil {
		return components.IngressCallResponse{}, err
	}

	url.Path = path

	scopes.Framework.Debugf("Sending call to ingress at: %s", url)

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return components.IngressCallResponse{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return components.IngressCallResponse{}, err
	}

	defer func() { _ = resp.Body.Close() }()

	var ba []byte
	ba, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		scopes.Framework.Warnf("Unable to connect to read from %s: %v", c.url, err)
		return components.IngressCallResponse{}, err
	}
	contents := string(ba)
	status := resp.StatusCode

	response := components.IngressCallResponse{
		Code: status,
		Body: contents,
	}

	// scopes.Framework.Debugf("Received response to ingress call (url: %s): %+v", url, response)

	return response, nil
}

func (a *kubeComponent) ConfigureSecretAndWaitForExistence(secret *v1.Secret) (*v1.Secret, error) {
	secret.Name = secretName
	secretApi := a.accessor.GetSecret(a.istioSystemNamespace)
	_, err := secretApi.Create(secret)
	if err != nil {
		switch t := err.(type) {
		case *errors2.StatusError:
			if t.ErrStatus.Reason == v12.StatusReasonAlreadyExists {
				_, err := secretApi.Update(secret)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, err
		}
	}
	return a.accessor.WaitForSecretExist(secretApi, secretName, secretWaitTime)

}
