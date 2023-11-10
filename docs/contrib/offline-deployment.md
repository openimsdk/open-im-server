

# OpenIM Offline Deployment Design

## 1. Base Images

Below are the base images and their versions you'll need:

- wurstmeister/kafka
- redis:7.0.0
- mongo:6.0.2
- mysql:5.7
- wurstmeister/zookeeper
- minio/minio

Use the following commands to pull these base images:

```
docker pull wurstmeister/kafka
docker pull redis:7.0.0
docker pull mongo:6.0.2
docker pull mysql:5.7
docker pull wurstmeister/zookeeper
docker pull minio/minio
```

## 2. OpenIM & Chat Images

**For detailed understanding of version management and storage of OpenIM and Chat**: [version.md](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/version.md)

### OpenIM Image

- Get image version info: [images.md](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md)
- Depending on the required version, execute the following command:

```bash
docker pull ghcr.io/openimsdk/openim-server:<version-name>
```

### Chat Image

- Execute the following command to pull the image:

```bash
docker pull ghcr.io/openimsdk/openim-server:<version-name>
```

## 3. Image Storage Selection

**Repositories**:

- Alibaba Cloud: `registry.cn-hangzhou.aliyuncs.com/openimsdk/openim-server`
- Docker Hub: `openim/openim-server`

**Version Selection**:

- Stable: e.g. release-v3.2 (or 3.1, 3.3)
- Latest: latest
- Latest of main: main

## 4. Version Selection

You can select from the following versions:

- Stable: e.g. release-v3.2
- Latest: latest
- Latest from main branch: main

## 5. Offline Deployment Steps

1. **Pull images**: Execute the above `docker pull` commands to pull all required images locally.
2. **Save images**:

```
docker save -o <tar-file-name>.tar <image-name>
```

1. **Fetch code**: Clone the repository:

```
git clone https://github.com/OpenIMSDK/openim-docker.git
```

Or download the code from [Releases](https://github.com/OpenIMSDK/openim-docker/releases/).

1. **Transfer files**: Use `scp` to transfer all images and code to the intranet server.

```
scp <tar-file-name>.tar user@remote-ip:/path/on/remote/server
```

Or choose other transfer methods such as a hard drive.

1. **Import images**: On the intranet server:

```
docker load -i <tar-file-name>.tar
```

1. **Deploy**: Navigate to the `openim-docker` repository directory and follow the README guide for deployment.
2. **Deploy using Docker-compose**:

```
docker-compose up -d

# Verify
docker-compose ps
```

> **Note**: If you're using a version of Docker prior to 20, make sure you've installed `docker-compose`.

## 6. Reference Links

- [OpenIMSDK Issue #432](https://github.com/openimsdk/open-im-server/issues/432)
- [Notion Link](https://nsddd.notion.site/435ee747c0bc44048da9300a2d745ad3?pvs=25)
- [OpenIMSDK Issue #474](https://github.com/openimsdk/open-im-server/issues/474)