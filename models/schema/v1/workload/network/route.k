schema Route:
    """ Route enables Service accessible through Internet by HTTP/HTTPS routing paths.

    Attributes
    ----------
    protocol: "HTTP" | "HTTPS", default is "HTTPS", required.
        The protocol used to access the route.
    paths: [RoutePath], default is Undefined, required.
        The list of exposed paths.

    Examples
    --------
    import models.schema.v1.workload.network as n

    route = n.Route {
        protocol: "HTTPS"
        paths: [
            n.RoutePath {
                path: "/example"
                port: 80
            }
        ]
    }
    """

    # The protocol used to access the route.
    protocol:                   "HTTP" | "HTTPS" = "HTTPS"

    # List of HTTP/HTTPS route paths.
    paths:                      [RoutePath]

schema RoutePath:
    """ RoutePath defines the HTTP/HTTPS path and corresponding Service port.

    Attributes
    ----------
    path: str, default is "/", required.
        The path exposed to HTTP/HTTPS routing.
    pathType: "Exact" | "Prefix", default is Undefined, optional.
        PathType defines how the IngressController matches the URL path. If not specified, the
        matching is up to the IngressController. The following values are supported:
        - "Exact": Matches the URL path exactly and with case sensitivity;
        - "Prefix": Matches based on a URL path prefix split by /.
    port: int, default is Undefined, required.
        The corresponding Service port that the traffic will forwarding to.

    Examples
    --------
    import models.schema.v1.workload.network as n

    routePath = n.RoutePath {
        path: "/example"
        pathType: "Exact"
        port: 80
    }
    """

    # The path exposed to HTTP/HTTPS routing.
    path:                       str = "/"

    # PathType defines how the URL matches path.
    pathType?:                  "Exact" | "Prefix"

    # The corresponding Service port.
    port:                       int

    check:
        1 <= port <= 65535, "port must be between 1 and 65535, inclusive"