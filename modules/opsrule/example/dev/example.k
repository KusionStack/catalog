import opsrule as o
import kam.v1.app_configuration as ac
import service
import service.container as c

helloworld : ac.AppConfiguration {
    workload: service.Service {
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