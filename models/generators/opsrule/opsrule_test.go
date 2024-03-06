package main

import (
	"reflect"
	"testing"

	v1 "kusionstack.io/kusion/pkg/apis/core/v1"
	"kusionstack.io/kusion/pkg/apis/core/v1/workload"
	"kusionstack.io/kusion/pkg/modules/proto"
	jsonutil "kusionstack.io/kusion/pkg/util/json"

	generators "kusion-modules"
)

func TestOpsRuleModule_Generate(t *testing.T) {
	resConfig30 := v1.Resource{
		ID:   "apps.kusionstack.io/v1alpha1:PodTransitionRule:foo:default-dev-foo",
		Type: "Kubernetes",
		Attributes: map[string]interface{}{
			"apiVersion": "apps.kusionstack.io/v1alpha1",
			"kind":       "PodTransitionRule",
			"metadata": map[string]interface{}{
				"creationTimestamp": interface{}(nil),
				"name":              "default-dev-foo",
				"namespace":         "foo",
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
		ID:   "apps.kusionstack.io/v1alpha1:PodTransitionRule:foo:default-dev-foo",
		Type: "Kubernetes",
		Attributes: map[string]interface{}{
			"apiVersion": "apps.kusionstack.io/v1alpha1",
			"kind":       "PodTransitionRule",
			"metadata": map[string]interface{}{
				"creationTimestamp": interface{}(nil),
				"name":              "default-dev-foo",
				"namespace":         "foo",
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

	jobWorkloadConfig := &workload.Workload{
		Header: workload.Header{
			Type: workload.TypeJob,
		},
	}
	serviceWorkloadConfig := &workload.Workload{
		Header: workload.Header{
			Type: workload.TypeService,
		},
		Service: &workload.Service{
			Type: workload.Collaset,
		},
	}
	ops30 := map[string]interface{}{
		"maxUnavailable": "30%",
	}
	ops40 := map[string]interface{}{
		"maxUnavailable": 40,
	}

	devConfig := jsonutil.Marshal2String(ops30)
	platformConfig := jsonutil.Marshal2String(ops40)
	jobWorkload := jsonutil.Marshal2String(jobWorkloadConfig)
	serviceWorkload := jsonutil.Marshal2String(serviceWorkloadConfig)
	res30 := jsonutil.Marshal2String(resConfig30)
	res40 := jsonutil.Marshal2String(resConfig40)
	response30 := &proto.GeneratorResponse{
		Resources: [][]byte{[]byte(res30)},
	}
	response40 := &proto.GeneratorResponse{
		Resources: [][]byte{[]byte(res40)},
	}

	project := "default"
	stack := "dev"
	app := "foo"

	type args struct {
		r *proto.GeneratorRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *proto.GeneratorResponse
		wantErr bool
	}{
		{
			name: "test Job",
			args: args{
				r: &proto.GeneratorRequest{
					Project:              project,
					Stack:                stack,
					App:                  app,
					Workload:             []byte(jobWorkload),
					DevModuleConfig:      []byte(devConfig),
					PlatformModuleConfig: []byte(platformConfig),
					RuntimeConfig:        nil,
				},
			},
			want: generators.EmptyResponse(),
		},
		{
			name: "test CollaSet with opsRule in appConfig",
			args: args{
				r: &proto.GeneratorRequest{
					Project:              project,
					Stack:                stack,
					App:                  app,
					Workload:             []byte(serviceWorkload),
					DevModuleConfig:      []byte(devConfig),
					PlatformModuleConfig: []byte(platformConfig),
				},
			},
			wantErr: false,
			want:    response30,
		},
		{
			name: "test CollaSet with opsRule in workspace",
			args: args{
				r: &proto.GeneratorRequest{
					Project:              project,
					Stack:                stack,
					App:                  app,
					Workload:             []byte(serviceWorkload),
					PlatformModuleConfig: []byte(platformConfig),
				},
			},
			wantErr: false,
			want:    response40,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OpsRuleModule{}
			got, err := o.Generate(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Generate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
