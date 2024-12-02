package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"

	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
)

func TestGenerate(t *testing.T) {
	jobConfig := &Job{
		Base: Base{
			Containers: map[string]Container{
				"busybox": {
					Image: "busybox:1.28",
					Command: []string{
						"/bin/sh",
						"-c",
						"echo hello",
					},
				},
			},
		},
		Schedule: "0 * * * *",
	}

	var devConfig map[string]interface{}
	temp, _ := yamlv2.Marshal(jobConfig)
	_ = yamlv2.Unmarshal(temp, &devConfig)

	tests := []struct {
		name    string
		request *module.GeneratorRequest
		want    *module.GeneratorResponse
		wantErr bool
	}{
		{
			name: "CronJob",
			request: &module.GeneratorRequest{
				Project:        "default",
				Stack:          "dev",
				App:            "foo",
				DevConfig:      devConfig,
				PlatformConfig: nil,
			},
			wantErr: false,
			want: &module.GeneratorResponse{
				Resources: []kusionapiv1.Resource{
					{
						ID:   "batch/v1:CronJob:default:default-dev-foo",
						Type: kusionapiv1.Kubernetes,
						Attributes: map[string]interface{}{
							"apiVersion": "batch/v1",
							"kind":       "CronJob",
							"metadata": map[string]interface{}{
								"creationTimestamp": nil,
								"labels": map[string]interface{}{
									"app.kubernetes.io/name":    "foo",
									"app.kubernetes.io/part-of": "default",
								},
								"name":      "default-dev-foo",
								"namespace": "default",
							},
							"spec": map[string]interface{}{
								"jobTemplate": map[string]interface{}{
									"metadata": map[string]interface{}{
										"creationTimestamp": nil,
									},
									"spec": map[string]interface{}{
										"template": map[string]interface{}{
											"metadata": map[string]interface{}{
												"labels": map[string]interface{}{
													"app.kubernetes.io/name":    "foo",
													"app.kubernetes.io/part-of": "default",
												},
												"creationTimestamp": nil,
											},
											"spec": map[string]interface{}{
												"containers": []interface{}{
													map[string]interface{}{
														"name":  "busybox",
														"image": "busybox:1.28",
														"command": []interface{}{
															"/bin/sh",
															"-c",
															"echo hello",
														},
														"resources": map[string]interface{}{},
													},
												},
												"restartPolicy": "Never",
											},
										},
									},
								},
								"schedule": "0 * * * *",
							},
							"status": map[string]interface{}{},
						},
						Extensions: map[string]interface{}{
							"GVK": "batch/v1, Kind=CronJob",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &Job{}
			got, err := o.Generate(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() err = %v, wanted error %v", err, tt.wantErr)
				return
			}

			for i, resource := range got.Resources {
				// Fixme: consider the case that more than one resource.
				assert.Equal(t, tt.want.Resources[i], resource)
			}
		})
	}
}
