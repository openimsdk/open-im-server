# Mac Developer Deployment Guide for OpenIM

## Introduction

This guide aims to assist Mac-based developers in contributing effectively to OpenIM. It covers the setup of a development environment tailored for Mac, including the use of GitHub for development workflow and `devcontainer` for a consistent development experience.

Before contributing to OpenIM through issues and pull requests, make sure you are familiar with GitHub and the [pull request workflow](https://docs.github.com/en/get-started/quickstart/github-flow).

## Prerequisites

### System Requirements

- macOS (latest stable version recommended)
- Internet connection
- Administrator access

### Knowledge Requirements

- Basic understanding of Git and GitHub
- Familiarity with Docker and containerization
- Experience with Go programming language

## Setting up the Development Environment

### Installing Homebrew

Homebrew is an essential package manager for macOS. Install it using:

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

### Installing and Configuring Git

1. Install Git:

   ```sh
   brew install git
   ```

2. Configure Git with your user details:

   ```sh
   git config --global user.name "Your Name"
   git config --global user.email "your.email@example.com"
   ```

### Setting Up the Devcontainer

`Devcontainers` provide a Docker-based isolated development environment. 

Read [README.md](https://github.com/openimsdk/open-im-server/tree/main/.devcontainer) in the `.devcontainer` directory of the project to learn more about the devcontainer.

To set it up:

1. Install Docker Desktop for Mac from [Docker Hub](https://docs.docker.com/desktop/install/mac-install/).
2. Install Visual Studio Code and the Remote - Containers extension.
3. Open the cloned OpenIM repository in VS Code.
4. VS Code will prompt to reopen the project in a container. Accept this to set up the environment automatically.

### Installing Go and Dependencies

Use Homebrew to install Go:

```sh
brew install go
```

Ensure the version of Go is compatible with the version required by OpenIM (refer to the main documentation for version requirements).

### Additional Tools

Install other required tools like Docker, Vagrant, and necessary GNU utils as described in the main documentation.

## Mac Deployment openim-chat and openim-server

To integrate the Chinese document into an English document for Linux deployment, we will first translate the content and then adapt it to suit the Linux environment. Here's how the translated and adapted content might look:

### Ensure a Clean Environment

- It's recommended to execute in a new directory.
- Run `ps -ef | grep openim` to ensure no OpenIM processes are running.
- Run `ps -ef | grep chat` to check for absence of chat-related processes.
- Execute `docker ps` to verify there are no related containers running.

### Source Code Deployment

#### Deploying openim-server

Source code deployment is slightly more complex because Docker's networking on Linux differs from Mac.

```bash
git clone https://github.com/openimsdk/open-im-server
cd open-im-server

export OPENIM_IP="Your IP" # If it's a cloud server, setting might not be needed
make init # Generates configuration files
```

Before deploying openim-server, modify the Kafka logic in the docker-compose.yml file. Replace:

```yaml
- KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092,EXTERNAL://${DOCKER_BRIDGE_GATEWAY:-172.28.0.1}:${KAFKA_PORT:-19094}
```

With:

```yaml
- KAFKA_CFG_ADVERTISED_LISTENERS=PLAINTEXT://kafka:9092,EXTERNAL://127.0.0.1:${KAFKA_PORT:-19094}
```

Then start the service:

```bash
docker compose up -d
```

Before starting the openim-server source, set `config/config.yaml` by replacing all instances of `172.28.0.1` with `127.0.0.1`:

```bash
vim config/config.yaml -c "%s/172\.28\.0\.1/127.0.0.1/g" -c "wq"
```

Then start openim-server:

```bash
make start
```

To check the startup:

```bash
make check
```

<aside>
🚧 To avoid mishaps, it's best to wait five minutes before running `make check` again.

</aside>

#### Deploying openim-chat

There are several ways to deploy openim-chat, either by source code or using Docker.

Navigate back to the parent directory:

```bash
cd ..
```

First, let's look at deploying chat from source:

```bash
git clone https://github.com/openimsdk/chat
cd chat
make init # Generates configuration files
```

If openim-chat has not deployed MySQL, you will need to deploy it. Note that the official Docker Hub for MySQL does not support architectures like ARM, so you can use the newer version of the open-source edition:

```bash
docker run -d \
  --name mysql \
  -p 13306:3306 \
  -p 3306:33060 \
  -v "$(pwd)/components/mysql/data:/var/lib/mysql" \
  -v "/etc/localtime:/etc/localtime" \
  -e MYSQL_ROOT_PASSWORD="openIM123" \
  --restart always \
  mariadb:10.6
```

Before starting the source code of openim-chat, set `config/config.yaml` by replacing all instances of `172.28.0.1` with `127.0.0.1`:

```bash
vim config/config.yaml -c "%s/172\.28\.0\.1/127.0.0.1/g" -c "wq"
```

Then start openim-chat from source:

```bash
make start
```

To check, ensure the following four processes start successfully:

```bash
make check 
```

### Docker Deployment

Refer to https://github.com/openimsdk/openim-docker for Docker deployment instructions, which can be followed similarly on Linux.

```bash
git clone https://github.com/openimsdk/openim-docker
cd openim-docker
export OPENIM_IP="Your IP"
make init
docker compose up -d 
docker compose logs -f openim-server
docker compose logs -f openim-chat
```

## GitHub Development Workflow

### Creating a New Branch

For new features or fixes, create a new branch:

```sh
git checkout -b feat/your-feature-name
```

### Making Changes and Committing

1. Make your changes in the code.
2. Stage your changes:

   ```sh
   git add .
   ```

3. Commit with a meaningful message:

   ```sh
   git commit -m "Add a brief description of your changes"
   ```

### Pushing Changes and Creating Pull Requests

1. Push your branch to GitHub:

   ```sh
   git push origin feat/your-feature-name
   ```

2. Go to your fork on GitHub and create a pull request to the main OpenIM repository.

### Keeping Your Fork Updated

Regularly sync your fork with the main repository:

```sh
git fetch upstream
git checkout main
git rebase upstream/main
```

More read: [https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md](https://github.com/openimsdk/open-im-server/blob/main/CONTRIBUTING.md)

## Testing and Quality Assurance

Run tests as described in the OpenIM documentation to ensure your changes do not break existing functionality.

## Conclusion

This guide provides a comprehensive overview for Mac developers to set up and contribute to OpenIM. By following these steps, you can ensure a smooth and efficient development experience. Happy coding!