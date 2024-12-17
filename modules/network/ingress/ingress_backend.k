schema IngressBackend:
    """ IngressBackend describes all endpoints for a given service and port.

    Attributes
    ----------
    resource: TypedLocalObjectReference, default is Undefined, optional.
        Resource is an ObjectRef to another Kubernetes resource in the namespace of the Ingress object. If resource is
        specified, a service.Name and service.Port must not be specified. This is a mutually exclusive setting with
        "Service".
    service: IngressServiceBackend, default is Undefined, optional.
        Service references a service as a backend. This is a mutually exclusive setting with "Resource".
    """

    # Resource is an ObjectRef to another Kubernetes resource in the namespace of the Ingress object.
    resource?:                 TypedLocalObjectReference

    # Service references a service as a backend.
    service?:                  IngressServiceBackend

    check:
        not resource or not service, "resource and number are mutually exclusive"


schema IngressServiceBackend:
    """ IngressServiceBackend references a Kubernetes Service as a Backend.

    Attributes
    ----------
    name: str, default is Undefined, optional.
         Name is the referenced service. The service must exist in the same namespace as the Ingress object.
         If the name is not set, the generated public service name will be used.
    port: ServiceBackendPort, default is Undefined, optional.
         Port of the referenced service. A port name or port number is required for a IngressServiceBackend.
    """

    # Name is the referenced service. The service must exist in the same namespace as the Ingress object.
    # If the name is not set, the generated public service name will be used.
    name?:                 str

    # Port of the referenced service. A port name or port number is required for a IngressServiceBackend.
    port?:                ServiceBackendPort


schema ServiceBackendPort:
    """ ServiceBackendPort is the service port being referenced. A port name or port number is required
    for a IngressServiceBackend.

    Attributes
    ----------
    name: str, default is Undefined, optional.
        Name is the name of the port on the Service. This is a mutually exclusive setting with "Number".
    number: int, default is Undefined, optional.
        Number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with "Name".
    """

    # Name is the name of the port on the Service. This is a mutually exclusive setting with "Number".
    name?:                 str

    # Number is the numerical port number (e.g. 80) on the Service. This is a mutually exclusive setting with "Name".
    number?:               int

    check:
        not name or not number, "name and number are mutually exclusive"