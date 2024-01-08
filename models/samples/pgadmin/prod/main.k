import models.schema.v1 as ac

# main.k declares customized configurations for dev stack. 

pgadmin: ac.AppConfiguration {
    workload.containers.pgadmin: {
        # dev stack has different image
        image = "dpage/pgadmin4:8.0"
    }
}