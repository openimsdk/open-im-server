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

### Forking and Cloning the Repository

1. Fork the OpenIM repository on GitHub to your account.
2. Clone your fork to your local machine:

   ```sh
   git clone https://github.com/<your-username>/open-im-server.git
   # The Docker bridging network mode for Mac is slightly different and needs to be set:
   export DOCKER_BRIDGE_SUBNET=127.0.0.0/16
   # Set OpenIM IP
   export OPENIM_IP=<your-ip>
   # Init Config
   make init

   # Start Components
   docker compose up -d

   # Start OpenIM Server
   make start
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