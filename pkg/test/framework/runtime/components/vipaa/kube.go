package vipaa

import (
	"testing"

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/context"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/api"
	kubeEnv "istio.io/istio/pkg/test/framework/runtime/components/environment/kube"
	"istio.io/istio/pkg/test/kube"

	"istio.io/istio/pkg/test/framework/tmpl"
)

const (
	service = `---
apiVersion: v1
kind: Service
metadata:
  name: {{ .name }}
  namespace: {{ .namespace }}
spec:
  ports:
  - port: {{ .port }}
    protocol: TCP
    targetPort: {{ .port }}
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
`
)

var (
	_ components.VirtualIPAddressAllocator = &kubeVipaa{}
	_ api.Component                        = &kubeVipaa{}
)

func NewKubeComponent() (api.Component, error) {
	return &kubeVipaa{}, nil
}

type kubeVipaa struct {
	scope     lifecycle.Scope
	accessor  *kube.Accessor
	namespace string
}

func (v *kubeVipaa) Descriptor() component.Descriptor {
	return descriptors.VirtualIPAddressAllocator
}

func (v *kubeVipaa) Scope() lifecycle.Scope {
	return v.scope
}

func (v *kubeVipaa) Start(ctx context.Instance, scope lifecycle.Scope) (err error) {
	v.scope = scope
	env, err := kubeEnv.GetEnvironment(ctx)
	if err != nil {
		return err
	}

	v.accessor = env.Accessor
	v.namespace = env.SuiteNamespace()
	return nil
}

func (v *kubeVipaa) AllocateIPAddress(port int, name string) (string, error) {
	_, err := v.accessor.ApplyContents(v.namespace, tmpl.EvaluateOrFail(service, map[string]interface{}{"name": name, "namespace": v.namespace, "port": port}, nil))
	if err != nil {
		return "", err
	}
	service, err := v.accessor.GetService(v.namespace, name)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}

func (v *kubeVipaa) AllocateIPAddressOrFail(port int, name string, t testing.TB) string {
	ip, err := v.AllocateIPAddress(port, name)
	if err != nil {
		t.Fatal(err)
	}
	return ip
}
