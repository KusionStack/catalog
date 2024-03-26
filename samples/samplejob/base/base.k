import models.schema.v1 as ac
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c

# base.k declares reusable configurations for all stacks.

samplejob: ac.AppConfiguration {
    workload: wl.Job {
        containers: {
            "busybox": c.Container {
                image: "busybox:1.28"
                # Run the following command as defined
                command: ["/bin/sh", "-c", "echo hello"]
            }
        }
        # Run every hour.
        schedule: "0 * * * *"
    }
}