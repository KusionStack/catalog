package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestNetworkModule_Generator(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	testcases := []struct {
		name                 string
		devModuleConfig      kusionapiv1.Accessory
		platformModuleConfig kusionapiv1.GenericConfig
		expectedErr          error
	}{
		{
			name: "Generate private port service",
			devModuleConfig: kusionapiv1.Accessory{
				"ports": []interface{}{
					map[string]any{
						"port":     8080,
						"protocol": "TCP",
					},
				},
			},
			platformModuleConfig: nil,
			expectedErr:          nil,
		},
		{
			name: "Generate public port service",
			devModuleConfig: kusionapiv1.Accessory{
				"ports": []interface{}{
					map[string]any{
						"port":     8080,
						"public":   true,
						"protocol": "TCP",
					},
				},
			},
			platformModuleConfig: kusionapiv1.GenericConfig{
				"port": map[string]any{
					"type": "alicloud",
					"annotations": map[string]string{
						"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-spec": "slb.s1.small",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		network := &Network{}
		t.Run(tc.name, func(t *testing.T) {
			r.DevConfig = tc.devModuleConfig
			r.PlatformConfig = tc.platformModuleConfig

			res, err := network.Generate(context.Background(), r)
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestNetworkModule_GetCompleteConfig(t *testing.T) {
	testcases := []struct {
		name                 string
		devModuleConfig      kusionapiv1.Accessory
		platformModuleConfig kusionapiv1.GenericConfig
		expectedErr          error
	}{
		{
			name: "Generate private port service",
			devModuleConfig: kusionapiv1.Accessory{
				"ports": []interface{}{
					map[string]any{
						"port":     8080,
						"protocol": "TCP",
					},
				},
			},
			platformModuleConfig: nil,
			expectedErr:          nil,
		},
		{
			name: "Generate public port service",
			devModuleConfig: kusionapiv1.Accessory{
				"ports": []interface{}{
					map[string]any{
						"port":     8080,
						"public":   true,
						"protocol": "TCP",
					},
				},
			},
			platformModuleConfig: kusionapiv1.GenericConfig{
				"port": map[string]any{
					"type": "alicloud",
					"annotations": map[string]string{
						"service.beta.kubernetes.io/alibaba-cloud-loadbalancer-spec": "slb.s1.small",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		network := &Network{}
		t.Run(tc.name, func(t *testing.T) {
			err := network.GetCompleteConfig(tc.devModuleConfig, tc.platformModuleConfig)
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNetworkModule_Validate(t *testing.T) {
	testcases := []struct {
		name        string
		network     *Network
		expectedErr error
	}{
		{
			name: "Invalid port",
			network: &Network{
				Ports: []Port{
					{
						Port: 0,
					},
				},
			},
			expectedErr: ErrInvalidPort,
		},
		{
			name: "Invalid target port",
			network: &Network{
				Ports: []Port{
					{
						Port:       80,
						TargetPort: 0,
					},
				},
			},
			expectedErr: ErrInvalidTargetPort,
		},
		{
			name: "Invalid protocol",
			network: &Network{
				Ports: []Port{
					{
						Port:       80,
						TargetPort: 80,
						Protocol:   "InvalidProtocol",
					},
				},
			},
			expectedErr: ErrInvalidProtocol,
		},
		{
			name: "Valid port",
			network: &Network{
				Ports: []Port{
					{
						Port:       80,
						TargetPort: 80,
						Protocol:   "TCP",
					},
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.network.Validate()
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNetworkModule_GenerateIngressResource(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "network.Network",
		},
	}

	ingresClassName := "nginx-example"
	testcases := []struct {
		name        string
		ingress     Ingress
		expectedErr bool
	}{
		{
			name: "Generate ingress resources",
			ingress: Ingress{
				Annotations: map[string]string{
					"nginx.ingress.kubernetes.io/rewrite-target": "/",
				},
				IngressClassName: &ingresClassName,
				DefaultBackend: &IngressBackend{
					Service: &IngressServiceBackend{
						Name: "test",
						Port: ServiceBackendPort{
							Number: 80,
						},
					},
				},
				TLS: []IngressTLS{
					{
						Hosts:      []string{"https-example.foo.com"},
						SecretName: "testsecret-tls",
					},
				},
				Rules: []IngressRule{
					{
						Host: "foo.bar.com",
						HTTP: &HTTPIngressRuleValue{
							Paths: []HTTPIngressPath{
								{
									Path:     "/foo",
									PathType: "Prefix",
									Backend: IngressBackend{
										Service: &IngressServiceBackend{
											Name: "service1",
											Port: ServiceBackendPort{
												Number: 4200,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		network := &Network{
			Ingress: &tc.ingress,
		}
		t.Run(tc.name, func(t *testing.T) {
			_, err := network.GenerateIngressResource(r)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}

func TestNetworkModule_GenerateIngressClassResource(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "network.Network",
		},
	}

	apiGroup := "k8s.example.net"
	testcases := []struct {
		name         string
		ingressClass IngressClass
		expectedErr  bool
	}{
		{
			name: "Generate ingressClass resources",
			ingressClass: IngressClass{
				Controller: "example.com/ingress-controller",
				Parameters: &IngressClassParametersReference{
					APIGroup: &apiGroup,
					Kind:     "ClusterIngressParameter",
					Name:     "external-config-1",
				},
			},
			expectedErr: false,
		},
	}

	for _, tc := range testcases {
		network := &Network{
			IngressClass: &tc.ingressClass,
		}
		t.Run(tc.name, func(t *testing.T) {
			_, err := network.GenerateIngressClassResource(r)
			assert.Equal(t, tc.expectedErr, err != nil)
		})
	}
}
