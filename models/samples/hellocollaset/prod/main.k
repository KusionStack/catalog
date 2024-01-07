import models.schema.v1 as ac

# main.k declares customized configurations for prod stack.

hellocollaset: ac.AppConfiguration {
    workload.containers.nginx: {
        # prod stack has different image
        image = "nginx:v2"
    }
}
