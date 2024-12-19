package main

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	k8snetworking "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"

	"kusionstack.io/kusion-module-framework/pkg/util/workspace"
)

func (network *Network) Generate(ctx context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	// Get the module logger with the generator context.
	logger := log.GetModuleLogger(ctx)
	logger.Info("Generating resources...")

	defer func() {
		if r := recover(); r != nil {
			logger.Debug("failed to generate network module: %v", r)
			response = nil
			rawRequest, _ := json.Marshal(request)
			err = fmt.Errorf("panic in network module generator but recovered with error: [%v] and stack %v and request %v",
				r, string(debug.Stack()), string(rawRequest))
		}
	}()

	// Network does not exist in AppConfiguration configs.
	if request.DevConfig == nil {
		logger.Info("Network does not exist in AppConfig config")

		return nil, nil
	}

	// Get the complete configs of the Network accessory.
	if err := network.GetCompleteConfig(request.DevConfig, request.PlatformConfig); err != nil {
		return nil, err
	}
	if len(network.Ports) != 0 && request.Workload == nil {
		return nil, ErrEmptySvcWorkload
	}

	var resources []kusionapiv1.Resource
	// Generate network port related resources.
	res, err := network.GeneratePortResources(request)
	if err != nil {
		return nil, err
	}
	resources = append(resources, res...)

	// Generate network ingress related resources.
	ingressRes, err := network.GenerateIngressResource(request)
	if err != nil {
		return nil, err
	}
	if ingressRes != nil {
		resources = append(resources, *ingressRes)
	}

	// Generate network ingressClass related resources.
	ingressClassRes, err := network.GenerateIngressClassResource(request)
	if err != nil {
		return nil, err
	}
	if ingressClassRes != nil {
		resources = append(resources, *ingressClassRes)
	}

	return &module.GeneratorResponse{
		Resources: resources,
	}, nil
}

// GetCompleteConfig combines the configs in devModuleConfig and platformModuleConfig to form a complete
// configuration for the Network accessory.
func (network *Network) GetCompleteConfig(devConfig kusionapiv1.Accessory, platformConfig kusionapiv1.GenericConfig) error {
	// Get the complete port config.
	if err := network.CompletePortConfig(devConfig, platformConfig); err != nil {
		return err
	}

	if err := network.CompleteIngressConfig(devConfig); err != nil {
		return err
	}

	return network.Validate()
}

// CompletePortConfig completes the network port related config.
func (network *Network) CompletePortConfig(devConfig kusionapiv1.Accessory, platformConfig kusionapiv1.GenericConfig) error {
	if devConfig != nil {
		ports, ok := devConfig["ports"]
		if ok {
			for _, port := range ports.([]interface{}) {
				// Retrieve port configs from the devConfig based on the result
				// of the type assertion.
				mp, err := toMapStringInterface(port)
				if err != nil {
					return fmt.Errorf("failed to retrieve port from dev config: %v", err)
				}

				yamlStr, err := yaml.Marshal(mp)
				if err != nil {
					return err
				}

				var p Port
				if err := yaml.Unmarshal(yamlStr, &p); err != nil {
					return err
				}

				network.Ports = append(network.Ports, p)
			}
		}

		ingressConf, ok := devConfig["ingress"]
		if ok {
			ingressYaml, err := yaml.Marshal(ingressConf)
			if err != nil {
				return err
			}
			var ingress Ingress
			err = yaml.Unmarshal(ingressYaml, &ingress)
			if err != nil {
				return err
			}
			network.Ingress = &ingress
		}

		ingressClassConf, ok := devConfig["ingressClass"]
		if ok {
			ingressClassYaml, err := yaml.Marshal(ingressClassConf)
			if err != nil {
				return err
			}
			var ingressClass IngressClass
			err = yaml.Unmarshal(ingressClassYaml, &ingressClass)
			if err != nil {
				return err
			}
			network.IngressClass = &ingressClass
		}
	}

	var portConfig kusionapiv1.GenericConfig
	if platformConfig != nil {
		pc, ok := platformConfig["port"]
		if ok {
			// Retrieve port configs from the platformConfig based on the result
			// of the type assertion.
			mpc, err := toMapStringInterface(pc)
			if err != nil {
				return fmt.Errorf("failed to retrieve port from platform config: %v", err)
			}

			yamlStr, err := yaml.Marshal(mpc)
			if err != nil {
				return err
			}

			if err := yaml.Unmarshal(yamlStr, &portConfig); err != nil {
				return err
			}
		}
	}

	for i := range network.Ports {
		if network.Ports[i].TargetPort == 0 {
			network.Ports[i].TargetPort = network.Ports[i].Port
		}
		if network.Ports[i].Public {
			// Get port type from platform config.
			if portConfig == nil {
				return ErrEmptyPortConfig
			}
			portType, err := workspace.GetStringFromGenericConfig(portConfig, FieldType)
			if err != nil {
				return err
			}
			if portType == "" {
				return ErrEmptyType
			}
			if portType != CSPAWS && portType != CSPAliCloud {
				return ErrUnsupportedType
			}
			network.Ports[i].Type = portType

			// Get labels from platform config.
			labels, err := workspace.GetStringMapFromGenericConfig(portConfig, FieldLabels)
			if err != nil {
				return err
			}
			network.Ports[i].Labels = labels

			// Get annotations from platform config.
			annotations, err := workspace.GetStringMapFromGenericConfig(portConfig, FieldAnnotations)
			if err != nil {
				return err
			}
			network.Ports[i].Annotations = annotations
		}
	}

	return nil
}

