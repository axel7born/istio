package vipaa

import (
	"fmt"

	"testing"

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/context"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/api"
)

var (
	_ components.VirtualIPAddressAllocator = &kubeVipaa{}
	_ api.Component                        = &nativeVipaa{}
)

func NewNativeComponent() (api.Component, error) {
	return &nativeVipaa{}, nil
}

type nativeVipaa struct {
	scope   lifecycle.Scope
	counter int
}

func (v *nativeVipaa) Descriptor() component.Descriptor {
	return descriptors.VirtualIPAddressAllocator
}

func (v *nativeVipaa) Scope() lifecycle.Scope {
	return v.scope
}

func (v *nativeVipaa) Start(ctx context.Instance, scope lifecycle.Scope) (err error) {
	v.scope = scope
	v.counter = 1
	return nil
}

func (v *nativeVipaa) AllocateIPAddress(port int, name string) (string, error) {
	v.counter++
	return fmt.Sprintf("127.0.0.%d", v.counter), nil
}

func (v *nativeVipaa) AllocateIPAddressOrFail(port int, name string, t testing.TB) string {
	s, _ := v.AllocateIPAddress(port, name)
	return s
}
