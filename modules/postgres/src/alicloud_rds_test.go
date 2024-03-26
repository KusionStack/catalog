package main

import (
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion/pkg/apis/core/v1/workload"
)

func TestPostgreSQLModule_GenerateAlicloudResources(t *testing.T) {
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
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "postgres.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	mockey.PatchConvey("set alicloud region env", t, func() {
		mockey.Mock(os.Getenv).Return("test-region").Build()

		resources, patchers, err := postgres.GenerateAlicloudResources(r)

		assert.Equal(t, 5, len(resources))
		assert.NotNil(t, patchers)
		assert.NoError(t, err)
	})
}

func TestPostgreSQLModule_GenerateAlicloudDBInstance(t *testing.T) {
	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "postgres.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	res, id, err := postgres.generateAlicloudDBInstance(defaultAlicloudProviderCfg, "test-region")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateAlicloudDBConnection(t *testing.T) {
	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "postgres.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	res, id, err := postgres.generateAlicloudDBConnection(defaultAlicloudProviderCfg, "test-region", "db_instance_id")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateAlicloudRDSAccount(t *testing.T) {
	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "postgres.n2.serverless.1c",
		Category:       "serverless_basic",
		SubnetID:       "test-subnet-id",
	}

	res, err := postgres.generateAlicloudRDSAccount(defaultAlicloudProviderCfg, "test-region",
		"account_name", "random_password_id", "db_instance_id")

	assert.NotNil(t, res)
	assert.NoError(t, err)
}
