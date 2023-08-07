import models.schema.v1 as ac

# main.k declares customized configurations for dev stack.

appConfiguration: ac.AppConfiguration {
    components.proxy.longRunningService.containers.nginx: {
        # dev stack has different image
        image = "nginx:v2"
    }
}
