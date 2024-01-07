import models.schema.v1 as ac

# main.k declares customized configurations for prod stack.

wordpress: ac.AppConfiguration {
    workload.containers.wordpress: {
        # prod stack has different image
        image = "wordpress:4.8-apache"
    }
}
