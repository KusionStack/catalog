[package]
name = "hellokpm"
edition = "0.0.1"
version = "0.0.1"

[dependencies]
catalog = { path = "./../../../../catalog" }

[profile]
entries = [
    "base/base.k",
    "dev/main.k"
]