schema IngressClass:
    """ IngressClass represents the class of the Ingress, referenced by the Ingress Spec. The
    `ingressclass.kubernetes.io/is-default-class` annotation can be used to indicate that an IngressClass should be
    considered default. When a single IngressClass resource has this annotation set to true, new Ingress resources
    without a class specified will be assigned this default class.

    Attributes
    ----------
    controller: str, default is Undefined, optional.
        Controller refers to the name of the controller that should handle this class. This allows for different "flavors"
        that are controlled by the same controller. For example, you may have different parameters for the same implementing
        controller. This should be specified as a domain-prefixed path no more than 250 characters in length,
        e.g. "acme.io/ingress-controller". This field is immutable.
    parameters: IngressClassParametersReference, default is Undefined, optional.
        Parameters is a link to a custom resource containing additional configuration for the controller. This is optional
        if the controller does not require extra parameters.
    labels: {str:str}, default is Undefined, optional.
        Labels are key/value pairs that are attached to the workload.
    annotations: {str:str}, default is Undefined, optional.
        Annotations are key/value pairs that attach arbitrary non-identifying metadata to the workload.
    """

    # Controller refers to the name of the controller that should handle this class.
    controller?:                 str

    # Parameters is a link to a custom resource containing additional configuration for the controller.
    parameters?:                 IngressClassParametersReference

    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    labels?:                     {str:str}
    annotations?:                {str:str}

schema IngressClassParametersReference:
    """ IngressClassParametersReference identifies an API object. This can be used to specify a cluster or
    namespace-scoped resource.

    Attributes
    ----------
    kind: str, default is Undefined, required.
        Kind is the type of resource being referenced.
    name: str, default is Undefined, required.
        Name is the name of resource being referenced.
    apiGroup: str, default is Undefined, optional.
        ApiGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be
        in the core API group. For any other third-party types, APIGroup is required.
    namespace: str, default is Undefined, optional.
        Namespace is the namespace of the resource being referenced. This field is required when scope is set to "Namespace"
        and must be unset when scope is set to "Cluster".
    scope: str, default is Undefined, optional.
        Scope represents if this refers to a cluster or namespace scoped resource. This may be set to "Cluster" (default)
        or "Namespace".
    """

    # Kind is the type of resource being referenced.
    kind:                 str

    # Name is the name of resource being referenced.
    name:                 str

    # ApiGroup is the group for the resource being referenced.
    apiGroup?:            str

    # Namespace is the namespace of the resource being referenced.
    namespace?:           str

    # Scope represents if this refers to a cluster or namespace scoped resource.
    scope?:               str

    check:
        scope in ["Namespace", "Cluster"] if scope, "scope value is invalid"