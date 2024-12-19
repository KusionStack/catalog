schema IngressRule:
    """ IngressRule represents the rules mapping the paths under a specified host to the related backend services.
    Incoming requests are first evaluated for a host match, then routed to the backend associated with the matching IngressRuleValue.

    Attributes
    ----------
    host: str, default is Undefined, optional.
        Host is the fully qualified domain name of a network host, as defined by RFC 3986. Note the following deviations
        from the "host" part of the URI as defined in RFC 3986: 1. IPs are not allowed. Currently an IngressRuleValue can
        only apply to the IP in the Spec of the parent Ingress. 2. The : delimiter is not respected because ports are not
        allowed. Currently the port of an Ingress is implicitly :80 for http and :443 for https. Both these may change in
        the future. Incoming requests are matched against the host before the IngressRuleValue. If the host is unspecified,
        the Ingress routes all traffic based on the specified IngressRuleValue.
        Host can be "precise" which is a domain name without the terminating dot of a network host (e.g. "foo.bar.com")
        or "wildcard", which is a domain name prefixed with a single wildcard label (e.g. ".foo.com"). The wildcard
        character '' must appear by itself as the first DNS label and matches only a single label. You cannot have a
        wildcard label by itself (e.g. Host == "*"). Requests will be matched against the Host field in the following
        way: 1. If host is precise, the request matches this rule if the http host header is equal to Host. 2. If host is
        a wildcard, then the request matches this rule if the http host header is to equal to the suffix (removing the
        first label) of the wildcard rule.
    http: HTTPIngressRuleValue, default is Undefined, optional.
        HTTPIngressRuleValue is a list of http selectors pointing to backends. In the example: http:///? -> backend where
        parts of the url correspond to RFC 3986, this resource will be used to match against everything after the last '/'
        and before the first '?' or '#'.
    """

    # Host is the fully qualified domain name of a network host, as defined by RFC 3986.
    host?:                 str

    # HTTPIngressRuleValue is a list of http selectors pointing to backends.
    http?:                 HTTPIngressRuleValue


schema HTTPIngressRuleValue:
    """ HTTPIngressRuleValue is a list of http selectors pointing to backends. In the example:
    http://<host>/<path>?<searchpart> -> backend where where parts of the url correspond to RFC 3986, this resource will
    be used to match against everything after the last '/' and before the first '?' or '#'.

    Attributes
    ----------
    paths: [HTTPIngressPath], default is Undefined, required.
        Paths is a collection of paths that map requests to backends.
    """

    # Paths is a collection of paths that map requests to backends.
    paths:                 [HTTPIngressPath]


schema HTTPIngressPath:
    """ HTTPIngressPath associates a path with a backend. Incoming urls matching the path are forwarded to the backend.

    Attributes
    ----------
    backend: IngressBackend, default is Undefined, required.
        Backend defines the referenced service endpoint to which the traffic will be forwarded to.
    pathType: str, default is Undefined, required.
        PathType determines the interpretation of the path matching. PathType can be one of the following values:
        * Exact: Matches the URL path exactly. * Prefix: Matches based on a URL path prefix split by '/'. Matching is
        done on a path element by element basis. A path element refers is the list of labels in the path split by the '/'
        separator. A request is a match for path p if every p is an element-wise prefix of p of the request path. Note
        that if the last element of the path is a substring of the last element in request path, it is not a match
        (e.g. /foo/bar matches /foo/bar/baz, but does not match /foo/barbaz).
        ImplementationSpecific: Interpretation of the Path matching is up to the IngressClass. Implementations can treat
        this as a separate PathType or treat it identically to Prefix or Exact path types. Implementations are required
        to support all path types.
    path: str, default is Undefined, optional.
        Path is matched against the path of an incoming request. Currently it can contain characters disallowed from the
        conventional "path" part of a URL as defined by RFC 3986. Paths must begin with a '/' and must be present when
        using PathType with value "Exact" or "Prefix".
    """

    # Backend defines the referenced service endpoint to which the traffic will be forwarded to.
    backend:                 IngressBackend

    # PathType determines the interpretation of the path matching.
    pathType:                str

    # Path is matched against the path of an incoming request.
    path?: str

    check:
        pathType in ["Exact", "Prefix", "ImplementationSpecific"] if pathType, "pathType value is invalid"