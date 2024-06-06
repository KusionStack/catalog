# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import network as n
import postgres

pgadmin: ac.AppConfiguration {
    workload: service.Service {
        containers: {
            pgadmin: c.Container {
                image: "dpage/pgadmin4:latest"
                env: {
                    "PGADMIN_DEFAULT_EMAIL": "admin@email.com"
                    "PGADMIN_DEFAULT_PASSWORD": "123456"
                    "PGADMIN_PORT": "80"
                }
            }
        }
    }
    accessories: {
        "network": n.Network {
            ports: [
                n.Port {
                    port: 80
                }
            ]
        }
        "postgres": postgres.PostgreSQL {
            type:   "local"
            version: "14.0"
        }
    }
}
