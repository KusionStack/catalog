package main

import (
	"context"
	"fmt"
	"time"

	prometheusv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/log"
	"kusionstack.io/kusion/pkg/modules"
)

func (g *MonitoringModule) Generate(_ context.Context, request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	// Monitoring does not exist in AppConfig and workspace config
	if request.DevConfig == nil && request.PlatformConfig == nil {
		log.Info("Monitoring does not exist in either AppConfig and workspace config")
		return nil, nil
	}

	log.Debug("request.DevModuleConfig: ", request.DevConfig)
	log.Debug("request.PlatformModuleConfig: ", request.PlatformConfig)
	// Parse workspace configurations for monitoring generator.
	if err := g.parseWorkspaceConfig(request.DevConfig, request.PlatformConfig); err != nil {
		return nil, err
	}

	// If operator mode is enabled, create monitor objects.
	if g != nil && g.OperatorMode {
		log.Info("Operator mode is enabled. Creating monitor objects...")
		if g.MonitorType == ServiceMonitorType {
			// Create ServiceMonitor if MonitorType is Service
			serviceMonitor, err := g.buildMonitorObject(request, g.MonitorType)
			if err != nil {
				return nil, err
			}
			resourceID := module.KubernetesResourceID(serviceMonitor.(*prometheusv1.ServiceMonitor).TypeMeta, serviceMonitor.(*prometheusv1.ServiceMonitor).ObjectMeta)
			resource, err := module.WrapK8sResourceToKusionResource(resourceID, serviceMonitor)
			if err != nil {
				return nil, err
			}
			patcher := &v1.Patcher{
				Labels: map[string]string{
					"kusion_monitoring_appname": request.App,
				},
			}
			return &module.GeneratorResponse{
				Resources: []v1.Resource{*resource},
				Patcher:   patcher,
			}, nil
		} else if g.MonitorType == PodMonitorType {
			// Create PodMonitor if MonitorType is Pod
			podMonitor, err := g.buildMonitorObject(request, g.MonitorType)
			if err != nil {
				return nil, err
			}
			resourceID := module.KubernetesResourceID(podMonitor.(*prometheusv1.PodMonitor).TypeMeta, podMonitor.(*prometheusv1.PodMonitor).ObjectMeta)
			resource, err := module.WrapK8sResourceToKusionResource(resourceID, podMonitor)
			if err != nil {
				return nil, err
			}
			patcher := &v1.Patcher{
				Labels: map[string]string{
					"kusion_monitoring_appname": request.App,
				},
			}
			return &module.GeneratorResponse{
				Resources: []v1.Resource{*resource},
				Patcher:   patcher,
			}, nil
		} else {
			return nil, fmt.Errorf("MonitorType should either be service or pod %s", g.MonitorType)
		}
	} else {
		// Operator mode is disabled. Patching workload annotations
		log.Info("Operator mode is disabled. Patching workload annotations...")
		// Patch workload annotations
		annotations := map[string]string{
			"prometheus.io/scrape": "true",
			"prometheus.io/path":   g.Path,
			"prometheus.io/port":   g.Port,
			"prometheus.io/scheme": g.Scheme,
		}
		patchers := &v1.Patcher{
			Annotations: annotations,
		}
		return &module.GeneratorResponse{
			Patcher: patchers,
		}, nil
	}
}

func main() {
	server.Start(&MonitoringModule{})
}

