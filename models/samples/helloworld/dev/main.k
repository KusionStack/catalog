import models.schema.v1 as ac

# main.k declares customized configurations for dev stack.

helloworld: ac.AppConfiguration {
    workload.containers.nginx: {
        # dev stack has different image
        image = "nginx:v2"
    }
}
