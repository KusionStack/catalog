import models.schema.v1 as ac
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.network as n
import models.schema.v1.accessories.postgres
import models.schema.v1.workload.secret as sec

# base.k declares reusable configurations for all stacks.
pgadmin: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            pgadmin: c.Container {
                image: "dpage/pgadmin4:latest"
                env: {
                    "PGADMIN_DEFAULT_EMAIL": "admin@admin.com"
                    "PGADMIN_DEFAULT_PASSWORD": "secret://pgadmin-secret/pgadmin-default-password"
                    "PGADMIN_PORT": "80"
                }
                resources = {
                    "cpu": "500m"
                    "memory": "512Mi"
                }
            }
        }
        secrets: {
            "pgadmin-secret": sec.Secret {
                type: "opaque"
                data: {
                    "pgadmin-default-password": "*******"
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
    database: {
        pgadmin: postgres.PostgreSQL {
            type: "cloud"
            version: "14.0"
        }
    }
}