// parseWorkspaceConfig parses the config items for monitoring generator in workspace configurations.
func (g *MonitoringModule) parseWorkspaceConfig(devConfig v1.Accessory, workspaceConfig v1.GenericConfig) error {
	// get path and port from devConfig
	if path, ok := devConfig[PathKey]; ok {
		g.Path = path.(string)
	}
	if port, ok := devConfig[PortKey]; ok {
		g.Port = port.(string)
	}

	// get operatorMode, monitorType, interval, timeout, scheme from workspaceConfig
	if operatorMode, ok := workspaceConfig[OperatorModeKey]; ok {
		g.OperatorMode = operatorMode.(bool)
	}

	if monitorType, ok := workspaceConfig[MonitorTypeKey]; ok {
		g.MonitorType = MonitorType(monitorType.(string))
	} else {
		g.MonitorType = DefaultMonitorType
	}

	if interval, ok := workspaceConfig[IntervalKey]; ok {
		g.Interval = prometheusv1.Duration(interval.(string))
	} else {
		g.Interval = DefaultInterval
	}

	if timeout, ok := workspaceConfig[TimeoutKey]; ok {
		g.Timeout = prometheusv1.Duration(timeout.(string))
	} else {
		g.Timeout = DefaultTimeout
	}

	if scheme, ok := workspaceConfig[SchemeKey]; ok {
		g.Scheme = scheme.(string)
	} else {
		g.Scheme = DefaultScheme
	}

	// validate the monitoring configuration
	parsedTimeout, err := time.ParseDuration(string(g.Timeout))
	if err != nil {
		return err
	}
	parsedInterval, err := time.ParseDuration(string(g.Interval))
	if err != nil {
		return err
	}

	// validate the monitoring configuration
	if parsedTimeout > parsedInterval {
		return ErrTimeoutGreaterThanInterval
	}

	return nil
}

func (g *MonitoringModule) buildMonitorObject(request *module.GeneratorRequest, monitorType MonitorType) (runtime.Object, error) {
	// If Prometheus runs as an operator, it relies on Custom Resources to
	// manage the scrape configs. CRs (ServiceMonitors and PodMonitors) rely on
	// corresponding resources (Services and Pods) to have labels that can be
	// used as part of the label selector for the CR to determine which
	// service/pods to scrape from.
	// Here we choose the label name kusion_monitoring_appname for two reasons:
	// 1. Unlike the label validation in Kubernetes, the label name accepted by
	// Prometheus cannot contain non-alphanumeric characters except underscore:
	// https://github.com/prometheus/common/blob/main/model/labels.go#L94
	// 2. The name should be unique enough that is only created by Kusion and
	// used to identify a certain application
	monitoringLabels := map[string]string{
		"kusion_monitoring_appname": request.App,
	}

	// Create ServiceMonitor or PodMonitor based on the monitorType
	if monitorType == ServiceMonitorType {
		// Create ServiceMonitor
		serviceEndpoint := prometheusv1.Endpoint{
			Interval:      g.Interval,
			ScrapeTimeout: g.Timeout,
			Port:          g.Port,
			Path:          g.Path,
			Scheme:        g.Scheme,
			BearerTokenSecret: &corev1.SecretKeySelector{
				Key: "",
			},
		}
		serviceEndpointList := []prometheusv1.Endpoint{serviceEndpoint}
		serviceMonitor := &prometheusv1.ServiceMonitor{
			TypeMeta: metav1.TypeMeta{
				Kind:       "ServiceMonitor",
				APIVersion: prometheusv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-service-monitor", modules.UniqueAppName(request.Project, request.Stack, request.App)),
				Namespace: request.Project,
			},
			Spec: prometheusv1.ServiceMonitorSpec{
				Selector: metav1.LabelSelector{
					MatchLabels: monitoringLabels,
				},
				Endpoints: serviceEndpointList,
			},
		}
		return serviceMonitor, nil
	} else if monitorType == PodMonitorType {
		// Create PodMonitor
		podMetricsEndpoint := prometheusv1.PodMetricsEndpoint{
			Interval:      g.Interval,
			ScrapeTimeout: g.Timeout,
			Port:          g.Port,
			Path:          g.Path,
			Scheme:        g.Scheme,
		}
		podMetricsEndpointList := []prometheusv1.PodMetricsEndpoint{podMetricsEndpoint}

		podMonitor := &prometheusv1.PodMonitor{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PodMonitor",
				APIVersion: prometheusv1.SchemeGroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-pod-monitor", modules.UniqueAppName(request.Project, request.Stack, request.App)),
				Namespace: request.Project,
			},
			Spec: prometheusv1.PodMonitorSpec{
				Selector: metav1.LabelSelector{
					MatchLabels: monitoringLabels,
				},
				PodMetricsEndpoints: podMetricsEndpointList,
			},
		}
		return podMonitor, nil
	}

	return nil, fmt.Errorf("MonitorType should either be service or pod %s", monitorType)
}
