[package]
name = "dev"
edition = "*"
version = "0.0.1"

[dependencies]
catalog = { path = "../../../../../catalog" }
[profile]
entries = ["../base/base.k", "main.k"]
