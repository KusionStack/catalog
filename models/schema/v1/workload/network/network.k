schema Network:
    """ Network describes how the Application exposes its service in network, it works when the
    workload is Service.
    Currently, Network only supports exposing TCP/HTTP/HTTPS service.

    Attributes
    ----------
    ports: [int], default is Undefined, required.
        The list of exposed ports of the Service's containers.
    exposeInternet: bool, default is Flase, required.
        ExposedInternet determines whether the Service can be accessed through Internet.
    routes: {str: Route}, default is Undefined, optional.
        The hostnames and corresponding HTTP/HTTPS routes to expose the Service through Internet,
        it works only when exposeInternet is True.

    Examples
    --------
    import models.schema.v1.workload.network as n

    network = n.Network {
        ports: [80]
        exposeInternet: True
        routes: {
            "example.helloworld.com": n.Route {
                protocol: "HTTPS"
                paths: [
                    n.RoutePath {
                        path: "/example"
                        port: 80
                    }
                ]
            }
        }
    }
    """

    # The list of ports exposed.
    ports:                      [int]

    # ExposeInternet indicates whether the service can be accessed through Internet.
    exposeInternet:             bool = False

    # Routes is the hostnames and corresponding HTTP/HTTPS routes exposed to Internet.
    routes?:                    {str: Route}

    check:
        sum([0 if (1 <= port <= 65535) else 0 for port in ports]) == 0, "port must be between 1 and 65535, inclusive"
