# Kusion Catalog

Kusion is an Intent-Driven Platform Orchestrator that enables developers to specify their desired intent in a declarative way and then use a consistent workflow to drive continuous delivery through the entire application lifecycle.

To achieve that, we've introduced the concept of [Kusion Modules](https://www.kusionstack.io/docs/concepts/module/overview) for users to prescribe their intent in a structured way. Kusion Modules are modular building blocks that represent common and re-usable capabilities required during an application delivery.

One of the core goals of Kusion is to build an open, inclusive and prosperous open-source community focused on solving real-world application delivery and management problems, in the meantime sharing the reusable building blocks and best practices.

This repository contains the source code for all Kusion Modules that can be used publicly. If your module is open to the public, we **welcome and highly encourage** you to contribute it to this repository, so that more people can benefit from the module. Submit a pull request to this repository, once it is merged, it will be published to the [KusionStack GitHub container registry](https://github.com/orgs/KusionStack/packages).

We also provide a module [developer guide](https://www.kusionstack.io/docs/concepts/module/develop-guide) on our website, if you have any questions, please don't hesitate to contact us directly.

Some of the modules in this repository are maintained by the KusionStack team, representing our understanding of a "golden path" and are designed to be used out-of-the-box. All examples can be found in the [User Guide](https://www.kusionstack.io/docs/user-guides/working-with-k8s/deploy-application) on our website.


## Catalog Structure

The `modules` directory contains all the out-of-the-box Kusion Module definitions, with the following directory structure.

```
├── modules
│   ├── monitoring          👈 Module for Promethues
│   │   ├── example         👈 Example for using the Promethues module
│   │   ├── kcl.mod         👈 kcl.mod includes the KCL package metadata
│   │   ├── prometheus.k    👈 Schema definition for Promethues configuration
│   │   └── src             👈 gRPC interfaces implementation for Promethues module in Go
│   ├── mysql               👈 Module for Mysql database
│   │   ├── ...
│   ├── network             👈 Module for Network
│   │   └── ...
│   ├── opsrule             👈 Module for Operational Rule
│   │   └── ...
│   └── postgres            👈 Module for Postgres database
│       └── ...
```

## Using the Catalog Modules

The modules defined in the `catalog` repository are published to the [KusionStack GitHub container registry](https://github.com/orgs/KusionStack/packages).

### Platform Engineers

1. Please visit [module references](https://www.kusionstack.io/docs/reference/modules/) on the website or example/readme.md in each module directory to understand the capabilities and usage of each module.
2. Register this module in your workspace and set default values to standardize the module's behavior

Please visit the [platform engineer development guide](https://www.kusionstack.io/docs/concepts/module/develop-guide) for more details.

### App Developers

As an application developer, the workflow of using a Kusion module looks like this:

1. Browse available modules registered by platform engineers in the workspace
2. Add modules you need to your Stack
3. Initialize modules
4. Apply the AppConfiguration

Please visit the [application developer user guide](https://www.kusionstack.io/docs/concepts/module/app-dev-guide) for more details.