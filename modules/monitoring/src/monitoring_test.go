package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/modules"

	"kusionstack.io/kusion-module-framework/pkg/module"
)

type TestCase struct {
	name    string
	request module.GeneratorRequest
	want    *module.GeneratorResponse
	wantErr bool
}

func BuildMonitoringTestCase(
	testName, projectName, stackName, appName string,
	interval, timeout, path, port, scheme, monitorType string,
	operatorMode, wantErr bool,
) *TestCase {
	var endpointType string
	var monitorKind MonitorType
	if monitorType == string(ServiceMonitorType) {
		monitorKind = "ServiceMonitor"
		endpointType = "endpoints"
	} else if monitorType == string(PodMonitorType) {
		monitorKind = "PodMonitor"
		endpointType = "podMetricsEndpoints"
	}
	var expectedResources []apiv1.Resource
	var expectedPatcher apiv1.Patcher
	uniqueName := modules.UniqueAppName(projectName, stackName, appName)
	if operatorMode {
		expectedResources = []apiv1.Resource{
			{
				ID:   fmt.Sprintf("monitoring.coreos.com/v1:%s:%s:%s-%s-monitor", monitorKind, projectName, uniqueName, strings.ToLower(monitorType)),
				Type: "Kubernetes",
				Attributes: map[string]interface{}{
					"apiVersion": "monitoring.coreos.com/v1",
					"kind":       string(monitorKind),
					"metadata": map[string]interface{}{
						"creationTimestamp": nil,
						"name":              fmt.Sprintf("%s-%s-monitor", uniqueName, strings.ToLower(monitorType)),
						"namespace":         projectName,
					},
					"spec": map[string]interface{}{
						endpointType: []interface{}{
							map[string]interface{}{
								"bearerTokenSecret": map[string]interface{}{
									"key": "",
								},
								"interval":      interval,
								"scrapeTimeout": timeout,
								"path":          path,
								"port":          port,
								"scheme":        scheme,
							},
						},
						"namespaceSelector": make(map[string]interface{}),
						"selector": map[string]interface{}{
							"matchLabels": map[string]interface{}{
								"kusion_monitoring_appname": appName,
							},
						},
					},
				},
				DependsOn: nil,
				Extensions: map[string]interface{}{
					"GVK": fmt.Sprintf("monitoring.coreos.com/v1, Kind=%s", string(monitorKind)),
				},
			},
		}
	} else {
		expectedPatcher = apiv1.Patcher{
			Annotations: map[string]string{
				"prometheus.io/scrape": "true",
				"prometheus.io/path":   path,
				"prometheus.io/port":   port,
				"prometheus.io/scheme": scheme,
			},
		}
	}
	testCase := &TestCase{
		name: testName,
		request: module.GeneratorRequest{
			Project: projectName,
			Stack:   stackName,
			App:     appName,
			PlatformConfig: apiv1.GenericConfig{
				TimeoutKey:      timeout,
				IntervalKey:     interval,
				SchemeKey:       scheme,
				OperatorModeKey: operatorMode,
				MonitorTypeKey:  monitorType,
			},
			DevConfig: apiv1.Accessory{
				PathKey: path,
				PortKey: port,
			},
		},
		want: &module.GeneratorResponse{
			Resources: expectedResources,
			Patcher:   &expectedPatcher,
		},
		wantErr: wantErr,
	}
	return testCase
}

func TestMonitoringGenerator_Generate(t *testing.T) {
	ctx := context.TODO()
	tests := []TestCase{
		*BuildMonitoringTestCase("ServiceMonitorTest", "test-project", "test-stack", "test-app", "15s", "5s", "/metrics", "web", "http", "Service", true, false),
		*BuildMonitoringTestCase("ServiceMonitorTest2", "test-project-2", "test-stack-2", "test-app-2", "15s", "5s", "/metrics", "web", "http", "Service", true, false),
		*BuildMonitoringTestCase("PodMonitorTest", "test-project", "test-stack", "test-app", "15s", "5s", "/metrics", "web", "http", "Pod", true, false),
		*BuildMonitoringTestCase("ServiceAnnotationTest", "test-project", "test-stack", "test-app", "30s", "15s", "/metrics", "8080", "http", "Service", false, false),
		*BuildMonitoringTestCase("PodAnnotationTest", "test-project", "test-stack", "test-app", "30s", "15s", "/metrics", "8080", "http", "Pod", false, false),
		*BuildMonitoringTestCase("InvalidDurationTest", "test-project", "test-stack", "test-app", "15s", "5ssss", "/metrics", "8080", "http", "Pod", false, true),
		*BuildMonitoringTestCase("InvalidTimeoutTest", "test-project", "test-stack", "test-app", "15s", "30s", "/metrics", "8080", "http", "Pod", false, true),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &MonitoringModule{}
			response, err := g.Generate(ctx, &tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Generate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				require.Equal(t, tt.want.Resources, response.Resources)
				require.Equal(t, *tt.want.Patcher, *response.Patcher)
			}
		})
	}
}
