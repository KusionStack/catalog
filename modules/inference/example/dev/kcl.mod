[package]
name = "example"

[dependencies]
inference = { oci = "oci://ghcr.io/kusionstack/inference", tag = "0.1.0-beta.5" }
service = {oci = "oci://ghcr.io/kusionstack/service", tag = "0.1.0" }
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.2.0" }
network = { oci = "oci://ghcr.io/kusionstack/network", tag = "0.2.0" }

[profile]
entries = ["main.k"]
