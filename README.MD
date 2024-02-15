# TR-Webdav

TR-Webdav is a very simplistic webdav server written in go, tailered to serve as a webdav server
for [Procurat](https://procurat.de) servers.

### Configuration

In order to run the server, you need the following configuration files:

#### `conf.json`

```
{
"dataPath": "/path/to/data", // The path to the data directory
"adminUserName": "admin", // The username for the admin user (can access all files)
"credentialsFilePath": "/path/to/user.passwd", // The path to the credentials file (explained below)
"permissionsFilePath": "/path/to/permissions.json" // The path to the permissions file (explained below)
"realm": "WebDAV" // The realm for the webdav server
"usersRoot": "users" // The name after the user directories come (for example /users/my-user) - it is only accessible by the user itself and the admin user
"certPath": "/path/to/cert.pem" // The path to the certificate file (for HTTPS)
"keyPath": "/path/to/key.pem" // The path to the key file (for HTTPS)
}
```

#### `user.passwd`

This file contains the credentials for the users, similar to a user.passwd file for HTTP Digest auth. It is a simple
file with the following format: `username:realm:passwordHash`

##### Example

`my-user:WebDAV:cba5dd28add5f5a89bd7b890a8e9db79`

`permissions.json`

```
{
  "restrictedDirectories": [
    "admin"
  ],
  "userPermissions": [
    {
      "username": "test",
      "allowed": [
        "admin",
        "datev"
      ],
      "restricted": [
        "finanzen",
        "datev"
      ]
    }
  ]
}
```

This file contains an array of restricted directories (which need explicit permission to access) and an array of
user-specific permissions that may override them.
Here you need to specify which permissions on globally restricted directories a user has and which directories are
restricted for the user. If a directory is both allowed and restricted, the restricted permission takes precedence.
This is useful if you have a non-restricted folder (for example "images") but want to restrict a certain user, for
example bob from accessing it. You just need to add "images" to the restricted array for bob and he won't be able to
access it.

### Usage

You can either run the service as a precompiled binary, use the Docker image or compile it yourself.

#### Compiling

To compile the server, you need to have go installed. Then, you can simply run `go build` in the root directory of the
repository.

#### Docker

You can either pull the docker image with `docker pull hub.docker.com/triargos/tr-webdav` or use a docker-compose file
like this:

```yaml
services:
  webdav:
    image: triargos/tr-webdav
    volumes:
      - /path/to/conf.json:/app/conf.json
      - /path/to/user.passwd:/app/user.passwd
      - /path/to/access-control.json:/app/access-control.json
      - /path/to/data:/app/data
      - /path/to/cert.pem:/app/cert.pem
      - /path/to/key.pem:/app/key.pem
    ports:
      - 8080:8080
      - 8443:8443
    restart: always

```

#### Precompiled binary (coming soon)

Get the latest release from the releases page and run it with `./tr-webdav` (or `tr-webdav.exe` on Windows).

### Contributing

Any contributions are welcome.