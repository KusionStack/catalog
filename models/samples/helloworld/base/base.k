import models.schema.v1 as ac
import models.schema.v1.trait as t
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.container.probe as p
import models.schema.v1.workload.network as n

# base.k declares reuseable configurations for all stacks.
helloworld: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
                # Run the following command as defined
                command: ["/bin/sh", "-c", "echo hi"]
                # Extra arguments append to command defined above
                args: ["/bin/sh", "-c", "echo hi"]
                env: {
                    # An environment variable of name "env1" and value "VALUE" will be set
                    "env1": "VALUE"
                    # An environment variable of name "env2" and value of the key "key" in the
                    # secret named "sec-name" will be set.
                    "env2": "secret://sec-name/key"
                }
                # Run the command "/bin/sh -c echo hi", as defined above, in the directory "/tmp"
                workingDir: "/tmp"
                # Expose container port 80
                ports: [
                    n.ContainerPort {
                        port: 80
                        accessMode: "Exposed"
                        accessPort: 80
                    }
                ]
                # Configure a HTTP readiness probe
                readinessProbe: p.Probe {
                    probeHandler: p.Http {
                        url: "http://localhost:80"
                    }
                    initialDelaySeconds: 10
                }
            }
        }
        routes: {
            "nginx.example.com": n.Route {
                paths: [
                    n.RoutePath {
                        accessPort: 80
                    }
                ]
            }
        }
    }
}

helloworldcollaset: ac.AppConfiguration {
    workload: wl.Service {
        type: "CollaSet"
        containers: {
            "nginx": c.Container {
                image: "nginx:v1"
                # Run the following command as defined
                command: ["/bin/sh", "-c", "echo hi"]
                # Extra arguments append to command defined above
                args: ["/bin/sh", "-c", "echo hi"]
                # Run the command "/bin/sh -c echo hi", as defined above, in the directory "/tmp"
                workingDir: "/tmp"
            }  
        }
    }
    opsRule: t.OpsRule {
        maxUnavailable: "30%"
    }
}

helloworldjob: ac.AppConfiguration {
    workload: wl.Job {
        containers: {
            "busybox": c.Container {
                image: "busybox:1.28"
                # Run the following command as defined
                command: ["/bin/sh", "-c", "echo hello"]
            }
        }
        # Run every hour.
        schedule: "0 * * * *"
    }
}