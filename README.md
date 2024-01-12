# Kusion Catalog

Kusion is a modern application delivery and management toolchain that enables developers to specify desired intent in a declarative way and then using a consistent workflow to drive continuous deployment through application lifecycle.

One of the core goals of Kusion is to build an open, inclusive and vibrant OSS developer community focused on solving real-world application delivery and management problems, sharing the reusable building blocks and best practices.

This repository contains a catalog of community maintained, shared `Kusion Module` resources, which are designed to be usable out of box.

## Catalog Structure

The `models/schema` directory contains all KCL schema definitions for application developers, and follows the following structure.

```
    ./schema/v1                   

        /accessories            👈 schema definition for various accessory resources

        /monitoring             👈 schema definition for monitoring e.g. Promethues
        
        /trait                  👈 schema definition for various operation capabilities e.g. OpsRule
        
        /workload               👈 default workload schema definition
        
        app_configuration.k     👈 root AppConfiguration schema definition
```

Based on the schema definitions, `models/samples` directory contains plenty of useful sample code. Here is a simple explanation of those samples.

* `hellocollaset` - demonstrates how to declare a long-running service, and the workload implementation of this service is `Collaset`, which is provided by [KusionStack Operating](https://github.com/KusionStack/operating).
* `helloworld` - also declare a long-running service, with default Kubernetes Deployment workload.
* `pgadmin` - declares a cloud provider managed Postgres resource, as well as a long-running service with `dpage/pgadmin4:latest` image.
* `samplejob` - demonstrates how to declare a periodic job.
* `wordpress` - declares a cloud provider managed MySQL resource, as well as a long-running service with `wordpress:6.3` image. 