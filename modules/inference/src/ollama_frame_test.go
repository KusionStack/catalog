package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"kusionstack.io/kusion-module-framework/pkg/module"
	v1 "kusionstack.io/kusion/pkg/apis/api.kusion.io/v1"
)

func TestInferenceModule_GenerateOllamaResource(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &v1.Workload{
			Header: v1.Header{
				Type: "Service",
			},
			Service: &v1.Service{},
		},
	}

	infer := &Inference{
		Model:       "qwen",
		Framework:   "Ollama",
		System:      "",
		Template:    "",
		TopK:        40,
		TopP:        0.9,
		Temperature: 0.8,
		NumPredict:  128,
		NumCtx:      2048,
	}

	res, patch, err := infer.GenerateOllamaResource(r)

	assert.NotNil(t, res)
	assert.NotNil(t, patch)
	assert.NoError(t, err)
}

func TestInferenceModule_GenerateOllamaPodSpec(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &v1.Workload{
			Header: v1.Header{
				Type: "Service",
			},
			Service: &v1.Service{},
		},
	}

	infer := &Inference{
		Model:       "qwen",
		Framework:   "Ollama",
		System:      "",
		Template:    "",
		TopK:        40,
		TopP:        0.9,
		Temperature: 0.8,
		NumPredict:  128,
		NumCtx:      2048,
	}

	res, err := infer.generateOllamaPodSpec(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestInferenceModule_GenerateOllamaDeployment(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &v1.Workload{
			Header: v1.Header{
				Type: "Service",
			},
			Service: &v1.Service{},
		},
	}

	infer := &Inference{
		Model:       "qwen",
		Framework:   "Ollama",
		System:      "",
		Template:    "",
		TopK:        40,
		TopP:        0.9,
		Temperature: 0.8,
		NumPredict:  128,
		NumCtx:      2048,
	}

	res, err := infer.generateOllamaDeployment(r)

	assert.NotNil(t, res)
	assert.NoError(t, err)
}

func TestInferenceModule_GenerateOllamaService(t *testing.T) {
	r := &module.GeneratorRequest{
		Project: "test-project",
		Stack:   "test-stack",
		App:     "test-app",
		Workload: &v1.Workload{
			Header: v1.Header{
				Type: "Service",
			},
			Service: &v1.Service{},
		},
	}

	infer := &Inference{
		Model:       "qwen",
		Framework:   "Ollama",
		System:      "",
		Template:    "",
		TopK:        40,
		TopP:        0.9,
		Temperature: 0.8,
		NumPredict:  128,
		NumCtx:      2048,
	}

	res, svcName, err := infer.generateOllamaService(r)

	assert.NotNil(t, res)
	assert.NotNil(t, svcName)
	assert.Equal(t, strings.ToLower(infer.Framework)+inferServiceSuffix, svcName)
	assert.NoError(t, err)
}
