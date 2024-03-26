package main

import (
	"os"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion/pkg/apis/core/v1/workload"
)

func TestPostgreSQLModule_GenerateAWSResources(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "db.t3.micro",
	}

	mockey.PatchConvey("set aws region env", t, func() {
		mockey.Mock(os.Getenv).Return("test-region").Build()

		resources, patchers, err := postgres.GenerateAWSResources(r)

		assert.Equal(t, 4, len(resources))
		assert.NotNil(t, patchers)
		assert.NoError(t, err)
	})
}

func TestPostgreSQLModule_GenerateAWSSecurityGroup(t *testing.T) {
	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "db.t3.micro",
	}

	res, id, err := postgres.generateAWSSecurityGroup(defaultAWSProviderCfg, "test-region")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateAWSDBInstance(t *testing.T) {
	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: false,
		Size:           defaultSize,
		InstanceType:   "db.t3.micro",
	}

	res, id, err := postgres.generateAWSDBInstance(defaultAWSProviderCfg, "test-region",
		"random_password_id", "aws_security_group_id")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}
