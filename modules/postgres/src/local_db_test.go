package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestPostgreSQLModule_GenerateLocalResources(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	resources, patchers, err := postgres.GenerateLocalResources(r)

	assert.Equal(t, 5, len(resources))
	assert.NotNil(t, patchers)
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateLocalSecret(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := postgres.generateLocalSecret(r, "123456")

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateLocalDeployment(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := postgres.generateLocalDeployment(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateLocalPodSpec(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := postgres.generateLocalPodSpec(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateLocalPVC(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, err := postgres.generateLocalPVC(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestPostgreSQLModule_GenerateLocalService(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: kusionapiv1.Accessory{
			"_type": "service.Service",
			"type":  "service",
		},
	}

	postgres := &PostgreSQL{
		Type:           "local",
		Version:        "14.0",
		DatabaseName:   "test-database",
		Username:       defaultUsername,
		Category:       defaultCategory,
		SecurityIPs:    defaultSecurityIPs,
		PrivateRouting: defaultPrivateRouting,
		Size:           defaultSize,
	}

	res, svcName, err := postgres.generateLocalService(r)

	assert.NotNil(t, res)
	assert.NotNil(t, svcName)
	assert.NoError(t, err)
}