// Validate validates whether the input of a Network accessory is valid.
func (network *Network) Validate() error {
	// Validate the port config.
	if err := network.ValidatePortConfig(); err != nil {
		return err
	}

	return nil
}

// CompleteIngressConfig completes the network ingress related config.
func (network *Network) CompleteIngressConfig(devConfig kusionapiv1.Accessory) error {
	if devConfig != nil {
		ingressConf, ok := devConfig["ingress"]
		if ok {
			ingressYaml, err := yaml.Marshal(ingressConf)
			if err != nil {
				return err
			}
			var ingress Ingress
			err = yaml.Unmarshal(ingressYaml, &ingress)
			if err != nil {
				return err
			}
			network.Ingress = &ingress
		}

		ingressClassConf, ok := devConfig["ingressClass"]
		if ok {
			ingressClassYaml, err := yaml.Marshal(ingressClassConf)
			if err != nil {
				return err
			}
			var ingressClass IngressClass
			err = yaml.Unmarshal(ingressClassYaml, &ingressClass)
			if err != nil {
				return err
			}
			network.IngressClass = &ingressClass
		}
	}
	return nil
}

// ValidatePortConfig validates whether the port configs are valid or not.
func (network *Network) ValidatePortConfig() error {
	for _, port := range network.Ports {
		if port.Port < 1 || port.Port > 65535 {
			return ErrInvalidPort
		}
		if port.TargetPort < 1 || port.TargetPort > 65535 {
			return ErrInvalidTargetPort
		}
		if port.Protocol != ProtocolTCP && port.Protocol != ProtocolUDP {
			return ErrInvalidProtocol
		}
	}

	return nil
}

// GeneratePortResources generates the resources related to the network port.
func (network *Network) GeneratePortResources(request *module.GeneratorRequest) ([]kusionapiv1.Resource, error) {
	var resources []kusionapiv1.Resource
	privatePorts, publicPorts := splitPorts(network.Ports)
	if len(privatePorts) != 0 {
		svc := generatePortK8sSvc(request, false, privatePorts)
		resourceID := module.KubernetesResourceID(svc.TypeMeta, svc.ObjectMeta)
		resource, err := module.WrapK8sResourceToKusionResource(resourceID, svc)
		if err != nil {
			return nil, err
		}
		resources = append(resources, *resource)
	}
	if len(publicPorts) != 0 {
		svc := generatePortK8sSvc(request, true, publicPorts)
		resourceID := module.KubernetesResourceID(svc.TypeMeta, svc.ObjectMeta)
		resource, err := module.WrapK8sResourceToKusionResource(resourceID, svc)
		if err != nil {
			return nil, err
		}
		resources = append(resources, *resource)
	}

	return resources, nil
}

