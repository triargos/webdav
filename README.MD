# webdav-go - lightweight webdav server for simple use cases

## Introduction

_webdav-go_ is a simple WebDAV server that provides the following features:

- Simple user management based on a configuration file (yaml)
- Predefined user subdirectories
- Global admin role for users

## Table of Contents

- [Configuration](#configuration)
    * [First steps](#first-steps)
    * [Additional configuration](#additional-configuration)
    * [User management](#user-management)
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
    password: admin (this will be hashed on startup)
    root: /Users/admin
    sub_directories:
      - documents
    admin: true
```

You can always regenerate this by executing `webdav-go` genconfig. If you specifiy the `--reset`, flag, it will
overwrite your existing configuration file.

If you add users, you can specify the password in plain text. The server will hash them on startup.
If you want to add them programmatically, add them
via `webdav-go adduser -u [USERNAME] -p [PASSWORD] -d [ROOT_DIRECTORY] [...OPTIONAL ARGS]`. The password will be hashed
on the fly.

With this configuration you'll grant access for an `admin` user (with the password `admin`) and the WebDAV
server is available under `http://127.0.0.1:8080/`.

### Additional configuration

You can modify the configuration file generated on first startup by specifying the following environment variables:

- `WEBDAV_PORT` - the port the server will listen on
- `CREATE_ADMIN_USER` - Set this to "0" or "false" to prevent auto-creation of the admin user
- `WEBDAV_DATA_DIR` - the directory where the user data will be stored

Please note that if you specify a custom configuration file, the environment variables will be ignored.


### User management

When adding users, you can specify the following options in the configuration

- `root` - the root directory of the user (mandatory)
- `admin` - a boolean flag that specifies if the user is an admin. No permissions will be checked (optional)
- `password` - the password of the user (mandatory). This will be hashed on startup
- `jailed` - a boolean value that specifies if the user should be jailed to his root directory and subdirectories (
  optional)
- `sub_directories` - a list of subdirectories that will be created for the user (optional)

Users that are non-admin are forbidden to access other users root directories and subdirectories by default.

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