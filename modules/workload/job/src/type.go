package main

import "gopkg.in/yaml.v2"

// Job is a kind of workload profile that describes how to run your application code. This is typically used for tasks that take from
// a few seconds to a few days to complete.
type Job struct {
	Base `yaml:",inline" json:",inline"`
	// The scheduling strategy in Cron format: https://en.wikipedia.org/wiki/Cron.
	Schedule string `yaml:"schedule,omitempty" json:"schedule,omitempty"`
}

const (
	BuiltinModulePrefix = ""
	ProbePrefix         = "service.container.probe."
	TypeHTTP            = BuiltinModulePrefix + ProbePrefix + "Http"
	TypeExec            = BuiltinModulePrefix + ProbePrefix + "Exec"
	TypeTCP             = BuiltinModulePrefix + ProbePrefix + "Tcp"
)

// Container describes how the App's tasks are expected to be run.
type Container struct {
	// Image to run for this container
	Image string `yaml:"image" json:"image"`
	// Entrypoint array.
	// The image's ENTRYPOINT is used if this is not provided.
	Command []string `yaml:"command,omitempty" json:"command,omitempty"`
	// Arguments to the entrypoint.
	// The image's CMD is used if this is not provided.
	Args []string `yaml:"args,omitempty" json:"args,omitempty"`
	// Collection of environment variables to set in the container.
	// The value of environment variable may be static text or a value from a secret.
	Env yaml.MapSlice `yaml:"env,omitempty" json:"env,omitempty"`
	// The current working directory of the running process defined in entrypoint.
	WorkingDir string `yaml:"workingDir,omitempty" json:"workingDir,omitempty"`
	// Resource requirements for this container.
	Resources map[string]string `yaml:"resources,omitempty" json:"resources,omitempty"`
	// Files configures one or more files to be created in the container.
	Files map[string]FileSpec `yaml:"files,omitempty" json:"files,omitempty"`
	// Dirs configures one or more volumes to be mounted to the specified folder.
	Dirs map[string]string `yaml:"dirs,omitempty" json:"dirs,omitempty"`
	// Periodic probe of container liveness.
	LivenessProbe *Probe `yaml:"livenessProbe,omitempty" json:"livenessProbe,omitempty"`
	// Periodic probe of container service readiness.
	ReadinessProbe *Probe `yaml:"readinessProbe,omitempty" json:"readinessProbe,omitempty"`
	// StartupProbe indicates that the Pod has successfully initialized.
	StartupProbe *Probe `yaml:"startupProbe,omitempty" json:"startupProbe,omitempty"`
	// Actions that the management system should take in response to container lifecycle events.
	Lifecycle *Lifecycle `yaml:"lifecycle,omitempty" json:"lifecycle,omitempty"`
}

// FileSpec defines the target file in a Container
type FileSpec struct {
	// The content of target file in plain text.
	Content string `yaml:"content,omitempty" json:"content,omitempty"`
	// Source for the file content, might be a reference to a secret value.
	ContentFrom string `yaml:"contentFrom,omitempty" json:"contentFrom,omitempty"`
	// Mode bits used to set permissions on this file.
	Mode string `yaml:"mode" json:"mode"`
}

// TypeWrapper is a thin wrapper to make YAML decoder happy.
type TypeWrapper struct {
	// Type of action to be taken.
	Type string `yaml:"_type" json:"_type"`
}

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
type Probe struct {
	// The action taken to determine the health of a container.
	ProbeHandler *ProbeHandler `yaml:"probeHandler" json:"probeHandler"`
	// Number of seconds after the container has started before liveness probes are initiated.
	InitialDelaySeconds int32 `yaml:"initialDelaySeconds,omitempty" json:"initialDelaySeconds,omitempty"`
	// Number of seconds after which the probe times out.
	TimeoutSeconds int32 `yaml:"timeoutSeconds,omitempty" json:"timeoutSeconds,omitempty"`
	// How often (in seconds) to perform the probe.
	PeriodSeconds int32 `yaml:"periodSeconds,omitempty" json:"periodSeconds,omitempty"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	SuccessThreshold int32 `yaml:"successThreshold,omitempty" json:"successThreshold,omitempty"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	FailureThreshold int32 `yaml:"failureThreshold,omitempty" json:"failureThreshold,omitempty"`
}

