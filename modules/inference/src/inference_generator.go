package main

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/log"
)

var (
	ErrUnsupportFramework = errors.New("framework must be Ollama or KubeRay")
	ErrRangeTopK          = errors.New("topK must be greater than 0 if exist")
	ErrRangeTopP          = errors.New("topP must be greater than 0 and less than or equal to 1 if exist")
	ErrRangeTemperature   = errors.New("temperature must be greater than 0 if exist")
	ErrRangeNumPredict    = errors.New("numPredict must be greater than or equal to -2")
	ErrRangeNumCtx        = errors.New("numCtx must be greater than 0 if exist")
)

var (
	inferDeploymentSuffix = "-infer-deployment"
	inferStorageSuffix    = "-infer-storage"
	inferServiceSuffix    = "-infer-service"
)

var (
	defaultTopK        int     = 40
	defaultTopP        float64 = 0.9
	defaultTemperature float64 = 0.8
	defaultNumPredict  int     = 128
	defaultNumCtx      int     = 2048
)

var (
	OllamaType = "ollama"
)

var (
	OllamaImage = "ollama"
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

func (infer *Inference) Generate(_ context.Context, request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Debugf("failed to generate inference module: %v", r)
		}
	}()

	// Inference module does not exist in AppConfiguration configs.
	if request.DevConfig == nil {
		log.Info("Inference does not exist in AppConfig config")
		return nil, nil
	}

	// Get the complete inference module configs.
	if err := infer.CompleteConfig(request.DevConfig, request.PlatformConfig); err != nil {
		log.Debugf("failed to get complete inference module configs: %v", err)
		return nil, err
	}

	// Validate the completed inference module configs.
	if err := infer.ValidateConfig(); err != nil {
		log.Debugf("failed to validate the inference module configs: %v", err)
		return nil, err
	}

	// var resources []apiv1.Resource
	// var patcher *apiv1.Patcher
	// var err error

	// switch strings.ToLower(infer.Framework) {
	// case OllamaType:
	// 	resources, patcher, err = infer.GenerateInferenceResource(request)
	// default:
	// 	return nil, ErrUnsupportFramework
	// }

	// Generate the Kubernetes Service related resource.
	resources, patcher, err := infer.GenerateInferenceResource(request)
	if err != nil {
		return nil, err
	}

	// Return the Kusion generator response.
	return &module.GeneratorResponse{
		Resources: resources,
		Patcher:   patcher,
	}, nil
}

