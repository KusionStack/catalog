import models.schema.v1.component.workload as wl

schema Component:
    """ Component defines the delivery artifact of one application. Each component
    can be composed of a Workload and its supporting Accessories.

    Attributes
    ----------
    job: wl.Job, default is Undefined, optional.
        job describes how the Application's tasks are expected to be run.
        A job will contiune to run until it succeeds.
    longRunningService: ws.LongRunningService, default is Undefined, optional.
        longRunningService defines long-running, scalable workload types that should "never" go down, 
        and handle short-lived latency-sensitive requests.
    labels: {str:str}, default is Undefined, optional.
        Labels are key/value pairs that are attached to the component.
    annotations: {str:str}, default is Undefined, optional.
        Annotations are key/value pairs that attach arbitrary non-identifying metadata to the component.

    Examples
    --------
    Instantiate a long-running service Component with "nginx:v1" image and customized labels.

    component = cmp.Component {
        longRunningService: wl.LongRunningService{
            containers: {
                "nginx": c.Container {
                    image: "nginx:v1"
                }
            }
        }
        labels: {
            "key": "value"
        }
    }
    """

    # Job describes how the Application's tasks are expected to be run. 
    # A job will contiune to run until it succeeds.
    job?: wl.Job
    # longRunningService defines long-running, scalable workload types that should "never" go down, 
    # and handle short-lived latency-sensitive requests.
    longRunningService?: wl.LongRunningService

    ###### Other metadata info
    # Labels and annotations can be used to attach arbitrary metadata as key-value pairs to resources.
    labels?:                    {str:str}
    annotations?:               {str:str}

    check:
        not longRunningService if job, "longRunningService and job cannot exist simultaneously."
        not job if longRunningService, "longRunningService and job cannot exist simultaneously."
