package tunnel

import (
	"fmt"
	"io/ioutil"

	"net/url"
	"testing"

	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/ids"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/components/environment/kube"

	"istio.io/istio/pkg/test/framework/tmpl"

	"istio.io/istio/pilot/pkg/model"
	"istio.io/istio/pkg/test/framework"
	"k8s.io/api/core/v1"
)

const (
	clientSideEgressConfig = `
---
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  labels:
    service: client
  name: istio-egressgateway-client
spec:
  selector:
    istio: egressgateway
  servers:
  - hosts:
    - {{ .sidecarSNI }}
    port:
      name: tcp-port-443
      number: 443
      protocol: TLS
    tls:
      caCertificates: /etc/certs/root-cert.pem
      mode: MUTUAL
      privateKey: /etc/certs/key.pem
      serverCertificate: /etc/certs/cert-chain.pem
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: egress-gateway-client
spec:
  gateways:
  - istio-egressgateway-client
  hosts:
  - {{ .sidecarSNI }}
  tcp:
  - match:
    - gateways:
      - istio-egressgateway-client
      port: 443
    route:
    - destination:
        host: {{ .ingressDNS }}
        port:
          number: 443
        subset: client-2
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: egressgateway-client
spec:
  host: {{ .ingressDNS }}
  configScope: PRIVATE
  subsets:
  - name: client-2
    trafficPolicy:
      portLevelSettings:
      - port:
          number: 443
        tls:
          caCertificates: /etc/istio/egressgateway-certs/ca.crt
          clientCertificate: /etc/istio/egressgateway-certs/client.crt
          mode: MUTUAL
          privateKey: /etc/istio/egressgateway-certs/client.key
          sni: {{ .ingressDNS }}
          subjectAltNames:
          - {{ .ingressDNS }}
---
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  creationTimestamp: null
  name: ingress-service-entry
spec:
  endpoints:
  - address: {{.ingressAddress}}
  hosts:
  - {{ .ingressDNS }}
  ports:
  - name: index-binding-id-7777
    number: {{.ingressPort}}
    protocol: TCP
  resolution: STATIC
`

	clientSideConfig = `
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: mesh-to-egress-client
spec:
  gateways:
  - mesh
  hosts:
  - {{ .serviceName }}
  tcp:
  - match:
    - destinationSubnets:
      - {{ .vip }}
      gateways:
      - mesh
    route:
    - destination:
        host: istio-egressgateway.{{ .systemNamespace }}.svc.cluster.local
        port:
          number: 443
        subset: client
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: sidecar-to-egress-client
spec:
  host: istio-egressgateway.{{ .systemNamespace }}.svc.cluster.local
  subsets:
  - name: client
    trafficPolicy:
      tls:
        mode: ISTIO_MUTUAL
        sni: {{ .sidecarSNI }}
`
	serverSideConfig = `---
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  creationTimestamp: null
  name: index-binding-id-gateway-tls
spec:
  selector:
    istio: ingressgateway # use Istio default gateway implementation
  servers:
  - hosts:
    - "*"
    port:
      name: tls
      number: 443
      protocol: TLS
    tls:
      caCertificates: /etc/istio/ingressgateway-certs/ca.crt
      mode: MUTUAL
      privateKey: /etc/istio/ingressgateway-certs/service.key
      serverCertificate: /etc/istio/ingressgateway-certs/service.crt
      subjectAltNames:
      - {{ .clientSAN }}
---
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  creationTimestamp: null
  name: index-binding-id-virtual-service-tls
spec:
  gateways:
  - index-binding-id-gateway-tls
  hosts:
  - {{ .ingressDNS }}
  tcp:
  - route:
    - destination:
        host: index-binding-id.service-fabrik
        port:
          number: {{ .port }}
---
apiVersion: networking.istio.io/v1alpha3
kind: ServiceEntry
metadata:
  creationTimestamp: null
  name: index-binding-id-service-entry
spec:
  endpoints:
  - address: {{ .address }}
  hosts:
  - index-binding-id.service-fabrik
  ports:
  - name: index-binding-id-7777
    number: {{ .port }}
    protocol: TCP
  resolution: STATIC

`
)

