schema Service(WorkloadBase):
    """ Service is a kind of workload profile that describes how to run your application code. This
    is typically used for long-running web applications that should "never" go down, and handle
    short-lived latency-sensitive web requests, or events.

    Attributes
    ----------
    type: str = "Deployment" | "CollaSet", default is Deployment, optianl.
        type represents the type of workload used by this Service. Currently, it supports several
        types, including Deployment and CollaSet.

    Examples
    --------
    Instantiate a long-running service and its image is "nginx:v1"

    import models.schema.v1.workload as wl
    import models.schema.v1.workload.container as c

    svc = wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
            }
        }
    }
    """

    # Type represents the type of workload used by this Service. Currently, it supports several
    types, including Deployment and CollaSet.
    type:       "Deployment" | "CollaSet" = "Deployment"
