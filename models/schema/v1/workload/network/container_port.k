schema ContainerPort:
    """ ContainerPort defines the available port on the container, and how other containers
    access the port.

    Attributes
    ----------
    port: int, default is 80, required.
        The available port on the container.
    accessMode: "Internal" | "Exposed", default is "Internal", required.
        AccessMode determines which containers can access the port. The following values are
        supported:
        - "Internal": accessed by the containers in the same Application;
        - "Exposed": accessed by the containers within the cluster.
    accessPort: int, default is Undefined, optional.
        AccessPort works when accessMode is "Exposed", defines the port that the container
        exposes within the cluster.
        If accessMode is "Exposed", accessPort defaults to the same as the port.
    accessProtocol: "TCP" | "UDP", default is "TCP", optional.
        accessProtocol defines the protocol of accessPort, works when accessMode is "Exposed".

    Examples
    --------
    import models.schema.v1.workload.network as n

    port = n.ContainerPort {
        port: 80
        accessMode: "Exposed"
        accessPort: 80
    }
    """

    # The available port on the container.
    port:                       int = 80

    # AccessMode defines the accessibility of the container port.
    accessMode:                 "Internal" | "Exposed" = "Internal"

    # AccessPort is the port exposed within the cluster.
    accessPort?:                int

    # The protocol of accessPort.
    accessProtocol?:            "TCP" | "UDP" = "TCP"

    check:
        1 <= port <= 65535, "port must be between 1 and 65535, inclusive"
        1 <= accessPort <= 65535 if accessPort, "accessPort must be between 1 and 65535, inclusive"