# OpenIM Suite Scripts

The OpenIM Suite represents a comprehensive collection of scripts, each tailored to manage and operate specific services within the OpenIM ecosystem. These scripts offer consistent, reliable, and efficient tools for initializing, controlling, and managing various OpenIM services on a Linux platform.

## Features

- **Robustness:** Built with Bash's error handling mechanisms (`errexit`, `nounset`, and `pipefail`), ensuring scripts fail fast and provide relevant error messages.
- **Modularity:** Each script is dedicated to a particular service, promoting clarity and ease of maintenance.
- **Comprehensive Logging:** Integrated logging utilities offer real-time insights into operations, enhancing transparency and debuggability.
- **Systemd Integration:** Where applicable, scripts integrate with the systemd service manager, offering standardized service controls like start, stop, restart, and status checks.

## Scripts Overview

1. **openim-api:** Control interface for the OpenIM API service.
2. **openim-cmdutils:** Utility toolkit for common OpenIM command-line operations.
3. **openim-crontask:** Manages the OpenIM CronTask service, with both direct and systemctl installation methods.
4. **openim-msggateway:** Script to operate the OpenIM Message Gateway service.
5. **openim-msgtransfer:** Manages the OpenIM Message Transfer functionalities.
6. **openim-push:** Interface for controlling the OpenIM Push Notification service.
7. **openim-rpc-auth:** Script dedicated to the OpenIM RPC Authentication service.
8. **openim-rpc-conversation:** Manages operations related to the OpenIM RPC Conversation service.
9. **openim-rpc-friend:** Control interface for the OpenIM RPC Friend functionalities.
10. **openim-rpc-group:** Script for managing the OpenIM RPC Group service.
11. **openim-rpc-msg:** Operates the OpenIM RPC Messaging service.
12. **openim-rpc-third:** Script dedicated to third-party integrations with OpenIM RPC.
13. **openim-rpc-user:** Control interface for OpenIM RPC User operations.

## Usage

The scripts within the OpenIM Suite generally adhere to two primary execution methodologies. To illustrate these methodologies, we'll use `openim-crontask` as a representative example.

1. **Direct Script Execution:** Running the script directly, typically for straightforward start/stop operations.

   ```bash
   ./[script-name].sh
   ```

2. **Function-based Execution:** Invoking specific functions within the script for more specialized operations (e.g., install, uninstall).

   ```bash
   ./scripts/install/install.sh [function-name]
   ```

### 1. Direct Script Execution

This method involves invoking the script directly, initiating its default behavior. For instance, with `openim-crontask`, direct execution will start the OpenIM CronTask as a background process.

**Example:**

```bash
./openim-crontask.sh
```

Upon execution, the script will source any necessary configurations, log the start of the CronTask, and finally run the CronTask in the background. The log messages will provide feedback about the process, ensuring the user is informed of the task's progress.

### 2. Function-based Execution

This approach is more specialized, enabling users to call specific functions defined within the script. This is particularly useful for tasks like installation, uninstallation, and status checks.

For the `openim-crontask` script:

- **Installation**: It includes building the service, generating configuration files, setting up systemd services, and starting the service.

  ```bash
  ./openim-crontask.sh openim::crontask::install
  ```

- **Uninstallation**: Stops the service, removes associated binaries, configuration files, and systemd service files.

  ```bash
  ./openim-crontask.sh openim::crontask::uninstall
  ```

- **Status Check**: Verifies the running status of the service, checking for active processes and listening ports.

  ```bash
  ./openim-crontask.sh openim::crontask::status
  ```

It's crucial to familiarize oneself with the available functions within each script. This ensures optimal utilization of the provided tools and a deeper understanding of the underlying operations.



## Notes

- Always ensure you have the correct permissions before executing any script.
- Environment variables may need to be set or sourced depending on your installation and configuration.