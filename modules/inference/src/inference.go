package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
)

// error type
var (
	ErrUnsupportFramework = errors.New("framework must be Ollama or KubeRay")
	ErrRangeTopK          = errors.New("topK must be greater than 0 if exist")
	ErrRangeTopP          = errors.New("topP must be greater than 0 and less than or equal to 1 if exist")
	ErrRangeTemperature   = errors.New("temperature must be greater than 0 if exist")
	ErrRangeNumPredict    = errors.New("numPredict must be greater than or equal to -2")
	ErrRangeNumCtx        = errors.New("numCtx must be greater than 0 if exist")
)

// resource naming
var (
	inferDeploymentSuffix    = "-infer-deployment"
	inferStorageSuffix       = "-infer-storage"
	inferServiceSuffix       = "-infer-service"
	inferContainerPortSuffix = "-port"
	inferContainerSuffix     = "-infer-container"
)

// default config
var (
	defaultTopK        int     = 40
	defaultTopP        float64 = 0.9
	defaultTemperature float64 = 0.8
	defaultNumPredict  int     = 128
	defaultNumCtx      int     = 2048
)

// port
var (
	CalledPort = 80
	OllamaPort = 11434
)

// framework type
var (
	OllamaType = "ollama"
)

// framework image
var (
	OllamaImage = "ollama/ollama"
)

// proxy
var (
	ProxyName  = "proxy"
	ProxyPort  = 5000
	ProxyImage = "kangy126/proxy"
)

func main() {
	server.Start(&Inference{})
}

// Inference implements the Kusion Module generator interface.
type Inference struct {
	Model       string  `yaml:"model,omitempty" json:"model,omitempty"`
	Framework   string  `yaml:"framework,omitempty" json:"framework,omitempty"`
	System      string  `yaml:"system,omitempty" json:"system,omitempty"`
	Template    string  `yaml:"template,omitempty" json:"template,omitempty"`
	TopK        int     `yaml:"top_k,omitempty" json:"top_k,omitempty"`
	TopP        float64 `yaml:"top_p,omitempty" json:"top_p,omitempty"`
	Temperature float64 `yaml:"temperature,omitempty" json:"temperature,omitempty"`
	NumPredict  int     `yaml:"num_predict,omitempty" json:"num_predict,omitempty"`
	NumCtx      int     `yaml:"num_ctx,omitempty" json:"num_ctx,omitempty"`
}

func (infer *Inference) Generate(ctx context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	// Get the module logger with the generator context.
	logger := log.GetModuleLogger(ctx)
	logger.Info("Generating resources...")

	// Handle the panic and update the returned error.
	defer func() {
		if r := recover(); r != nil {
			logger.Debug("failed to generate inference module: %v", r)
			response = nil
			rawRequest, _ := json.Marshal(request)
			err = fmt.Errorf("panic in inference module generator but recovered with error: [%v] and stack %v and request %v", r, string(debug.Stack()), string(rawRequest))
		}
	}()

	// Inference module does not exist in AppConfiguration configs.
	if request.DevConfig == nil {
		logger.Info("Inference does not exist in AppConfig config")
		return nil, nil
	}

	// Get the complete inference module configs.
	if err := infer.CompleteConfig(request.DevConfig, request.PlatformConfig); err != nil {
		logger.Debug("failed to get complete inference module configs: %v", err)
		return nil, err
	}

	// Validate the completed inference module configs.
	if err := infer.ValidateConfig(); err != nil {
		logger.Debug("failed to validate the inference module configs: %v", err)
		return nil, err
	}

	var resources []kusionapiv1.Resource
	var patcher *kusionapiv1.Patcher

	switch strings.ToLower(infer.Framework) {
	case OllamaType:
		resources, patcher, err = infer.GenerateOllamaResource(request)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrUnsupportFramework
	}

	// Return the Kusion generator response.
	return &module.GeneratorResponse{
		Resources: resources,
		Patcher:   patcher,
	}, nil
}

// CompleteConfig completes the inference module configs with both devModuleConfig and platformModuleConfig.
func (infer *Inference) CompleteConfig(devConfig kusionapiv1.Accessory, platformConfig kusionapiv1.GenericConfig) error {
	infer.TopK = defaultTopK
	infer.TopP = defaultTopP
	infer.Temperature = defaultTemperature
	infer.NumPredict = defaultNumPredict
	infer.NumCtx = defaultNumCtx

	// Retrieve the config items the developers are concerned about.
	if devConfig != nil {
		devCfgYamlStr, err := yaml.Marshal(devConfig)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(devCfgYamlStr, infer); err != nil {
			return err
		}
	}
	// Retrieve the config items the platform engineers care about.
	if platformConfig != nil {
		platformCfgYamlStr, err := yaml.Marshal(platformConfig)
		if err != nil {
			return err
		}

		if err = yaml.Unmarshal(platformCfgYamlStr, infer); err != nil {
			return err
		}
	}
	return nil
}

// ValidateConfig validates the completed inference configs are valid or not.
func (infer *Inference) ValidateConfig() error {
	if infer.TopK <= 0 {
		return ErrRangeTopK
	}
	if infer.TopP <= 0 || infer.TopP > 1 {
		return ErrRangeTopP
	}
	if infer.Temperature <= 0 {
		return ErrRangeTemperature
	}
	if infer.NumPredict < -2 {
		return ErrRangeNumPredict
	}
	if infer.NumCtx <= 0 {
		return ErrRangeNumCtx
	}
	return nil
}

func (infer *Inference) GenerateEnv(svcName string) (*kusionapiv1.Patcher, error) {
	envVars := []v1.EnvVar{
		{
			Name:  "INFERENCE_URL",
			Value: svcName,
		},
	}
	patcher := &kusionapiv1.Patcher{
		Environments: envVars,
	}

	return patcher, nil
}
