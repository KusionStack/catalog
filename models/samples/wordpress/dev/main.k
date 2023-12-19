import models.schema.v1 as ac
import models.schema.v1.trait as t
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.container.probe as p
import models.schema.v1.workload.secret as sec
import models.schema.v1.workload.network as n
import models.schema.v1.monitoring as m
import models.schema.v1.accessories.mysql

# main.k declares reusable configurations for dev stacks.
wordpress: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            wordpress: c.Container {
                image = "wordpress:6.3"
                env: {
                    "WORDPRESS_DB_HOST": "$(KUSION_DB_HOST)"
                    "WORDPRESS_DB_USER": "$(KUSION_DB_USERNAME)"
                    "WORDPRESS_DB_PASSWORD": "$(KUSION_DB_PASSWORD)"
                    "WORDPRESS_DB_NAME": "mysql"
                }
                resources: {
                    "cpu": "500m"
                    "memory": "512Mi"
                }
            }
        }
        replicas: 1
        ports: [
            n.Port {
                port: 80
            }
        ]
    }

    mysql: mysql.MySQL {
        type: "cloud"
        version: "8.0"
    }
}
