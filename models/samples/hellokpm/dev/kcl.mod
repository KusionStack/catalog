[package]
name = "dev"
edition = "*"
version = "0.0.1"

[dependencies]
base = { path = "../base" }
catalog = { path = "../../../../../catalog" }

[profile]
entries = [
    "${base:KCL_MOD}/base.k",
    "main.k"
]
