package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"kusionstack.io/kusion-module-framework/pkg/module"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
)

func TestOpenSearch_CompleteConfig(t *testing.T) {
	devConfig := v1.Accessory{
		"engineVersion": "OpenSearch_1.0",
		"domainName":    "test-domain",
	}

	platformConfig := v1.GenericConfig{
		"clusterConfig": map[string]interface{}{
			"instanceType": "t2.micro.search",
		},
		"ebsOptions": map[string]interface{}{
			"ebsEnabled": true,
			"volumeSize": 10,
		},
		"statement": []map[string]interface{}{
			{
				"effect": "Allow",
				"principal": []Principal{
					{
						Type:        "AWS",
						Identifiers: []string{"arn:aws:iam::12345678901:role/yak-role"},
					},
				},
				"action": []string{"ec2:RunInstances", "s3:*"},
			},
		},
	}

	type args struct {
		devConfig      v1.Accessory
		platformConfig v1.GenericConfig
	}
	tests := []struct {
		name    string
		os      *OpenSearch
		args    args
		wantErr bool
	}{
		{
			name: "CompleteConfig with valid devConfig and platformConfig",
			os:   &OpenSearch{},
			args: args{
				devConfig:      devConfig,
				platformConfig: platformConfig,
			},
			wantErr: false,
		},
		{
			name: "CompleteConfig with invalid devConfig",
			os:   &OpenSearch{},
			args: args{
				devConfig: v1.Accessory{
					"engineVersion": 123, // Invalid type
					"domainName":    "test-domain",
				},
				platformConfig: platformConfig,
			},
			wantErr: true,
		},
		{
			name: "CompleteConfig with empty statements",
			os:   &OpenSearch{},
			args: args{
				devConfig: devConfig,
				platformConfig: v1.GenericConfig{
					"clusterConfig": map[string]interface{}{
						"instanceType": "t2.micro.search",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.os.CompleteConfig(tt.args.devConfig, tt.args.platformConfig); (err != nil) != tt.wantErr {
				t.Errorf("CompleteConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenSearch_ValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		os      *OpenSearch
		wantErr bool
	}{
		{
			name: "ValidateConfig with valid statements",
			os: &OpenSearch{
				Statement: []Statement{
					{
						Effect: "Allow",
						Principals: []Principal{
							{
								Type:        "AWS",
								Identifiers: []string{"arn:aws:iam::12345678901:role/yak-role"},
							},
						},
						Action: []string{"ec2:RunInstances", "s3:*"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "ValidateConfig with invalid effect",
			os: &OpenSearch{
				Statement: []Statement{
					{
						Effect: "InvalidEffect", // Invalid effect
						Principals: []Principal{
							{
								Type:        "AWS",
								Identifiers: []string{"arn:aws:iam::12345678901:role/yak-role"},
							},
						},
						Action: []string{"ec2:RunInstances", "s3:*"},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.os.ValidateConfig(); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpenSearch_GenerateOpenSearchDomain(t *testing.T) {
	type args struct {
		request *module.GeneratorRequest
	}
	os := &OpenSearch{
		EngineVersion: "OpenSearch_1.0",
		DomainName:    "test-domain",
		ClusterConfig: ClusterConfig{
			InstanceType: "t2.micro.search",
		},
		Statement: []Statement{
			{
				Effect: "Allow",
				Principals: []Principal{
					{
						Type:        "AWS",
						Identifiers: []string{"arn:aws:iam::12345678901:role/yak-role"},
					},
				},
				Action: []string{"ec2:RunInstances", "s3:*"},
			},
		},
	}
	tests := []struct {
		name    string
		os      *OpenSearch
		args    args
		wantErr bool
	}{
		{
			name: "GenerateOpenSearchDomain with valid request",
			os:   os,
			args: args{
				request: &module.GeneratorRequest{
					Project: "test-project",
					Stack:   "test-stack",
					App:     "test-app",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _, err := tt.os.GenerateOpenSearchDomain(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateOpenSearchDomain() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.NotEmpty(t, res)
		})
	}
}

func TestOpenSearch_Generate(t *testing.T) {
	os := &OpenSearch{
		EngineVersion: "OpenSearch_1.0",
		DomainName:    "test-domain",
		ClusterConfig: ClusterConfig{
			InstanceType: "t2.micro.search",
		},
		Statement: []Statement{
			{
				Effect: "Allow",
				Principals: []Principal{
					{
						Type:        "AWS",
						Identifiers: []string{"arn:aws:iam::12345678901:role/yak-role"},
					},
				},
				Action: []string{"ec2:RunInstances", "s3:*"},
			},
		},
	}
	tests := []struct {
		name    string
		os      *OpenSearch
		request *module.GeneratorRequest
		wantErr bool
	}{
		{
			name: "Generate with valid request",
			os:   os,
			request: &module.GeneratorRequest{
				Project: "test-project",
				Stack:   "test-stack",
				App:     "test-app",
			},
			wantErr: false,
		},
		{
			name: "Generate with invalid request",
			os:   os,
			request: &module.GeneratorRequest{
				Project: "", // Invalid project
				Stack:   "test-stack",
				App:     "test-app",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.os.Generate(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
