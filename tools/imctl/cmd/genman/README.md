# OpenIM `man` Module README

Welcome to the `man` module of OpenIM, the comprehensive guide for using OpenIM's range of powerful commands. Here, you'll find in-depth details for each command, its options, and examples to help you harness the full power of the OpenIM suite.

## Overview

OpenIM is a robust instant messaging solution. To ensure users can effectively harness its capabilities, OpenIM provides a suite of commands that serve different functionalities, from the API level to RPC calls and utilities.

The `man` module ensures that users, both new and experienced, have a reliable source of information and documentation to use these commands effectively.

## Available Commands

The OpenIM commands are divided into core services and tools. Below is a brief overview of each:

### Core Services

- **openim-api**: Interface to the main functionalities of OpenIM.
- **openim-cmdutils**: Utilities for executing common tasks.
- **openim-crontask**: Schedule and manage routine tasks within OpenIM.
- **openim-msggateway**: Gateway for managing messages within the OpenIM system.
- **openim-msgtransfer**: Handle message transfers across different parts of OpenIM.
- **openim-push**: Service for pushing notifications and updates.
- **openim-rpc-auth**: RPC interface for authentication tasks.
- **openim-rpc-conversation**: RPC service for handling conversations.
- **openim-rpc-friend**: Manage friend lists and related functionalities through RPC.
- **openim-rpc-group**: Group management via RPC.
- **openim-rpc-msg**: Message handling at the RPC level.
- **openim-rpc-third**: Third-party integrations and related tasks through RPC.
- **openim-rpc-user**: User management and tasks via RPC.

### Tools

- **changelog**: Track and manage changes in OpenIM.
- **component**: Utilities related to different components within OpenIM.
- **infra**: Infrastructure and backend management tools.
- **ncpu**: Monitor and manage CPU usage and related tasks.
- **yamlfmt**: A tool for formatting and linting YAML files within the OpenIM configuration.

## How to Use

To view the manual page for any of the OpenIM commands, use the `man` command followed by the command name. For example:

```
man openim-api
```

## Contributions

We welcome contributions to enhance the `man` pages. If you discover inconsistencies, errors, or areas where further details are required, feel free to raise an issue or submit a pull request.