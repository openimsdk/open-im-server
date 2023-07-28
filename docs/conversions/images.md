# OpenIM Image Management Strategy and Pulling Guide

OpenIM is an efficient, stable, and scalable instant messaging framework that provides convenient deployment methods through Docker images. OpenIM manages multiple image sources, hosted respectively on GitHub (ghcr), Alibaba Cloud, and Docker Hub. This document is aimed at detailing the image management strategy of OpenIM and providing the steps for pulling these images.

## Image Management Strategy

OpenIM's versions correspond to GitHub's tag versions. Each time we release a new version and tag it on GitHub, an automated process is triggered that pushes the new Docker image version to the following three platforms:

1. **GitHub (ghcr.io):** We use GitHub Container Registry (ghcr.io) to host OpenIM's Docker images. This allows us to better integrate with the GitHub source code repository, providing better version control and continuous integration/deployment (CI/CD) features. You can view all GitHub images [here](https://github.com/orgs/OpenIMSDK/packages).
2. **Alibaba Cloud (registry.cn-hangzhou.aliyuncs.com):** For users in Mainland China, we also host OpenIM's Docker images on Alibaba Cloud to provide faster pull speeds. You can view all Alibaba Cloud images on this [page](https://cr.console.aliyun.com/cn-hangzhou/instances/repositories) of Alibaba Cloud Image Service (note that you need to log in to your Alibaba Cloud account first).
3. **Docker Hub (docker.io):** Docker Hub is the most commonly used Docker image hosting platform, and we also host OpenIM's images there to facilitate developers worldwide. You can view all Docker Hub images on the [OpenIM's Docker Hub page](https://hub.docker.com/r/openim).

## Methods and Steps for Pulling Images

When pulling OpenIM's Docker images, you can choose the most suitable source based on your geographic location and network conditions. Here are the steps to pull OpenIM images from each source:

1. First, make sure Docker is installed on your machine. If not, you can refer to the [Docker official documentation](https://docs.docker.com/get-docker/) for installation.

2. Open the terminal and run the following commands to pull the images:

   For OpenIM Server:

   - Pull from GitHub:

     ```
     bashCopy code
     docker pull ghcr.io/openimsdk/openim-server:latest
     ```

   - Pull from Alibaba Cloud:

     ```
     bashCopy code
     docker pull registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-server:latest
     ```

   - Pull from Docker Hub:

     ```
     bashCopy code
     docker pull docker.io/openim/openim-server:latest
     ```

   For OpenIM Chat:

   - Pull from GitHub:

     ```
     bashCopy code
     docker pull ghcr.io/openimsdk/openim-chat:latest
     ```

   - Pull from Alibaba Cloud:

     ```
     bashCopy code
     docker pull registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-chat:latest
     ```

   - Pull from Docker Hub:

     ```
     bashCopy code
     docker pull docker.io/openim/openim-chat:latest
     ```

3. Run the `docker images` command to confirm that the image has been successfully pulled.

This concludes OpenIM's image management strategy and the steps for pulling images. If you have any questions, please feel free to ask.