package main

import (
	"context"
	"errors"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kusionstack.io/kube-api/apps/v1alpha1"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
	"kusionstack.io/kusion/pkg/log"
)

type OpsRuleModule struct{}

func (o *OpsRuleModule) Generate(_ context.Context, request *module.GeneratorRequest) (*module.GeneratorResponse, error) {
	// opsRule does not exist in AppConfig and workspace config
	if request.DevConfig == nil && request.PlatformConfig == nil {
		log.Info("OpsRule does not exist in AppConfig and workspace config")
		return nil, nil
	}

	// Job does not support maxUnavailable
	if request.Workload.Header.Type == v1.TypeJob {
		log.Infof("Job does not support opsRule")
		return nil, nil
	}

	if request.Workload.Service.Type == v1.Collaset {
		maxUnavailable, err := GetMaxUnavailable(request.DevConfig, request.PlatformConfig)
		if err != nil {
			return nil, err
		}
		ptr := &v1alpha1.PodTransitionRule{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.GroupVersion.String(),
				Kind:       "PodTransitionRule",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      module.UniqueAppName(request.Project, request.Stack, request.App),
				Namespace: request.Project,
			},
			Spec: v1alpha1.PodTransitionRuleSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: module.UniqueAppLabels(request.Project, request.App),
				},
				Rules: []v1alpha1.TransitionRule{
					{
						Name: "maxUnavailable",
						TransitionRuleDefinition: v1alpha1.TransitionRuleDefinition{
							AvailablePolicy: &v1alpha1.AvailableRule{
								MaxUnavailableValue: &maxUnavailable,
							},
						},
					},
				},
			},
		}
		resourceID := module.KubernetesResourceID(ptr.TypeMeta, ptr.ObjectMeta)
		resource, err := module.WrapK8sResourceToKusionResource(resourceID, ptr)
		if err != nil {
			return nil, err
		}
		return &module.GeneratorResponse{
			Resources: []v1.Resource{*resource},
		}, nil
	}
	return nil, nil
}

func GetMaxUnavailable(devConfig v1.Accessory, platformConfig v1.GenericConfig) (intstr.IntOrString, error) {
	var maxUnavailable interface{}
	key := "maxUnavailable"

	// developer config
	// kusionstack/opsrule@v0.1 : t.OpsRule {
	//    maxUnavailable: "30%"
	// }
	if devConfig != nil && devConfig[key] != "" {
		maxUnavailable = devConfig[key]
	} else if platformConfig == nil {
		return intstr.IntOrString{}, nil
	} else {
		// platformConfig example
		// kusionstack/opsrule@v0.1:
		//   maxUnavailable: 1 # or 10%
		maxUnavailable = platformConfig[key]
	}
	var mu string
	mu, isString := maxUnavailable.(string)
	if !isString {
		temp, isInt := maxUnavailable.(int)
		if isInt {
			mu = strconv.Itoa(temp)
		} else {
			return intstr.IntOrString{}, errors.New("illegal opsRule config. opsRule.maxUnavailable is not string or int")
		}
	}
	return intstr.Parse(mu), nil
}

func main() {
	server.Start(&OpsRuleModule{})
}
