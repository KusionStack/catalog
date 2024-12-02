package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"kusionstack.io/kube-api/apps/v1alpha1"
	kusionapiv1 "kusionstack.io/kusion-api-go/api.kusion.io/v1"
	"kusionstack.io/kusion-module-framework/pkg/log"
	"kusionstack.io/kusion-module-framework/pkg/module"
	"kusionstack.io/kusion-module-framework/pkg/server"
)

type OpsRuleModule struct{}

func (o *OpsRuleModule) Generate(ctx context.Context, request *module.GeneratorRequest) (response *module.GeneratorResponse, err error) {
	// Get the module logger with the generator context.
	logger := log.GetModuleLogger(ctx)
	logger.Info("Generating resources...")

	defer func() {
		if r := recover(); r != nil {
			logger.Debug("failed to generate opsrule module: %v", r)
			response = nil
			rawRequest, _ := json.Marshal(request)
			err = fmt.Errorf("panic in opsrule module generator but recovered with error: [%v] and stack %v and request %v",
				r, string(debug.Stack()), string(rawRequest))
		}
	}()

	// opsRule does not exist in AppConfig and workspace config
	if request.DevConfig == nil && request.PlatformConfig == nil {
		log.Info("OpsRule does not exist in AppConfig and workspace config")
		return nil, nil
	}

	// Job does not support maxUnavailable
	if workloadType, ok := request.Workload["_type"]; ok && strings.Contains(workloadType.(string), ".Job") {
		log.Infof("Job does not support opsRule")
		return nil, nil
	}

	if workloadType, ok := request.Workload["type"]; ok && strings.ToLower(workloadType.(string)) == "collaset" {
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
			Resources: []kusionapiv1.Resource{*resource},
		}, nil
	}
	return nil, nil
}

func GetMaxUnavailable(devConfig kusionapiv1.Accessory, platformConfig kusionapiv1.GenericConfig) (intstr.IntOrString, error) {
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
