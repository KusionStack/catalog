import models.schema.v1 as ac
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.container.probe as p
import models.schema.v1.monitoring as m
import models.schema.v1.workload.network as n

# base.k declares reusable configurations for all stacks.
prometheus: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
                env: {
                    # An environment variable of name "env1" and value "VALUE" will be set
                    "env1": "VALUE"
                }
                # Run the command "/bin/sh -c echo hi", as defined above, in the directory "/tmp"
                workingDir: "/tmp"
                # Configure a HTTP readiness probe
                readinessProbe: p.Probe {
                    probeHandler: p.Http {
                        url: "http://localhost:8080"
                    }
                    initialDelaySeconds: 10
                }
            }
        }
        # Ports to expose
        ports: [
            n.Port {
                port: 8080
            }
        ]
        # Set the replicas
        replicas: 2
    }
    # Add the monitoring configuration backed by Prometheus
    # monitoring: m.Prometheus{
    #     path:           "/metrics"
    # }
}