package main

import (
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestMySQLModule_GenerateAlicloudResources(t *testing.T) {
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
		Type:           "local",
		Version:        "8.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "mysql.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	mockey.PatchConvey("set alicloud region env", t, func() {
		mockey.Mock(os.Getenv).Return("test-region").Build()

		resources, patchers, err := mysql.GenerateAlicloudResources(r)

		assert.Equal(t, 5, len(resources))
		assert.NotNil(t, patchers)
		assert.NoError(t, err)
	})
}

func TestMySQLModule_GenerateAlicloudDBInstance(t *testing.T) {
	mysql := &MySQL{
		Type:           "local",
		Version:        "8.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "mysql.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	res, id, err := mysql.generateAlicloudDBInstance(defaultAlicloudProviderCfg, "test-region")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateAlicloudDBConnection(t *testing.T) {
	mysql := &MySQL{
		Type:           "local",
		Version:        "8.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "mysql.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	res, id, err := mysql.generateAlicloudDBConnection(defaultAlicloudProviderCfg, "test-region", "db_instance_id")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateAlicloudRDSAccount(t *testing.T) {
	mysql := &MySQL{
		Type:           "local",
		Version:        "8.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "mysql.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	res, err := mysql.generateAlicloudRDSAccount(defaultAlicloudProviderCfg, "test-region",
		"account_name", "random_password_id", "db_instance_id")

	assert.NotNil(t, res)
	assert.NoError(t, err)
}
