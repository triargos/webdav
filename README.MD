# webdav-go - lightweight webdav server for simple use cases

## Introduction

_webdav-go_ is a simple WebDAV server that provides the following features:

- Simple user management based on a configuration file (yaml)
- Predefined user subdirectories
- Global admin role for users

## Table of Contents

- [Configuration](#configuration)
    * [First steps](#first-steps)
    * [Persisting data](#persisting-data)
    * [TLS](#tls)
    * [Build and run with Docker](#build-and-run-with-docker)
- [Acknowledgements](#acknowledgements)

## Configuration

The configuration is done in form of a yaml file. The program will scan for the following locations

- `/etc/webdav/config/config.yaml` when the `DOCKER_ENABLED` environment variable is set to 1
- `./config/config.yaml` when the `DOCKER_ENABLED` environment variable is not set or is 0

### First steps

If you specify no configuration file, the server will autogenerate a sample configuration file for you:

```yaml
network:
  address: 0.0.0.0
  port: "8080"
  prefix: /
content:
  dir: /var/webdav/data
users:
  admin:
    password: [ HASHED_PASSWORD ]
    root: /document/Admin
    sub_directories:
      - documents
    admin: true
```

If you add users, you can specify the password in plain text. The server will hash them on startup.

With this configuration you'll grant access for two users and the WebDAV
server is available under `http://127.0.0.1:8080/`.

### Persisting data

The server will write every `user data` (no configuration!) to the directory specified in content -> dir. You can mount
this directory to a volume to persist the data.

### TLS

This service is designed to be used behind a reverse proxy, which is responsible for the TLS.

### Build and run with Docker

The image of webdav-go is available on Docker Hub
as [`triargos/webdav-go`](https://hub.docker.com/r/triargos/webdav-go). And example docker-compose file is provided in
the repository.

## Acknowledgements

This project was partly inspired by [dave](https://github.com/micromata/dave).

Please be aware of the licenses of the components we use in this project. Everything else that has
been developed by the contributions to this project is under the [Apache 2 License](LICENSE.txt).