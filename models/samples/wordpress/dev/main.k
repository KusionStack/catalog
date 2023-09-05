import models.schema.v1 as ac

# main.k declares customized configurations for dev stack.

wordpress: ac.AppConfiguration {
    workload.containers.wordpress: {
        # dev stack has different image
        image = "wordpress:4.8-apache"
    }
}
