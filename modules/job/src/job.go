package main

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"gopkg.in/yaml.v2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
)

func (j *Job) Generate(ctx context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	// Get the module logger with the generator context.
	logger := log.GetModuleLogger(ctx)
	logger.Info("Generating resources...")

	defer func() {
		if r := recover(); r != nil {
			logger.Debug("failed to generate Job module: %v", r)
			response = nil
			rawRequest, _ := json.Marshal(request)
			err = fmt.Errorf("panic in job module generator but recovered with error: [%v] and stack %v and request %v",
				r, string(debug.Stack()), string(rawRequest))
		}
	}()

	if request.DevConfig == nil {
		logger.Info("Job does not exist in AppConfig config")
		return nil, nil
	}
	out, err := yaml.Marshal(request.DevConfig)
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(out, j); err != nil {
		return nil, fmt.Errorf("complete Job by dev config failed, %w", err)
	}

	if err = completeBaseWorkload(&j.Base, request.PlatformConfig); err != nil {
		return nil, fmt.Errorf("complete Job by platform config failed, %w", err)
	}

	uniqueAppName := module.UniqueAppName(request.Project, request.Stack, request.App)

	meta := metav1.ObjectMeta{
		Namespace: request.Project,
		Name:      uniqueAppName,
		Labels: module.MergeMaps(
			module.UniqueAppLabels(request.Project, request.App),
			j.Labels,
		),
		Annotations: module.MergeMaps(
			j.Annotations,
		),
	}

	containers, volumes, configMaps, err := toOrderedContainers(j.Containers, uniqueAppName)
	if err != nil {
		return nil, err
	}

	res := make([]kusionapiv1.Resource, 0)
	for _, cm := range configMaps {
		cm.Namespace = request.Project
		resourceID := module.KubernetesResourceID(cm.TypeMeta, cm.ObjectMeta)
		resource, err := module.WrapK8sResourceToKusionResource(resourceID, &cm)
		if err != nil {
			return nil, err
		}
		res = append(res, *resource)
	}

	jobSpec := batchv1.JobSpec{
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      module.MergeMaps(module.UniqueAppLabels(request.Project, request.App), j.Labels),
				Annotations: module.MergeMaps(j.Annotations),
			},
			Spec: corev1.PodSpec{
				Containers:    containers,
				RestartPolicy: corev1.RestartPolicyNever,
				Volumes:       volumes,
			},
		},
	}

	if j.Schedule == "" {
		k8sJob := &batchv1.Job{
			ObjectMeta: meta,
			TypeMeta: metav1.TypeMeta{
				Kind:       "Job",
				APIVersion: batchv1.SchemeGroupVersion.String(),
			},
			Spec: jobSpec,
		}

		resourceID := module.KubernetesResourceID(k8sJob.TypeMeta, k8sJob.ObjectMeta)
		resource, err := module.WrapK8sResourceToKusionResource(resourceID, k8sJob)
		if err != nil {
			return nil, err
		}
		res = append(res, *resource)

		return &module.GeneratorResponse{
			Resources: res,
		}, nil
	}

	cronJob := &batchv1.CronJob{
		ObjectMeta: meta,
		TypeMeta: metav1.TypeMeta{
			Kind:       "CronJob",
			APIVersion: batchv1.SchemeGroupVersion.String(),
		},
		Spec: batchv1.CronJobSpec{
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: jobSpec,
			},
			Schedule: j.Schedule,
		},
	}

	resourceID := module.KubernetesResourceID(cronJob.TypeMeta, cronJob.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, cronJob)
	if err != nil {
		return nil, err
	}
	res = append(res, *resource)
	return &module.GeneratorResponse{
		Resources: res,
	}, nil
}

func main() {
	server.Start(&Job{})
}
