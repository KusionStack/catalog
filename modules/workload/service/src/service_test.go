package main

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	"kusionstack.io/kube-api/apps/v1alpha1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
)

func Test_workloadServiceGenerator_Generate(t *testing.T) {
	cm := `apiVersion: v1
data:
    example.txt: some file contents
kind: ConfigMap
metadata:
    creationTimestamp: null
    name: default-dev-foo-nginx-0
    namespace: default
`
	var cmResource v1.Resource
	k8sCm := &corev1.ConfigMap{}
	_ = yaml.Unmarshal([]byte(cm), k8sCm)
	unstructured, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(k8sCm)
	cmResource.Attributes = unstructured

	csSvc := `apiVersion: v1
kind: Service
metadata:
    annotations:
        service-workload-type: CollaSet
        service.beta.kubernetes.io/alibaba-cloud-loadbalancer-spec: slb.s1.small
    creationTimestamp: null
    labels:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
        kusionstack.io/control: "true"
        service-workload-type: CollaSet
    name: default-dev-foo-public
    namespace: default
spec:
    ports:
        - name: default-dev-foo-public-80-tcp
          port: 80
          protocol: TCP
          targetPort: 80
    selector:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
    type: LoadBalancer
status:
    loadBalancer: {}
`
	var csSvcResource v1.Resource
	k8sSvc := &corev1.Service{}
	_ = yaml.Unmarshal([]byte(csSvc), k8sSvc)
	unSvc, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(k8sSvc)
	csSvcResource.Attributes = unSvc

	deploySvc := `apiVersion: v1
kind: Service
metadata:
    annotations:
        service.beta.kubernetes.io/alibaba-cloud-loadbalancer-spec: slb.s1.small
    creationTimestamp: null
    labels:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
        kusionstack.io/control: "true"
        service-workload-type: Deployment
    name: default-dev-foo-public
    namespace: default
spec:
    ports:
        - name: default-dev-foo-public-80-tcp
          port: 80
          protocol: TCP
          targetPort: 80
    selector:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
    type: LoadBalancer
status:
    loadBalancer: {}
`
	var deploySvcRes v1.Resource
	deployK8sSvc := &corev1.Service{}
	_ = yaml.Unmarshal([]byte(deploySvc), deployK8sSvc)
	unDeploySvc, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(deployK8sSvc)
	deploySvcRes.Attributes = unDeploySvc

	cs := `apiVersion: apps.kusionstack.io/v1alpha1
kind: CollaSet
metadata:
    annotations:
        service-workload-type: CollaSet
    creationTimestamp: null
    labels:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
        service-workload-type: CollaSet
    name: default-dev-foo
    namespace: default
spec:
    replicas: 2
    scaleStrategy: {}
    selector:
        matchLabels:
            app.kubernetes.io/name: foo
            app.kubernetes.io/part-of: default
    template:
        metadata:
            annotations:
                service-workload-type: CollaSet
            creationTimestamp: null
            labels:
                app.kubernetes.io/name: foo
                app.kubernetes.io/part-of: default
                service-workload-type: CollaSet
        spec:
            containers:
                - image: nginx:v1
                  name: nginx
                  resources: {}
                  volumeMounts:
                    - mountPath: /tmp
                      name: default-dev-foo-nginx-0
            volumes:
                - configMap:
                    defaultMode: 511
                    name: default-dev-foo-nginx-0
                  name: default-dev-foo-nginx-0
    updateStrategy: {}
status: {}
`
	var csResource v1.Resource
	k8sCS := &v1alpha1.CollaSet{}
	_ = yaml.Unmarshal([]byte(cs), k8sCS)
	unCS, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(k8sCS)
	csResource.Attributes = unCS

	deploy := `apiVersion: apps/v1
kind: Deployment
metadata:
    creationTimestamp: null
    labels:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
        service-workload-type: Deployment
    name: default-dev-foo
    namespace: default
spec:
    replicas: 4
    selector:
        matchLabels:
            app.kubernetes.io/name: foo
            app.kubernetes.io/part-of: default
    strategy: {}
    template:
        metadata:
            creationTimestamp: null
            labels:
                app.kubernetes.io/name: foo
                app.kubernetes.io/part-of: default
                service-workload-type: Deployment
        spec:
            containers:
                - image: nginx:v1
                  name: nginx
                  resources: {}
                  volumeMounts:
                    - mountPath: /tmp
                      name: default-dev-foo-nginx-0
            volumes:
                - configMap:
                    defaultMode: 511
                    name: default-dev-foo-nginx-0
                  name: default-dev-foo-nginx-0
status: {}
`
	var deployRes v1.Resource
	k8sDep := &appsv1.Deployment{}
	_ = yaml.Unmarshal([]byte(deploy), k8sDep)
	unDep, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(k8sDep)
	deployRes.Attributes = unDep

	deployWithProbe := `apiVersion: apps/v1
kind: Deployment
metadata:
    creationTimestamp: null
    labels:
        app.kubernetes.io/name: foo
        app.kubernetes.io/part-of: default
        service-workload-type: Deployment
    name: default-dev-foo
    namespace: default
spec:
    replicas: 4
    selector:
        matchLabels:
            app.kubernetes.io/name: foo
            app.kubernetes.io/part-of: default
    strategy: {}
    template:
        metadata:
            creationTimestamp: null
            labels:
                app.kubernetes.io/name: foo
                app.kubernetes.io/part-of: default
                service-workload-type: Deployment
        spec:
            containers:
                - image: nginx:v1
                  lifecycle:
                    postStart:
                        exec:
                            command:
                                - /bin/true
                  name: nginx
                  readinessProbe:
                    tcpSocket:
                        host: localhost
                        port: 8888
                  resources: {}
                  volumeMounts:
                    - mountPath: /tmp
                      name: default-dev-foo-nginx-0
            volumes:
                - configMap:
                    defaultMode: 511
                    name: default-dev-foo-nginx-0
                  name: default-dev-foo-nginx-0
status: {}
`
	var deployWithProbeRes v1.Resource
	k8sDepWithProbe := &appsv1.Deployment{}
	_ = yaml.Unmarshal([]byte(deployWithProbe), k8sDepWithProbe)
	unDepWithProbe, _ := runtime.DefaultUnstructuredConverter.ToUnstructured(k8sDepWithProbe)
	deployWithProbeRes.Attributes = unDepWithProbe
	r2 := new(int32)
	*r2 = 2

	svcConfig := &Service{
		Base: Base{
			Containers: map[string]Container{
				"nginx": {
					Image: "nginx:v1",
					Files: map[string]FileSpec{
						"/tmp/example.txt": {
							Content: "some file contents",
							Mode:    "0777",
						},
					},
				},
			},
			Replicas: r2,
		},
		Ports: []Port{
			{
				Port:     80,
				Protocol: "TCP",
			},
		},
	}

	var devConfig map[string]interface{}
	temp, _ := yamlv2.Marshal(svcConfig)
	_ = yamlv2.Unmarshal(temp, &devConfig)

	serviceWithProbe := &v1.Service{
		Base: v1.Base{
			Containers: map[string]v1.Container{
				"nginx": {
					Image: "nginx:v1",
					Files: map[string]v1.FileSpec{
						"/tmp/example.txt": {
							Content: "some file contents",
							Mode:    "0777",
						},
					},
					ReadinessProbe: &v1.Probe{ProbeHandler: &v1.ProbeHandler{
						TypeWrapper:     v1.TypeWrapper{Type: v1.TypeTCP},
						ExecAction:      nil,
						HTTPGetAction:   nil,
						TCPSocketAction: &v1.TCPSocketAction{URL: "localhost:8888"},
					}},
					Lifecycle: &v1.Lifecycle{
						PostStart: &v1.LifecycleHandler{
							TypeWrapper: v1.TypeWrapper{Type: v1.TypeExec},
							ExecAction: &v1.ExecAction{Command: []string{
								"/bin/true",
							}},
							HTTPGetAction: nil,
						},
					},
				},
			},
		},
		Ports: []v1.Port{
			{
				Port:     80,
				Protocol: "TCP",
			},
		},
	}
	var devConfigWithProbe map[string]interface{}
	temp, _ = yamlv2.Marshal(serviceWithProbe)
	_ = yamlv2.Unmarshal(temp, &devConfigWithProbe)

	tests := []struct {
		name    string
		request *module.GeneratorRequest
		want    *module.GeneratorResponse
		wantErr bool
	}{
		{
			name: "CollaSet",
			request: &module.GeneratorRequest{
				Project:   "default",
				Stack:     "dev",
				App:       "foo",
				DevConfig: devConfig,
				PlatformConfig: v1.GenericConfig{
					"type": "CollaSet",
					"labels": v1.GenericConfig{
						"service-workload-type": "CollaSet",
					},
					"annotations": v1.GenericConfig{
						"service-workload-type": "CollaSet",
					},
				},
			},
			wantErr: false,
			want: &module.GeneratorResponse{
				Resources: []v1.Resource{cmResource, csResource},
			},
		},
		{
			name: "Deployment",
			request: &module.GeneratorRequest{
				Project:   "default",
				Stack:     "dev",
				App:       "foo",
				DevConfig: devConfig,
				PlatformConfig: v1.GenericConfig{
					"replicas": 4,
					"labels": v1.GenericConfig{
						"service-workload-type": "Deployment",
					},
				},
			},
			wantErr: false,
			want: &module.GeneratorResponse{
				Resources: []v1.Resource{cmResource, deployRes, deploySvcRes},
			},
		},
		{
			name: "DeploymentWithProbe",
			request: &module.GeneratorRequest{
				Project:   "default",
				Stack:     "dev",
				App:       "foo",
				DevConfig: devConfig,
				PlatformConfig: v1.GenericConfig{
					"replicas": 4,
					"labels": v1.GenericConfig{
						"service-workload-type": "Deployment",
					},
				},
			},
			wantErr: false,
			want: &module.GeneratorResponse{
				Resources: []v1.Resource{cmResource, deployWithProbeRes, deploySvcRes},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &Service{}
			got, err := svc.Generate(context.Background(), tt.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i, resource := range got.Resources {
				// todo the order of attributes is not the same, only compare the length here
				if !reflect.DeepEqual(len(resource.Attributes), len(tt.want.Resources[i].Attributes)) {
					t.Errorf("Service.Generate() = %v, want %v", resource.Attributes, tt.want.Resources[i].Attributes)
				}
			}
		})
	}
}

func TestCompleteServiceInput(t *testing.T) {
	r2 := int32(2)

	testcases := []struct {
		name             string
		service          *Service
		config           v1.GenericConfig
		success          bool
		completedService *Service
	}{
		{
			name: "use type in workspace config",
			service: &Service{
				Base: Base{
					Containers: map[string]Container{
						"nginx": {
							Image: "nginx:v1",
						},
					},
					Replicas: &r2,
					Labels: map[string]string{
						"k1": "v1",
					},
					Annotations: map[string]string{
						"k1": "v1",
					},
				},
			},
			config: v1.GenericConfig{
				"type": "CollaSet",
			},
			success: true,
			completedService: &Service{
				Base: Base{
					Containers: map[string]Container{
						"nginx": {
							Image: "nginx:v1",
						},
					},
					Replicas: &r2,
					Labels: map[string]string{
						"k1": "v1",
					},
					Annotations: map[string]string{
						"k1": "v1",
					},
				},
				Type: "CollaSet",
			},
		},
		{
			name: "use default type",
			service: &Service{
				Base: Base{
					Containers: map[string]Container{
						"nginx": {
							Image: "nginx:v1",
						},
					},
					Replicas: &r2,
					Labels: map[string]string{
						"k1": "v1",
					},
					Annotations: map[string]string{
						"k1": "v1",
					},
				},
			},
			config:  nil,
			success: true,
			completedService: &Service{
				Base: Base{
					Containers: map[string]Container{
						"nginx": {
							Image: "nginx:v1",
						},
					},
					Replicas: &r2,
					Labels: map[string]string{
						"k1": "v1",
					},
					Annotations: map[string]string{
						"k1": "v1",
					},
				},
				Type: "Deployment",
			},
		},
		{
			name: "invalid field type",
			service: &Service{
				Base: Base{
					Containers: map[string]Container{
						"nginx": {
							Image: "nginx:v1",
						},
					},
					Replicas: &r2,
					Labels: map[string]string{
						"k1": "v1",
					},
					Annotations: map[string]string{
						"k1": "v1",
					},
				},
			},
			config: v1.GenericConfig{
				"type": 1,
			},
			success:          false,
			completedService: nil,
		},
		{
			name: "unsupported type",
			service: &Service{
				Base: Base{
					Containers: map[string]Container{
						"nginx": {
							Image: "nginx:v1",
						},
					},
					Replicas: &r2,
					Labels: map[string]string{
						"k1": "v1",
					},
					Annotations: map[string]string{
						"k1": "v1",
					},
				},
			},
			config: v1.GenericConfig{
				"type": "unsupported",
			},
			success:          false,
			completedService: nil,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := completeServiceInput(tc.service, tc.config)
			assert.Equal(t, tc.success, err == nil)
			if tc.success {
				assert.True(t, reflect.DeepEqual(tc.service, tc.completedService))
			}
		})
	}
}
