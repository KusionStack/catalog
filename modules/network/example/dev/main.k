# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import network as n

nginx: ac.AppConfiguration {
    workload: service.Service {
        containers: {
            nginx: c.Container {
                image: "nginx:1.25.2"
            }
        }
    }
    accessories: {
        "network": n.Network {
            ports: [
                n.Port {
                    port: 80
                    public: True
                }
            ]
        }
    }
}
