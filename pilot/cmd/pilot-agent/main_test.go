package main

import (
	"github.com/onsi/gomega"
	meshconfig "istio.io/api/mesh/v1alpha1"
	"istio.io/istio/pilot/pkg/serviceregistry"
	"os"
	"testing"
)

func TestNoPilotSanIfAuthenticationNone(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	role.IdentityDomain = ""
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_NONE.String()

<<<<<<< HEAD
	pilotSAN := getPilotSAN(role.Domain, "anything")
||||||| merged common ancestors
	pilotSAN := determinePilotSAN("anything")

=======
	pilotSAN, _ := determinePilotSanAndDomain("anything")

>>>>>>> Add parameter --identity-domain to start pilot
	g.Expect(pilotSAN).To(gomega.BeNil())
}

<<<<<<< HEAD
func TestPilotSanIfAuthenticationMutualDomainEmptyKubernetes(t *testing.T) {
||||||| merged common ancestors
func TestNoPilotSanIfAuthentificationMutualDomainEmptyKubernetes(t *testing.T) {
=======
func TestPilotSanIfAuthentificationMutualDomainEmptyKubernetes(t *testing.T) {
>>>>>>> Add parameter --identity-domain to start pilot
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	role.IdentityDomain = ""
	registry = serviceregistry.KubernetesRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

<<<<<<< HEAD
	pilotSAN:= getPilotSAN(role.Domain,"anything")
||||||| merged common ancestors
	pilotSAN := determinePilotSAN("anything")
=======
	pilotSAN, _ := determinePilotSanAndDomain("anything")
>>>>>>> Add parameter --identity-domain to start pilot

	g.Expect(pilotSAN).To(gomega.Equal("cluster.local" ))
}

<<<<<<< HEAD
func TestPilotSanIfAuthenticationMutualDomainNotEmptyKubernetes(t *testing.T) {
||||||| merged common ancestors
func TestNoPilotSanIfAuthentificationMutualDomainNotEmptyKubernetes(t *testing.T) {
=======
func TestPilotSanIfAuthentificationMutualDomainNotEmptyKubernetes(t *testing.T) {
>>>>>>> Add parameter --identity-domain to start pilot
	g := gomega.NewGomegaWithT(t)
	role.Domain = "my.domain"
	role.IdentityDomain = ""
	registry = serviceregistry.KubernetesRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

<<<<<<< HEAD
	pilotSAN:= getPilotSAN(role.Domain,"anything")
||||||| merged common ancestors
	pilotSAN := determinePilotSAN("anything")
=======
	pilotSAN, _ := determinePilotSanAndDomain("anything")
>>>>>>> Add parameter --identity-domain to start pilot

	g.Expect(pilotSAN).To(gomega.Equal("my.domain"))
}

//TODO Is this really correct?
<<<<<<< HEAD
func TestPilotSanIfAuthenticationMutualDomainEmptyConsul(t *testing.T) {
||||||| merged common ancestors
func TestNoPilotSanIfAuthentificationMutualDomainEmptyConsul(t *testing.T) {
=======
func TestPilotSanIfAuthentificationMutualDomainEmptyConsul(t *testing.T) {
>>>>>>> Add parameter --identity-domain to start pilot
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	role.IdentityDomain = ""
	registry = serviceregistry.ConsulRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

<<<<<<< HEAD
	pilotSAN:= getPilotSAN(role.Domain,"anything")
||||||| merged common ancestors
	pilotSAN := determinePilotSAN("anything")
=======
	pilotSAN, _ := determinePilotSanAndDomain("anything")
>>>>>>> Add parameter --identity-domain to start pilot

	g.Expect(pilotSAN).To(gomega.Equal("" ))
}

<<<<<<< HEAD
func TestPilotSanIfAuthenticationMutualIdentityDomain(t *testing.T) {
||||||| merged common ancestors
func TestNoPilotSanIfAuthentificationMutualIdentityDomain(t *testing.T) {
=======
func TestPilotSanIfAuthentificationMutualIdentityDomain(t *testing.T) {
>>>>>>> Add parameter --identity-domain to start pilot
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	role.IdentityDomain = "secured"
	registry = serviceregistry.KubernetesRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

<<<<<<< HEAD
	pilotSAN:= getPilotSAN(role.Domain,"anything")
||||||| merged common ancestors
	pilotSAN := determinePilotSAN("anything")
=======
	pilotSAN, _ := determinePilotSanAndDomain("anything")
>>>>>>> Add parameter --identity-domain to start pilot

	g.Expect(pilotSAN).To(gomega.Equal("secured" ))
}

<<<<<<< HEAD
func TestPilotSanIfAuthenticationMutualIdentityDomainAndDomain(t *testing.T) {
||||||| merged common ancestors
func TestNoPilotSanIfAuthentificationMutualIdentityDomainAndDomain(t *testing.T) {
=======
func TestPilotSanIfAuthentificationMutualIdentityDomainAndDomain(t *testing.T) {
>>>>>>> Add parameter --identity-domain to start pilot
	g := gomega.NewGomegaWithT(t)
	role.Domain = "my.domain"
	role.IdentityDomain = "secured"
	registry = serviceregistry.KubernetesRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

<<<<<<< HEAD
	pilotSAN:= getPilotSAN(role.Domain,"anything")

	g.Expect(pilotSAN).To(gomega.Equal("secured" ))
}

func TestPilotDefaultDomainKubernetes(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	registry = serviceregistry.KubernetesRegistry
	os.Setenv("POD_NAMESPACE", "default")

	domain := getDomain(role.Domain)

	g.Expect(domain).To(gomega.Equal("default.svc.cluster.local"))
	os.Unsetenv("POD_NAMESPACE")
}

func TestPilotDefaultDomainConsul(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	role.IdentityDomain = ""
	registry = serviceregistry.ConsulRegistry

	domain := getDomain(role.Domain)

	g.Expect(domain).To(gomega.Equal("service.consul"))
}


func TestPilotDefaultDomainOthers(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	registry = serviceregistry.MockRegistry

	domain := getDomain(role.Domain)
||||||| merged common ancestors
	pilotSAN := determinePilotSAN("anything")
=======
	pilotSAN, _ := determinePilotSanAndDomain("anything")
>>>>>>> Add parameter --identity-domain to start pilot

<<<<<<< HEAD
	g.Expect(domain).To(gomega.Equal(""))
}

func TestPilotDomain(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = "my.domain"
	registry = serviceregistry.MockRegistry

	domain := getDomain(role.Domain)

	g.Expect(domain).To(gomega.Equal("my.domain"))
}

func TestPilotSanIfAuthenticationMutualStdDomainKubernetes(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ".svc.cluster.local"
	role.IdentityDomain = ""
	registry = serviceregistry.KubernetesRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

	pilotSAN:= getPilotSAN(role.Domain,"anything")

	g.Expect(pilotSAN).To(gomega.Equal([]string{"spiffe://cluster.local/ns/anything/sa/istio-pilot-service-account"} ))
}

func TestPilotSanIfAuthenticationMutualStdDomainConsul(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = "service.consul"
	role.IdentityDomain = ""
	registry = serviceregistry.ConsulRegistry
	controlPlaneAuthPolicy = meshconfig.AuthenticationPolicy_MUTUAL_TLS.String()

	pilotSAN:= getPilotSAN(role.Domain,"anything")

	g.Expect(pilotSAN).To(gomega.Equal([]string{"spiffe:///ns/anything/sa/istio-pilot-service-account"} ))
||||||| merged common ancestors
	g.Expect(pilotSAN).To(gomega.Equal([]string{"spiffe://secured/ns/anything/sa/istio-pilot-service-account"} ))
}
=======
	g.Expect(pilotSAN).To(gomega.Equal([]string{"spiffe://secured/ns/anything/sa/istio-pilot-service-account"} ))
}

func TestPilotDefaultDomainKubernetes(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	registry = serviceregistry.KubernetesRegistry
	os.Setenv("POD_NAMESPACE", "default")

	_, domain := determinePilotSanAndDomain("anything")

	g.Expect(domain).To(gomega.Equal("default.svc.cluster.local"))
}

func TestPilotDefaultDomainConsul(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	registry = serviceregistry.ConsulRegistry

	_, domain := determinePilotSanAndDomain("anything")

	g.Expect(domain).To(gomega.Equal("service.consul"))
}


func TestPilotDefaultDomainOthers(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = ""
	registry = serviceregistry.MockRegistry

	_, domain := determinePilotSanAndDomain("anything")

	g.Expect(domain).To(gomega.Equal(""))
}

func TestPilotDomain(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	role.Domain = "my.domain"
	registry = serviceregistry.MockRegistry

	_, domain := determinePilotSanAndDomain("anything")

	g.Expect(domain).To(gomega.Equal("my.domain"))
}

>>>>>>> Add parameter --identity-domain to start pilot
