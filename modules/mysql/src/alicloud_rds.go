package main

import (
	"errors"
	"os"
	"strings"

	"kusionstack.io/kusion-module-framework/pkg/module"
	apiv1 "kusionstack.io/kusion/pkg/apis/core/v1"
	"kusionstack.io/kusion/pkg/modules"
)

var ErrEmptyAlicloudProviderRegion = errors.New("empty alicloud provider region")

var (
	alicloudRegionEnv    = "ALICLOUD_REGION"
	alicloudDBInstance   = "alicloud_db_instance"
	alicloudDBConnection = "alicloud_db_connection"
	alicloudRDSAccount   = "alicloud_rds_account"
)

var defaultAlicloudProviderCfg = apiv1.ProviderConfig{
	Source:  "aliyun/alicloud",
	Version: "1.209.1",
}

type alicloudServerlessConfig struct {
	AutoPause   bool `yaml:"auto_pause" json:"auto_pause"`
	SwitchForce bool `yaml:"switch_force" json:"switch_force"`
	MaxCapacity int  `yaml:"max_capacity,omitempty" json:"max_capacity,omitempty"`
	MinCapacity int  `yaml:"min_capacity,omitempty" json:"min_capacity,omitempty"`
}

// GenerateAlicloudResources generates Alicloud provided MySQL database instance.
func (mysql *MySQL) GenerateAlicloudResources(request *module.GeneratorRequest) ([]apiv1.Resource, []apiv1.Patcher, error) {
	var resources []apiv1.Resource
	var patchers []apiv1.Patcher

	// Set the Alicloud provider with the default provider config.
	alicloudProviderCfg := defaultAlicloudProviderCfg

	// Get the Alicloud Terraform provider region, which should not be empty.
	var region string
	if region = module.TerraformProviderRegion(alicloudProviderCfg); region == "" {
		region = os.Getenv(alicloudRegionEnv)
	}
	if region == "" {
		return nil, nil, ErrEmptyAlicloudProviderRegion
	}

	// Build random_password resource.
	randomPasswordRes, randomPasswordID, err := mysql.GenerateTFRandomPassword(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *randomPasswordRes)

	// Build alicloud_db_instance resource.
	alicloudDBInstanceRes, alicloudDBInstanceID, err := mysql.generateAlicloudDBInstance(
		alicloudProviderCfg, region,
	)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *alicloudDBInstanceRes)

	// Build alicloud_db_connection resource.
	var alicloudDBConnectionRes *apiv1.Resource
	var alicloudDBConnectionID string
	if IsPublicAccessible(mysql.SecurityIPs) {
		alicloudDBConnectionRes, alicloudDBConnectionID, err = mysql.generateAlicloudDBConnection(
			alicloudProviderCfg,
			region, alicloudDBInstanceID,
		)
		if err != nil {
			return nil, nil, err
		}

		resources = append(resources, *alicloudDBConnectionRes)
	}

	// Build alicloud_rds_account resuorce.
	alicloudRDSAccountRes, err := mysql.generateAlicloudRDSAccount(
		alicloudProviderCfg,
		region, mysql.Username, randomPasswordID, alicloudDBInstanceID,
	)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *alicloudRDSAccountRes)

	hostAddress := modules.KusionPathDependency(alicloudDBInstanceID, "connection_string")
	if !mysql.PrivateRouting {
		// Set the public network connection string as the host address.
		hostAddress = modules.KusionPathDependency(alicloudDBConnectionID, "connection_string")
	}
	password := modules.KusionPathDependency(randomPasswordID, "result")

	// Build Kubernetes Secret with the hostAddress, username and password of the Alicloud provided MySQL instance,
	// and inject the credentials as the environment variable patcher.
	dbSecret, pathcer, err := mysql.GenerateDBSecret(request, hostAddress, mysql.Username, password)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *dbSecret)
	patchers = append(patchers, *pathcer)

	return resources, patchers, nil
}

// generateAlicloudDBInstance generates alicloud_db_instance resource
// for the Alicloud provided MySQL database instance.
func (mysql *MySQL) generateAlicloudDBInstance(alicloudProviderCfg apiv1.ProviderConfig,
	region string,
) (*apiv1.Resource, string, error) {
	resAttrs := map[string]interface{}{
		"category":         mysql.Category,
		"engine":           "MySQL",
		"engine_version":   mysql.Version,
		"instance_storage": mysql.Size,
		"instance_type":    mysql.InstanceType,
		"security_ips":     mysql.SecurityIPs,
		"vswitch_id":       mysql.SubnetID,
		"instance_name":    mysql.DatabaseName,
	}

	// Set the serverless-specific attributes of the alicloud_db_instance resource.
	if strings.Contains(mysql.Category, "serverless") {
		resAttrs["db_instance_storage_type"] = "cloud_essd"
		resAttrs["instance_charge_type"] = "Serverless"

		serverlessConfig := alicloudServerlessConfig{
			MaxCapacity: 8,
			MinCapacity: 1,
		}
		serverlessConfig.AutoPause = false
		serverlessConfig.SwitchForce = false

		resAttrs["serverless_config"] = []alicloudServerlessConfig{
			serverlessConfig,
		}
	}

	id, err := module.TerraformResourceID(alicloudProviderCfg, alicloudDBInstance, mysql.DatabaseName)
	if err != nil {
		return nil, "", err
	}

	resExts, err := module.TerraformProviderExtensions(alicloudProviderCfg, map[string]any{"region": region}, alicloudDBInstance)
	if err != nil {
		return nil, "", err
	}

	resource, err := module.WrapTFResourceToKusionResource(id, resAttrs, resExts, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}

// generateAlicloudDBConnection generates alicloud_db_connection resource
// for the Alicloud provided MySQL database instance.
func (mysql *MySQL) generateAlicloudDBConnection(alicloudProviderCfg apiv1.ProviderConfig,
	region, dbInstanceID string,
) (*apiv1.Resource, string, error) {
	resAttrs := map[string]interface{}{
		"instance_id": modules.KusionPathDependency(dbInstanceID, "id"),
	}

	id, err := module.TerraformResourceID(alicloudProviderCfg, alicloudDBConnection, mysql.DatabaseName)
	if err != nil {
		return nil, "", err
	}

	resExts, err := module.TerraformProviderExtensions(alicloudProviderCfg, map[string]any{"region": region}, alicloudDBConnection)
	if err != nil {
		return nil, "", err
	}

	resource, err := module.WrapTFResourceToKusionResource(id, resAttrs, resExts, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}

// generateAlicloudRDSAccount generates alicloud_rds_account resource
// for the Alicloud provided MySQL database instance.
func (mysql *MySQL) generateAlicloudRDSAccount(alicloudProviderCfg apiv1.ProviderConfig,
	region, accountName, randomPasswordID, dbInstanceID string,
) (*apiv1.Resource, error) {
	resAttrs := map[string]interface{}{
		"account_name":     accountName,
		"account_password": modules.KusionPathDependency(randomPasswordID, "result"),
		"account_type":     "Super",
		"db_instance_id":   modules.KusionPathDependency(dbInstanceID, "id"),
	}

	id, err := module.TerraformResourceID(alicloudProviderCfg, alicloudRDSAccount, mysql.DatabaseName)
	if err != nil {
		return nil, err
	}

	resExts, err := module.TerraformProviderExtensions(alicloudProviderCfg, map[string]any{"region": region}, alicloudRDSAccount)
	if err != nil {
		return nil, err
	}

	resource, err := module.WrapTFResourceToKusionResource(id, resAttrs, resExts, nil)
	if err != nil {
		return nil, err
	}

	return resource, nil
}
