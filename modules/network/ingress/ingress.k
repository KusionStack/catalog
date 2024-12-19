schema Ingress:
    """ Ingress is a collection of rules that allow inbound connections to reach the endpoints defined by a backend.
    An Ingress can be configured to give services externally-reachable urls, load balance traffic, terminate SSL,
    offer name based virtual hosting etc.

    Attributes
    ----------
    defaultBackend: IngressBackend, default is Undefined, optional.
        DefaultBackend is the backend that should handle requests that don't match any rule. If Rules are not specified,
        DefaultBackend must be specified. If DefaultBackend is not set, the handling of requests that do not match any
        of the rules will be up to the Ingress controller.
    ingressClassName: str, default is Undefined, optional.
        IngressClassName is the name of an IngressClass cluster resource. Ingress controller implementations use this
        field to know whether they should be serving this Ingress resource, by a transitive connection
        (controller -> IngressClass -> Ingress resource). Although the `kubernetes.io/ingress.class` annotation
        (simple constant name) was never formally defined, it was widely supported by Ingress controllers to create a
        direct binding between Ingress controller and Ingress resources. Newly created Ingress resources should prefer
        using the field. However, even though the annotation is officially deprecated, for backwards compatibility
        reasons, ingress controllers should still honor that annotation if present.
    rules: [IngressRule], default is Undefined, optional.
        Rules is a list of host rules used to configure the Ingress. If unspecified, or no rule matches, all traffic is
        sent to the default backend.
    tls: [IngressTLS], default is Undefined, optional.
        TLS represents the TLS configuration. Currently the Ingress only supports a single TLS port, 443. If multiple
        members of this list specify different hosts, they will be multiplexed on the same port according to the hostname
        specified through the SNI TLS extension, if the ingress controller fulfilling the ingress supports SNI.
    labels: {str:str}, default is Undefined, optional.
        Labels are key/value pairs that are attached to the workload.
    annotations: {str:str}, default is Undefined, optional.
        Annotations are key/value pairs that attach arbitrary non-identifying metadata to the workload.
    """

    # DefaultBackend is the backend that should handle requests that don't match any rule.
    defaultBackend?:                 IngressBackend

    # IngressClassName is the name of an IngressClass cluster resource.
    ingressClassName?:               str

    # Rules is a list of host rules used to configure the Ingress.
    rules?:                          [IngressRule]

    # TLS represents the TLS configuration.
    tls?:                            [IngressTLS]

    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    labels?:                         {str:str}
    annotations?:                    {str:str}