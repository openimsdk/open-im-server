# Building OpenIM

Building OpenIM is easy if you take advantage of the containerized build environment. This document will help guide you through understanding this build process.

## Requirements

1. Docker, using one of the following configurations:
  * **macOS** Install Docker for Mac. See installation instructions [here](https://docs.docker.com/docker-for-mac/).
     **Note**: You will want to set the Docker VM to have at least 4GB of initial memory or building will likely fail.
  * **Linux with local Docker**  Install Docker according to the [instructions](https://docs.docker.com/installation/#installation) for your OS.
  * **Windows with Docker Desktop WSL2 backend**  Install Docker according to the [instructions](https://docs.docker.com/docker-for-windows/wsl-tech-preview/). Be sure to store your sources in the local Linux file system, not the Windows remote mount at `/mnt/c`.
  
  **Note**: You will need to check if Docker CLI plugin buildx is properly installed (`docker-buildx` file should be present in `~/.docker/cli-plugins`). You can install buildx according to the [instructions](https://github.com/docker/buildx/blob/master/README.md#installing).

2. **Optional** [Google Cloud SDK](https://developers.google.com/cloud/sdk/)

You must install and configure Google Cloud SDK if you want to upload your release to Google Cloud Storage and may safely omit this otherwise.

## Actions

About [Images packages](https://github.com/orgs/OpenIMSDK/packages?repo_name=Open-IM-Server)

All files in the `build/images` directory are not templated and are instead rendered by Github Actions, which is an automated process.

Trigger condition:
1. create a new tag with the format `vX.Y.Z` (e.g. `v1.0.0`)
2. push the tag to the remote repository
3. wait for the build to finish
4. download the artifacts from the release page

## Make images

**help info:**

```bash
$ make image.help
```

**build images:**

```bash
$ make image
```

## Overview

While it is possible to build OpenIM using a local golang installation, we have a build process that runs in a Docker container.  This simplifies initial set up and provides for a very consistent build and test environment.


## Basic Flow

The scripts directly under [`build/`](.) are used to build and test.  They will ensure that the `openim-build` Docker image is built (based on [`build/build-image/Dockerfile`](../Dockerfile) and after base image's `OPENIM_BUILD_IMAGE_CROSS_TAG` from Dockerfile is replaced with one of those actual tags of the base image, like `v1.13.9-2`) and then execute the appropriate command in that container.  These scripts will both ensure that the right data is cached from run to run for incremental builds and will copy the results back out of the container. You can specify a different registry/name and version for `openim-cross` by setting `OPENIM_CROSS_IMAGE` and `OPENIM_CROSS_VERSION`, see [`common.sh`](common.sh) for more details.

The `openim-build` container image is built by first creating a "context" directory in `_output/images/build-image`.  It is done there instead of at the root of the OpenIM repo to minimize the amount of data we need to package up when building the image.

There are 3 different containers instances that are run from this image.  The first is a "data" container to store all data that needs to persist across to support incremental builds. Next there is an "rsync" container that is used to transfer data in and out to the data container.  Lastly there is a "build" container that is used for actually doing build actions.  The data container persists across runs while the rsync and build containers are deleted after each use.

`rsync` is used transparently behind the scenes to efficiently move data in and out of the container.  This will use an ephemeral port picked by Docker.  You can modify this by setting the `OPENIM_RSYNC_PORT` env variable.

All Docker names are suffixed with a hash derived from the file path (to allow concurrent usage on things like CI machines) and a version number.  When the version number changes all state is cleared and clean build is started.  This allows the build infrastructure to be changed and signal to CI systems that old artifacts need to be deleted.

## Build artifacts
The build system output all its products to a top level directory in the source repository named `_output`.
These include the binary compiled packages (e.g. imctl, openim-api etc.) and archived Docker images.
If you intend to run a component with a docker image you will need to import it from this directory with
