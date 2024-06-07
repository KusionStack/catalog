package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/log"
	"kusionstack.io/kusion/pkg/modules"
)

func main() {
	server.Start(&OpenSearch{})
}

// OpenSearch implements the Kusion Module generator interface.
type OpenSearch struct {
	// DevConfigs
	// Either Elasticsearch_X.Y or OpenSearch_X.Y to specify the engine version for the Amazon OpenSearch Service domain.
	// For example, OpenSearch_1.0 or Elasticsearch_7.9. Defaults to the latest version of OpenSearch.
	EngineVersion string `yaml:"engineVersion,omitempty" json:"engineVersion,omitempty"`
	// Name of the domain
	DomainName string `json:"domainName" yaml:"domainName"`

	// Platform Configs
	// ClusterConfig contains configurations for the cluster of the domain.
	ClusterConfig ClusterConfig `json:"clusterConfig" yaml:"clusterConfig"`
	// EbsOptions contains options for EBS volumes attached to data nodes in the domain.
	EbsOptions EbsOptions `json:"ebsOptions" yaml:"ebsOptions"`
	// Region represent the aws region
	Region    string      `json:"region" yaml:"region"`
	Statement []Statement `json:"statement" yaml:"statement"`
}

type ClusterConfig struct {
	// Instance type of data nodes in the cluster.
	// Argument values end in search for OpenSearch vs. elasticsearch for Elasticsearch (e.g., t2.micro.search vs. t2.micro.elasticsearch).
	InstanceType string `json:"instanceType" yaml:"instanceType"`
}

type EbsOptions struct {
	// Whether EBS volumes are attached to data nodes in the domain.
	EbsEnabled bool `json:"ebsEnabled,omitempty" yaml:"ebsEnabled,omitempty"`
	// Size of EBS volumes attached to data nodes (in GiB). Required if ebs_enabled is set to true.
	VolumeSize int `json:"volumeSize,omitempty" yaml:"volumeSize,omitempty"`
}

type EffectType string

const Allow EffectType = "Allow"
const Deny EffectType = "Deny"

type Statement struct {
	// Whether this statement allows or denies the given actions. Valid values are Allow and Deny. Defaults to Allow.
	Effect EffectType `json:"effect" yaml:"effect"`
	// Configuration block for principals.
	Principals []Principal `json:"principals" yaml:"principals"`
	// List of actions that this statement either allows or denies. For example, ["ec2:RunInstances", "s3:*"].
	Action []string `json:"action" yaml:"action"`
}
type Principal struct {
	// Type of principal. Valid values include AWS, Service, Federated, CanonicalUser and *.
	Type string `json:"type" yaml:"type"`
	// List of identifiers for principals. When type is AWS, these are IAM principal ARNs, e.g., arn:aws:iam::12345678901:role/yak-role. When type is Service, these are AWS Service roles, e.g., lambda.amazonaws.com. When type is Federated, these are web identity users or SAML provider ARNs, e.g., accounts.google.com or arn:aws:iam::12345678901:saml-provider/yak-saml-provider.
	// When type is CanonicalUser, these are canonical user IDs, e.g., 79a59df900b949e55d96a1e698fbacedfd6e09d98eacf8f8d5218e7cd47ef2be.
	Identifiers []string `json:"identifiers" yaml:"identifiers"`
}

// Generate implements the generation logic of openSearch module
func (k *OpenSearch) Generate(_ context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case error:
				err = x
			default:
				err = errors.New("unknown panic")
			}
			log.Errorf("failed to generate openSearch module: %v", r)
		}
	}()

	generatorResponse, err := validateRequest(request)
	if err != nil {
		return generatorResponse, err
	}

	// OpenSearch module does not exist in AppConfiguration configs.
	if request.DevConfig == nil {
		log.Info("OpenSearch module does not exist in AppConfiguration configs")
	}

	// Get the complete openSearch module configs.
	if err = k.CompleteConfig(request.DevConfig, request.PlatformConfig); err != nil {
		log.Debugf("failed to get complete openSearch module configs: %v", err)
		return nil, err
	}

	// Validate the completed openSearch module configs.
	if err := k.ValidateConfig(); err != nil {
		log.Errorf("failed to validate the openSearch module configs: %s", err.Error())
		return nil, err
	}

	var resources []v1.Resource

	// Generate the Terraform aws_opensearch_domain
	resource, patcher, err := k.GenerateOpenSearchDomain(request)
	if err != nil {
		return nil, err
	}
	resources = append(resources, *resource)

	// Return the Kusion generator response.
	return &module.GeneratorResponse{
		Resources: resources,
		Patcher:   patcher,
	}, nil
}

