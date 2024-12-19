# The configuration codes in perspective of developers. 
import kam.v1.app_configuration as ac
import service
import service.container as c
import network as n
import network.ingress as ni

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
                    port: 8080
                    targetPort: 80
                    public: False
                }
            ]
           ingress: {
               defaultBackend: {
                    service: {
                        port: {
                            number: 8080
                        }
                    }
               }
               rules: [
                    {
                        http: {
                            paths: [
                                {
                                    path: "/"
                                    pathType: "Prefix"
                                    backend: {
                                        service: {
                                            port: {
                                                number: 8080
                                            }
                                        }
                                    }
                                }
                            ]
                        }
                    }
               ]
           }
        }
    }
}
