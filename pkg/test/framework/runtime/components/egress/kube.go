package egress

import (
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"istio.io/istio/pkg/test/framework/api/component"
	"istio.io/istio/pkg/test/framework/api/components"
	"istio.io/istio/pkg/test/framework/api/context"
	"istio.io/istio/pkg/test/framework/api/descriptors"
	"istio.io/istio/pkg/test/framework/api/lifecycle"
	"istio.io/istio/pkg/test/framework/runtime/api"
	"istio.io/istio/pkg/test/framework/runtime/components/environment/kube"
	kube2 "istio.io/istio/pkg/test/kube"
)

const (
	// Specifies how long we wait before a secret becomes existent.
	secretWaitTime = 120 * time.Second
	// Name of secret used by egress
	secretName = "istio-egressgateway-certs"
	istioLabel = "istio-egressgateway"
)

var (
	_ components.Egress = &kubeEgress{}
	_ api.Component     = &kubeEgress{}
)

type kubeEgress struct {
	scope                lifecycle.Scope
	accessor             *kube2.Accessor
	istioSystemNamespace string
}

func NewKubeComponent() (api.Component, error) {
	return &kubeEgress{}, nil
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

	err = c.accessor.WaitForFilesExistence(c.istioSystemNamespace, fmt.Sprintf("istio=%s", istioLabel), []string{"/etc/certs/cert-chain.pem"}, secretWaitTime)
	if err != nil {
		return err
	}

	return nil
}

func (c *kubeEgress) ConfigureSecretAndWaitForExistence(secret *corev1.Secret) (*corev1.Secret, error) {
	secret.Name = secretName
	secretAPI := c.accessor.GetSecret(c.istioSystemNamespace)
	_, err := secretAPI.Create(secret)
	if err != nil {
		switch t := err.(type) {
		case *errors2.StatusError:
			if t.ErrStatus.Reason == v1.StatusReasonAlreadyExists {
				_, err := secretAPI.Update(secret)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, err
		}
	}
	secret, err = c.accessor.WaitForSecretExist(secretAPI, secretName, secretWaitTime)
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for key := range secret.Data {
		files = append(files, "/etc/istio/egressgateway-certs/"+key)
	}
	err = c.accessor.WaitForFilesExistence(c.istioSystemNamespace, fmt.Sprintf("istio=%s", istioLabel), files, secretWaitTime)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
