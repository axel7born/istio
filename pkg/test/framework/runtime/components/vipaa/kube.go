package vipaa

import (
	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/context"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/api"
	"istio.io/istio/pkg/test/kube"
	"testing"
	kubeEnv "istio.io/istio/pkg/test/framework/runtime/components/environment/kube"

	"github.com/google/uuid"
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
  - port: 5555
    protocol: TCP
    targetPort: 5555
  sessionAffinity: None
  type: ClusterIP
status:
  loadBalancer: {}
`
)

var (

	_ components.VirtualIPAddressAllocator = &kubeVipaa{}
	_ api.Component      = &kubeVipaa{}
)

type kubeVipaa struct {
	scope   lifecycle.Scope
	accessor *kube.Accessor
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
	v.namespace = env.TestNamespace()
	return  nil
}


func (v *kubeVipaa) AllocateIPAddress(port int) (string, error) {
	name := "kube-" + uuid.New().String()
	_, err := v.accessor.ApplyContents(v.namespace,tmpl.EvaluateOrFail(service, map[string]interface{}{"name": name, "namespace": v.namespace, "port": port}, nil))
	if err != nil {
		return "", err
	}
	service, err := v.accessor.GetService(v.namespace, name)
	if err != nil {
		return "", err
	}
	return service.Spec.ClusterIP, nil
}

func (v *kubeVipaa) AllocateIPAddressOrFail(port int, t testing.TB) string {
	ip, err := v.AllocateIPAddress(port)
	if err != nil {
		t.Fatal(err)
	}
	return ip
}
