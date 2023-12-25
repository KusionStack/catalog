import models.schema.v1 as ac
import models.schema.v1.trait as t
import models.schema.v1.workload as wl
import models.schema.v1.workload.container as c
import models.schema.v1.workload.secret as sec
import models.schema.v1.workload.network as n
import models.schema.v1.monitoring as m
import models.schema.v1.accessories.database as db

# base.k declares reusable configurations for all stacks.
wordpress: ac.AppConfiguration {
    workload: wl.Service {
        containers: {
            "wordpress": c.Container {
                image: "wordpress:latest"
                env: {
                    "WORDPRESS_DB_HOST": "secret://wordpress-db/hostAddress"
                    "WORDPRESS_DB_PASSWORD": "secret://wordpress-db/password"
                }
                resources: {
                    "cpu": "1"
                    "memory": "2Gi"
                }
            }
        }
        replicas: 2
        ports: [
            n.Port {
                port: 80
                public: True
            }
        ]
    }
    database: db.Database {
        type: "alicloud"
        engine: "MySQL"
        version: "5.7"
        size: 20
        instanceType: "mysql.n2.serverless.1c"
        category: "serverless_basic"

        # SubnetID defines the virtual subnet ID associated with the VPC that the rds 
        # instance will be created in. The rds instance won't be created in user's own VPC 
        # if this field is not provided. 
        subnetID: "your_subnet_id"
    }
}
