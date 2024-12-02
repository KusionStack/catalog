package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestMySQLModule_GenerateLocalResources(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	resources, patchers, err := mysql.GenerateLocalResources(r)

	assert.Equal(t, 5, len(resources))
	assert.NotNil(t, patchers)
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateLocalSecret(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := mysql.generateLocalSecret(r, "123456")

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateLocalDeployment(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := mysql.generateLocalDeployment(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateLocalPodSpec(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := mysql.generateLocalPodSpec(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateLocalPVC(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := mysql.generateLocalPVC(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestMySQLModule_GenerateLocalService(t *testing.T) {
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
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, svcName, err := mysql.generateLocalService(r)

	assert.NotNil(t, res)
	assert.NotNil(t, svcName)
	assert.NoError(t, err)
}
