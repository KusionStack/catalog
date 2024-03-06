package generators

import (
	"fmt"

	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	v1 "kusionstack.io/kusion/pkg/apis/core/v1"
	"kusionstack.io/kusion/pkg/apis/core/v1/workload"
	"kusionstack.io/kusion/pkg/log"
	"kusionstack.io/kusion/pkg/modules"
	"kusionstack.io/kusion/pkg/modules/proto"
)

// HandshakeConfig is a common handshake that is shared by plugin and host.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MODULE_PLUGIN",
	MagicCookieValue: "ON",
}

func StartModule(module modules.Module) {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			modules.PluginKey: &modules.GRPCPlugin{Impl: module},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

type GeneratorRequest struct {
	// Project represents the project name
	Project string `json:"project,omitempty"`
	// Stack represents the stack name
	Stack string `json:"stack,omitempty"`
	// App represents the application name, which is typically the same as the namespace of Kubernetes resources
	App string `json:"app,omitempty"`
	// Workload represents the workload configuration
	Workload *workload.Workload `json:"workload,omitempty"`
	// DevModuleConfig is the developer's inputs of this module
	DevModuleConfig v1.Accessory `json:"dev_module_config,omitempty"`
	// PlatformModuleConfig is the platform engineer's inputs of this module
	PlatformModuleConfig v1.GenericConfig `json:"platform_module_config,omitempty"`
	// RuntimeConfig is the runtime configurations defined in the workspace config
	RuntimeConfig *v1.RuntimeConfigs `json:"runtime_config,omitempty"`
}

func NewGeneratorRequest(req *proto.GeneratorRequest) (*GeneratorRequest, error) {

	log.Infof("module proto request received:%s", req.String())

	// validate workload
	if req.Workload == nil {
		return nil, fmt.Errorf("workload in the request is nil")
	}
	w := &workload.Workload{}
	if err := yaml.Unmarshal(req.Workload, w); err != nil {
		return nil, fmt.Errorf("unmarshal workload failed. %w", err)
	}

	var dc v1.Accessory
	if req.DevModuleConfig != nil {
		if err := yaml.Unmarshal(req.DevModuleConfig, &dc); err != nil {
			return nil, fmt.Errorf("unmarshal dev module config failed. %w", err)
		}
	}

	var pc v1.GenericConfig
	if req.PlatformModuleConfig != nil {
		if err := yaml.Unmarshal(req.PlatformModuleConfig, &pc); err != nil {
			return nil, fmt.Errorf("unmarshal platform module config failed. %w", err)
		}
	}

	var rc *v1.RuntimeConfigs
	if req.RuntimeConfig != nil {
		if err := yaml.Unmarshal(req.RuntimeConfig, rc); err != nil {
			return nil, fmt.Errorf("unmarshal runtime config failed. %w", err)
		}
	}

	result := &GeneratorRequest{
		Project:              req.Project,
		Stack:                req.Stack,
		App:                  req.App,
		Workload:             w,
		DevModuleConfig:      dc,
		PlatformModuleConfig: pc,
		RuntimeConfig:        rc,
	}
	out, err := yaml.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("marshal new generator request failed. %w", err)
	}
	log.Infof("new generator request:%s", string(out))
	return result, nil
}

func EmptyResponse() *proto.GeneratorResponse {
	return &proto.GeneratorResponse{}
}

func WrapK8sResourceToKusionResource(id string, resource any) (*v1.Resource, error) {
	gvk := resource.(runtime.Object).GetObjectKind().GroupVersionKind().String()

	// fixme: this function converts int to int64 by default
	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(resource)
	if err != nil {
		return nil, err
	}
	return &v1.Resource{
		ID:         id,
		Type:       v1.Kubernetes,
		Attributes: unstructured,
		DependsOn:  nil,
		Extensions: map[string]any{
			v1.ResourceExtensionGVK: gvk,
		},
	}, nil
}
