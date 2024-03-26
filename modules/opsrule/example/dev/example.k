import opsrule as o
import kam.v1.app_configuration as ac
import kam.v1.workload as wl
import kam.v1.workload.container as c

helloworld : ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
            }
        }
    }
    accessories: {
        "opsRule": o.OpsRule {
            maxUnavailable: "30%"
        }
    }
}