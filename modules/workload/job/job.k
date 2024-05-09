import common as c

schema Job(c.WorkloadBase):
    """ Job is a kind of workload profile that describes how to run your application code. This
    is typically used for tasks that take from a few seconds to a few days to complete.

    Attributes
    ----------
    schedule: str, default is Undefined, required.
        The scheduling strategy in Cron format. More info: https://en.wikipedia.org/wiki/Cron.

    Examples
    --------
    Instantiate a job with busybox image and runs every hour

    import catalog.workload as wl
    import catalog.workload.container as c

    echoJob : wl.Job {
        containers: {
            "busybox": c.Container{
                image:   "busybox:1.28"
                command: ["/bin/sh", "-c", "echo hello"]
            }
        }
        schedule: "0 * * * *"
    }
    """

    # The scheduling strategy in Cron format.
    # More info: https://en.wikipedia.org/wiki/Cron.
    schedule:                   str
