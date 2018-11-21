package egress

import (
	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/context"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/api"
	"istio.io/istio/pkg/test/framework/runtime/components/environment/kube"
	kube2 "istio.io/istio/pkg/test/kube"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

const (
	// Specifies how long we wait before a secret becomes existent.
	secretWaitTime = 20 * time.Second
	// Name of secret used by egress
	secretName = "istio-egressgateway-certs"
)

var (

	_ components.Egress = &kubeEgress{}
	_ api.Component      = &kubeEgress{}
)

type kubeEgress struct {
	scope   lifecycle.Scope
	accessor *kube2.Accessor
	istioSystemNamespace string
}


func (c *kubeEgress) Descriptor() component.Descriptor {
	return descriptors.Egress
}

func (c *kubeEgress) Scope() lifecycle.Scope {
	return c.scope
}


func (c *kubeEgress) Start(ctx context.Instance, scope lifecycle.Scope) (err error) {
	c.scope = scope
	env, err := kube.GetEnvironment(ctx)
	if err != nil {
		return err
	}

	c.accessor = env.Accessor
	c.istioSystemNamespace = env.SystemNamespace()
	return nil
}


func (a *kubeEgress) ConfigureSecretAndWaitForExistence(secret *corev1.Secret) (*corev1.Secret, error) {
	secret.Name = secretName
	secretApi := a.accessor.GetSecret(a.istioSystemNamespace)
	_, err := secretApi.Create(secret)
	if err != nil {
		switch t := err.(type) {
		case *errors.StatusError:
			if t.ErrStatus.Reason == v1.StatusReasonAlreadyExists {
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
