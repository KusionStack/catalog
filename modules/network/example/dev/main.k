# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import kam.v1.workload as wl
import kam.v1.workload.container as c
import network.network as n

nginx: ac.AppConfiguration {
    workload: wl.Service {
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
