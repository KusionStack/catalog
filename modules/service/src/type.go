package main

import (
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	BuiltinModulePrefix = ""
	ProbePrefix         = "service.container.probe."
	TypeHTTP            = BuiltinModulePrefix + ProbePrefix + "Http"
	TypeExec            = BuiltinModulePrefix + ProbePrefix + "Exec"
	TypeTCP             = BuiltinModulePrefix + ProbePrefix + "Tcp"
)

// LabelSelector is a label query over a set of resources.
type LabelSelector struct {
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels
	// map is equivalent to an element of matchExpressions, whose key field is "key", the
	// operator is "In", and the values array contains only "value".
	MatchLabels map[string]string `yaml:"matchLabels,omitempty" json:"matchLabels,omitempty"`
	// matchExpressions is a list of label selector requirements.
	MatchExpressions []LabelSelectorRequirement `yaml:"matchExpressions,omitempty" json:"matchExpressions,omitempty"`
}

// LabelSelectorRequirement is a selector that contains values, a key, and an operator that relates the key and values.
type LabelSelectorRequirement struct {
	// key is the label key that the selector applies to.
	Key string `yaml:"key" json:"key"`
	// operator represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists and DoesNotExist.
	Operator metav1.LabelSelectorOperator `yaml:"operator" json:"operator"`
	// values is an array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. This array is replaced during a strategic merge patch.
	Values []string `yaml:"values,omitempty" json:"values,omitempty"`
}

// TopologySpreadConstraint specifies how to spread matching pods among the given topology.
type TopologySpreadConstraint struct {
	// MaxSkew describes the degree to which pods may be unevenly distributed.
	MaxSkew int32 `yaml:"maxSkew" json:"maxSkew"`
	// TopologyKey is the key of node labels.
	TopologyKey string `yaml:"topologyKey" json:"topologyKey"`
	// WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy the spread constraint.
	WhenUnsatisfiable corev1.UnsatisfiableConstraintAction `yaml:"whenUnsatisfiable" json:"whenUnsatisfiable"`
	// LabelSelector is used to find matching pods.
	LabelSelector *LabelSelector `yaml:"labelSelector,omitempty" json:"labelSelector,omitempty"`
	// MinDomains indicates a minimum number of eligible domains.
	MinDomains *int32 `yaml:"minDomains,omitempty" json:"minDomains,omitempty"`
	// NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew.
	NodeAffinityPolicy *corev1.NodeInclusionPolicy `yaml:"nodeAffinityPolicy,omitempty" json:"nodeAffinityPolicy,omitempty"`
	// NodeTaintsPolicy indicates how we will treat node taints when calculating pod topology spread skew.
	NodeTaintsPolicy *corev1.NodeInclusionPolicy `yaml:"nodeTaintsPolicy,omitempty" json:"nodeTaintsPolicy,omitempty"`
	// MatchLabelKeys is a set of pod label keys to select the pods over which spreading will be calculated.
	MatchLabelKeys []string `yaml:"matchLabelKeys,omitempty" json:"matchLabelKeys,omitempty"`
}

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

// Port defines the exposed port of Service.
type Port struct {
	// Port is the exposed port of the Service.
	Port int `yaml:"port,omitempty" json:"port,omitempty"`
	// TargetPort is the backend .Container port.
	TargetPort int `yaml:"targetPort,omitempty" json:"targetPort,omitempty"`
	// Protocol is protocol used to expose the port, support ProtocolTCP and ProtocolUDP.
	Protocol Protocol `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

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
	// TopologySpreadConstraints describes how a group of pods ought to spread across topology domains.
	// Scheduler will schedule pods in a way which abides by the constraints. All topologySpreadConstraints are ANDed.
	TopologySpreadConstraints map[string]TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty" yaml:"topologySpreadConstraints,omitempty"`
}

type ServiceType string

const (
	ModuleService                 = "service"
	ModuleServiceType             = "type"
	Deployment        ServiceType = "Deployment"
	Collaset          ServiceType = "CollaSet"
)

// Service is a kind of workload profile that describes how to run your application code.
// This is typically used for long-running web applications that should "never" go down, and handle short-lived latency-sensitive
// web requests, or events.
type Service struct {
	Base `yaml:",inline" json:",inline"`
	// Type represents the type of workload.Service, support Deployment and CollaSet.
	Type ServiceType `yaml:"type" json:"type"`
	// Ports describe the list of ports need getting exposed.
	Ports []Port `yaml:"ports,omitempty" json:"ports,omitempty"`
}
