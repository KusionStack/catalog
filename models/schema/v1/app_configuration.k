import models.schema.v1.workload as wl
import models.schema.v1.trait as t

schema AppConfiguration:
    """ AppConfiguration is a developer-centric definition that describes how to run an Application.
    This application model builds upon a decade of experience at AntGroup running super large scale
    internal developer platform, combined with best-of-breed ideas and practices from the community.

    Attributes
    ----------
    workload: wl.Service | wl.Job, default is Undefined, required.
        Workload defines how to run your application code. Currently supported workload profile
        includes Service and Job.
    opsRule: t.OpsRule, default is Undefined, optional.
        OpsRule specifies collection of rules that will be checked for Day-2 operation.

    Examples
    --------
    Instantiate an App with a long-running service and its image is "nginx:v1"

    import models.schema.v1 as ac
    import models.schema.v1.workload as wl
    import models.schema.v1.workload.container as c

    appConfiguration = ac.AppConfiguration {
        workload: wl.Service {
            containers: {
                "nginx": c.Container {
                    image: "nginx:v1"
                }
            }
        }
    }
    """

    # Workload defines how to run your application code.
    workload:                   wl.Service | wl.Job
    
    # OpsRule specifies collection of rules that will be checked for Day-2 operation.
    opsRule?:                   t.OpsRule

    ###### Other metadata info

    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    labels?:                    {str:str}
    annotations?:               {str:str}