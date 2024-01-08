import models.schema.v1 as ac
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.network as n
import models.schema.v1.accessories.mysql

# main.k declares customized configurations for prod stack.

wordpress: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            wordpress: c.Container {
                image = "wordpress:6.3"
                env = {
                    "WORDPRESS_DB_HOST": "$(KUSION_DB_HOST_WORDPRESS)"
                    "WORDPRESS_DB_USER": "$(KUSION_DB_USERNAME_WORDPRESS)"
                    "WORDPRESS_DB_PASSWORD": "$(KUSION_DB_PASSWORD_WORDPRESS)"
                    "WORDPRESS_DB_NAME": "mysql"
                }
                resources = {
                    "cpu": "500m"
                    "memory": "512Mi"
                }
            }
        }
        replicas = 1
        ports = [
            n.Port {
                port: 80
            }
        ]
    }
    database = {
        wordpress: mysql.MySQL {
            type: "cloud"
            version: "8.0"
        }
    }
}
