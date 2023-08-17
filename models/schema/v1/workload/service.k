import models.schema.v1.workload.network as n

schema Service(WorkloadBase):
    """ Service is a kind of workload profile that describes how to run your application code. This
    is typically used for long-running web applications that should "never" go down, and handle
    short-lived latency-sensitive web requests, or events.

    Attributes
    ----------
    routes: {str:n.Route}, default is Undefined, optional.
        Routes describes the paths based on HTTP/HTTPS routing which can access the Service's exposed
        containers from outside the cluster.
        There must be one or more containers whose ports' accessMode are "Exposed", so that the
        configuration of routes can work.
        Multiple hostnames are supported, which are set in the keys. If the key is empty, then the
        route applies the cluster specified IP address.
    type: str = "Deployment" | "CollaSet", default is Deployment, optianl.
        type represents the type of workload used by this Service. Currently, it supports several types, including Deployment and CollaSet.

    Examples
    --------
    # Instantiate a long-running service and its image is "nginx:v1"

    import models.schema.v1.workload as wl
    import models.schema.v1.workload.container as c
    import models.schema.v1.workload.network as n

    svc = wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
                ports: [
                    n.ContainerPort {
                        port: 80
                        accessMode: "Exposed"
                        accessPort: 80
                    }
                ]
            }
        }
        routes: {
            "nginx.example.com": n.Route {
                paths: [
                    n.RoutePath {
                        path: "/example"
                        containerAccessPort: 80
                    }
                ]
            }
        }
    }
    """

    # More service workload attributes here

    # type represents the type of workload used by this Service. Currently, it supports several types, including Deployment and CollaSet.
    type:       "Deployment" | "CollaSet" = "Deployment"

    # The hostnames and corresponding HTTP/HTTPS routes to expose the Service outside the cluster.
    routes?:                     {str : n.Route}