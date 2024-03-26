schema Network: 
    """ Network describes the network accessories of Workload, which typically contains the exposed ports, load balancer 
    and other related resource configs. 

    Attributes
    ----------
    ports: [n.Port], default is Undefined, optional. 
        The list of ports which the Workload should get exposed. 

    Examples
    --------
    import catalog.models.schema.v1.network as n

    accessories: {
        "network": n.Network {
            ports: [
                n.Port {
                    port: 80
                    public: True
                }
                n.Port {
                    port: 8080
                }
            ]
        }
    }
    """

    # The list of ports getting exposed. 
    ports?:                         [Port]

schema Port:
    """ Port defines the exposed port of Workload, which can be used to describe how the Workload
    get accessed.

    Attributes
    ----------
    port: int, default is 80, required.
        The exposed port of the Workload.
    targetPort: int, default is Undefined, optional.
        The backend container port. If empty, set it the same as the port.
    protocol: "TCP" | "UDP", default is "TCP", required.
        The protocol to access the port.
    public: bool, default is False, required.
        Public defines whether the port can be accessed through Internet.

    Examples
    --------
    import catalog.models.schema.v1.network as n

    port = n.Port {
        port: 80
        targetPort: 8080
        protocol: "TCP"
        public: True
    }
    """

    # The exposed port of the Service.
    port:                       int = 80

    # The backend container port.
    targetPort?:                int

    # The protocol of port.
    protocol:                   "TCP" | "UDP" = "TCP"

    # Public defines whether to expose the port through Internet.
    public:                     bool = False

    check:
        1 <= port <= 65535, "port must be between 1 and 65535, inclusive"
        1 <= targetPort <= 65535 if targetPort, "targetPort must be between 1 and 65535, inclusive"
