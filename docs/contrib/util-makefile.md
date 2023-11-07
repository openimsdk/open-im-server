# Open-IM-Server Development Tools Guide

- [Open-IM-Server Development Tools Guide](#open-im-server-development-tools-guide)
  - [Introduction](#introduction)
  - [Getting Started](#getting-started)
  - [Toolset Categories](#toolset-categories)
  - [Installation Commands](#installation-commands)
    - [Basic Installation](#basic-installation)
    - [Individual Tool Installation](#individual-tool-installation)
    - [Tool Verification](#tool-verification)
  - [Detailed Tool Descriptions](#detailed-tool-descriptions)
  - [Best Practices](#best-practices)
  - [Conclusion](#conclusion)


## Introduction

Open-IM-Server boasts a robust set of tools encapsulated within its Makefile system, designed to ease development, code formatting, and tool management. This guide aims to familiarize developers with the features and usage of the Makefile toolset provided within the Open-IM-Server project.

## Getting Started

Executing `make tools` ensures verification and installation of the default tools:

- golangci-lint
- goimports
- addlicense
- deepcopy-gen
- conversion-gen
- ginkgo
- go-junit-report
- go-gitlint

The installation path is situated at `/root/workspaces/openim/Open-IM-Server/_output/tools/`.

## Toolset Categories

The Makefile logically groups tools into different categories for better management:

- **Development Tools**: `BUILD_TOOLS`
- **Code Analysis Tools**: `ANALYSIS_TOOLS`
- **Code Generation Tools**: `GENERATION_TOOLS`
- **Testing Tools**: `TEST_TOOLS`
- **Version Control Tools**: `VERSION_CONTROL_TOOLS`
- **Utility Tools**: `UTILITY_TOOLS`
- **Tencent Cloud Object Storage Tools**: `COS_TOOLS`

## Installation Commands

1. **golangci-lint**: high performance Go code linter with integration of multiple inspection tools.
2. **goimports**: Used to format Go source files and automatically add or remove imports.
3. **addlicense**: Adds a license header to the source file.
4. **deepcopy-gen and conversions-gen **: Generate deepcopy and conversion functionality.
5. **ginkgo**: Testing framework for Go.
6. **go-junit-report**: Converts Go test output to junit xml format.
7. **go-gitlint**: For checking git commit information. ... (And so on, list all the tools according to the `make tools.help` output above)...

 

### Basic Installation

- `tools.install`: Installs tools mentioned in the `BUILD_TOOLS` list.
- `tools.install-all`: Installs all tools from the `ALL_TOOLS` list.

### Individual Tool Installation

- `tools.install.%`: Installs a single tool in the `$GOBIN/` directory.
- `tools.install-all.%`: Parallelly installs an individual tool located in `./tools/*`.

### Tool Verification

- `tools.verify.%`: Checks if a specific tool is installed, and if not, installs it.

## Detailed Tool Descriptions

The following commands serve the purpose of installing particular development tools:

- `install.golangci-lint`: Installs `golangci-lint`.
- `install.addlicense`: Installs `addlicense`. ... (and so on for every tool as mentioned in the provided Makefile source)...

The commands primarily leverage Go's `install` operation, fetching and installing tools from their respective repositories. This method is especially convenient as it auto-handles dependencies and installation paths. For tools not written directly with Go (like `install.coscli`), other installation methods like wget or pip are employed.

## Best Practices

1. **Regular Updates**: To ensure tools are up-to-date, periodically run the `make tools` command.
2. **Individual Tools**: If only specific tools are required, employ the `make install.<tool-name>` command for individual installations.
3. **Verification**: Before code submissions, use the `make tools.verify.%` command to guarantee that all necessary tools are present and up-to-date.

## Conclusion

The Makefile provided by Open-IM-Server presents a centralized approach to manage and install all necessary tools during the development process. It ensures that all developers employ consistent tool versions, reducing potential issues due to version disparities. Whether you're a maintainer or a contributor to the Open-IM-Server project, understanding the workings of this Makefile will significantly enhance your developmental efficiency.
