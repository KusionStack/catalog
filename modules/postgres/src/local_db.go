package main

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"

	"kusionstack.io/kusion-module-framework/pkg/module"

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiv1 "kusionstack.io/kusion/pkg/apis/core/v1"
)

// GenerateLocalResources generates the resources of locally deployed PostgreSQL database instance.
func (postgres *PostgreSQL) GenerateLocalResources(request *module.GeneratorRequest) ([]apiv1.Resource, []apiv1.Patcher, error) {
	var resources []apiv1.Resource
	var patchers []apiv1.Patcher

	// Build Kubernetes Secret for the random password of the local PostgreSQL instance.
	password := postgres.generateLocalPassword(request)
	localSecret, err := postgres.generateLocalSecret(request, password)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *localSecret)

	// Build Kubernetes Deployment for the local PostgreSQL instance.
	localDeployment, err := postgres.generateLocalDeployment(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *localDeployment)

	// Build Kubernetes Persistent Volume Claim for the lcoal PostgreSQL instance.
	localPVC, err := postgres.generateLocalPVC(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *localPVC)

	// Build Kubernetes Service for the local PostgreSQL instance.
	localSvc, hostAddress, err := postgres.generateLocalService(request)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *localSvc)

	// Build Kubernetes Secret with the hostAddress, username and password of the local PostgreSQL instance,
	// and inject the credentials as the environment variable patcher.
	dbSecret, patcher, err := postgres.GenerateDBSecret(request, hostAddress, postgres.Username, password)
	if err != nil {
		return nil, nil, err
	}
	resources = append(resources, *dbSecret)
	patchers = append(patchers, *patcher)

	return resources, patchers, nil
}

// generateLocalSecret generates the Kubernetes Secret resource for the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalSecret(request *module.GeneratorRequest, password string) (*apiv1.Resource, error) {
	// Set the password string.
	data := make(map[string]string)
	data["password"] = password
	data["username"] = postgres.Username
	data["database"] = postgres.DatabaseName

	// Construct the Kubernetes Secret resource.
	secret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgres.DatabaseName + localSecretSuffix,
			Namespace: request.Project,
		},
		StringData: data,
	}

	resourceID := module.KubernetesResourceID(secret.TypeMeta, secret.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, secret)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// generateLocalDeployment generates the Kubernetes Deployment resource for the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalDeployment(request *module.GeneratorRequest) (*apiv1.Resource, error) {
	// Prepare the Pod Spec for the local PostgreSQL instance.
	podSpec, err := postgres.generateLocalPodSpec(request)
	if err != nil {
		return nil, nil
	}

	// Create the Kubernetes Deployment for the local PostgreSQL instance.
	deployment := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgres.DatabaseName + localDeploymentSuffix,
			Namespace: request.Project,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: postgres.generateLocalMatchLabels(),
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: postgres.generateLocalMatchLabels(),
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

// generateLocalPodSpec generates the Kubernetes PodSpec for the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalPodSpec(_ *module.GeneratorRequest) (v1.PodSpec, error) {
	image := dbEngine + ":" + postgres.Version
	secretName := postgres.DatabaseName + localSecretSuffix

	var portName string
	if len(postgres.DatabaseName) > 15 {
		portName = postgres.DatabaseName[:15]
	} else {
		portName = postgres.DatabaseName
	}
	ports := []v1.ContainerPort{
		{
			Name:          portName,
			ContainerPort: int32(5432),
		},
	}

	volumes := []v1.Volume{
		{
			Name: postgres.DatabaseName,
			VolumeSource: v1.VolumeSource{
				PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
					ClaimName: postgres.DatabaseName + localPVCSuffix,
				},
			},
		},
	}
	volumeMounts := []v1.VolumeMount{
		{
			Name:      postgres.DatabaseName,
			MountPath: "/var/lib/postgresql/data",
		},
	}

	env := []v1.EnvVar{
		{
			Name: "POSTGRES_USER",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
					Key: "username",
				},
			},
		},
		{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
					Key: "password",
				},
			},
		},
		{
			Name: "POSTGRES_DB",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{
						Name: secretName,
					},
					Key: "database",
				},
			},
		},
	}

	podSpec := v1.PodSpec{
		Containers: []v1.Container{
			{
				Name:         postgres.DatabaseName,
				Image:        image,
				Env:          env,
				Ports:        ports,
				VolumeMounts: volumeMounts,
			},
		},
		Volumes: volumes,
	}

	return podSpec, nil
}

// generateLocalPVC generates the Kubernetes Persistent Volume Claim resource for the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalPVC(request *module.GeneratorRequest) (*apiv1.Resource, error) {
	// Create the Kubernetes PVC with the storage size of `postgres.Size`.
	pvc := &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      postgres.DatabaseName + localPVCSuffix,
			Namespace: request.Project,
			Labels:    postgres.generateLocalMatchLabels(),
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{
				v1.ReadWriteOnce,
			},
			Resources: v1.VolumeResourceRequirements{
				Requests: map[v1.ResourceName]resource.Quantity{
					v1.ResourceStorage: resource.MustParse(strconv.Itoa(postgres.Size) + "Gi"),
				},
			},
		},
	}

	resourceID := module.KubernetesResourceID(pvc.TypeMeta, pvc.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, pvc)
	if err != nil {
		return nil, err
	}

	return resource, nil
}

// generateLocalService generates the Kubernetes Service resource for the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalService(request *module.GeneratorRequest) (*apiv1.Resource, string, error) {
	// Prepare the service port for the local PostgreSQL instance.
	svcPort := postgres.generateLocalSvcPort()
	svcName := postgres.DatabaseName + localServiceSuffix

	// Create the Kubernetes service for local PostgreSQL instance.
	service := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: v1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: request.Project,
			Labels:    postgres.generateLocalMatchLabels(),
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "None",
			Ports:     svcPort,
			Selector:  postgres.generateLocalMatchLabels(),
		},
	}

	resourceID := module.KubernetesResourceID(service.TypeMeta, service.ObjectMeta)
	resource, err := module.WrapK8sResourceToKusionResource(resourceID, service)
	if err != nil {
		return nil, "", err
	}

	return resource, svcName, nil
}

// generateLocalSvcPort generates the Kubernetes ServicePort resource of the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalSvcPort() []v1.ServicePort {
	svcPort := []v1.ServicePort{
		{
			Port: int32(5432),
		},
	}

	return svcPort
}

// generateLocalMatchLabels generates the match labels for the Kubernetes resources of the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalMatchLabels() map[string]string {
	return map[string]string{
		"accessory": postgres.DatabaseName,
	}
}

// generateLocalPassword generates a fixed password string with the specified length for the local PostgreSQL instance.
func (postgres *PostgreSQL) generateLocalPassword(request *module.GeneratorRequest) string {
	hashInput := request.Project + request.Stack + request.App + postgres.DatabaseName
	hash := md5.Sum([]byte(hashInput))

	hashString := hex.EncodeToString(hash[:])

	return hashString[:16]
}
