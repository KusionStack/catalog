package main

import (
	"context"
	"errors"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	apiv1 "kusionstack.io/kusion/pkg/apis/core/v1"
	"kusionstack.io/kusion/pkg/apis/core/v1/workload"
)

func TestPostgreSQLModule_Generator(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &workload.Workload{
			Header: workload.Header{
				Type: "Service",
			},
			Service: &workload.Service{},
		},
	}

	testcases := []struct {
		name            string
		devModuleConfig apiv1.Accessory
		platformConfig  apiv1.GenericConfig
		expectedErr     error
	}{
		{
			name: "Generate local PostgreSQL database",
			devModuleConfig: apiv1.Accessory{
				"type":    "local",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"databaseName": "test-postgres",
			},
			expectedErr: nil,
		},
		{
			name: "Generate AWS PostgreSQL RDS",
			devModuleConfig: apiv1.Accessory{
				"type":    "cloud",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"cloud":          "aws",
				"size":           20,
				"instanceType":   "db.t3.micro",
				"privateRouting": false,
			},
			expectedErr: nil,
		},
		{
			name: "Generate Alicloud PostgreSQL RDS",
			devModuleConfig: apiv1.Accessory{
				"type":    "cloud",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"cloud":          "alicloud",
				"size":           20,
				"instanceType":   "postgres.n2.serverless.1c",
				"category":       "serverless_basic",
				"privateRouting": false,
				"subnetID":       "test-subnet-id",
			},
			expectedErr: nil,
		},
		{
			name: "Unsupported PostgreSQL type",
			devModuleConfig: apiv1.Accessory{
				"type":    "unsupported-type",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"databaseName": "test-postgres",
			},
			expectedErr: errors.New("unsupported postgres type"),
		},
		{
			name: "Unsupported Terraform provider type",
			devModuleConfig: apiv1.Accessory{
				"type":    "cloud",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"cloud":        "unsupported-type",
				"instanceType": "db.t3.micro",
			},
			expectedErr: errors.New("unsupported cloud provider type"),
		},
		{
			name: "Empty cloud PostgreSQL instance type",
			devModuleConfig: apiv1.Accessory{
				"type":    "cloud",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"cloud": "aws",
			},
			expectedErr: ErrEmptyInstanceTypeForCloudDB,
		},
	}

	for _, tc := range testcases {
		postgres := &PostgreSQL{}
		t.Run(tc.name, func(t *testing.T) {
			r.DevModuleConfig = tc.devModuleConfig
			r.PlatformModuleConfig = tc.platformConfig

			res, err := postgres.Generate(context.Background(), r)
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestPostgreSQLModule_GetCompleteConfig(t *testing.T) {
	testcases := []struct {
		name               string
		devModuleConfig    apiv1.Accessory
		platformConfig     apiv1.GenericConfig
		expectedPostgreSQL *PostgreSQL
	}{
		{
			name: "Empty platform config",
			devModuleConfig: apiv1.Accessory{
				"type":    "local",
				"version": "14.0",
			},
			platformConfig: nil,
			expectedPostgreSQL: &PostgreSQL{
				Type:           "local",
				Version:        "14.0",
				Username:       defaultUsername,
				Category:       defaultCategory,
				SecurityIPs:    defaultSecurityIPs,
				PrivateRouting: defaultPrivateRouting,
				Size:           defaultSize,
			},
		},
		{
			name: "Default config with specified platform config",
			devModuleConfig: apiv1.Accessory{
				"type":    "cloud",
				"version": "14.0",
			},
			platformConfig: apiv1.GenericConfig{
				"size":           100,
				"privateRouting": true,
				"instanceType":   "test-instance-type",
				"subnetID":       "test-subnet-id",
				"databaseName":   "test-database",
			},
			expectedPostgreSQL: &PostgreSQL{
				Type:           "cloud",
				Version:        "14.0",
				Username:       defaultUsername,
				Category:       defaultCategory,
				SecurityIPs:    defaultSecurityIPs,
				PrivateRouting: true,
				Size:           100,
				InstanceType:   "test-instance-type",
				SubnetID:       "test-subnet-id",
				DatabaseName:   "test-database",
			},
		},
	}

	for _, tc := range testcases {
		postgres := &PostgreSQL{}
		t.Run(tc.name, func(t *testing.T) {
			mockey.PatchConvey("mock postgres validate", t, func() {
				mockey.Mock(postgres.Validate).Return(nil).Build()

				_ = postgres.GetCompleteConfig(tc.devModuleConfig, tc.platformConfig)
				assert.Equal(t, tc.expectedPostgreSQL, postgres)
			})
		})
	}
}

func TestPostgreSQLModule_GenerateDBSecret(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &workload.Workload{
			Header: workload.Header{
				Type: "Service",
			},
			Service: &workload.Service{},
		},
	}

	postgres := &PostgreSQL{
		Type:         "local",
		Version:      "14.0",
		DatabaseName: "test-database",
	}

	hostAddress := "test-host-address"
	username := "test-username"
	password := "test-password"

	sec := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-database-postgres",
			Namespace: "test-project",
		},
		StringData: map[string]string{
			"hostAddress": "test-host-address",
			"username":    "test-username",
			"password":    "test-password",
		},
	}

	resID := module.KubernetesResourceID(sec.TypeMeta, sec.ObjectMeta)
	expectedResource, err := module.WrapK8sResourceToKusionResource(resID, sec)
	if err != nil {
		t.Fatalf("failed to wrap secret resource for unit test: %v", err)
	}

	expectedPatcher := &apiv1.Patcher{
		Environments: []v1.EnvVar{
			{
				Name: "KUSION_DB_HOST_TEST_DATABASE",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "test-database-postgres",
						},
						Key: "hostAddress",
					},
				},
			},
			{
				Name: "KUSION_DB_USERNAME_TEST_DATABASE",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "test-database-postgres",
						},
						Key: "username",
					},
				},
			},
			{
				Name: "KUSION_DB_PASSWORD_TEST_DATABASE",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "test-database-postgres",
						},
						Key: "password",
					},
				},
			},
		},
	}

	actualResource, actualPatchers, err := postgres.GenerateDBSecret(r, hostAddress, username, password)

	assert.Nil(t, err)
	assert.Equal(t, expectedPatcher, actualPatchers)
	assert.Equal(t, expectedResource, actualResource)
}

