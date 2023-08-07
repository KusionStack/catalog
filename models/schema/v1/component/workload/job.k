import models.schema.v1.component.container as c

schema Job:
    """ Job describes how the Application's tasks are expected to be run. A job will contiune to run until it succeeds.

    Attributes
    ----------
    containers: {str:c.Container}, default is Undefined, required.
        Containers defines the templates of containers to be ran.
        More info: https://kubernetes.io/docs/concepts/containers
    schedule: str, default is Undefined, optional.
        The scheduling strategy in Cron format, see https://en.wikipedia.org/wiki/Cron.

    Examples
    --------
    import models.schema.v1.component.workload as wl
    import models.schema.v1.component.container as c

    job: wl.Job {
        containers: {
            "busybox": c.Container{
                image:   "busybox:1.28"
                command: ["/bin/sh", "-c", "echo hello"]
            }
            schedule: "0 0 13 * 5"
        }
    }
    """

    # The template of container to be ran.
    containers: {str:c.Container}

    # The scheduling strategy in Cron format.
    schedule?: str