func TestTunnel(t *testing.T) {
	ctx := framework.GetContext(t)
	ctx.RequireOrSkip(t, lifecycle.Suite, &descriptors.KubernetesEnvironment, &ids.Egress, &ids.Ingress, &ids.Apps, &ids.VirtualIPAddressAllocator)

	egress := components.GetEgress(ctx, t)

	_, err := egress.ConfigureSecretAndWaitForExistence(&v1.Secret{
		Data: map[string][]byte{
			"ca.crt":     readFileOrFail("certs/ca.crt", t),
			"client.crt": readFileOrFail("certs/client.crt", t),
			"client.key": readFileOrFail("certs/client.key", t),
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	ingress := components.GetIngress(ctx, t)

	_, err = ingress.ConfigureSecretAndWaitForExistence(&v1.Secret{
		Data: map[string][]byte{
			"ca.crt":      readFileOrFail("certs/ca.crt", t),
			"service.crt": readFileOrFail("certs/service.crt", t),
			"service.key": readFileOrFail("certs/service.key", t),
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	apps := components.GetApps(ctx, t)
	a := apps.GetAppOrFail("a", t)
	b := apps.GetAppOrFail("t", t)

	be := b.EndpointsForProtocol(model.ProtocolHTTP)[0]

	ingressURL, err := ingress.URL(model.ProtocolHTTPS)
	if err != nil {
		t.Fatal(err)
	}

	ingressPort := ingressURL.Port()

	vipaa := components.GetVirtualIPAddressAllocator(ctx, t)
	beURL := be.URL()
	virtualPort := 5555
	serviceName := "client"
	virtualIP := vipaa.AllocateIPAddressOrFail(virtualPort, serviceName, t)
	env := kube.GetEnvironmentOrFail(ctx, t)

	_, err = env.ApplyContents(env.SystemNamespace(),
		dump(tmpl.EvaluateOrFail(clientSideEgressConfig, map[string]interface{}{
			"vip":             virtualIP,
			"serviceName":     serviceName,
			"ingressAddress":  ingressURL.Hostname(),
			"ingressPort":     ingressPort,
			"ingressDNS":      "service.istio.test.local", // Must match CN in certs/server.crt
			"sidecarSNI":      "sni.of.destination.rule.in.sidecar",
			"systemNamespace": env.SystemNamespace(),
		}, t)))

	if err != nil {
		t.Fatal(err)
	}

	_, err = env.ApplyContents(env.SuiteNamespace(),
		dump(tmpl.EvaluateOrFail(clientSideConfig, map[string]interface{}{
			"vip":             virtualIP,
			"serviceName":     serviceName,
			"ingressAddress":  ingressURL.Hostname(),
			"ingressPort":     ingressPort,
			"ingressDNS":      "service.istio.test.local", // Must match CN in certs/server.crt
			"sidecarSNI":      "sni.of.destination.rule.in.sidecar",
			"systemNamespace": env.SystemNamespace(),
		}, t)),
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = env.ApplyContents(env.TestNamespace(),
		dump(tmpl.EvaluateOrFail(serverSideConfig, map[string]interface{}{
			"address":    b.Service().ClusterIP(),
			"port":       beURL.Port(),
			"ingressDNS": "service.istio.test.local", // Must match CN in certs/server.crt
			"clientSAN":  "client.istio.test.local",  // Must match CN and SAN in certs/client.crt
		}, t)),
	)

	if err != nil {
		t.Fatal(err)
	}
	tunnelURL := &url.URL{Host: fmt.Sprintf("%s:%d", virtualIP, virtualPort), Path: beURL.Path, Scheme: beURL.Scheme}
	result := a.CallURLOrFail(tunnelURL, b, components.AppCallOptions{}, t)[0]

	//result := a.CallOrFail( be, components.AppCallOptions{}, t)[0]

	if !result.IsOK() {
		t.Fatalf("HTTP Request unsuccessful: %s", result.Body)
	}
}

func dump(yaml string) string {
	fmt.Println(yaml)
	return yaml
}

func readFileOrFail(filename string, t testing.TB) []byte {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	return content
}
