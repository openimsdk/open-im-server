# OpenIM System: Setup and Usage Guide

<!-- vscode-markdown-toc -->
* 1. [1. Introduction](#Introduction)
* 2. [2. Prerequisites (Requires root permissions)](#PrerequisitesRequiresrootpermissions)
* 3. [3. Create `openim-api` systemd unit template file](#Createopenim-apisystemdunittemplatefile)
* 4. [4. Copy systemd unit template file to systemd config directory (Requires root permissions)](#CopysystemdunittemplatefiletosystemdconfigdirectoryRequiresrootpermissions)
* 5. [5. Start systemd service](#Startsystemdservice)


##  0. <a name='Introduction'></a>0. Introduction

Systemd is the default service management form for the latest Linux distributions, replacing the original init.

The OpenIM system is a comprehensive suite of services tailored to address a wide variety of messaging needs. This guide will walk you through the steps of setting up the OpenIM system services and provide insights into its usage.

**Prerequisites:**

+ A Linux server with necessary privileges.
+ Ensure you have `systemctl` installed and running.


##  1. <a name='Deployment'></a>1. Deployment

1. **Retrieve the Installation Script**:

   Begin by obtaining the OpenIM installation script which will be utilized to deploy the entire OpenIM system.

2. **Install OpenIM**:

   To install all the components of OpenIM, run:

   ```bash
   ./scripts/install/install.sh -i  
   ```

   or

   ```bash
   ./scripts/install/install.sh --install  
   ```

   This will initiate the installation process for all OpenIM components.

3. **Check the Status**:

   Post installation, it is good practice to verify if all the services are running as expected:

   ```bash
   systemctl status openim.target
   ```

   This will list the status of all related services of OpenIM.

**Maintenance & Management:**

1. **Checking Individual Service Status**:

   You can monitor the status of individual services with the following command:

   ```bash
   systemctl status <service-name>
   ```

   For instance:

   ```bash
   systemctl status openim-api.service
   ``

2. **Starting and Stopping Services**:

   If you wish to start or stop any specific service, you can do so with `systemctl start` or `systemctl stop` followed by the service name:

   ```bash
   systemctl start openim-api.service
   systemctl stop openim-api.service
   ```

3. **Uninstalling OpenIM**:

   In case you wish to remove the OpenIM components from your server, utilize:

   ```bash
   ./scripts/install/install.sh -u
   ```

   or

   ```bash
   ./scripts/install/install.sh --uninstall
   ```

   Ensure you take a backup of any important data before executing the uninstall command.

4. **Logs & Troubleshooting**:

   Logs play a pivotal role in understanding the system's operation and troubleshooting any issues. OpenIM logs can typically be found in the directory specified during installation, usually `${OPENIM_LOG_DIR}`.

   Always refer to the logs when troubleshooting. Look for any error messages or warnings that might give insights into the issue at hand.


**Note:**

+ `openim-api.service`: Manages the main API gateways for OpenIM communication.
+ `openim-crontask.service`: Manages scheduled tasks and jobs.
+ `openim-msggateway.service`: Takes care of message gateway operations.
+ `openim-msgtransfer.service`: Handles message transfer functionalities.
+ `openim-push.service`: Responsible for push notification services.
+ `openim-rpc-auth.service`: Manages RPC (Remote Procedure Call) for authentication.
+ `openim-rpc-conversation.service`: Manages RPC for conversations.
+ `openim-rpc-friend.service`: Handles RPC for friend-related operations.
+ `openim-rpc-group.service`: Manages group-related RPC operations.
+ `openim-rpc-msg.service`: Takes care of message RPCs.
+ `openim-rpc-third.service`: Deals with third-party integrations using RPC.
+ `openim-rpc-user.service`: Manages user-related RPC operations.
+ `openim.target`: A target that bundles all the above services for collective operations.


**Viewing Logs with `journalctl`:**

`systemctl` services usually log their output to the systemd journal, which you can access using the `journalctl` command.

1. **View Logs for a Specific Service**:

   To view the logs for a particular service, you can use:

   ```bash
   journalctl -u <service-name>
   ```

   For example, to see the logs for the `openim-api.service`, you would use:

   ```bash
   journalctl -u openim-api.service
   ```

2. **Filtering Logs**:

   + By Time

     : If you wish to see logs since a specific time:

     ```bash
     journalctl -u openim-api.service --since "2023-10-28 12:00:00"
     ```

   + Most Recent Logs

     : To view the most recent logs, you can combine 
`tail` functionality with `journalctl`:

     ```bash
     journalctl -u openim-api.service -n 100
     ```

3. **Continuous Monitoring of Logs**:

   To see new log messages in real-time, you can use the `-f` flag, which mimics the behavior of `tail -f`:

   ```bash
   journalctl -u openim-api.service -f
   ```

### Continued Maintenance:

1. **Regularly Check Service Status**:

   It's good practice to routinely verify that all services are active and running. This can be done with:

   ```bash
   systemctl status openim-api.service openim-push.service openim-rpc-group.service openim-crontask.service openim-rpc-auth.service openim-rpc-msg.service openim-msggateway.service openim-rpc-conversation.service openim-rpc-third.service openim-msgtransfer.service openim-rpc-friend.service openim-rpc-user.service
   ```

2. **Update Services**:

   Periodically, there might be updates or patches to the OpenIM system or its components. Make sure you keep the system updated. After updating any service, always reload the daemon and restart the service:

   ```bash
   systemctl daemon-reload
   systemctl restart openim-api.service
   ```

3. **Backup Important Data**:

   Regularly backup any configuration files, user data, and other essential data. This ensures that you can restore the system to a working state in case of failures.

### Important `systemctl` and Logging Commands to Learn:

1. **Start/Stop/Restart Services**:

   ```bash
   systemctl start <service-name>
   systemctl stop <service-name>
   systemctl restart <service-name>
   ```

2. **Enable/Disable Services**:

   If you want a service to start automatically at boot:

   ```bash
   systemctl enable <service-name>
   ```

   To prevent it from starting at boot:

   ```bash
   systemctl disable <service-name>
   ```

3. **Check Failed Services**:

   To quickly check if any service has failed:

   ```bash
   systemctl --failed
   ```

4. **Log Rotation**:

   `journalctl` logs can grow large. To clear all archived journal entries, use:

   ```bash
   journalctl --vacuum-time=1d
   ```


**Advanced requirements:**

- Convenient service runtime log recording for problem analysis
- Service management logs
- Option to restart upon abnormal exit

The daemon does not meet these advanced requirements.

`nohup` only logs the service's runtime outputs and errors.

Only systemd can fulfill all of the above requirements.

> The default logs are enhanced with timestamps, usernames, service names, PIDs, etc., making them user-friendly. You can view logs of abnormal service exits. Advanced customization is possible through the configuration files in `/lib/systemd/system/`.

In short, systemd is the current mainstream way to manage backend services on Linux, so I've abandoned `nohup` in my new versions of bash scripts, opting instead for systemd.

##  2. <a name='PrerequisitesRequiresrootpermissions'></a>Prerequisites (Requires root permissions)

1. Configure `environment.sh` based on the comments.
2. Create a data directory:

```bash
mkdir -p ${OPENIM_DATA_DIR}/{openim-api,openim-crontask}
```

3. Create a bin directory and copy `openim-api` and `openim-crontask` executable files:

```bash
source ./environment.sh
mkdir -p ${OPENIM_INSTALL_DIR}/bin
cp openim-api openim-crontask ${OPENIM_INSTALL_DIR}/bin
```

4. Copy the configuration files of `openim-api` and `openim-crontask` to the `${OPENIM_CONFIG_DIR}` directory:

```bash
mkdir -p ${OPENIM_CONFIG_DIR}
cp openim-api.yaml openim-crontask.yaml ${OPENIM_CONFIG_DIR}
```

##  3. <a name='Createopenim-apisystemdunittemplatefile'></a> Create `openim-api` systemd unit template file

For each OpenIM service, we will create a systemd unit template. Follow the steps below for each service:

Run the following shell script to generate the `openim-api.service.template`:

```bash
source ./environment.sh
cat > openim-api.service.template <<EOF
[Unit]
Description=OpenIM Server API
Documentation=https://github.com/oepnimsdk/open-im-server/blob/master/init/README.md

[Service]
WorkingDirectory=${OPENIM_DATA_DIR}/openim-api
ExecStart=${OPENIM_INSTALL_DIR}/bin/openim-api --config=${OPENIM_CONFIG_DIR}/openim-api.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
```

Following the above style, create the respective template files or generate them in bulk:

First, make sure you've sourced the environment variables:

```bash
source ./environment.sh
```

Use the shell script to generate the systemd unit template for each service:

```bash
declare -a services=(
"openim-api"
... [other services]
)

for service in "${services[@]}"
do
   cat > $service.service.template <<EOF
[Unit]
Description=OpenIM Server - $service
Documentation=https://github.com/oepnimsdk/open-im-server/blob/master/init/README.md

[Service]
WorkingDirectory=${OPENIM_DATA_DIR}/$service
ExecStart=${OPENIM_INSTALL_DIR}/bin/$service --config=${OPENIM_CONFIG_DIR}/$service.yaml
Restart=always
RestartSec=5
StartLimitInterval=0

[Install]
WantedBy=multi-user.target
EOF
done
```

##  4. <a name='CopysystemdunittemplatefiletosystemdconfigdirectoryRequiresrootpermissions'></a>Copy systemd unit template file to systemd config directory (Requires root permissions)

Ensure you have root permissions to perform this operation:

```bash
for service in "${services[@]}"
do
   sudo cp $service.service.template /etc/systemd/system/$service.service
done
...
```

##  5. <a name='Startsystemdservice'></a>Start systemd service

To start the OpenIM services:

```bash
for service in "${services[@]}"
do
   sudo systemctl daemon-reload 
   sudo systemctl enable $service 
   sudo systemctl restart $service
done
```