func validateRequest(request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if request.Project == "" {
		return nil, fmt.Errorf("empty project")
	}
	if request.Stack == "" {
		return nil, fmt.Errorf("empty stack")
	}
	if request.App == "" {
		return nil, fmt.Errorf("empty app")
	}
	return nil, nil
}

// CompleteConfig completes the openSearch module configs with both DevConfig and platformModuleConfig.
func (k *OpenSearch) CompleteConfig(devConfig v1.Accessory, platformConfig v1.GenericConfig) error {
	if devConfig != nil {
		devCfgYamlStr, err := json.Marshal(devConfig)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(devCfgYamlStr, k); err != nil {
			return err
		}
	}

	if platformConfig != nil {
		platformCfgYamlStr, err := json.Marshal(platformConfig)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(platformCfgYamlStr, k); err != nil {
			return err
		}
	}

	return nil
}

// ValidateConfig validates the completed openSearch configs are valid or not.
func (k *OpenSearch) ValidateConfig() error {
	if k.Region == "" {
		return fmt.Errorf("empty region")
	}

	statements := k.Statement
	if statements != nil {
		for _, statement := range statements {
			if statement.Effect != Allow && statement.Effect != Deny {
				return fmt.Errorf("invalid effect type: %s. Only 'Allow' and 'Deny' are allowed", statement.Effect)
			}
		}
	}
	return nil
}

// GenerateOpenSearchDomain generates the AWS Terraform provider OpenSearch resource
func (k *OpenSearch) GenerateOpenSearchDomain(request *module.GeneratorRequest) (*v1.Resource, *v1.Patcher, error) {
	// Set the random_password provider config.

	providerConfig := module.ProviderConfig{
		Source:  "hashicorp/aws",
		Version: "5.51.1",
		ProviderMeta: map[string]any{
			"region": k.Region,
		},
	}

	// Set attributes.
	attrs := map[string]any{
		"domain_name":    k.DomainName,
		"engine_version": k.EngineVersion,
		"cluster_config": map[string]string{
			"instance_type": k.ClusterConfig.InstanceType,
		},
		"ebs_options": map[string]any{
			"ebs_enabled": k.EbsOptions.EbsEnabled,
			"volume_size": k.EbsOptions.VolumeSize,
		},
	}

	if k.Statement != nil {
		marshal, err := json.Marshal(k.Statement)
		if err != nil {
			return nil, nil, err
		}
		attrs["access_policies"] = string(marshal)
	}

	// Generate Kusion resource ID and extensions
	appUniqueName := modules.UniqueAppName(request.Project, request.Stack, request.App)
	resType := "aws_opensearch_domain"
	resourceID, err := module.TerraformResourceID(providerConfig, resType, appUniqueName)
	if err != nil {
		return nil, nil, err
	}

	resource, err := module.WrapTFResourceToKusionResource(providerConfig, resType, resourceID, attrs, nil)
	if err != nil {
		return nil, nil, err
	}

	endpoint := modules.KusionPathDependency(resourceID, "endpoint")
	patcher := &v1.Patcher{
		Environments: []corev1.EnvVar{{
			Name:  "OPEN_SEARCH_ENDPOINT",
			Value: endpoint,
		}, {
			Name:  "OPEN_SEARCH_REGION",
			Value: k.Region,
		},
		},
	}

	return resource, patcher, nil
}
