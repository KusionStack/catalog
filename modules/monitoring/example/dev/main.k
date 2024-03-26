import kam.v1.app_configuration as ac
import kam.v1.workload as wl
import kam.v1.workload.container as c
import monitoring as m

# base.k declares reusable configurations for all stacks.
prometheus: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            nginx: c.Container {
                image: "nginx:v1"
                env: {
                    # An environment variable of name "env1" and value "VALUE" will be set
                    "env1": "VALUE"
                }
            }
        }
    }
    # Add the monitoring configuration backed by Prometheus
    accessories: {
        # Add the monitoring configuration backed by Prometheus
        "monitoring": m.Prometheus {
            path:           "/metrics"
            port:           "web"
        }
    }
}
