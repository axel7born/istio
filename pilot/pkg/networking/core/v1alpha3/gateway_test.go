package v1alpha3

import (
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
	"github.com/envoyproxy/go-control-plane/envoy/api/v2/core"
	networking "istio.io/api/networking/v1alpha3"
	"testing"
)

func TestBuildGatewayListenerTLSContextPassthrough(t *testing.T) {

	tlsContext := buildGatewayListenerTLSContext(&networking.Server{}, []string{})
	if tlsContext != nil {
		t.Error("Not nil")
	}
}

func TestBuildGatewayListenerTLSContextMutial(t *testing.T) {

	tlsContext := buildGatewayListenerTLSContext(&networking.Server{Tls: &networking.Server_TLSOptions{
		Mode:              networking.Server_TLSOptions_MUTUAL,
		ServerCertificate: "test",
		PrivateKey:        "private",
		SubjectAltNames:   []string{"san1", "san2"},
		CaCertificates:    "ca.crt",
	}}, []string{})
	if tlsContext == nil {
		t.Error("nil")
	}

	if tlsContext.CommonTlsContext.TlsCertificates[0].CertificateChain.Specifier.(*core.DataSource_Filename).Filename != "test" {
		t.Error("invalid server certificate")
	}

	if tlsContext.CommonTlsContext.TlsCertificates[0].PrivateKey.Specifier.(*core.DataSource_Filename).Filename != "private" {
		t.Error("invalid private key")
	}
	if len(tlsContext.CommonTlsContext.ValidationContextType.(*auth.CommonTlsContext_ValidationContext).ValidationContext.VerifySubjectAltName) != 2 {
		t.Error("invalid subject alternate names")
	}
	if tlsContext.CommonTlsContext.ValidationContextType.(*auth.CommonTlsContext_ValidationContext).ValidationContext.TrustedCa.Specifier.(*core.DataSource_Filename).Filename != "ca.crt" {
		t.Error("invalid trust ca")
	}

}

func TestBuildGatewayListenerTLSContextIstioMutial(t *testing.T) {

	tlsContext := buildGatewayListenerTLSContext(&networking.Server{Tls: &networking.Server_TLSOptions{
		Mode: networking.Server_TLSOptions_ISTIO_MUTUAL,
	}}, []string{"spiffe://cluster.local/ns/default/sa/default"})
	if tlsContext == nil {
		t.Error("nil")
	}

	if tlsContext.CommonTlsContext.TlsCertificates[0].CertificateChain.Specifier.(*core.DataSource_Filename).Filename != "/etc/certs/cert-chain.pem" {
		t.Error("invalid server certificate")
	}

	if tlsContext.CommonTlsContext.TlsCertificates[0].PrivateKey.Specifier.(*core.DataSource_Filename).Filename != "/etc/certs/key.pem" {
		t.Error("invalid private key")
	}
	if tlsContext.CommonTlsContext.ValidationContextType.(*auth.CommonTlsContext_ValidationContext).ValidationContext.VerifySubjectAltName[0] != "spiffe://cluster.local/ns/default/sa/default" {
		t.Error("invalid subject alternate names")
	}
	if tlsContext.CommonTlsContext.ValidationContextType.(*auth.CommonTlsContext_ValidationContext).ValidationContext.TrustedCa.Specifier.(*core.DataSource_Filename).Filename != "/etc/certs/root-cert.pem" {
		t.Error("invalid trust ca")
	}

}
