[package]
name = "example"

[dependencies]
monitoring = { oci = "oci://ghcr.io/kusionstack/monitoring", tag = "0.1.0" }
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.1.0" }

[profile]
entries = ["main.k"]

