import models.schema.v1.component.container as c

schema LongRunningService:
    """ LongRunningService defines long-running, scalable workload types that should "never"
    go down, and handle short-lived latency-sensitive requests. Such services are commonly
    used for web applications and services that expose APIs.

    Attributes
    ----------
    containers: {str:c.Container}, default is Undefined, required.
        Containers defines the templates of containers to be ran.
        More info: https://kubernetes.io/docs/concepts/containers
    replicas: int, default is Undefined, optional.
        Replicas is the number of desired replicas.

    Examples
    --------
    import models.schema.v1.component.workload as wl
    import models.schema.v1.component.container as c

    webservice: wl.LongRunningService {
        containers: {
            "web": c.Container {
               image:   "nginx:latest"
               command: ["/bin/sh", "-c", "echo hi"]
               env: {
                   "name": "value"
               }
            }
        }
        replicas: 2
    }
    """
    # Type of long running service.
    $type:          "Deployment" | "StatefulSet" = "Deployment"

    # The templates of containers to be ran.
    containers:     {str:c.Container}

    # The number of containers that should be ran.
    # Default is 2 to meet high availability requirements.
    # Not supported when type is Job.
    replicas:       int = 2
