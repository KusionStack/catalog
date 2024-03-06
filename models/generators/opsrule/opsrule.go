package main

import (
	"errors"
	"strconv"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kusionstack.io/kube-api/apps/v1alpha1"
	v1 "kusionstack.io/kusion/pkg/apis/core/v1"
	"kusionstack.io/kusion/pkg/apis/core/v1/workload"
	"kusionstack.io/kusion/pkg/log"
	"kusionstack.io/kusion/pkg/modules"
	"kusionstack.io/kusion/pkg/modules/proto"
	jsonutil "kusionstack.io/kusion/pkg/util/json"

	generators "kusion-modules"
)

type OpsRuleModule struct{}

func (o *OpsRuleModule) Generate(r *proto.GeneratorRequest) (*proto.GeneratorResponse, error) {
	emptyResponse := generators.EmptyResponse()
	request, err := generators.NewGeneratorRequest(r)
	if err != nil {
		return nil, err
	}

	// opsRule does not exist in AppConfig and workspace config
	if request.DevModuleConfig == nil && request.PlatformModuleConfig == nil {
		log.Info("OpsRule does not exist in AppConfig and workspace config")
		return emptyResponse, nil
	}

	// Job does not support maxUnavailable
	if request.Workload.Header.Type == workload.TypeJob {
		log.Infof("Job does not support opsRule")
		return emptyResponse, nil
	}

	if request.Workload.Service.Type == workload.Collaset {
		maxUnavailable, err := GetMaxUnavailable(request.DevModuleConfig, request.PlatformModuleConfig)
		if err != nil {
			return nil, err
		}
		ptr := &v1alpha1.PodTransitionRule{
			TypeMeta: metav1.TypeMeta{
				APIVersion: v1alpha1.GroupVersion.String(),
				Kind:       "PodTransitionRule",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      modules.UniqueAppName(request.Project, request.Stack, request.App),
				Namespace: request.App,
			},
			Spec: v1alpha1.PodTransitionRuleSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: modules.UniqueAppLabels(request.Project, request.App),
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
		resourceID := modules.KubernetesResourceID(ptr.TypeMeta, ptr.ObjectMeta)
		resource, err := generators.WrapK8sResourceToKusionResource(resourceID, ptr)
		if err != nil {
			return nil, err
		}
		str := jsonutil.Marshal2String(resource)
		b := []byte(str)
		return &proto.GeneratorResponse{
			Resources: [][]byte{b},
		}, nil
	}
	return emptyResponse, nil
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
	modules.StartModule(&OpsRuleModule{})
}