// generatePortK8sSvc generates the Kubernetes Service resource for the network port.
func generatePortK8sSvc(request *module.GeneratorRequest, public bool, ports []Port) *v1.Service {
	appUname := module.UniqueAppName(request.Project, request.Stack, request.App)
	var name string
	if public {
		name = fmt.Sprintf("%s-%s", appUname, suffixPublic)
	} else {
		name = fmt.Sprintf("%s-%s", appUname, suffixPrivate)
	}
	svcType := v1.ServiceTypeClusterIP
	if public {
		svcType = v1.ServiceTypeLoadBalancer
	}

	svcLabels, ok := request.Workload["labels"]
	if !ok {
		svcLabels = make(map[string]string)
	}

	svcAnnotations, ok := request.Workload["annotations"]
	if !ok {
		svcAnnotations = make(map[string]string)
	}

	labels := module.MergeMaps(module.UniqueAppLabels(request.Project, request.App), svcLabels.(map[string]string))
	annotations := module.MergeMaps(svcAnnotations.(map[string]string))
	selector := module.UniqueAppLabels(request.Project, request.App)

	svc := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			APIVersion: v1.SchemeGroupVersion.String(),
			Kind:       k8sKindService,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   request.Project,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: v1.ServiceSpec{
			Ports:    toSvcPorts(name, ports),
			Selector: selector,
			Type:     svcType,
		},
	}

	if public {
		if len(svc.Labels) == 0 {
			svc.Labels = make(map[string]string)
		}
		if len(svc.Annotations) == 0 {
			svc.Annotations = make(map[string]string)
		}

		labels := ports[0].Labels
		for k, v := range labels {
			svc.Labels[k] = v
		}
		annotations := ports[0].Annotations
		for k, v := range annotations {
			svc.Annotations[k] = v
		}
	}

	return svc
}

// splitPorts splits the network ports into private ports and public ports.
func splitPorts(ports []Port) ([]Port, []Port) {
	var privatePorts, publicPorts []Port
	for _, port := range ports {
		if port.Public {
			publicPorts = append(publicPorts, port)
		} else {
			privatePorts = append(privatePorts, port)
		}
	}
	return privatePorts, publicPorts
}

// toSvcPorts returns the Kubernetes ServicePort resource.
func toSvcPorts(name string, ports []Port) []v1.ServicePort {
	svcPorts := make([]v1.ServicePort, len(ports))
	for i, port := range ports {
		svcPorts[i] = v1.ServicePort{
			Name:       fmt.Sprintf("%s-%d-%s", name, port.Port, strings.ToLower(port.Protocol)),
			Port:       int32(port.Port),
			TargetPort: intstr.FromInt(port.TargetPort),
			Protocol:   v1.Protocol(port.Protocol),
		}
	}
	return svcPorts
}

// toMapStringInterface changes the input interface (usually map[interface{}]interface{})
// into map[string]interface{}.
func toMapStringInterface(i any) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if p, ok := i.(map[interface{}]interface{}); ok {
		for k, v := range p {
			m[fmt.Sprintf("%v", k)] = v
		}
	} else if p, ok := i.(map[string]interface{}); ok {
		m = p
	} else if p, ok := i.(kusionapiv1.Accessory); ok {
		m = map[string]interface{}(p)
	} else if p, ok := i.(kusionapiv1.GenericConfig); ok {
		m = map[string]interface{}(p)
	} else {
		return nil, fmt.Errorf("unexpected type: %T", i)
	}

	return m, nil
}

