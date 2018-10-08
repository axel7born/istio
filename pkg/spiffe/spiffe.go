// Copyright 2018 Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spiffe

import (
	"fmt"
	"os"
)

const (
	Scheme = "spiffe"
)

var globalDomain Domain

type Domain struct {
	Suffix   string
	Identity string
}

var (
	KubernetesDefaultDomain = Domain{orDefault(os.Getenv("POD_NAMESPACE"), "default") + ".svc.cluster.local", "cluster.local"}
	ConsulDefaultDomain     = Domain{"service.consul", ""}
	DefaultDomain           = Domain{"", ""}
)

func SetDomain(domain Domain, dflt Domain, mutualTLS bool) Domain {
	globalDomain = determineDomain(domain, dflt, mutualTLS)
	return globalDomain
}

func determineDomain(domain Domain, dflt Domain, mutualTLS bool) Domain {
	if mutualTLS {
		domain.Identity = orDefault(domain.Identity, domain.Suffix)
		domain.Identity = orDefault(domain.Identity, os.Getenv("ISTIO_SA_DOMAIN_CANONICAL"))
		domain.Identity = orDefault(domain.Identity, dflt.Identity)
	} else {
		domain.Identity = ""
	}
	domain.Suffix = orDefault(domain.Suffix, dflt.Suffix)

	return domain
}

func orDefault(value string, dflt string) string {
	if len(value) == 0 {
		return dflt
	}
	return value
}

// GenSpiffeURI returns the formatted uri(SPIFFEE format for now) for the certificate.
func GenSpiffeURI(ns, serviceAccount string) (string, error) {
	if globalDomain.Identity == "" {
		return "", nil
	}
	if ns == "" || serviceAccount == "" {
		return "", fmt.Errorf(
			"namespace or service account can't be empty ns=%v serviceAccount=%v", ns, serviceAccount)
	}
	return fmt.Sprintf(Scheme+"://%s/ns/%s/sa/%s", globalDomain.Identity, ns, serviceAccount), nil
}

func MustGenSpiffeURI(ns, serviceAccount string) string {
	uri, err := GenSpiffeURI(ns, serviceAccount)
	if err != nil {
		panic(err)
	}
	return uri
}
