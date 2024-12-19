package main

import (
	"errors"

	k8snetworking "k8s.io/api/networking/v1"
)

const (
	FieldType        = "type"
	FieldLabels      = "labels"
	FieldAnnotations = "annotations"
)

const (
	CSPAWS      = "aws"
	CSPAliCloud = "alicloud"
)

const (
	ProtocolTCP = "TCP"
	ProtocolUDP = "UDP"
)

const (
	K8sKindIngress      = "Ingress"
	K8sKindIngressClass = "IngressClass"
	k8sKindService      = "Service"
	suffixPublic        = "public"
	suffixPrivate       = "private"
	ingressSuffix       = "ingress"
	ingressClassSuffix  = "ingressclass"
)

var (
	ErrEmptyPortConfig   = errors.New("empty port config")
	ErrEmptyType         = errors.New("type must not be empty when public")
	ErrUnsupportedType   = errors.New("type only support alicloud and aws for now")
	ErrInvalidPort       = errors.New("port must be between 1 and 65535")
	ErrInvalidTargetPort = errors.New("targetPort must be between 1 and 65535 if exist")
	ErrInvalidProtocol   = errors.New("protocol must be TCP or UDP")
	ErrEmptySvcWorkload  = errors.New("network port should be binded to a service workload")
)

// Network describes the network accessories of workload, which typically contains the exposed
// ports, load balancer and other related resource configs.
type Network struct {
	Ports        []Port        `yaml:"ports,omitempty" json:"ports,omitempty"`
	Ingress      *Ingress      `yaml:"ingress,omitempty" json:"ingress,omitempty"`
	IngressClass *IngressClass `yaml:"ingressClass,omitempty" json:"ingressClass,omitempty"`
}

