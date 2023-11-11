# OpenIM Image Management Strategy and Pulling Guide

OpenIM is an efficient, stable, and scalable instant messaging framework that provides convenient deployment methods through Docker images. OpenIM manages multiple image sources, hosted respectively on GitHub (ghcr), Alibaba Cloud, and Docker Hub. This document is aimed at detailing the image management strategy of OpenIM and providing the steps for pulling these images.


## Image Management Strategy

OpenIM's versions correspond to GitHub's tag versions. Each time we release a new version and tag it on GitHub, an automated process is triggered that pushes the new Docker image version to the following three platforms:

1. **GitHub (ghcr.io):** We use GitHub Container Registry (ghcr.io) to host OpenIM's Docker images. This allows us to better integrate with the GitHub source code repository, providing better version control and continuous integration/deployment (CI/CD) features. You can view all GitHub images [here](https://github.com/orgs/OpenIMSDK/packages).
2. **Alibaba Cloud (registry.cn-hangzhou.aliyuncs.com):** For users in Mainland China, we also host OpenIM's Docker images on Alibaba Cloud to provide faster pull speeds. You can view all Alibaba Cloud images on this [page](https://cr.console.aliyun.com/cn-hangzhou/instances/repositories) of Alibaba Cloud Image Service (note that you need to log in to your Alibaba Cloud account first).
3. **Docker Hub (docker.io):** Docker Hub is the most commonly used Docker image hosting platform, and we also host OpenIM's images there to facilitate developers worldwide. You can view all Docker Hub images on the [OpenIM's Docker Hub page](https://hub.docker.com/r/openim).

## Base images design

+ [https://github.com/openim-sigs/openim-base-image](https://github.com/openim-sigs/openim-base-image)

## OpenIM Image Design and Usage Guide

OpenIM offers a comprehensive and flexible system of Docker images, available across multiple repositories. We actively maintain these images across different platforms, namely GitHub's ghcr.io, Alibaba Cloud, and Docker Hub. However, we highly recommend ghcr.io for deployment.

### Available Versions

We provide multiple versions of our images to meet different project requirements. Here's a quick overview of what you can expect:

1. `main`: This image corresponds to the latest version of the main branch in OpenIM. It is updated frequently, making it perfect for users who want to stay at the cutting edge of our features.
2. `release-v3.*`: This is the image that corresponds to the latest version of OpenIM's stable release branch. It's ideal for users who prefer a balance between new features and stability.
3. `v3.*.*`: These images are specific to each tag in OpenIM. They are preserved in their original state and are never overwritten. These are the go-to images for users who need a specific, unchanging version of OpenIM.

### Multi-Architecture Images

In order to cater to a wider range of needs, some of our images are provided with multiple architectures under `OS / Arch`. These images offer greater compatibility across different operating systems and hardware architectures, ensuring that OpenIM can be deployed virtually anywhere.

**Example:**

+ [https://github.com/OpenIMSDK/chat/pkgs/container/openim-chat/113925695?tag=v1.1.0](https://github.com/OpenIMSDK/chat/pkgs/container/openim-chat/113925695?tag=v1.1.0)


## Methods and Steps for Pulling Images

When pulling OpenIM's Docker images, you can choose the most suitable source based on your geographic location and network conditions. Here are the steps to pull OpenIM images from each source:

### Select image

1. Choose the image repository platform you prefer. As previously mentioned, we recommend [OpenIM ghcr.io](https://github.com/orgs/OpenIMSDK/packages).

2. Choose the image name and image version that suits your needs. Refer to the description above for more details.


### Install image

1. First, make sure Docker is installed on your machine. If not, you can refer to the [Docker official documentation](https://docs.docker.com/get-docker/) for installation.

2. Open the terminal and run the following commands to pull the images:

   For OpenIM Server:

   - Pull from GitHub:

     ```bash
     docker pull ghcr.io/openimsdk/openim-server:latest
     ```

   - Pull from Alibaba Cloud:

     ```bash
     docker pull registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-server:latest
     ```

   - Pull from Docker Hub:

     ```bash
     docker pull docker.io/openim/openim-server:latest
     ```

   For OpenIM Chat:

   - Pull from GitHub:

     ```bash
     docker pull ghcr.io/openimsdk/openim-chat:latest
     ```

   - Pull from Alibaba Cloud:

     ```bash
     docker pull registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-chat:latest
     ```

   - Pull from Docker Hub:

     ```bash
     docker pull docker.io/openim/openim-chat:latest
     ```

3. Run the `docker images` command to confirm that the image has been successfully pulled.

This concludes OpenIM's image management strategy and the steps for pulling images. If you have any questions, please feel free to ask.
