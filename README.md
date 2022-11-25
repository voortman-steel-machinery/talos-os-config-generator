# Talos OS config generator API

[![License](https://img.shields.io/badge/License-MIT-brightgreen?style=for-the-badge)](https://raw.githubusercontent.com/voortman/talos-os-config-generator/main/LICENSE)

Serves an Rest-API that can generate Talos OS config files.

## What is Talos OS

https://github.com/siderolabs/talos

**Talos** is a modern OS for running Kubernetes: secure, immutable, and minimal.
Talos is fully open source, production-ready, and supported by the people at [Sidero Labs](https://www.SideroLabs.com/)
All system management is done via an API - there is no shell or interactive console.
Benefits include:

- **Security**: Talos reduces your attack surface: It's minimal, hardened, and immutable.
  All API access is secured with mutual TLS (mTLS) authentication.
- **Predictability**: Talos eliminates configuration drift, reduces unknown factors by employing immutable infrastructure ideology, and delivers atomic updates.
- **Evolvability**: Talos simplifies your architecture, increases your agility, and always delivers current stable Kubernetes and Linux versions.

## Feature overview
- Create ControlPlane and Worker node config with default values with minimal input
- Create a Talosconfig to use in combination with talosctl
- Able to add a patch file to merge with the created config files

