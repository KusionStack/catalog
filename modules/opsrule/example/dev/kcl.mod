[package]
name = "example"
version = "0.1.0"

[dependencies]
opsrule = { oci = "oci://ghcr.io/kusionstack/opsrule", tag = "0.2.0" }
kam = { git = "https://github.com/KusionStack/kam.git", tag = "0.2.0" }
service = {oci = "oci://ghcr.io/kusionstack/service", tag = "0.1.0" }