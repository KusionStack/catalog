package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"kusionstack.io/kube-api/apps/v1alpha1"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	"kusionstack.io/kusion-module-framework/pkg/util/workspace"
)

var (
	ErrEmptySelectors        = errors.New("selectors must not be empty")
	ErrInvalidPort           = errors.New("port must be between 1 and 65535")
	ErrInvalidTargetPort     = errors.New("targetPort must be between 1 and 65535 if exist")
	ErrInvalidProtocol       = errors.New("protocol must be TCP or UDP")
	ErrDuplicatePortProtocol = errors.New("port-protocol pair must not be duplicate")
)

func (svc *Service) Generate(ctx context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	// Get the module logger with the generator context.
	logger := log.GetModuleLogger(ctx)
	logger.Info("Generating resources...")

	defer func() {
		if r := recover(); r != nil {
			logger.Debug("failed to generate Service module: %v", r)
			response = nil
			rawRequest, _ := json.Marshal(request)
			err = fmt.Errorf("panic in service module generator but recovered with error: [%v] and stack %v and request %v",
				r, string(debug.Stack()), string(rawRequest))
		}
	}()

	if request.DevConfig == nil {
		logger.Info("Service does not exist in AppConfig config")
		return nil, nil
	}
	out, err := yaml.Marshal(request.DevConfig)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(out, svc); err != nil {
		return nil, fmt.Errorf("complete Service by dev config failed, %w", err)
	}

	if err = completeServiceInput(svc, request.PlatformConfig); err != nil {
		return nil, fmt.Errorf("complete Service by platform config failed, %w", err)
	}

	uniqueAppName := module.UniqueAppName(request.Project, request.Stack, request.App)

	// Create a slice of containers based on the App's containers along with related volumes and configMaps.
	containers, volumes, configMaps, err := toOrderedContainers(svc.Containers, uniqueAppName)
	if err != nil {
		return nil, err
	}

	topologySpreadConstraints := handleTopologySpreadConstraints(svc.TopologySpreadConstraints)

	res := make([]kusionapiv1.Resource, 0)
	// Create ConfigMap objects based on the App's configuration.
	for _, cm := range configMaps {
		cm.Namespace = request.Project
		resourceID := module.KubernetesResourceID(cm.TypeMeta, cm.ObjectMeta)
		resource, err := module.WrapK8sResourceToKusionResource(resourceID, &cm)
		if err != nil {
			return nil, err
		}
		res = append(res, *resource)
	}

	labels := module.MergeMaps(module.UniqueAppLabels(request.Project, request.App), svc.Labels)
	annotations := module.MergeMaps(svc.Annotations)
	selectors := module.UniqueAppLabels(request.Project, request.App)

	// Create a K8s Workload object based on the App's configuration.
	// common parts
	objectMeta := metav1.ObjectMeta{
		Labels:      labels,
		Annotations: annotations,
		Name:        uniqueAppName,
		Namespace:   request.Project,
	}
	podTemplateSpec := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: corev1.PodSpec{
			TopologySpreadConstraints: topologySpreadConstraints,
			Containers:                containers,
			Volumes:                   volumes,
		},
	}

	var k8sResource runtime.Object
	typeMeta := metav1.TypeMeta{}

	switch svc.Type {
	case Deployment:
		typeMeta = metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       string(Deployment),
		}
		spec := appsv1.DeploymentSpec{
			Replicas: svc.Replicas,
			Selector: &metav1.LabelSelector{MatchLabels: selectors},
			Template: podTemplateSpec,
		}
		k8sResource = &appsv1.Deployment{
			TypeMeta:   typeMeta,
			ObjectMeta: objectMeta,
			Spec:       spec,
		}
	case Collaset:
		typeMeta = metav1.TypeMeta{
			APIVersion: v1alpha1.GroupVersion.String(),
			Kind:       string(Collaset),
		}
		k8sResource = &v1alpha1.CollaSet{
			TypeMeta:   typeMeta,
			ObjectMeta: objectMeta,
			Spec: v1alpha1.CollaSetSpec{
				Replicas: svc.Replicas,
				Selector: &metav1.LabelSelector{MatchLabels: selectors},
				Template: podTemplateSpec,
			},
		}
	}

	// append the Deployment/Collaset resource to res.
	resourceID := module.KubernetesResourceID(typeMeta, objectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, k8sResource)
	if err != nil {
		return nil, err
	}
	res = append(res, *resource)

	// validate and complete service ports
	if len(svc.Ports) != 0 {
		if err = validate(selectors, svc.Ports); err != nil {
			return nil, err
		}
		if err = complete(svc.Ports); err != nil {
			return nil, err
		}
	}
	response = &module.GeneratorResponse{
		Resources: res,
	}

	return response, nil
}

func validatePorts(ports []Port) error {
	portProtocolRecord := make(map[string]struct{})
	for _, port := range ports {
		if err := validatePort(&port); err != nil {
			return fmt.Errorf("invalid port config %+v, %w", port, err)
		}

		// duplicate "port-protocol" pairs are not allowed.
		portProtocol := fmt.Sprintf("%d-%s", port.Port, port.Protocol)
		if _, ok := portProtocolRecord[portProtocol]; ok {
			return fmt.Errorf("invalid port config %+v, %v", port, ErrDuplicatePortProtocol)
		}
		portProtocolRecord[portProtocol] = struct{}{}
	}
	return nil
}

func validatePort(port *Port) error {
	if port.Port < 1 || port.Port > 65535 {
		return ErrInvalidPort
	}
	if port.TargetPort < 0 || port.Port > 65535 {
		return ErrInvalidTargetPort
	}
	if port.Protocol != TCP && port.Protocol != UDP {
		return ErrInvalidProtocol
	}
	return nil
}

func validate(selectors map[string]string, ports []Port) error {
	if len(selectors) == 0 {
		return ErrEmptySelectors
	}
	if err := validatePorts(ports); err != nil {
		return err
	}
	return nil
}

func complete(ports []Port) error {
	for i := range ports {
		if ports[i].TargetPort == 0 {
			ports[i].TargetPort = ports[i].Port
		}
	}
	return nil
}

func completeServiceInput(service *Service, config kusionapiv1.GenericConfig) error {
	if err := completeBaseWorkload(&service.Base, config); err != nil {
		return err
	}
	serviceTypeStr, err := workspace.GetStringFromGenericConfig(config, ModuleServiceType)
	platformServiceType := ServiceType(serviceTypeStr)
	if err != nil {
		return err
	}
	// if not set in workspace, use Deployment as default type
	if platformServiceType == "" {
		platformServiceType = Deployment
	}
	if platformServiceType != Deployment && platformServiceType != Collaset {
		return fmt.Errorf("unsupported Service type %s", platformServiceType)
	}
	if service.Type == "" {
		service.Type = platformServiceType
	}
	return nil
}

func main() {
	server.Start(&Service{})
}