func TestPostgreSQLModule_GenerateTFRandomPassword(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &workload.Workload{
			Header: workload.Header{
				Type: "Service",
			},
			Service: &workload.Service{},
		},
	}

	postgres := &PostgreSQL{
		Type:         "local",
		Version:      "14.0",
		DatabaseName: "test-database",
	}

	t.Run("failed to generate tf resource id", func(t *testing.T) {
		mockey.PatchConvey("failed to generate tf resource id", t, func() {
			mockey.Mock(module.TerraformResourceID).Return("", errors.New("failed to generate tf resource id")).Build()

			res, id, err := postgres.GenerateTFRandomPassword(r)

			assert.Nil(t, res)
			assert.Equal(t, id, "")
			assert.ErrorContains(t, err, "failed to generate tf resource id")
		})
	})

	t.Run("failed to generate provider extensions", func(t *testing.T) {
		mockey.PatchConvey("failed to generate provider extensions", t, func() {
			mockey.Mock(module.TerraformResourceID).Return("", nil).Build()
			mockey.Mock(module.TerraformProviderExtensions).Return(nil, errors.New("failed to generate provider extensions")).Build()

			res, id, err := postgres.GenerateTFRandomPassword(r)

			assert.Nil(t, res)
			assert.Equal(t, id, "")
			assert.ErrorContains(t, err, "failed to generate provider extensions")
		})
	})

	t.Run("failed to wrap tf resource to kusion resource", func(t *testing.T) {
		mockey.PatchConvey("failed to wrap tf resource to kusion resource", t, func() {
			mockey.Mock(module.TerraformResourceID).Return("", nil).Build()
			mockey.Mock(module.TerraformProviderExtensions).Return(nil, nil).Build()
			mockey.Mock(module.WrapTFResourceToKusionResource).Return(nil, errors.New("failed to wrap tf resource to kusion resource")).Build()

			res, id, err := postgres.GenerateTFRandomPassword(r)

			assert.Nil(t, res)
			assert.Equal(t, id, "")
			assert.ErrorContains(t, err, "failed to wrap tf resource to kusion resource")
		})
	})

	t.Run("successfully generate random_password resource", func(t *testing.T) {
		res, id, err := postgres.GenerateTFRandomPassword(r)

		assert.NotNil(t, res)
		assert.NotEqual(t, id, "")
		assert.NoError(t, err)
	})
}

func TestPostgreSQLModule_Validate(t *testing.T) {
	t.Run("cloud db with empty instanceType", func(t *testing.T) {
		postgres := &PostgreSQL{
			Type:    "cloud",
			Version: "14.0",
		}

		err := postgres.Validate()

		assert.ErrorContains(t, err, ErrEmptyInstanceTypeForCloudDB.Error())
	})

	t.Run("valid postgres config", func(t *testing.T) {
		postgres := &PostgreSQL{
			Type:         "cloud",
			Version:      "14.0",
			InstanceType: "test-instance-type",
		}

		err := postgres.Validate()

		assert.NoError(t, err)
	})
}

func TestIsPublicAccessible(t *testing.T) {
	testcases := []struct {
		name        string
		securityIPs []string
		expected    bool
	}{
		{
			name: "Public CIDR",
			securityIPs: []string{
				"0.0.0.0/0",
			},
			expected: true,
		},
		{
			name: "Private CIDR",
			securityIPs: []string{
				"172.16.0.0/24",
			},
			expected: false,
		},
		{
			name: "Private IP Address",
			securityIPs: []string{
				"172.16.0.1",
			},
			expected: false,
		},
	}

	for _, tc := range testcases {
		actual := IsPublicAccessible(tc.securityIPs)

		assert.Equal(t, tc.expected, actual)
	}
}
