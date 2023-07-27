# OpenIM Scripts Directory Structure

This document outlines the directory structure for scripts in the OpenIM Server project. These scripts play a critical role in various areas like building, deploying, running and managing the services of OpenIM.

```bash
scripts/
├── LICENSE                        # License related files
│   ├── LICENSE                    # The license file
│   └── LICENSE_TEMPLATES          # Template for license file
├── README.md                      # Readme file for scripts directory
├── advertise.sh                   # Script for advertisement services
├── batch_start_all.sh             # Script to start all services in batch
├── build.cmd                      # Windows build command script
├── build_all_service.sh           # Script to build all services
├── build_push_k8s_images.sh       # Script to build and push images for Kubernetes
├── check_all.sh                   # Script to check status of all services
├── common.sh                      # Contains common functions used by other scripts
├── coverage.awk                   # AWK script for coverage report generation
├── coverage.sh                    # Script for generating coverage reports
├── docker_check_service.sh        # Docker specific service check script
├── docker_start_all.sh            # Script to start all services in a docker environment
├── ensure_tag.sh                  # Script to ensure proper tagging of docker images
├── enterprise                     # Scripts specific to enterprise version
│   ├── check_all.sh               # Check status of all enterprise services
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
├── install_im_server.sh           # Script to install IM server
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
├── msg_gateway_start.sh           # Script to start message gateway service
├── msg_transfer_start.sh          # Script to start message transfer service
├── path_info.sh                   # Script containing path information
├── push_start.sh                  # Script to start push service
├── release.sh                     # Script to perform release process
├── start_all.sh                   # Script to start all services
├── start_cron.sh                  # Script to start cron jobs
├── start_rpc_service.sh           # Script to start RPC service
├── stop_all.sh                    # Script to stop all services
└── style_info.sh                  # Script containing style related information
```

The purpose of having a structured scripts directory like this is to make the operations of OpenIM Server clear and easy to manage. Each script has its own responsibility, making it easier to maintain and update. It's also helpful for newcomers who can easily understand what each part of the system is doing by just looking at this directory structure.

Each directory and script in the structure should be understood as a part of a larger whole. All scripts work together to ensure the smooth operation and maintenance of the OpenIM Server.

## Supported platforms

- Linux x86_64 (linux_amd64) : 64-bit Linux for most desktop and server systems.

- Windows x86_64 (windows_amd64) : 64-bit version for most Windows operating systems.

- macOS x86_64 (darwin_amd64) : 64-bit version for Apple Macintosh computers.

- Linux ARM64 (linux_arm64) : For ARM-based 64-bit Linux systems such as Raspberry Pi 4 and Jetson Nano.

- Linux s390x (linux_s390x) : 64-bit Linux for IBM System z hosts.

- Linux MIPS64 (linux_mips64) : 64-bit Linux for MIPS architecture.

- Linux MIPS64LE (linux_mips64le) : Suitable for 64-bit Linux systems with little endian MIPS architecture.

 

## examples
Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple.

Examples:

* https://github.com/kubernetes/helm/tree/master/scripts
* https://github.com/cockroachdb/cockroach/tree/master/scripts
* https://github.com/hashicorp/terraform/tree/master/scripts