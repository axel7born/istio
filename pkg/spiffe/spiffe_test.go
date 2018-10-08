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
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

func TestGenSpiffeURI(t *testing.T) {
	WithIdentityDomain("cluster.local", func() {
		testCases := []struct {
			namespace      string
			serviceAccount string
			expectedError  string
			expectedURI    string
		}{
			{
				serviceAccount: "sa",
				expectedError:  "namespace or service account can't be empty",
			},
			{
				namespace:     "ns",
				expectedError: "namespace or service account can't be empty",
			},
			{
				namespace:      "namespace-foo",
				serviceAccount: "service-bar",
				expectedURI:    "spiffe://cluster.local/ns/namespace-foo/sa/service-bar",
			},
		}
		for id, tc := range testCases {
			got, err := GenSpiffeURI(tc.namespace, tc.serviceAccount)
			if tc.expectedError == "" && err != nil {
				t.Errorf("teste case [%v] failed, error %v", id, tc)
			}
			if tc.expectedError != "" {
				if err == nil {
					t.Errorf("want get error %v, got nil", tc.expectedError)
				} else if !strings.Contains(err.Error(), tc.expectedError) {
					t.Errorf("want error contains %v,  got error %v", tc.expectedError, err)
				}
				continue
			}
			if got != tc.expectedURI {
				t.Errorf("unexpected subject name, want %v, got %v", tc.expectedURI, got)
			}

		}
	})
}

func TestMustGenSpiffeURI(t *testing.T) {
	WithIdentityDomain("cluster.local", func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expect that MustGenSpiffeURI panics in case of empty namespace")
			}
		}()

		MustGenSpiffeURI("", "")
	})
}

func TestMustGenSpiffeURIWitoutIdentityDomain(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	g.Expect(MustGenSpiffeURI("namespace", "account")).To(gomega.BeEmpty())
}

func TestDetermineDomainIfAuthenticationNone(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", ""}, DefaultDomain, false)
	g.Expect(domain.Identity).To(gomega.Equal(""))
}

func TestDetermineDomainKubernetes(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", ""}, KubernetesDefaultDomain, true)
	g.Expect(domain.Identity).To(gomega.Equal("cluster.local"))
}

func TestDetermineDomainIfAuthenticationMutualDomainNotEmptyKubernetes(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"my.domain", ""}, KubernetesDefaultDomain, true)
	g.Expect(domain.Identity).To(gomega.Equal("my.domain"))
}

func TestDetermineDomainIfAuthenticationMutualDomainEmptyConsul(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"my.domain", ""}, ConsulDefaultDomain, true)
	g.Expect(domain.Identity).To(gomega.Equal("my.domain"))
}

func TestDetermineDomainIfAuthenticationMutualIdentityDomain(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", ""}, ConsulDefaultDomain, true)
	g.Expect(domain.Identity).To(gomega.Equal(""))
}

func TestDetermineDomainIfAuthenticationMutualIdentityDomainAndDomain(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"my.domain", ""}, ConsulDefaultDomain, true)
	g.Expect(domain.Identity).To(gomega.Equal("my.domain"))
}

func TestDetermineDomainIfKubernetes(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"cluster.local", ""}, KubernetesDefaultDomain, true)
	g.Expect(domain.Suffix).To(gomega.Equal("cluster.local"))
	g.Expect(domain.Identity).To(gomega.Equal("cluster.local"))
}

func TestDetermineDomainWithEmptyDomain(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", ""}, KubernetesDefaultDomain, true)
	g.Expect(domain.Identity).To(gomega.Equal("cluster.local"))
	g.Expect(domain.Suffix).To(gomega.Equal(KubernetesDefaultDomain.Suffix))

}

func TestDetermineDomainMututalTLS(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", "secured"}, KubernetesDefaultDomain, true)
	g.Expect(domain.Suffix).To(gomega.Equal(KubernetesDefaultDomain.Suffix))
	g.Expect(domain.Identity).To(gomega.Equal("secured"))
}

func TestDetermineDomainConsul(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", ""}, ConsulDefaultDomain, true)
	g.Expect(domain.Suffix).To(gomega.Equal("service.consul"))
}

func TestDetermineDomainMututalTLSDefault(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"my.Domain", ""}, KubernetesDefaultDomain, true)
	g.Expect(domain.Suffix).To(gomega.Equal("my.Domain"))
}

func TestIdentityDomainNoMututalTLS(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	domain := determineDomain(Domain{"", "secured"}, KubernetesDefaultDomain, false)
	g.Expect(domain.Suffix).To(gomega.Equal(KubernetesDefaultDomain.Suffix))
}
