package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestMySQLModule_Generator(t *testing.T) {
	// Set provider envs.
	originAWSRegion := os.Getenv("AWS_REGION")
	originAlicloudRegion := os.Getenv("ALICLOUD_REGION")

	defer func() {
		os.Setenv("AWS_REGION", originAWSRegion)
		os.Setenv("ALICLOUD_REGION", originAlicloudRegion)
	}()

	os.Setenv("AWS_REGION", "us-east-1")
	os.Setenv("ALICLOUD_REGION", "cn-beijing")

	// TODO: set env for AWS & Alicloud Region.
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
		name            string
		devModuleConfig kusionapiv1.Accessory
		platformConfig  kusionapiv1.GenericConfig
		expectedErr     error
	}{
		{
			name: "Generate local MySQL database",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "local",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"databaseName": "test-mysql",
			},
			expectedErr: nil,
		},
		{
			name: "Generate AWS MySQL RDS",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "cloud",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"cloud":          "aws",
				"size":           20,
				"instanceType":   "db.t3.micro",
				"privateRouting": false,
			},
			expectedErr: nil,
		},
		{
			name: "Generate Alicloud MySQL RDS",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "cloud",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"cloud":          "alicloud",
				"size":           20,
				"instanceType":   "mysql.n2.serverless.1c",
				"category":       "serverless_basic",
				"privateRouting": false,
				"subnetID":       "test-subnet-id",
			},
			expectedErr: nil,
		},
		{
			name: "Unsupported MySQL type",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "unsupported-type",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"databaseName": "test-mysql",
			},
			expectedErr: errors.New("unsupported mysql type"),
		},
		{
			name: "Unsupported Terraform provider type",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "cloud",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"cloud":        "unsupported-type",
				"instanceType": "db.t3.micro",
			},
			expectedErr: errors.New("unsupported cloud provider type"),
		},
		{
			name: "Empty cloud MySQL instance type",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "cloud",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"cloud": "aws",
			},
			expectedErr: ErrEmptyInstanceTypeForCloudDB,
		},
	}

	for _, tc := range testcases {
		mysql := &MySQL{}
		t.Run(tc.name, func(t *testing.T) {
			r.DevConfig = tc.devModuleConfig
			r.PlatformConfig = tc.platformConfig

			res, err := mysql.Generate(context.Background(), r)
			if tc.expectedErr != nil {
				assert.ErrorContains(t, err, tc.expectedErr.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
			}
		})
	}
}

func TestMySQLModule_GetCompleteConfig(t *testing.T) {
	testcases := []struct {
		name            string
		devModuleConfig kusionapiv1.Accessory
		platformConfig  kusionapiv1.GenericConfig
		expectedMySQL   *MySQL
	}{
		{
			name: "Empty platform config",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "local",
				"version": "8.0",
			},
			platformConfig: nil,
			expectedMySQL: &MySQL{
				Type:           "local",
				Version:        "8.0",
				Username:       defaultUsername,
				Category:       defaultCategory,
				SecurityIPs:    defaultSecurityIPs,
				PrivateRouting: defaultPrivateRouting,
				Size:           defaultSize,
			},
		},
		{
			name: "Default config with specified platform config",
			devModuleConfig: kusionapiv1.Accessory{
				"type":    "cloud",
				"version": "8.0",
			},
			platformConfig: kusionapiv1.GenericConfig{
				"size":           100,
				"privateRouting": true,
				"instanceType":   "test-instance-type",
				"subnetID":       "test-subnet-id",
				"databaseName":   "test-database",
			},
			expectedMySQL: &MySQL{
				Type:           "cloud",
				Version:        "8.0",
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
		mysql := &MySQL{}
		t.Run(tc.name, func(t *testing.T) {
			mockey.PatchConvey("mock mysql validate", t, func() {
				mockey.Mock(mysql.Validate).Return(nil).Build()

				_ = mysql.GetCompleteConfig(tc.devModuleConfig, tc.platformConfig)
				assert.Equal(t, tc.expectedMySQL, mysql)
			})
		})
	}
}

func TestMySQLModule_GenerateDBSecret(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	mysql := &MySQL{
		Type:         "local",
		Version:      "8.0",
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
			Name:      "test-database-mysql",
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

	expectedPatcher := &kusionapiv1.Patcher{
		Environments: []v1.EnvVar{
			{
				Name: "KUSION_DB_HOST_TEST_DATABASE",
				ValueFrom: &v1.EnvVarSource{
					SecretKeyRef: &v1.SecretKeySelector{
						LocalObjectReference: v1.LocalObjectReference{
							Name: "test-database-mysql",
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
							Name: "test-database-mysql",
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
							Name: "test-database-mysql",
						},
						Key: "password",
					},
				},
			},
		},
	}

	actualResource, actualPatcher, err := mysql.GenerateDBSecret(r, hostAddress, username, password)

	assert.Nil(t, err)
	assert.Equal(t, expectedResource, actualResource)
	assert.Equal(t, expectedPatcher, actualPatcher)
}

func TestMySQLModule_GenerateTFRandomPassword(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	mysql := &MySQL{
		Type:         "local",
		Version:      "8.0",
		DatabaseName: "test-database",
	}

	t.Run("failed to generate tf resource id", func(t *testing.T) {
		mockey.PatchConvey("failed to generate tf resource id", t, func() {
			mockey.Mock(module.TerraformResourceID).Return("", errors.New("failed to generate tf resource id")).Build()

			res, id, err := mysql.GenerateTFRandomPassword(r)

			assert.Nil(t, res)
			assert.Equal(t, id, "")
			assert.ErrorContains(t, err, "failed to generate tf resource id")
		})
	})

	t.Run("failed to generate provider extensions", func(t *testing.T) {
		mockey.PatchConvey("failed to generate provider extensions", t, func() {
			mockey.Mock(module.TerraformResourceID).Return("", nil).Build()
			mockey.Mock(module.TerraformProviderExtensions).Return(nil, errors.New("failed to generate provider extensions")).Build()

			res, id, err := mysql.GenerateTFRandomPassword(r)

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

			res, id, err := mysql.GenerateTFRandomPassword(r)

			assert.Nil(t, res)
			assert.Equal(t, id, "")
			assert.ErrorContains(t, err, "failed to wrap tf resource to kusion resource")
		})
	})

	t.Run("successfully generate random_password resource", func(t *testing.T) {
		res, id, err := mysql.GenerateTFRandomPassword(r)

		assert.NotNil(t, res)
		assert.NotEqual(t, id, "")
		assert.NoError(t, err)
	})
}

func TestMySQLModule_Validate(t *testing.T) {
	t.Run("cloud db with empty instanceType", func(t *testing.T) {
		mysql := &MySQL{
			Type:    "cloud",
			Version: "8.0",
		}

		err := mysql.Validate()

		assert.ErrorContains(t, err, ErrEmptyInstanceTypeForCloudDB.Error())
	})

	t.Run("valid mysql config", func(t *testing.T) {
		mysql := &MySQL{
			Type:         "cloud",
			Version:      "8.0",
			InstanceType: "test-instance-type",
		}

		err := mysql.Validate()

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
