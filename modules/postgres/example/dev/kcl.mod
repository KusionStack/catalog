[package]
name = "example"

[dependencies]
postgres = { oci = "oci://ghcr.io/kusionstack/postgres", tag = "0.1.0" }
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.1.0" }
network = { oci = "oci://ghcr.io/kusionstack/network", tag = "0.1.0" }

[profile]
entries = ["main.k"]

