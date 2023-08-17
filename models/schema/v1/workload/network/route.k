Schema Route
    """ Route enables exposed container ports accessible by HTTP/HTTPS routing paths, so that the
     container can be accessed outside the cluster.
     There must be an IngressController in the cluster, then the configuration of Route can work.

    Attributes
    ----------
    paths: [RoutePath], default is Undefined, required.
        The list of exposed paths.
    tlsSecret: str, default is Undefined, optional.
        The name of Secret which contains a TLS private key and certificate.

    Examples
    --------
    import models.schema.v1.workload.network as n

    route = n.Route {
        paths: [
            n.RoutePath {
                path: "/example"
                containerAccessPort: 80
            }
        ]
        tlsSecret: "example-tls-secret"
    }
    """

    # List of HTTP/HTTPS route paths.
    paths:                      [RoutePath]

    # The Secret name which contains a TLS private key and certificate.
    tlsSecret?:                 str

Schema RoutePath:
    """ RoutePath defines the HTTP/HTTPS path and corresponding backend container accessPort.

    Attributes
    ----------
    path: str, default is "/", required.
        The path exposed to HTTP/HTTPS routing.
    pathType: "Exact" | "Prefix", default is Undefined, optional.
        PathType defines how the IngressController matches the URL path. If not specified, the
        matching is up to the IngressController. The following values are supported:
        - "Exact": Matches the URL path exactly and with case sensitivity;
        - "Prefix": Matches based on a URL path prefix split by /.
    containerAccessPort: int, default is Undefined, required.
        The corresponding container accessPort that the traffic will forwarding to.

    Examples
    --------
    import models.schema.v1.workload.network as n

    routePath = n.RoutePath {
        path: "/example"
        pathType: "Exact"
        containerAccessPort: 80
    }
    """

    # The path exposed to HTTP/HTTPS routing.
    path:                       str = "/"

    # PathType defines how the URL matches path.
    pathType?:                  "Exact" | "Prefix"

    # The backend container accessPort.
    containerAccessPort:        int

    check:
        1 <= containerAccessPort <= 65535, "containerAccessPort must be between 1 and 65535, inclusive"