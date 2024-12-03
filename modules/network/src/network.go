package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"gopkg.in/yaml.v3"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"

	"kusionstack.io/kusion-module-framework/pkg/util/workspace"
)

const (
	FieldType        = "type"
	FieldLabels      = "labels"
	FieldAnnotations = "annotations"
)

const (
	CSPAWS      = "aws"
	CSPAliCloud = "alicloud"
)

const (
	ProtocolTCP = "TCP"
	ProtocolUDP = "UDP"
)

const (
	k8sKindService = "Service"
	suffixPublic   = "public"
	suffixPrivate  = "private"
)

var (
	ErrEmptyPortConfig   = errors.New("empty port config")
	ErrEmptyType         = errors.New("type must not be empty when public")
	ErrUnsupportedType   = errors.New("type only support alicloud and aws for now")
	ErrInvalidPort       = errors.New("port must be between 1 and 65535")
	ErrInvalidTargetPort = errors.New("targetPort must be between 1 and 65535 if exist")
	ErrInvalidProtocol   = errors.New("protocol must be TCP or UDP")
	ErrEmptySvcWorkload  = errors.New("network port should be binded to a service workload")
)

// Network describes the network accessories of workload, which typically contains the exposed
// ports, load balancer and other related resource configs.
type Network struct {
	Ports []Port `yaml:"ports,omitempty" json:"ports,omitempty"`
}

// Port defines the exposed port of workload, which can be used to describe how
// the workload get accessed.
type Port struct {
	// Type is the specific cloud vendor that provides load balancer, works when Public
	// is true, supports CSPAliCloud and CSPAWS for now.
	Type string `yaml:"type,omitempty" json:"type,omitempty"`

	// Port is the exposed port of the workload.
	Port int `yaml:"port,omitempty" json:"port,omitempty"`

	// TargetPort is the backend container.Container port.
	TargetPort int `yaml:"targetPort,omitempty" json:"targetPort,omitempty"`

	// Protocol is protocol used to expose the port, support ProtocolTCP and ProtocolUDP.
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`

	// Public defines whether to expose the port through Internet.
	Public bool `yaml:"public,omitempty" json:"public,omitempty"`

	// Labels are the attached labels of the port, works only when the Public is true.
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`

	// Annotations are the attached annotations of the port, works only when the Public is true.
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

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

func main() {
	server.Start(&Network{})
}
