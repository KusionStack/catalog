package main

import (
	"context"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"
	"kusionstack.io/kusion-module-framework/pkg/module"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
)

func TestOpsRuleModule_Generate(t *testing.T) {
	resConfig30 := v1.Resource{
		ID:   "apps.kusionstack.io/v1alpha1:PodTransitionRule:default:default-dev-foo",
		Type: "Kubernetes",
		Attributes: map[string]interface{}{
			"apiVersion": "apps.kusionstack.io/v1alpha1",
			"kind":       "PodTransitionRule",
			"metadata": map[string]interface{}{
				"creationTimestamp": interface{}(nil),
				"name":              "default-dev-foo",
				"namespace":         "default",
			},
			"spec": map[string]interface{}{
				"rules": []interface{}{map[string]interface{}{
					"availablePolicy": map[string]interface{}{
						"maxUnavailableValue": "30%",
					},
					"name": "maxUnavailable",
				}},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app.kubernetes.io/name": "foo", "app.kubernetes.io/part-of": "default",
					},
				},
			}, "status": map[string]interface{}{},
		},
		DependsOn: []string(nil),
		Extensions: map[string]interface{}{
			"GVK": "apps.kusionstack.io/v1alpha1, Kind=PodTransitionRule",
		},
	}
	resConfig40 := v1.Resource{
		ID:   "apps.kusionstack.io/v1alpha1:PodTransitionRule:default:default-dev-foo",
		Type: "Kubernetes",
		Attributes: map[string]interface{}{
			"apiVersion": "apps.kusionstack.io/v1alpha1",
			"kind":       "PodTransitionRule",
			"metadata": map[string]interface{}{
				"creationTimestamp": interface{}(nil),
				"name":              "default-dev-foo",
				"namespace":         "default",
			},
			"spec": map[string]interface{}{
				"rules": []interface{}{map[string]interface{}{
					"availablePolicy": map[string]interface{}{
						"maxUnavailableValue": 40,
					},
					"name": "maxUnavailable",
				}},
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app.kubernetes.io/name": "foo", "app.kubernetes.io/part-of": "default",
					},
				},
			}, "status": map[string]interface{}{},
		},
		DependsOn: []string(nil),
		Extensions: map[string]interface{}{
			"GVK": "apps.kusionstack.io/v1alpha1, Kind=PodTransitionRule",
		},
	}

	jobWorkloadConfig := &v1.Workload{
		Header: v1.Header{
			Type: v1.TypeJob,
		},
	}
	serviceWorkloadConfig := &v1.Workload{
		Header: v1.Header{
			Type: v1.TypeService,
		},
		Service: &v1.Service{
			Type: v1.Collaset,
		},
	}
	devConfig := map[string]interface{}{
		"maxUnavailable": "30%",
	}
	platformConfig := map[string]interface{}{
		"maxUnavailable": 40,
	}

	response30 := &module.GeneratorResponse{
		Resources: []v1.Resource{resConfig30},
	}
	response40 := &module.GeneratorResponse{
		Resources: []v1.Resource{resConfig40},
	}

	project := "default"
	stack := "dev"
	app := "foo"

	type args struct {
		r *module.GeneratorRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *module.GeneratorResponse
		wantErr bool
	}{
		{
			name: "test Job",
			args: args{
				r: &module.GeneratorRequest{
					Project:        project,
					Stack:          stack,
					App:            app,
					Workload:       jobWorkloadConfig,
					DevConfig:      devConfig,
					PlatformConfig: platformConfig,
				},
			},
			want: nil,
		},
		{
			name: "test CollaSet with opsRule in appConfig",
			args: args{
				r: &module.GeneratorRequest{
					Project:        project,
					Stack:          stack,
					App:            app,
					Workload:       serviceWorkloadConfig,
					DevConfig:      devConfig,
					PlatformConfig: platformConfig,
				},
			},
			wantErr: false,
			want:    response30,
		},
		{
			name: "test CollaSet with opsRule in workspace",
			args: args{
				r: &module.GeneratorRequest{
					Project:        project,
					Stack:          stack,
					App:            app,
					Workload:       serviceWorkloadConfig,
					PlatformConfig: platformConfig,
				},
			},
			wantErr: false,
			want:    response40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OpsRuleModule{}
			got, err := o.Generate(context.Background(), tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			out, _ := yaml.Marshal(got)
			out2, _ := yaml.Marshal(tt.want)
			if !reflect.DeepEqual(string(out), string(out2)) {
				t.Errorf("Generate()\ngot = %v\nwant = %v", string(out), string(out2))
			}
		})
	}
}
