import models.schema.v1 as ac

# main.k declares customized configurations for prod stack.
prometheus: ac.AppConfiguration {
    workload.containers.nginx: {
        # prod stack has different image
        image = "quay.io/brancz/prometheus-example-app:v0.3.0"
    }
}
