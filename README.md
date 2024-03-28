# Kusion Catalog

Kusion is an Intent-Driven Platform Orchestrator that enables developers to specify their desired intent in a declarative way and then using a consistent workflow to drive continuous delivery through the entire application lifecycle.

To achieve that, we've introduced the concept of [Kusion Modules](https://www.kusionstack.io/docs/kusion/concepts/kusion-module) for users to prescribe their intent in a structured way. Kusion Modules are modular building blocks that represent common and re-usable capabilities required during an application delivery.

One of the core goals of Kusion is to build an open, inclusive and prosper open source community focused on solving real-world application delivery and management problems, in the meantime sharing the reusable building blocks and best practices.

This repository contains the source code for all the official `Kusion Modules` maintained by the KusionStack team (all contributions welcome), representing our understanding of a "golden path" and is designed to be used out-of-the-box.

## Catalog Structure

The `modules` directory contains all the out-of-the-box Kusion Module definitions, with the following directory structure.

```
├── modules
│   ├── monitoring          👈 Module definition for Promethues
│   │   ├── example         👈 Example for using the Promethues module
│   │   ├── kcl.mod         👈 kcl.mod includes the KCL package metadata
│   │   ├── prometheus.k    👈 Schema definition for Promethues configuration
│   │   └── src             👈 Generator implementation for Promethues module in Go
│   ├── mysql               👈 Module definition for Mysql database
│   │   ├── ...
│   ├── network             👈 Module definition for Network
│   │   └── ...
│   ├── opsrule             👈 Module definition for Operational Rule
│   │   └── ...
│   └── postgres            👈 Module definition for Postgres database
│       └── ...
```

## Using the Catalog Modules

The modules defined in the `catalog` repository are published to the [KusionStack GitHub container registry](https://github.com/orgs/KusionStack/packages).

To reference and import the official Kusion Modules defined in this catalog repository, you can declare the dependencies in the corresponding `kcl.mod` file (Pick and choose the ones you need):

```
[package]
name = "my-project"
edition = "0.5.0"
version = "0.1.0"

[dependencies]
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.1.0" }
monitoring = { oci = "oci://ghcr.io/kusionstack/monitoring", tag = "0.1.0" }
mysql = { oci = "oci://ghcr.io/kusionstack/mysql", tag = "0.1.0" }
postgres = { oci = "oci://ghcr.io/kusionstack/postgres", tag = "0.1.0" }
network = { oci = "oci://ghcr.io/kusionstack/network", tag = "0.1.0" }
opsrule = { oci = "oci://ghcr.io/kusionstack/opsrule", tag = "0.1.0" }

[profile]
entries = ["../base/base.k", "main.k"]
```

The `kam` repository referenced in the `kcl.mod` contains the definition for the `AppConfiguration` schema, which is a top layer concept for describing an application and may contains a collection of modules.