// ProbeHandler defines a specific action that should be taken in a probe.
// One and only one of the fields must be specified.
type ProbeHandler struct {
	// Type of action to be taken.
	TypeWrapper `yaml:"_type" json:"_type"`
	// Exec specifies the action to take.
	// +optional
	*ExecAction `yaml:",inline" json:",inline"`
	// HTTPGet specifies the http request to perform.
	// +optional
	*HTTPGetAction `yaml:",inline" json:",inline"`
	// TCPSocket specifies an action involving a TCP port.
	// +optional
	*TCPSocketAction `yaml:",inline" json:",inline"`
}

// ExecAction describes a "run in container" action.
type ExecAction struct {
	// Command is the command line to execute inside the container, the working directory for the
	// command  is root ('/') in the container's filesystem.
	// Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
	Command []string `yaml:"command,omitempty" json:"command,omitempty"`
}

// HTTPGetAction describes an action based on HTTP Get requests.
type HTTPGetAction struct {
	// URL is the full qualified url location to send HTTP requests.
	URL string `yaml:"url,omitempty" json:"url,omitempty"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
}

// TCPSocketAction describes an action based on opening a socket.
type TCPSocketAction struct {
	// URL is the full qualified url location to open a socket.
	URL string `yaml:"url,omitempty" json:"url,omitempty"`
}

// Lifecycle describes actions that the management system should take in response
// to container lifecycle events.
type Lifecycle struct {
	// PreStop is called immediately before a container is terminated due to an
	// API request or management event such as liveness/startup probe failure,
	// preemption, resource contention, etc.
	PreStop *LifecycleHandler `yaml:"preStop,omitempty" json:"preStop,omitempty"`
	// PostStart is called immediately after a container is created.
	PostStart *LifecycleHandler `yaml:"postStart,omitempty" json:"postStart,omitempty"`
}

// LifecycleHandler defines a specific action that should be taken in a lifecycle
// hook. One and only one of the fields, except TCPSocket must be specified.
type LifecycleHandler struct {
	// Type of action to be taken.
	TypeWrapper `yaml:"_type" json:"_type"`
	// Exec specifies the action to take.
	// +optional
	*ExecAction `yaml:",inline" json:",inline"`
	// HTTPGet specifies the http request to perform.
	// +optional
	*HTTPGetAction `yaml:",inline" json:",inline"`
}

type Protocol string

const (
	TCP Protocol = "TCP"
	UDP Protocol = "UDP"
)

type Secret struct {
	Type      string            `yaml:"type" json:"type"`
	Params    map[string]string `yaml:"params,omitempty" json:"params,omitempty"`
	Data      map[string]string `yaml:"data,omitempty" json:"data,omitempty"`
	Immutable bool              `yaml:"immutable,omitempty" json:"immutable,omitempty"`
}

const (
	FieldLabels      = "labels"
	FieldAnnotations = "annotations"
	FieldReplicas    = "replicas"
)

// Base defines set of attributes shared by different workload profile, e.g. Service and Job.
type Base struct {
	// The templates of containers to be run.
	Containers map[string]Container `yaml:"containers,omitempty" json:"containers,omitempty"`
	// The number of containers that should be run.
	Replicas *int32 `yaml:"replicas,omitempty" json:"replicas,omitempty"`
	// Secret
	Secrets map[string]Secret `json:"secrets,omitempty" yaml:"secrets,omitempty"`
	// Dirs configures one or more volumes to be mounted to the specified folder.
	Dirs map[string]string `json:"dirs,omitempty" yaml:"dirs,omitempty"`
	// Labels and Annotations can be used to attach arbitrary metadata as key-value pairs to resources.
	Labels      map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
}
