import models.schema.v1 as ac
import models.schema.v1.trait as t
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.secret as sec
import models.schema.v1.workload.network as n
import models.schema.v1.monitoring as m
import models.schema.v1.accessories.mysql

# base.k declares reusable configurations for all stacks.
wordpress: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            "wordpress": c.Container {
                image: "wordpress:latest"
                env: {
                    "WORDPRESS_DB_HOST": "secret://wordpress-db/hostAddress"
                    "WORDPRESS_DB_PASSWORD": "secret://wordpress-db/password"
                }
                resources: {
                    "cpu": "1"
                    "memory": "2Gi"
                }
            }
        }
        replicas: 2
        ports: [
            n.Port {
                type: "aliyun"
                port: 80
                public: True
            }
        ]
    }
    mysql: mysql.MySQL {
        type: "cloud"
        version: "5.7"
    }
}