// GenerateIngressResource generates the resources related to the network ingress.
func (network *Network) GenerateIngressResource(request *module.GeneratorRequest) (*kusionapiv1.Resource, error) {
	if network.Ingress == nil {
		return nil, nil
	}
	ingress, err := network.generateIngress(request)
	if err != nil {
		return nil, err
	}
	resourceID := module.KubernetesResourceID(ingress.TypeMeta, ingress.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, ingress)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (network *Network) generateIngress(request *module.GeneratorRequest) (*k8snetworking.Ingress, error) {
	appUname := module.UniqueAppName(request.Project, request.Stack, request.App)
	resourceName := fmt.Sprintf("%s-%s", appUname, ingressSuffix)
	k8sIngress := &k8snetworking.Ingress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8snetworking.SchemeGroupVersion.String(),
			Kind:       K8sKindIngress,
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:      network.Ingress.Labels,
			Annotations: network.Ingress.Annotations,
			Name:        resourceName,
			Namespace:   request.Project,
		},
		Spec: k8snetworking.IngressSpec{
			IngressClassName: network.Ingress.IngressClassName,
		},
	}

	if network.Ingress.DefaultBackend != nil {
		defaultBackend, err := network.toIngressBackend(*network.Ingress.DefaultBackend, appUname)
		if err != nil {
			return nil, err
		}
		k8sIngress.Spec.DefaultBackend = defaultBackend
	}

	for _, t := range network.Ingress.TLS {
		tls := k8snetworking.IngressTLS{
			Hosts:      t.Hosts,
			SecretName: t.SecretName,
		}
		k8sIngress.Spec.TLS = append(k8sIngress.Spec.TLS, tls)
	}

	var rules []k8snetworking.IngressRule
	for _, r := range network.Ingress.Rules {
		var rule k8snetworking.IngressRule
		if r.HTTP != nil {
			var paths []k8snetworking.HTTPIngressPath
			for _, p := range r.HTTP.Paths {
				httpPath := k8snetworking.HTTPIngressPath{Path: p.Path}
				if p.PathType != "" {
					httpPath.PathType = &p.PathType
				}

				backend, err := network.toIngressBackend(p.Backend, appUname)
				if err != nil {
					return nil, err
				}
				if backend != nil {
					httpPath.Backend = *backend
				}
				paths = append(paths, httpPath)
			}
			rule.HTTP = &k8snetworking.HTTPIngressRuleValue{
				Paths: paths,
			}
		}
		rule.Host = r.Host
		rules = append(rules, rule)
	}
	k8sIngress.Spec.Rules = rules
	return k8sIngress, nil
}

func (network *Network) toIngressBackend(b IngressBackend, appUname string) (*k8snetworking.IngressBackend, error) {
	var backend k8snetworking.IngressBackend
	if b.Service != nil {
		svcName := b.Service.Name
		if b.Service.Name == "" {
			foundPort := false
			for _, port := range network.Ports {
				if b.Service.Port.Number == int32(port.Port) {
					if port.Public {
						svcName = fmt.Sprintf("%s-%s", appUname, suffixPublic)
					} else {
						svcName = fmt.Sprintf("%s-%s", appUname, suffixPrivate)
					}
					foundPort = true
					break
				}
			}
			if !foundPort {
				return nil, fmt.Errorf("not found available service for backend, please check service name or port")
			}
		}

		backend.Service = &k8snetworking.IngressServiceBackend{
			Name: svcName,
			Port: k8snetworking.ServiceBackendPort{
				Name:   b.Service.Port.Name,
				Number: b.Service.Port.Number,
			},
		}
	}

	if b.Resource != nil {
		backend.Resource = &v1.TypedLocalObjectReference{
			APIGroup: b.Resource.APIGroup,
			Kind:     b.Resource.Kind,
			Name:     b.Resource.Name,
		}
	}
	return &backend, nil
}

// GenerateIngressClassResource generates the resources related to the network ingressClass.
func (network *Network) GenerateIngressClassResource(request *module.GeneratorRequest) (*kusionapiv1.Resource, error) {
	if network.IngressClass == nil {
		return nil, nil
	}
	ingressClass := network.generateIngressClass(request)
	resourceID := module.KubernetesResourceID(ingressClass.TypeMeta, ingressClass.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, ingressClass)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (network *Network) generateIngressClass(request *module.GeneratorRequest) *k8snetworking.IngressClass {
	appUname := module.UniqueAppName(request.Project, request.Stack, request.App)
	resourceName := fmt.Sprintf("%s-%s", appUname, ingressClassSuffix)
	k8sIngressClass := &k8snetworking.IngressClass{
		TypeMeta: metav1.TypeMeta{
			APIVersion: k8snetworking.SchemeGroupVersion.String(),
			Kind:       K8sKindIngressClass,
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:      network.IngressClass.Labels,
			Annotations: network.IngressClass.Annotations,
			Name:        resourceName,
		},
		Spec: k8snetworking.IngressClassSpec{
			Controller: network.IngressClass.Controller,
		},
	}

	if network.IngressClass.Parameters != nil {
		k8sIngressClass.Spec.Parameters = &k8snetworking.IngressClassParametersReference{
			APIGroup:  network.IngressClass.Parameters.APIGroup,
			Kind:      network.IngressClass.Parameters.Kind,
			Name:      network.IngressClass.Parameters.Name,
			Scope:     network.IngressClass.Parameters.Scope,
			Namespace: network.IngressClass.Parameters.Namespace,
		}
	}
	return k8sIngressClass
}

func main() {
	server.Start(&Network{})
}
