package main

import (
	"errors"
	"fmt"
	"os"

	"kusionstack.io/kusion-module-framework/pkg/module"
	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/modules"
)

var ErrEmptyAWSProviderRegion = errors.New("empty aws provider region")

var (
	awsRegionEnv     = "AWS_REGION"
	awsSecurityGroup = "aws_security_group"
	awsDBInstance    = "aws_db_instance"
)

var defaultAWSProviderCfg = module.ProviderConfig{
	Source:  "hashicorp/aws",
	Version: "5.0.1",
}

type awsSecurityGroupTraffic struct {
	CidrBlocks     []string `yaml:"cidr_blocks" json:"cidr_blocks"`
	Description    string   `yaml:"description" json:"description"`
	FromPort       int      `yaml:"from_port" json:"from_port"`
	IPv6CIDRBlocks []string `yaml:"ipv6_cidr_blocks" json:"ipv6_cidr_blocks"`
	PrefixListIDs  []string `yaml:"prefix_list_ids" json:"prefix_list_ids"`
	Protocol       string   `yaml:"protocol" json:"protocol"`
	SecurityGroups []string `yaml:"security_groups" json:"security_groups"`
	Self           bool     `yaml:"self" json:"self"`
	ToPort         int      `yaml:"to_port" json:"to_port"`
}

// GenerateAWSResources generates the AWS provided MySQL database instance.
func (mysql *MySQL) GenerateAWSResources(request *module.GeneratorRequest) ([]apiv1.Resource, *apiv1.Patcher, error) {
	var resources []apiv1.Resource

	// Set the AWS provider with the default provider config.
	awsProviderCfg := defaultAWSProviderCfg

	// Get the AWS Terraform provider region, which should not be empty.
	var region string
	if region = module.TerraformProviderRegion(awsProviderCfg); region == "" {
		region = os.Getenv(awsRegionEnv)
	}
	if region == "" {
		return nil, nil, ErrEmptyAWSProviderRegion
	}

	// Build random_password resource.
	randomPasswordRes, randomPasswordID, err := mysql.GenerateTFRandomPassword(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *randomPasswordRes)

	// Build aws_security_group resource.
	awsSecurityGroupRes, awsSecurityGroupID, err := mysql.generateAWSSecurityGroup(awsProviderCfg, region)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *awsSecurityGroupRes)

	// Build aws_db_instance resource.
	awsDBInstance, awsDBInstanceID, err := mysql.generateAWSDBInstance(awsProviderCfg, region, randomPasswordID, awsSecurityGroupID)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *awsDBInstance)

	hostAddress := modules.KusionPathDependency(awsDBInstanceID, "address")
	password := modules.KusionPathDependency(randomPasswordID, "result")

	// Build Kubernetes Secret with the hostAddress, username and password of the AWS provided MySQL instance,
	// and inject the credentials as the environment variable patcher.
	dbSecret, patcher, err := mysql.GenerateDBSecret(request, hostAddress, mysql.Username, password)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *dbSecret)

	return resources, patcher, nil
}

// generateAWSSecurityGroup generates aws_security_group resource for the AWS provided MySQL database instance.
func (mysql *MySQL) generateAWSSecurityGroup(awsProviderCfg module.ProviderConfig, region string) (*apiv1.Resource, string, error) {
	// SecurityIPs should be in the format of IP address or Classes Inter-Domain
	// Routing (CIDR) mode.
	for _, ip := range mysql.SecurityIPs {
		if !IsIPAddress(ip) && !IsCIDR(ip) {
			return nil, "", fmt.Errorf("illegal security ip format: %s", ip)
		}
	}

	resAttrs := map[string]interface{}{
		"egress": []awsSecurityGroupTraffic{
			{
				CidrBlocks: []string{"0.0.0.0/0"},
				Protocol:   "-1",
				FromPort:   0,
				ToPort:     0,
			},
		},
		"ingress": []awsSecurityGroupTraffic{
			{
				CidrBlocks: mysql.SecurityIPs,
				Protocol:   "tcp",
				FromPort:   3306,
				ToPort:     3306,
			},
		},
	}

	id, err := module.TerraformResourceID(awsProviderCfg, awsSecurityGroup, mysql.DatabaseName+dbResSuffix)
	if err != nil {
		return nil, "", err
	}

	awsProviderCfg.ProviderMeta = map[string]any{"region": region}
	resource, err := module.WrapTFResourceToKusionResource(awsProviderCfg, awsSecurityGroup, id, resAttrs, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}

// generateAWSDBInstance generates aws_db_instance resource for the AWS provided MySQL database instance.
func (mysql *MySQL) generateAWSDBInstance(awsProviderCfg module.ProviderConfig, region, randomPasswordID, awsSecurityGroupID string) (*apiv1.Resource, string, error) {
	resAttrs := map[string]interface{}{
		"allocated_storage":   mysql.Size,
		"engine":              dbEngine,
		"engine_version":      mysql.Version,
		"identifier":          mysql.DatabaseName,
		"instance_class":      mysql.InstanceType,
		"password":            modules.KusionPathDependency(randomPasswordID, "result"),
		"publicly_accessible": IsPublicAccessible(mysql.SecurityIPs),
		"skip_final_snapshot": true,
		"username":            mysql.Username,
		"vpc_security_group_ids": []string{
			modules.KusionPathDependency(awsSecurityGroupID, "id"),
		},
	}

	if mysql.SubnetID != "" {
		resAttrs["db_subnet_group_name"] = mysql.SubnetID
	}

	id, err := module.TerraformResourceID(awsProviderCfg, awsDBInstance, mysql.DatabaseName)
	if err != nil {
		return nil, "", err
	}

	awsProviderCfg.ProviderMeta = map[string]any{"region": region}
	resource, err := module.WrapTFResourceToKusionResource(awsProviderCfg, awsDBInstance, id, resAttrs, nil)
	if err != nil {
		return nil, "", err
	}

	return resource, id, nil
}
