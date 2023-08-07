import models.schema.v1 as ac
import models.schema.v1.component as cmp
import models.schema.v1.component.container as c
import models.schema.v1.component.workload as wl

# base.k declares reuseable configurations for all stacks.

appConfiguration: ac.AppConfiguration {
    components: {
        "proxy": cmp.Component {
            longRunningService: wl.LongRunningService {
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
                    }  
                }
                replicas: 2  
            }
        }
        "echo-hello": {
            job: wl.Job {
                containers: {
                    "busybox": {
                        image: "busybox:1.28"
                        # Run the following command as defined
                        command: ["/bin/sh", "-c", "echo hello"]
                    }
                }
                # The scheduling strategy in Cron format.
                schedule: "0 0 13 * 5"
            }
        }
    }
}