// CompleteConfig completes the inference module configs with both devModuleConfig and platformModuleConfig.
func (infer *Inference) CompleteConfig(devConfig apiv1.Accessory, platformConfig apiv1.GenericConfig) error {
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
	if infer.Framework != "Ollama" && infer.Framework != "KubeRay" {
		return ErrUnsupportFramework
	}
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

// GenerateInferenceResource generates the Kubernetes Service related to the inference module service.
//
// Note that we will use the SDK provided by the kusion module framework to wrap the Kubernetes resource
// into Kusion resource.
func (infer *Inference) GenerateInferenceResource(request *module.GeneratorRequest) ([]apiv1.Resource, *apiv1.Patcher, error) {
	var resources []apiv1.Resource

	// Build Kubernetes Deployment for the Inference instance.
	deployment, err := infer.generateDeployment(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *deployment)

	// Build Kubernetes Service for the Inference instance.
	svc, svcName, err := infer.generateService(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *svc)

	envVars := []v1.EnvVar{
		{
			Name:  "INFERENCE_PATH",
			Value: svcName,
		},
	}
	patcher := &apiv1.Patcher{
		Environments: envVars,
	}

	return resources, patcher, nil
}

// generatePodSpec generates the Kubernetes PodSpec for the Inference instance.
func (infer *Inference) generatePodSpec(_ *module.GeneratorRequest) (v1.PodSpec, error) {
	var mountPath string
	var modelPullCmd []string
	var containerPort int32
	switch strings.ToLower(infer.Framework) {
	case OllamaType:
		mountPath = "/root/.ollama"

		var builder strings.Builder
		builder.WriteString("'")
		builder.WriteString(fmt.Sprintf("FROM %s\n", infer.Model))
		if infer.System != "" {
			builder.WriteString(fmt.Sprintf(`SYSTEM """%s"""`, infer.System))
			builder.WriteString("\n")
		}
		if infer.Template != "" {
			builder.WriteString(fmt.Sprintf(`TEMPLATE """%s""""`, infer.Template))
			builder.WriteString("\n")
		}
		builder.WriteString(fmt.Sprintf("PARAMETER top_k %d\n", infer.TopK))
		builder.WriteString(fmt.Sprintf("PARAMETER top_p %f\n", infer.TopP))
		builder.WriteString(fmt.Sprintf("PARAMETER temperature %f\n", infer.Temperature))
		builder.WriteString(fmt.Sprintf("PARAMETER num_predict %d\n", infer.NumPredict))
		builder.WriteString(fmt.Sprintf("PARAMETER num_ctx %d\n", infer.NumCtx))
		builder.WriteString("'")

		var commandParts []string
		commandParts = append(commandParts, fmt.Sprintf("echo %s > Modelfile", builder.String()))
		commandParts = append(commandParts, fmt.Sprintf("ollama create %s -f Modelfile", infer.Model))

		modelPullCmd = append(modelPullCmd, "/bin/sh", "-c", strings.Join(commandParts, " && "))
		containerPort = 11434
	default:
	}

	image := OllamaImage

	volumes := []v1.Volume{
		{
			Name: infer.Framework + inferStorageSuffix,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}

	volumeMounts := []v1.VolumeMount{
		{
			Name:      infer.Framework + inferStorageSuffix,
			MountPath: mountPath,
		},
	}

	ports := []v1.ContainerPort{
		{
			Name:          infer.Framework,
			ContainerPort: containerPort,
		},
	}

	podSpec := v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:         infer.Framework,
				Image:        image,
				Ports:        ports,
				Command:      modelPullCmd,
				VolumeMounts: volumeMounts,
			},
		},
		Volumes: volumes,
	}
	return podSpec, nil
}

// generateDeployment generates the Kubernetes Deployment resource for the Inference instance.
func (infer *Inference) generateDeployment(request *module.GeneratorRequest) (*apiv1.Resource, error) {
	// Prepare the Pod Spec for the Inference instance.
	podSpec, err := infer.generatePodSpec(request)
	if err != nil {
		return nil, nil
	}

	// Create the Kubernetes Deployment for the Inference instance.
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      infer.Framework + inferDeploymentSuffix,
			Namespace: request.Project,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: infer.generateMatchLabels(),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: infer.generateMatchLabels(),
				},
				Spec: podSpec,
			},
		},
	}

	resourceID := module.KubernetesResourceID(deployment.TypeMeta, deployment.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, deployment)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// generateService generates the Kubernetes Service resource for the Inference instance.
func (infer *Inference) generateService(request *module.GeneratorRequest) (*apiv1.Resource, string, error) {
	// Prepare the service port for the Inference instance.
	svcName := infer.Framework + inferServiceSuffix
	svcPort := []v1.ServicePort{
		{
			Port: int32(80),
		},
	}

	// Create the Kubernetes service for Inference instance.
	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: request.Project,
			Labels:    infer.generateMatchLabels(),
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeClusterIP,
			Ports:    svcPort,
			Selector: infer.generateMatchLabels(),
		},
	}

	resourceID := module.KubernetesResourceID(service.TypeMeta, service.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, service)
	if err != nil {
		return nil, svcName, err
	}

	return resource, svcName, nil
}

// generateMatchLabels generates the match labels for the Kubernetes resources of the Inference instance.
func (infer *Inference) generateMatchLabels() map[string]string {
	return map[string]string{
		"accessory": infer.Framework,
	}
}
