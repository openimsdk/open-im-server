# OpenIM Offline Deployment Design

## 1. Base Images

Below are the base images and their versions you'll need:

- [ ] bitnami/kafka:3.5.1
- [ ] redis:7.0.0
- [ ] mongo:6.0.2
- [ ] bitnami/zookeeper:3.8
- [ ] minio/minio:RELEASE.2024-01-11T07-46-16Z

> [!IMPORTANT]
> It is important to note that OpenIM removed mysql components from versions v3.5.0 (release-v3.5) and above, so mysql can be deployed without this requirement or above

**If you need to install more IM components or monitoring products：**

OpenIM:

> [!TIP]
> If you need to install more IM components or monitoring products [images.md](https://github.com/openimsdk/open-im-server/blob/main/docs/contrib/images.md)

- [ ] ghcr.io/openimsdk/openim-web:<version-name>
- [ ] ghcr.io/openimsdk/openim-admin:<version-name>
- [ ] ghcr.io/openimsdk/openim-chat:<version-name>
- [ ] ghcr.io/openimsdk/openim-server:<version-name>


Monitoring:

- [ ] prom/prometheus：v2.48.1
- [ ] prom/alertmanager：v0.23.0
- [ ] grafana/grafana：10.2.2
- [ ] bitnami/node-exporter：1.7.0


Use the following commands to pull these base images:

```bash
docker pull bitnami/kafka:3.5.1
docker pull redis:7.0.0
docker pull mongo:6.0.2
docker pull mariadb:10.6
docker pull bitnami/zookeeper:3.8
docker pull minio/minio:2024-01-11T07-46-16Z
```

If you need to install more IM components or monitoring products:

```bash
docker pull prom/prometheus:v2.48.1
docker pull prom/alertmanager:v0.23.0
docker pull grafana/grafana:10.2.2
docker pull bitnami/node-exporter:1.7.0
```

## 2. OpenIM Images

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
docker pull ghcr.io/openimsdk/openim-chat:<version-name>
```

### Web Image

- Execute the following command to pull the image:

```bash
docker pull ghcr.io/openimsdk/openim-web:<version-name>
```

### Admin Image

- Execute the following command to pull the image:

```bash
docker pull ghcr.io/openimsdk/openim-admin:<version-name>
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

```bash
docker save -o <tar-file-name>.tar <image-name>
```

If you want to save all the images, use the following command:

```bash
docker save -o <tar-file-name>.tar $(docker images -q)
```

3. **Fetch code**: Clone the repository:

```bash
git clone https://github.com/openimsdk/openim-docker.git
```

Or download the code from [Releases](https://github.com/openimsdk/openim-docker/releases/).

> Because of the difference between win and linux newlines, please do not clone the repository with win and then synchronize scp to linux.

4. **Transfer files**: Use `scp` to transfer all images and code to the intranet server.

```bash
scp <tar-file-name>.tar user@remote-ip:/path/on/remote/server
```

Or choose other transfer methods such as a hard drive.

5. **Import images**: On the intranet server:

```bash
docker load -i <tar-file-name>.tar
```

Import directly with shortcut commands:

```bash
for i in `ls ./`;do docker load -i $i;done
```

6. **Deploy**: Navigate to the `openim-docker` repository directory and follow the [README guide](https://github.com/openimsdk/openim-docker) for deployment.

7. **Deploy using docker compose**:

```bash
export OPENIM_IP="your ip" # Set Ip
make init # Init config
docker compose up -d # Deployment
docker compose ps # Verify
```

> **Note**: If you're using a version of Docker prior to 20, make sure you've installed `docker-compose`.

## 6. Reference Links

- [openimsdk Issue #432](https://github.com/openimsdk/open-im-server/issues/432)
- [Notion Link](https://nsddd.notion.site/435ee747c0bc44048da9300a2d745ad3?pvs=25)
- [openimsdk Issue #474](https://github.com/openimsdk/open-im-server/issues/474)
