import models.schema.v1 as ac

# main.k declares customized configurations for prod stack. 

pgadmin: ac.AppConfiguration {
    workload.containers.pgadmin: {
        # prod stack has different image
        image = "dpage/pgadmin4:8.0"
    }
}