# OpenIM Scripts Directory Structure

- [OpenIM Scripts Directory Structure](#openim-scripts-directory-structure)
  - [log directory](#log-directory)
  - [Supported platforms](#supported-platforms)
  - [Get started quickly - demo.sh](#get-started-quickly---demosh)
  - [Guide: Using and Understanding OpenIM Utility Functions](#guide-using-and-understanding-openim-utility-functions)
    - [Table of Contents](#table-of-contents)
    - [1. Checking the Status of Services by Ports](#1-checking-the-status-of-services-by-ports)
      - [Function: `openim::util::check_ports`](#function-openimutilcheck_ports)
      - [Example:](#example)
    - [2. Checking the Status of Services by Process Names](#2-checking-the-status-of-services-by-process-names)
      - [Function: `openim::util::check_process_names`](#function-openimutilcheck_process_names)
      - [Example:](#example-1)
    - [3. Stopping Services by Ports](#3-stopping-services-by-ports)
      - [Function: `openim::util::stop_services_on_ports`](#function-openimutilstop_services_on_ports)
      - [Example:](#example-2)
    - [4. Stopping Services by Process Names](#4-stopping-services-by-process-names)
      - [Function: `openim::util::stop_services_with_name`](#function-openimutilstop_services_with_name)
      - [Example:](#example-3)
    - [system management and installation of openim via Linux system](#system-management-and-installation-of-openim-via-linux-system)
  - [examples](#examples)


This document outlines the directory structure for scripts in the OpenIM Server project. These scripts play a critical role in various areas like building, deploying, running and managing the services of OpenIM.

```bash
scripts/
├── README.md                           # Documentation for the scripts directory.
├── advertise.sh                        # Script to advertise or broadcast services.
├── batch_start_all.sh                  # Batch script to start all services.
├── build-all-service.sh                # Script to build all services.
├── build.cmd                           # Command script for building (usually for Windows).
├── check-all.sh                        # Check script for all components or services.
├── cherry-pick.sh                      # Helper script for git cherry-pick operations.
├── common.sh                           # Common utilities and shared functions.
├── coverage.awk                        # AWK script for processing code coverage data.
├── coverage.sh                         # Script to gather and report code coverage.
├── demo.sh                             # Demonstration or example script.
├── docker-check-service.sh             # Docker script to check services' status.
├── docker-start-all.sh                 # Docker script to start all containers/services.
├── ensure_tag.sh                       # Ensure correct tags or labeling.
├── env_check.sh                        # Environment verification and checking.
├── gen-swagger-docs.sh                 # Script to generate Swagger documentation.
├── genconfig.sh                        # Generate configuration files.
├── gendoc.sh                           # General documentation generation script.
├── githooks                            # Directory containing git hooks.
│   ├── commit-msg                      # Git hook for commit messages.
│   ├── pre-commit                      # Pre-commit git hook.
│   └── pre-push                        # Pre-push git hook.
├── golangci.yml                        # Configuration for GolangCI linting.
├── init-config.sh                      # Initialize configurations.
├── init-env.sh                         # Initialize the environment.
├── init-pwd.sh                         # Initialize or set password.
├── install                             # Installation scripts directory.
│   ├── README.md                       # Installation documentation.
│   ├── common.sh                       # Common utilities for installation.
│   ├── dependency.sh                   # Script to install dependencies.
│   ├── environment.sh                  # Set up the environment during installation.
│   ├── install-protobuf.sh             # Install Protocol Buffers.
│   ├── install.sh                      # Main installation script.
│   ├── openim-api.sh                   # Install OpenIM API.
│   ├── openim-crontask.sh              # Install OpenIM crontask.
│   ├── openim-man.sh                   # Install OpenIM management tool.
│   ├── openim-msggateway.sh            # Install OpenIM message gateway.
│   ├── openim-msgtransfer.sh           # Install OpenIM message transfer.
│   ├── openim-push.sh                  # Install OpenIM push service.
│   ├── openim-rpc.sh                   # Install OpenIM RPC.
│   ├── openim-tools.sh                 # Install OpenIM tools.
│   ├── test.sh                         # Installation testing script.
│   └── vimrc                           # Vim configuration file.
├── install-im-server.sh                # Script to install the OpenIM server.
├── install_im_compose.sh               # Install OpenIM using Docker Compose.
├── lib                                 # Library or utility scripts directory.
│   ├── chat.sh                         # Utilities related to chat.
│   ├── color.sh                        # Color-related utilities.
│   ├── golang.sh                       # Golang utilities.
│   ├── init.sh                         # Initialization utilities.
│   ├── logging.sh                      # Logging utilities.
│   ├── release.sh                      # Release related utilities.
│   ├── util.sh                         # General utility functions.
│   └── version.sh                      # Version management utilities.
├── list-feature-tests.sh               # Script to list feature tests.
├── make-rules                          # Makefile rule templates.
│   ├── common.mk                       # Common make rules.
│   ├── copyright.mk                    # Copyright related make rules.
│   ├── dependencies.mk                 # Dependency management rules.
│   ├── gen.mk                          # Generic or general rules.
│   ├── golang.mk                       # Golang-specific make rules.
│   ├── image.mk                        # Image or container-related rules.
│   ├── release.mk                      # Release specific rules.
│   ├── swagger.mk                      # Swagger documentation rules.
│   └── tools.mk                        # Tooling-related make rules.
├── mongo-init.sh                       # MongoDB initialization script.
├── release.sh                          # Script for releasing or deployment.
├── run-in-gopath.sh                    # Script to run commands within GOPATH.
├── start-all.sh                        # Script to start all services.
├── start.bat                           # Batch file to start services (usually for Windows).
├── stop-all.sh                         # Script to stop all services.
├── template                            # Directory containing template files.
│   ├── LICENSE                         # License template.
│   ├── LICENSE_TEMPLATES               # Collection of license templates.
│   ├── boilerplate.txt                 # Boilerplate template.
│   ├── footer.md.tmpl                  # Footer template for markdown.
│   ├── head.md.tmpl                    # Header template for markdown.
│   └── project_README.md               # Project README template.
├── update-generated-docs.sh            # Update generated documentation.
├── update-yamlfmt.sh                   # Update YAML formatting.
├── verify-pkg-names.sh                 # Verify package names.
├── verify-shellcheck.sh                # Shell script linting verification.
├── verify-spelling.sh                  # Spelling verification script.
├── verify-typecheck.sh                 # Type checking verification.
├── verify-yamlfmt.sh                   # Verify YAML format.
└── wait-for-it.sh                      # Script to wait for a condition or service to be ready.
```

The purpose of having a structured scripts directory like this is to make the operations of OpenIM Server clear and easy to manage. Each script has its own responsibility, making it easier to maintain and update. It's also helpful for newcomers who can easily understand what each part of the system is doing by just looking at this directory structure.

Each directory and script in the structure should be understood as a part of a larger whole. All scripts work together to ensure the smooth operation and maintenance of the OpenIM Server.


## log directory

**PATH:** `scripts/lib/logging.sh`

+ [log details](../docs/contrib/bash-log.md)

## Supported platforms

- Linux x86_64 (linux_amd64) : 64-bit Linux for most desktop and server systems.

- Windows x86_64 (windows_amd64) : 64-bit version for most Windows operating systems.

- macOS x86_64 (darwin_amd64) : 64-bit version for Apple Macintosh computers.

- Linux ARM64 (linux_arm64) : For ARM-based 64-bit Linux systems such as Raspberry Pi 4 and Jetson Nano.

- Linux s390x (linux_s390x) : 64-bit Linux for IBM System z hosts.

- Linux MIPS64 (linux_mips64) : 64-bit Linux for MIPS architecture.

- Linux MIPS64LE (linux_mips64le) : Suitable for 64-bit Linux systems with little endian MIPS architecture.

## Get started quickly - demo.sh

Is the `demo.sh` script teaching you how to quickly get started with OpenIM development and use


Steps to run demo:

```sh
$ make demo
```

More about `make` read:

+ [makefile](../docs/contrib/go-code.md)

Instructions for producing the demo movie:

```bash
# Create temporary directory
mkdir /tmp/kb-demo
cd /tmp/kb-demo

asciinema rec
<path-to-KB-repo>/scripts/demo/run.sh

<CTRL-C> to terminate the script
<CTRL-D> to terminate the asciinema recording
<CTRL-C> to save the recording locally

# Edit the recorded file by editing the controller-gen path
# Once you are happy with the recording, use svg-term program to generate the svg

svg-term --cast=<movie-id> --out _output/demo.svg --window
```

Here you will learn how to test a script, We take the four functions for starting and checking a service as an example.

## Guide: Using and Understanding OpenIM Utility Functions

This document provides an overview of the four utility functions designed for managing processes and services. These functions can check the status of services based on ports and process names, as well as stop services based on the same criteria.

### Table of Contents
- [1. Checking the Status of Services by Ports](#checking-the-status-of-services-by-ports)
- [2. Checking the Status of Services by Process Names](#checking-the-status-of-services-by-process-names)
- [3. Stopping Services by Ports](#stopping-services-by-ports)
- [4. Stopping Services by Process Names](#stopping-services-by-process-names)

### 1. Checking the Status of Services by Ports

#### Function: `openim::util::check_ports`

This function checks the status of services running on specified ports.

**Usage**:

```bash
$ openim::util::check_ports <port1> <port2> ...
```

**Design**:

- The function iterates through each provided port.
- It uses the `lsof` command to identify if there is a service running on the specified port.
- If a service is running, it logs the command, PID, and start time of the service.
- If a service is not running, it logs that the port is not started.
- If any service is not running, the function returns a status of 1.

#### Example:

```bash
$ openim::util::check_ports 8080 8081 8082
```

### 2. Checking the Status of Services by Process Names

#### Function: `openim::util::check_process_names`

This function checks the status of services based on their process names.

**Usage**:

```bash
$ openim::util::check_process_names <process_name1> <process_name2> ...
```

**Design**:

- The function uses `pgrep` to find process IDs associated with the given process names.
- If processes are found, it logs the command, PID, associated port, and start time.
- If no processes are found for a name, it logs that the process is not started.
- If any process is not running, the function returns a status of 1.

#### Example:

```bash
$ openim::util::check_process_names nginx mysql redis
```

### 3. Stopping Services by Ports

#### Function: `openim::util::stop_services_on_ports`

This function attempts to stop services running on the specified ports.

**Usage**:

```bash
$ openim::util::stop_services_on_ports <port1> <port2> ...
```

**Design**:

- The function uses the `lsof` command to identify services running on the specified ports.
- If a service is running on a port, it tries to terminate the associated process using the `kill` command.
- It logs successful terminations and any failures.
- If any service couldn't be stopped, the function returns a status of 1.

#### Example:

```bash
$ openim::util::stop_services_on_ports 8080 8081 8082
```

### 4. Stopping Services by Process Names

#### Function: `openim::util::stop_services_with_name`

This function attempts to stop services based on their process names.

**Usage**:

```bash
$ openim::util::stop_services_with_name <process_name1> <process_name2> ...
```

**Design**:

- The function uses `pgrep` to identify processes associated with the specified names.
- If processes are found, it tries to terminate them using the `kill` command.
- It logs successful terminations and any failures.
- If any service couldn't be stopped, the function returns a status of 1.

#### Example:

```bash
$ openim::util::stop_services_with_name nginx apache
```

### system management and installation of openim via Linux system

```bash
$ ./scripts/install/install.sh
```

## examples
Scripts to perform various build, install, analysis, etc operations.

The script directory design of OpenIM and the writing of scripts and tools refer to many excellent open source projects, such as helm, iam, kubernetes, docker, etc.

Maybe they'll give you inspiration for later maintenance...

These scripts keep the root level Makefile small and simple.

Examples:

* https://github.com/kubernetes/helm/tree/master/scripts
* https://github.com/cockroachdb/cockroach/tree/master/scripts
* https://github.com/hashicorp/terraform/tree/master/scripts