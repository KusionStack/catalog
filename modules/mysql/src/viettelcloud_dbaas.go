package main

import (
	"os"

	"kusionstack.io/kusion-module-framework/pkg/module"
	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/modules"
)

const (
	DefaultViettelCloudRegion = "vn-central-1"
)

var (
	viettelCloudRegionEnv  = "VIETTEL_CLOUD_REGION"
	viettelCloudDBInstance = "viettelcloud_db_instance"
)

var defaultViettelCloudProviderCfg = module.ProviderConfig{
	Source:  "hashicorp/viettelcloud",
	Version: "1.0.0-dev",
}

// GenerateViettelCloudResources generates the ViettelCloud provided MySQL database instance.
func (mysql *MySQL) GenerateViettelCloudResources(request *module.GeneratorRequest) ([]apiv1.Resource, *apiv1.Patcher, error) {
	var resources []apiv1.Resource

	// Set the ViettelCloud provider with the default provider config.
	viettelCloudProviderCfg := defaultViettelCloudProviderCfg

	// Get the ViettelCloud Terraform provider region, which should not be empty.
	var region string
	if region = module.TerraformProviderRegion(viettelCloudProviderCfg); region == "" {
		region = os.Getenv(viettelCloudRegionEnv)
	}
	if region == "" {
		region = DefaultViettelCloudRegion
	}

	// Build random_password resource.
	randomPasswordRes, randomPasswordID, err := mysql.GenerateTFRandomPassword(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *randomPasswordRes)

	// Build viettelCloud_db_instance resource.
	viettelCloudDBInstance, viettelCloudDBInstanceID, err := mysql.generateViettelCloudDBInstance(viettelCloudProviderCfg, region, randomPasswordID)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *viettelCloudDBInstance)

	hostAddress := modules.KusionPathDependency(viettelCloudDBInstanceID, "private_url")
	password := modules.KusionPathDependency(randomPasswordID, "result")

	// Build Kubernetes Secret with the hostAddress, username and password of the ViettelCloud provided MySQL instance,
	// and inject the credentials as the environment variable patcher.
	dbSecret, patcher, err := mysql.GenerateDBSecret(request, hostAddress, mysql.Username, password)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *dbSecret)

	return resources, patcher, nil
}

// generateViettelCloudDBInstance generates viettelCloud_db_instance resource for the ViettelCloud provided MySQL database instance.
func (mysql *MySQL) generateViettelCloudDBInstance(viettelCloudProviderCfg module.ProviderConfig, region, randomPasswordID string) (*apiv1.Resource, string, error) {
	resAttrs := map[string]interface{}{
		"database_type":      dbEngine,
		"region":             region,
		"name":               mysql.DatabaseName,
		"db_version":         mysql.Version,
		"flavor":             mysql.InstanceType,
		"volume_type":        mysql.VolumeType,
		"disk_size":          mysql.Size,
		"vpc_name":           mysql.VPC,
		"subnet":             mysql.SubnetID,
		"solution":           mysql.Category,
		"root_password":      modules.KusionPathDependency(randomPasswordID, "result"),
		"enable_auto_backup": false,
	}

	id, err := module.TerraformResourceID(viettelCloudProviderCfg, viettelCloudDBInstance, mysql.DatabaseName)
	if err != nil {
		return nil, "", err
	}

	viettelCloudProviderCfg.ProviderMeta = map[string]any{}
	resource, err := module.WrapTFResourceToKusionResource(viettelCloudProviderCfg, viettelCloudDBInstance, id, resAttrs, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}
