package main

import (
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kusionstack.io/kusion-module-framework/pkg/module"
	apiv1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
)

// GenerateOllamaResource generates the the resources of Ollama
func (infer *Inference) GenerateOllamaResource(request *module.GeneratorRequest) ([]apiv1.Resource, *apiv1.Patcher, error) {
	var resources []apiv1.Resource

	// Build Kubernetes Deployment for Ollama framework.
	deployment, err := infer.generateOllamaDeployment(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *deployment)

	// Build Kubernetes Service for Ollama framework.
	svc, svcName, err := infer.generateOllamaService(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *svc)

	// Build Kubernetes Deployment for proxy.
	deploymentProxy, err := infer.generateProxyDeployment(request, svcName)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *deploymentProxy)

	// Build Kubernetes Service for proxy.
	svcProxy, svcNameProxy, err := infer.generateProxyService(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *svcProxy)

	patcher, err := infer.GenerateEnv(svcNameProxy)
	if err != nil {
		return nil, nil, err
	}

	return resources, patcher, nil
}

// generatePodSpec generates the Kubernetes PodSpec for Ollama framework.
func (infer *Inference) generateOllamaPodSpec(_ *module.GeneratorRequest) (v1.PodSpec, error) {
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
	commandParts = append(commandParts, "ollama serve & OLLAMA_SERVE_PID=$!")
	commandParts = append(commandParts, "sleep 5")
	commandParts = append(commandParts, fmt.Sprintf("ollama create %s -f Modelfile", infer.Model))
	commandParts = append(commandParts, "wait $OLLAMA_SERVE_PID")

	var modelPullCmd []string
	modelPullCmd = append(modelPullCmd, "/bin/sh", "-c", strings.Join(commandParts, " && "))

	volumes := []v1.Volume{
		{
			Name: strings.ToLower(infer.Framework) + inferStorageSuffix,
			VolumeSource: v1.VolumeSource{
				EmptyDir: &v1.EmptyDirVolumeSource{},
			},
		},
	}

	mountPath := "/root/.ollama"
	volumeMounts := []v1.VolumeMount{
		{
			Name:      strings.ToLower(infer.Framework) + inferStorageSuffix,
			MountPath: mountPath,
		},
	}

	portName := strings.ToLower(infer.Framework) + inferContainerPortSuffix
	if len(portName) > 15 {
		portName = portName[:15]
	}
	containerPort := int32(OllamaPort)
	ports := []v1.ContainerPort{
		{
			Name:          portName,
			ContainerPort: containerPort,
		},
	}

	image := OllamaImage
	podSpec := v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:         strings.ToLower(infer.Framework) + inferContainerSuffix,
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

// generateDeployment generates the Kubernetes Deployment resource for Ollama framework.
func (infer *Inference) generateOllamaDeployment(request *module.GeneratorRequest) (*apiv1.Resource, error) {
	// Prepare the Pod Spec for Ollama framework.
	podSpec, err := infer.generateOllamaPodSpec(request)
	if err != nil {
		return nil, nil
	}

	// Create the Kubernetes Deployment for Ollama framework.
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(infer.Framework) + inferDeploymentSuffix,
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

// generateService generates the Kubernetes Service resource for Ollama framework.
func (infer *Inference) generateOllamaService(request *module.GeneratorRequest) (*apiv1.Resource, string, error) {
	// Prepare the service port for Ollama framework.
	svcName := strings.ToLower(infer.Framework) + inferServiceSuffix
	svcPort := []v1.ServicePort{
		{
			Port: int32(CalledPort),
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(OllamaPort),
			},
		},
	}

	// Create the Kubernetes service for Ollama framework.
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

// generateMatchLabels generates the match labels for the Kubernetes resources of Ollama framework.
func (infer *Inference) generateMatchLabels() map[string]string {
	return map[string]string{
		"accessory": strings.ToLower(infer.Framework),
	}
}

// generateMatchLabels generates the match labels for the Kubernetes resources of proxy.
func (infer *Inference) generateMatchLabelsForProxy() map[string]string {
	return map[string]string{
		"accessory": strings.ToLower(ProxyName),
	}
}

// generatePodSpec generates the Kubernetes PodSpec for proxy.
func (infer *Inference) generateProxyPodSpec(_ *module.GeneratorRequest, svcName string) (v1.PodSpec, error) {
	portName := strings.ToLower(ProxyName) + inferContainerPortSuffix
	if len(portName) > 15 {
		portName = portName[:15]
	}
	containerPort := int32(ProxyPort)
	ports := []v1.ContainerPort{
		{
			Name:          portName,
			ContainerPort: containerPort,
		},
	}

	envVars := []v1.EnvVar{
		{
			Name:  "MODEL",
			Value: infer.Model,
		},
		{
			Name:  "FRAMEWORK_URL",
			Value: svcName,
		},
	}

	image := ProxyImage
	podSpec := v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:  strings.ToLower(ProxyName) + inferContainerSuffix,
				Image: image,
				Ports: ports,
				Env:   envVars,
			},
		},
	}
	return podSpec, nil
}

// generateDeployment generates the Kubernetes Deployment resource for proxy.
func (infer *Inference) generateProxyDeployment(request *module.GeneratorRequest, svcName string) (*apiv1.Resource, error) {
	podSpec, err := infer.generateProxyPodSpec(request, svcName)
	if err != nil {
		return nil, nil
	}

	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      strings.ToLower(ProxyName) + inferDeploymentSuffix,
			Namespace: request.Project,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: infer.generateMatchLabelsForProxy(),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: infer.generateMatchLabelsForProxy(),
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

// generateService generates the Kubernetes Service resource for proxy.
func (infer *Inference) generateProxyService(request *module.GeneratorRequest) (*apiv1.Resource, string, error) {
	svcName := strings.ToLower(ProxyName) + inferServiceSuffix
	svcPort := []v1.ServicePort{
		{
			Port: int32(CalledPort),
			TargetPort: intstr.IntOrString{
				Type:   intstr.Int,
				IntVal: int32(ProxyPort),
			},
		},
	}

	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: request.Project,
			Labels:    infer.generateMatchLabelsForProxy(),
		},
		Spec: v1.ServiceSpec{
			Type:     v1.ServiceTypeClusterIP,
			Ports:    svcPort,
			Selector: infer.generateMatchLabelsForProxy(),
		},
	}

	resourceID := module.KubernetesResourceID(service.TypeMeta, service.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, service)
	if err != nil {
		return nil, svcName, err
	}

	return resource, svcName, nil
}