// Port defines the exposed port of workload, which can be used to describe how
// the workload get accessed.
type Port struct {
	// Type is the specific cloud vendor that provides load balancer, works when Public
	// is true, supports CSPAliCloud and CSPAWS for now.
	Type string `yaml:"type,omitempty" json:"type,omitempty"`

	// Port is the exposed port of the workload.
	Port int `yaml:"port,omitempty" json:"port,omitempty"`

	// TargetPort is the backend container.Container port.
	TargetPort int `yaml:"targetPort,omitempty" json:"targetPort,omitempty"`

	// Protocol is protocol used to expose the port, support ProtocolTCP and ProtocolUDP.
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`

	// Public defines whether to expose the port through Internet.
	Public bool `yaml:"public,omitempty" json:"public,omitempty"`

	// Labels are the attached labels of the port, works only when the Public is true.
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`

	// Annotations are the attached annotations of the port, works only when the Public is true.
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

// Ingress is a collection of rules that allow inbound connections to reach the
// endpoints defined by a backend. An Ingress can be configured to give services
// externally-reachable urls, load balance traffic, terminate SSL, offer name
// based virtual hosting etc.
type Ingress struct {
	// Labels are the attached labels of the port, works only when the Public is true.
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`

	// Annotations are the attached annotations of the port, works only when the Public is true.
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`

	// ingressClassName is the name of the IngressClass cluster resource. The
	// associated IngressClass defines which controller will implement the
	// resource. This replaces the deprecated `kubernetes.io/ingress.class`
	// annotation. For backwards compatibility, when that annotation is set, it
	// must be given precedence over this field. The controller may emit a
	// warning if the field and annotation have different values.
	// Implementations of this API should ignore Ingresses without a class
	// specified. An IngressClass resource may be marked as default, which can
	// be used to set a default value for this field. For more information,
	// refer to the IngressClass documentation.
	IngressClassName *string `yaml:"ingressClassName,omitempty" json:"ingressClassName,omitempty"`

	// defaultBackend is the backend that should handle requests that don't
	// match any rule. If Rules are not specified, DefaultBackend must be specified.
	// If DefaultBackend is not set, the handling of requests that do not match any
	// of the rules will be up to the Ingress controller.
	DefaultBackend *IngressBackend `yaml:"defaultBackend,omitempty" json:"defaultBackend,omitempty"`

	// tls represents the TLS configuration. Currently the ingress only supports a
	// single TLS port, 443. If multiple members of this list specify different hosts,
	// they will be multiplexed on the same port according to the hostname specified
	// through the SNI TLS extension, if the ingress controller fulfilling the
	// ingress supports SNI.
	TLS []IngressTLS `yaml:"tls,omitempty" json:"tls,omitempty"`

	// rules is a list of host rules used to configure the Ingress. If unspecified, or
	// no rule matches, all traffic is sent to the default backend.
	Rules []IngressRule `yaml:"rules,omitempty" json:"rules,omitempty"`
}

// IngressBackend describes all endpoints for a given service and port.
type IngressBackend struct {
	// service references a service as a backend.
	// This is a mutually exclusive setting with "Resource".
	Service *IngressServiceBackend `yaml:"service,omitempty" json:"service,omitempty"`

	// resource is an ObjectRef to another Kubernetes resource in the namespace
	// of the Ingress object. If resource is specified, a service.Name and
	// service.Port must not be specified.
	// This is a mutually exclusive setting with "Service".
	Resource *TypedLocalObjectReference `yaml:"resource,omitempty" json:"resource,omitempty"`
}

// IngressServiceBackend references a Kubernetes Service as a Backend.
type IngressServiceBackend struct {
	// name is the referenced service.
	// The service must exist in the same namespace as the Ingress object.
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// port of the referenced service.
	// A port name or port number is required for a IngressServiceBackend.
	Port ServiceBackendPort `yaml:"port,omitempty" json:"port,omitempty"`
}

// ServiceBackendPort is the service port being referenced.
type ServiceBackendPort struct {
	// name is the name of the port on the Service.
	// This must be an IANA_SVC_NAME (following RFC6335).
	// This is a mutually exclusive setting with "Number".
	Name string `yaml:"name,omitempty" json:"name,omitempty"`

	// number is the numerical port number (e.g. 80) on the Service.
	// This is a mutually exclusive setting with "Name".
	Number int32 `yaml:"number,omitempty" json:"number,omitempty"`
}

// TypedLocalObjectReference contains enough information to let you locate the typed referenced object inside the same namespace.
type TypedLocalObjectReference struct {
	// APIGroup is the group for the resource being referenced.
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	APIGroup *string `yaml:"apiGroup,omitempty" json:"apiGroup,omitempty"`

	// Kind is the type of resource being referenced
	Kind string `yaml:"kind" json:"kind"`

	// Name is the name of resource being referenced
	Name string `yaml:"name" json:"name"`
}

// IngressTLS describes the transport layer security associated with an ingress.
type IngressTLS struct {
	// hosts is a list of hosts included in the TLS certificate. The values in
	// this list must match the name/s used in the tlsSecret. Defaults to the
	// wildcard host setting for the loadbalancer controller fulfilling this
	// Ingress, if left unspecified.
	Hosts []string `yaml:"hosts,omitempty" json:"hosts,omitempty"`

	// secretName is the name of the secret used to terminate TLS traffic on
	// port 443. Field is left optional to allow TLS routing based on SNI
	// hostname alone. If the SNI host in a listener conflicts with the "Host"
	// header field used by an IngressRule, the SNI host is used for termination
	// and value of the "Host" header is used for routing.
	SecretName string `yaml:"secretName,omitempty" json:"secretName,omitempty"`
}

// IngressRule represents the rules mapping the paths under a specified host to
// the related backend services. Incoming requests are first evaluated for a
// host match, then routed to the backend associated with the matching
// IngressRuleValue.
type IngressRule struct {
	// host is the fully qualified domain name of a network host, as defined by RFC 3986.
	// Note the following deviations from the "host" part of the
	// URI as defined in RFC 3986:
	// 1. IPs are not allowed. Currently an IngressRuleValue can only apply to
	//    the IP in the Spec of the parent Ingress.
	// 2. The `:` delimiter is not respected because ports are not allowed.
	//	  Currently the port of an Ingress is implicitly :80 for http and
	//	  :443 for https.
	// Both these may change in the future.
	// Incoming requests are matched against the host before the
	// IngressRuleValue. If the host is unspecified, the Ingress routes all
	// traffic based on the specified IngressRuleValue.
	//
	// host can be "precise" which is a domain name without the terminating dot of
	// a network host (e.g. "foo.bar.com") or "wildcard", which is a domain name
	// prefixed with a single wildcard label (e.g. "*.foo.com").
	// The wildcard character '*' must appear by itself as the first DNS label and
	// matches only a single label. You cannot have a wildcard label by itself (e.g. Host == "*").
	// Requests will be matched against the host field in the following way:
	// 1. If host is precise, the request matches this rule if the http host header is equal to Host.
	// 2. If host is a wildcard, then the request matches this rule if the http host header
	// is to equal to the suffix (removing the first label) of the wildcard rule.
	Host string `yaml:"host,omitempty" json:"host,omitempty"`

	// HTTP is a list of http selectors pointing to backends.
	HTTP *HTTPIngressRuleValue `yaml:"http,omitempty" json:"http,omitempty"`
}

// HTTPIngressRuleValue is a list of http selectors pointing to backends.
type HTTPIngressRuleValue struct {
	// paths is a collection of paths that map requests to backends.
	Paths []HTTPIngressPath `yaml:"paths" json:"paths"`
}

// HTTPIngressPath associates a path with a backend. Incoming urls matching the path are forwarded to the backend.
type HTTPIngressPath struct {
	// path is matched against the path of an incoming request. Currently it can
	// contain characters disallowed from the conventional "path" part of a URL
	// as defined by RFC 3986. Paths must begin with a '/' and must be present
	// when using PathType with value "Exact" or "Prefix".
	Path string `yaml:"path,omitempty" json:"path,omitempty"`

	// pathType determines the interpretation of the path matching. PathType can
	// be one of Exact, Prefix, or ImplementationSpecific. Implementations are
	// required to support all path types.
	PathType k8snetworking.PathType `yaml:"pathType" json:"pathType"`

	// backend defines the referenced service endpoint to which the traffic
	// will be forwarded to.
	Backend IngressBackend `yaml:"backend" json:"backend"`
}

// IngressClass provides information about the class of an Ingress.
type IngressClass struct {
	// Labels are the attached labels of the port, works only when the Public is true.
	Labels map[string]string `yaml:"labels,omitempty" json:"labels,omitempty"`

	// Annotations are the attached annotations of the port, works only when the Public is true.
	Annotations map[string]string `yaml:"annotations,omitempty" json:"annotations,omitempty"`

	// controller refers to the name of the controller that should handle this
	// class. This allows for different "flavors" that are controlled by the
	// same controller. For example, you may have different parameters for the
	// same implementing controller. This should be specified as a
	// domain-prefixed path no more than 250 characters in length, e.g.
	// "acme.io/ingress-controller". This field is immutable.
	Controller string `yaml:"controller,omitempty" json:"controller,omitempty"`

	// parameters is a link to a custom resource containing additional
	// configuration for the controller. This is optional if the controller does
	// not require extra parameters.
	// +optional
	Parameters *IngressClassParametersReference `yaml:"parameters,omitempty" json:"parameters,omitempty"`
}

// IngressClassParametersReference identifies an API object. This can be used
// to specify a cluster or namespace-scoped resource.
type IngressClassParametersReference struct {
	// apiGroup is the group for the resource being referenced. If apiGroup is
	// not specified, the specified kind must be in the core API group. For any
	// other third-party types, apiGroup is required.
	APIGroup *string `yaml:"apiGroup,omitempty" json:"apiGroup,omitempty"`

	// kind is the type of resource being referenced.
	Kind string `yaml:"kind" json:"kind"`

	// name is the name of resource being referenced.
	Name string `yaml:"name" json:"name"`

	// scope represents if this refers to a cluster or namespace scoped resource.
	// This may be set to "Cluster" (default) or "Namespace".
	Scope *string `yaml:"scope,omitempty" json:"scope,omitempty"`

	// namespace is the namespace of the resource being referenced. This field is
	// required when scope is set to "Namespace" and must be unset when scope is set to
	// "Cluster".
	Namespace *string `yaml:"namespace,omitempty" json:"namespace,omitempty"`
}
