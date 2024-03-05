import models.schema.v1 as ac
import models.schema.v1.trait as t
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.network as n

# base.k declares reusable configurations for all stacks.

hellocollaset: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
                # Run the following command as defined
                command: ["/bin/sh", "-c", "echo hi"]
                # Extra arguments append to command defined above
                args: ["/bin/sh", "-c", "echo hi"]
                # Run the command "/bin/sh -c echo hi", as defined above, in the directory "/tmp"
                workingDir: "/tmp"
            }
        }
        # Ports to expose
        ports: [
            n.Port {
                port: 80
            }
        ]
    }
    # accaccessories :{
    #         "kusionstack/opsrule@v0.1" : t.OpsRule {
    #         maxUnavailable: "30%"
    #     }
    # }
}