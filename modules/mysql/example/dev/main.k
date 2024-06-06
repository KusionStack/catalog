# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import network as n
import mysql

quickstart: ac.AppConfiguration {
    workload: service.Service {
        containers: {
            quickstart: c.Container {
                image: "kusionstack/kusion-quickstart:latest"
                env: {
                    "DB_HOST": "$(KUSION_DB_HOST_TEST_DATABASE)"
                    "DB_USERNAME": "$(KUSION_DB_USERNAME_TEST_DATABASE)"
                    "DB_PASSWORD": "$(KUSION_DB_PASSWORD_TEST_DATABASE)"
                }
            }
        }
    }
    accessories: {
        "network": n.Network {
            ports: [
                n.Port {
                    port: 8080
                }
            ]
        }
        "mysql": mysql.MySQL {
            type:   "local"
            version: "8.0"
        }
    }
}
