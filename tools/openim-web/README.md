# OpenIM Web Service

- [OpenIM Web Service](#openim-web-service)
  - [Overview](#overview)
  - [User](#user)
  - [Docker Deployment](#docker-deployment)
    - [Build the Docker Image](#build-the-docker-image)
    - [Run the Docker Container](#run-the-docker-container)
  - [Configuration](#configuration)
  - [Contributions](#contributions)


OpenIM Web Service is a lightweight containerized service built with Go. The service serves static files and allows customization via environment variables.

## Overview

- Built using Go.
- Deployed as a Docker container.
- Serves static files from a directory which can be set via an environment variable.
- The default port for the service is `20001`, but it can be customized using an environment variable.

## User

example：

```bash
$ ./openim-web -h
Usage of ./openim-web:
  -distPath string
        Path to the distribution (default "/app/dist")
  -port string
        Port to run the server on (default "20001")
```

Variables can be set as above, Environment variables can also be set

example:

```bash
$ export OPENIM_WEB_PPRT="11001"
```

Initialize the env configuration file:

```bash
$ make init
```

## Docker Deployment

### Build the Docker Image

Even though we've implemented automation, it's to make the developer experience easier:

To build the Docker image for OpenIM Web Service:

```bash
$ docker build -t openim-web .
```

### Run the Docker Container

To run the service:

```bash
$ docker run -e DIST_PATH=/app/dist -e PORT=20001 -p 20001:20001 openim-web
```

## Configuration

You can configure the OpenIM Web Service using the following environment variables:

- **DIST_PATH**: The path to the directory containing the static files. Default: `/app/dist`.
- **PORT**: The port on which the service listens. Default: `11001`.

## Contributions

We welcome contributions from the community. If you find any bugs or have feature suggestions, please create an issue or send a pull request.