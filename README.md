# Dockit

[![Docker](https://img.shields.io/badge/Docker-2496ED?logo=docker&logoColor=white)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white)](https://go.dev/)

Dockit is a simple DevSecOps toolkit for Docker.

It helps developers scan Docker images, inspect image filesystems, and debug running containers.

## Planned features

- Scan Docker images for vulnerabilities
- Inspect files inside Docker images
- Generate image security reports
- Inspect running containers
- Debug container processes, networks, and mounts
- Integrate with CI/CD pipelines

## Example commands

```bash
dockit image inspect nginx:latest
dockit container inspect my-container
dockit container debug my-container
```

## Project status

Dockit is currently under development.

## Goals

Dockit aims to make Docker security and troubleshooting simple for developers and DevOps teams.
