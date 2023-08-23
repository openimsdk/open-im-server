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
  - [examples](#examples)


This document outlines the directory structure for scripts in the OpenIM Server project. These scripts play a critical role in various areas like building, deploying, running and managing the services of OpenIM.

```bash
scripts/
├── LICENSE                        # License related files
│   ├── LICENSE                    # The license file
│   └── LICENSE_TEMPLATES          # Template for license file
├── README.md                      # Readme file for scripts directory
├── advertise.sh                   # Script for advertisement services
├── batch_start-all.sh             # Script to start all services in batch
├── build.cmd                      # Windows build command script
├── build-all-service.sh           # Script to build all services
├── build_push_k8s_images.sh       # Script to build and push images for Kubernetes
├── check-all.sh                   # Script to check status of all services
├── common.sh                      # Contains common functions used by other scripts
├── coverage.awk                   # AWK script for coverage report generation
├── coverage.sh                    # Script for generating coverage reports
├── docker-check-service.sh        # Docker specific service check script
├── docker-start-all-all.sh            # Script to start all services in a docker environment
├── ensure_tag.sh                  # Script to ensure proper tagging of docker images
├── enterprise                     # Scripts specific to enterprise version
│   ├── check-all.sh               # Check status of all enterprise services
│   ├── function.sh                # Functions specific to enterprise version
│   └── path_info.cfg              # Path information configuration for enterprise version
├── env_check.sh                   # Script to check the environment
├── function.sh                    # Contains functions used by other scripts
├── githooks                       # Git hook scripts
│   ├── commit-msg                 # Script to validate commit message
│   ├── pre-commit                 # Script to run before each commit
│   └── pre-push                   # Script to run before each push
├── init_pwd.sh                    # Script to initialize password
├── install_im_compose.sh          # Script to install IM with Docker Compose
├── install-im-server.sh           # Script to install IM server
├── lib                            # Library scripts
│   ├── color.sh                   # Script for console color manipulation
│   ├── golang.sh                  # Script for golang related utility functions
│   ├── init.sh                    # Script for initialization tasks
│   ├── logging.sh                 # Script for logging related utility functions
│   ├── release.sh                 # Script for release related utility functions
│   ├── util.sh                    # Script for generic utility functions
│   └── version.sh                 # Script for versioning related tasks
├── make-rules                     # Makefile rules
│   ├── common.mk                  # Common Make rules
│   ├── copyright.mk               # Copyright related Make rules
│   ├── dependencies.mk            # Dependencies related Make rules
│   ├── gen.mk                     # Make rules for code generation
│   ├── golang.mk                  # Golang specific Make rules
│   ├── image.mk                   # Make rules for image building
│   ├── release.mk                 # Make rules for release process
│   ├── swagger.mk                 # Make rules for swagger documentation
│   └── tools.mk                   # Make rules for tools and utilities
├── mongo-init.sh                  # Script to initialize MongoDB
├── openim-msggateway.sh           # Script to start message gateway service
├── openim-msgtransfer.sh          # Script to start message transfer service
├── path_info.sh                   # Script containing path information
├── openim-push.sh                  # Script to start push service
├── release.sh                     # Script to perform release process
├── start-all.sh                   # Script to start all services
├── openim-crontask.sh                  # Script to start cron jobs
├── openim-rpc.sh           # Script to start RPC service
├── stop-all.sh                    # Script to stop all services
└── style_info.sh                  # Script containing style related information
```

The purpose of having a structured scripts directory like this is to make the operations of OpenIM Server clear and easy to manage. Each script has its own responsibility, making it easier to maintain and update. It's also helpful for newcomers who can easily understand what each part of the system is doing by just looking at this directory structure.

Each directory and script in the structure should be understood as a part of a larger whole. All scripts work together to ensure the smooth operation and maintenance of the OpenIM Server.


## log directory

**PATH:** `scripts/lib/logging.sh`

+ [log details](../docs/conversions/bash_log.md)

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
make demo
```

More about `make` read:

+ [makefile](../docs/conversions/go_code.md)

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
openim::util::check_ports <port1> <port2> ...
```

**Design**:

- The function iterates through each provided port.
- It uses the `lsof` command to identify if there is a service running on the specified port.
- If a service is running, it logs the command, PID, and start time of the service.
- If a service is not running, it logs that the port is not started.
- If any service is not running, the function returns a status of 1.

#### Example:

```bash
openim::util::check_ports 8080 8081 8082
```

### 2. Checking the Status of Services by Process Names

#### Function: `openim::util::check_process_names`

This function checks the status of services based on their process names.

**Usage**:

```bash
openim::util::check_process_names <process_name1> <process_name2> ...
```

**Design**:

- The function uses `pgrep` to find process IDs associated with the given process names.
- If processes are found, it logs the command, PID, associated port, and start time.
- If no processes are found for a name, it logs that the process is not started.
- If any process is not running, the function returns a status of 1.

#### Example:

```bash
openim::util::check_process_names nginx mysql redis
```

### 3. Stopping Services by Ports

#### Function: `openim::util::stop_services_on_ports`

This function attempts to stop services running on the specified ports.

**Usage**:

```bash
openim::util::stop_services_on_ports <port1> <port2> ...
```

**Design**:

- The function uses the `lsof` command to identify services running on the specified ports.
- If a service is running on a port, it tries to terminate the associated process using the `kill` command.
- It logs successful terminations and any failures.
- If any service couldn't be stopped, the function returns a status of 1.

#### Example:

```bash
openim::util::stop_services_on_ports 8080 8081 8082
```

### 4. Stopping Services by Process Names

#### Function: `openim::util::stop_services_with_name`

This function attempts to stop services based on their process names.

**Usage**:

```bash
openim::util::stop_services_with_name <process_name1> <process_name2> ...
```

**Design**:

- The function uses `pgrep` to identify processes associated with the specified names.
- If processes are found, it tries to terminate them using the `kill` command.
- It logs successful terminations and any failures.
- If any service couldn't be stopped, the function returns a status of 1.

#### Example:

```bash
openim::util::stop_services_with_name nginx apache
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