package main

import (
	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"kusionstack.io/kusion-module-framework/pkg/module"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"os"
	"testing"
)

func TestMySQLModule_GenerateViettelCloudResources(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &v1.Workload{
			Header: v1.Header{
				Type: "Service",
			},
			Service: &v1.Service{},
		},
	}

	mysql := &MySQL{
		Type:         "local",
		Version:      "8.0",
		DatabaseName: "test-database",
		Username:     defaultUsername,
		Category:     defaultCategory,
		Size:         defaultSize,
		InstanceType: "DBAAS_1vCPU_1_RAM",
		VolumeType:   "ssd",
		VPC:          "vpc-new",
		SubnetID:     "subnet",
	}

	mockey.PatchConvey("set viettelcloud region env", t, func() {
		mockey.Mock(os.Getenv).Return("test-region").Build()

		resources, patchers, err := mysql.GenerateViettelCloudResources(r)

		assert.Equal(t, 3, len(resources))
		assert.NotNil(t, patchers)
		assert.NoError(t, err)
	})
}

func TestMySQLModule_GenerateViettelCloudDBInstance(t *testing.T) {
	mysql := &MySQL{
		Type:         "local",
		Version:      "8.0",
		DatabaseName: "test-database",
		Username:     defaultUsername,
		Category:     defaultCategory,
		Size:         defaultSize,
		InstanceType: "DBAAS_1vCPU_1_RAM",
		VolumeType:   "ssd",
		VPC:          "vpc-new",
		SubnetID:     "subnet",
	}

	res, id, err := mysql.generateViettelCloudDBInstance(defaultViettelCloudProviderCfg, "test-region",
		"random_password_id")

	assert.NotNil(t, res)
	assert.NotEqual(t, id, "")
	assert.NoError(t, err)
}
