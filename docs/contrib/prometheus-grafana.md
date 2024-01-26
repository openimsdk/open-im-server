# Deployment and Design of OpenIM's Management Backend and Monitoring

<!-- vscode-markdown-toc -->
* 1. [Source Code & Docker](#SourceCodeDocker)
	* 1.1. [Deployment](#Deployment)
	* 1.2. [Configuration](#Configuration)
	* 1.3. [Monitoring Running in Docker Guide](#MonitoringRunninginDockerGuide)
		* 1.3.1. [Introduction](#Introduction)
		* 1.3.2. [Prerequisites](#Prerequisites)
		* 1.3.3. [Step 1: Clone the Repository](#Step1:ClonetheRepository)
		* 1.3.4. [Step 2: Start Docker Compose](#Step2:StartDockerCompose)
		* 1.3.5. [Step 3: Use the OpenIM Web Interface](#Step3:UsetheOpenIMWebInterface)
		* 1.3.6. [Running Effect](#RunningEffect)
		* 1.3.7. [Step 4: Access the Admin Panel](#Step4:AccesstheAdminPanel)
		* 1.3.8. [Step 5: Access the Monitoring Interface](#Step5:AccesstheMonitoringInterface)
		* 1.3.9. [Next Steps](#NextSteps)
		* 1.3.10. [Troubleshooting](#Troubleshooting)
* 2. [Kubernetes](#Kubernetes)
	* 2.1. [Middleware Monitoring](#MiddlewareMonitoring)
	* 2.2. [Custom OpenIM Metrics](#CustomOpenIMMetrics)
	* 2.3. [Node Exporter](#NodeExporter)
* 3. [Setting Up and Configuring AlertManager Using Environment Variables and `make init`](#SettingUpandConfiguringAlertManagerUsingEnvironmentVariablesandmakeinit)
	* 3.1. [Introduction](#Introduction-1)
	* 3.2. [Prerequisites](#Prerequisites-1)
	* 3.3. [Configuration Steps](#ConfigurationSteps)
		* 3.3.1. [Exporting Environment Variables](#ExportingEnvironmentVariables)
		* 3.3.2. [Initializing AlertManager](#InitializingAlertManager)
		* 3.3.3. [Key Configuration Fields](#KeyConfigurationFields)
		* 3.3.4. [Configuring SMTP Authentication Password](#ConfiguringSMTPAuthenticationPassword)
		* 3.3.5. [Useful Links for Common Email Servers](#UsefulLinksforCommonEmailServers)
	* 3.4. [Conclusion](#Conclusion)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

OpenIM offers various flexible deployment options to suit different environments and requirements. Here is a simplified and optimized description of these deployment options:

1. Source Code Deployment:
   + **Regular Source Code Deployment**: Deployment using the `nohup` method. This is a basic deployment method suitable for development and testing environments. For details, refer to the [Regular Source Code Deployment Guide](https://docs.openim.io/).
   + **Production-Level Deployment**: Deployment using the `system` method, more suitable for production environments. This method provides higher stability and reliability. For details, refer to the [Production-Level Deployment Guide](https://docs.openim.io/guides/gettingStarted/install-openim-linux-system).
2. Cluster Deployment:
   + **Kubernetes Deployment**: Provides two deployment methods, including deployment through Helm and sealos. This is suitable for environments that require high availability and scalability. Specific methods can be found in the [Kubernetes Deployment Guide](https://docs.openim.io/guides/gettingStarted/k8s-deployment).
3. Docker Deployment:
   + **Regular Docker Deployment**: Suitable for quick deployments and small projects. For detailed information, refer to the [Docker Deployment Guide](https://docs.openim.io/guides/gettingStarted/dockerCompose).
   + **Docker Compose Deployment**: Provides more convenient service management and configuration, suitable for complex multi-container applications.

Next, we will introduce the specific steps, monitoring, and management backend configuration for each of these deployment methods, as well as usage tips to help you choose the most suitable deployment option according to your needs.

##  1. <a name='SourceCodeDocker'></a>Source Code & Docker

###  1.1. <a name='Deployment'></a>Deployment

OpenIM deploys openim-server and openim-chat from source code, while other components are deployed via Docker.

For Docker deployment, you can deploy all components with a single command using the [openimsdk/openim-docker](https://github.com/openimsdk/openim-docker) repository. The deployment configuration can be found in the [environment.sh](https://github.com/openimsdk/open-im-server/blob/main/scripts/install/environment.sh) document, which provides information on how to learn and familiarize yourself with various environment variables.

For Prometheus, it is not enabled by default. To enable it, set the environment variable before executing `make init`:

```bash
export PROMETHEUS_ENABLE=true   # Default is false
```

Then, execute:

```bash
make init
docker compose up -d
```

###  1.2. <a name='Configuration'></a>Configuration

To configure Prometheus data sources in Grafana, follow these steps:

1. **Log in to Grafana**: First, open your web browser and access the Grafana URL. If you haven't changed the port, the address is typically [http://localhost:13000](http://localhost:13000/).

2. **Log in with default credentials**: Grafana's default username and password are both `admin`. You will be prompted to change the password on your first login.

3. **Access Data Sources Settings**:

   + In the left menu of Grafana, look for and click the "gear" icon representing "Configuration."
   + In the configuration menu, select "Data Sources."

4. **Add a New Data Source**:

   + On the Data Sources page, click the "Add data source" button.
   + In the list, find and select "Prometheus."

   ![image-20231114175117374](http://sm.nsddd.top/sm202311141751692.png)

   Click `Add New connection` to add more data sources, such as Loki (responsible for log storage and query processing).

5. **Configure the Prometheus Data Source**:

   + On the configuration page, fill in the details of the Prometheus server. This typically includes the URL of the Prometheus service (e.g., if Prometheus is running on the same machine as OpenIM, the URL might be `http://172.28.0.1:19090`, with the address matching the `DOCKER_BRIDGE_GATEWAY` variable address). OpenIM and the components are linked via a gateway. The default port used by OpenIM is `19090`.
   + Adjust other settings as needed, such as authentication and TLS settings.

   ![image-20231114180351923](http://sm.nsddd.top/sm202311141803076.png)

6. **Save and Test**:

   + After completing the configuration, click the "Save & Test" button to ensure that Grafana can successfully connect to Prometheus.

**Importing Dashboards in Grafana**

Importing Grafana Dashboards is a straightforward process and is applicable to OpenIM Server application services and Node Exporter. Here are detailed steps and necessary considerations:

**Key Metrics Overview and Deployment Steps**

To monitor OpenIM in Grafana, you need to focus on three categories of key metrics, each with its specific deployment and configuration steps:

**OpenIM Metrics (`prometheus-dashboard.yaml`)**:

- **Configuration File Path**: Find this at `config/prometheus-dashboard.yaml`.
- **Enabling Monitoring**: Activate Prometheus monitoring by setting the environment variable: `export PROMETHEUS_ENABLE=true`.
- **More Information**: For detailed instructions, see the [OpenIM Configuration Guide](https://docs.openim.io/configurations/prometheus-integration).

**Node Exporter**:

- **Container Deployment**: Use the container `quay.io/prometheus/node-exporter` for effective node monitoring.
- **Access Dashboard**: Visit the [Node Exporter Full Feature Dashboard](https://grafana.com/grafana/dashboards/1860-node-exporter-full/) for dashboard integration either through YAML file download or ID.
- **Deployment Guide**: For deployment steps, consult the [Node Exporter Deployment Documentation](https://prometheus.io/docs/guides/node-exporter/).

**Middleware Metrics**: Different middlewares require unique steps and configurations for monitoring:

- MySQL:
    - **Configuration**: Make sure MySQL is set up for performance monitoring.
    - **Guide**: See the [MySQL Monitoring Configuration Guide](https://grafana.com/docs/grafana/latest/datasources/mysql/).
- Redis:
    - **Configuration**: Adjust Redis settings to enable monitoring data export.
    - **Guide**: Consult the [Redis Monitoring Guide](https://grafana.com/docs/grafana/latest/datasources/redis/).
- MongoDB:
    - **Configuration**: Configure MongoDB for monitoring metrics.
    - **Guide**: Visit the [MongoDB Monitoring Guide](https://grafana.com/grafana/plugins/grafana-mongodb-datasource/).
- Kafka:
    - **Configuration**: Set up Kafka for Prometheus monitoring integration.
    - **Guide**: Refer to the [Kafka Monitoring Guide](https://grafana.com/grafana/plugins/grafana-kafka-datasource/).
- Zookeeper:
    - **Configuration**: Ensure Prometheus can monitor Zookeeper.
    - **Guide**: Check out the [Zookeeper Monitoring Configuration](https://grafana.com/docs/grafana/latest/datasources/zookeeper/).

**Importing Steps**:

1. Access the Dashboard Import Interface:

   + Click the `+` icon on the left menu or in the top right corner of Grafana, then select "Create."
   + Choose "Import" to access the dashboard import interface.

2. **Perform Dashboard Import**:
   + **Upload via File**: Directly upload your YAML file.
   + **Paste Content**: Open the YAML file, copy its content, and paste it into the import interface.
   + **Import via Grafana.com Dashboard**: Visit [Grafana Dashboards](https://grafana.com/grafana/dashboards/), search for the desired dashboard, and import it using its ID.
3. **Configure the Dashboard**:
   + Select the appropriate data source, such as the previously configured Prometheus.
   + Adjust other settings, such as the dashboard name or folder.
4. **Save and View the Dashboard**:
   + After configuring, click "Import" to complete the process.
   + Immediately view the new dashboard after successful import.

**Graph Examples:**

![image-20231114194451673](http://sm.nsddd.top/sm202311141944953.png)



###  1.3. <a name='MonitoringRunninginDockerGuide'></a>Monitoring Running in Docker Guide

####  1.3.1. <a name='Introduction'></a>Introduction

This guide provides the steps to run OpenIM using Docker. OpenIM is an open-source instant messaging solution that can be quickly deployed using Docker. For more information, please refer to the [OpenIM Docker GitHub](https://github.com/openimsdk/openim-docker).

####  1.3.2. <a name='Prerequisites'></a>Prerequisites

+ Ensure that Docker and Docker Compose are installed.
+ Basic understanding of Docker and containerization technology.

####  1.3.3. <a name='Step1:ClonetheRepository'></a>Step 1: Clone the Repository

First, clone the OpenIM Docker repository:

```bash
git clone https://github.com/openimsdk/openim-docker.git
```

Navigate to the repository directory and check the `README` file for more information and configuration options.

####  1.3.4. <a name='Step2:StartDockerCompose'></a>Step 2: Start Docker Compose

In the repository directory, run the following command to start the service:

```bash
docker-compose up -d
```

This will download the required Docker images and start the OpenIM service.

####  1.3.5. <a name='Step3:UsetheOpenIMWebInterface'></a>Step 3: Use the OpenIM Web Interface

+ Open a browser in private mode and access [OpenIM Web](http://localhost:11001/).
+ Register two users and try adding friends.
+ Test sending messages and pictures.

####  1.3.6. <a name='RunningEffect'></a>Running Effect

![image-20231115100811208](http://sm.nsddd.top/sm202311151008639.png)

####  1.3.7. <a name='Step4:AccesstheAdminPanel'></a>Step 4: Access the Admin Panel

+ Access the [OpenIM Admin Panel](http://localhost:11002/).
+ Log in using the default username and password (`admin1:admin1`).

Running Effect Image:

![image-20231115101039837](http://sm.nsddd.top/sm202311151010116.png)

####  1.3.8. <a name='Step5:AccesstheMonitoringInterface'></a>Step 5: Access the Monitoring Interface

+ Log in to the [Monitoring Interface](http://localhost:3000/login) using the credentials (`admin:admin`).

####  1.3.9. <a name='NextSteps'></a>Next Steps

+ Configure and manage the services following the steps provided in the OpenIM source code.
+ Refer to the `README` file for advanced configuration and management.

####  1.3.10. <a name='Troubleshooting'></a>Troubleshooting

+ If you encounter any issues, please check the documentation on [OpenIM Docker GitHub](https://github.com/openimsdk/openim-docker) or search for related issues in the Issues section.
+ If the problem persists, you can create an issue on the [openim-docker](https://github.com/openimsdk/openim-docker/issues/new/choose) repository or the [openim-server](https://github.com/openimsdk/open-im-server/issues/new/choose) repository.



##  2. <a name='Kubernetes'></a>Kubernetes

Refer to [openimsdk/helm-charts](https://github.com/openimsdk/helm-charts).

When deploying and monitoring OpenIM in a Kubernetes environment, you will focus on three main metrics: middleware, custom OpenIM metrics, and Node Exporter. Here are detailed steps and guidelines:

###  2.1. <a name='MiddlewareMonitoring'></a>Middleware Monitoring

Middleware monitoring is crucial to ensure the overall system's stability. Typically, this includes monitoring the following components:

+ **MySQL**: Monitor database performance, query latency, and more.
+ **Redis**: Track operation latency, memory usage, and more.
+ **MongoDB**: Observe database operations, resource usage, and more.
+ **Kafka**: Monitor message throughput, latency, and more.
+ **Zookeeper**: Keep an eye on cluster status, performance metrics, and more.

For Kubernetes environments, you can use the corresponding Prometheus Exporters to collect monitoring data for these middleware components.

###  2.2. <a name='CustomOpenIMMetrics'></a>Custom OpenIM Metrics

Custom OpenIM metrics provide essential information about the OpenIM application itself, such as user activity, message traffic, system performance, and more. To monitor these metrics in Kubernetes:

+ Ensure OpenIM application configurations expose Prometheus metrics.
+ When deploying using Helm charts (refer to [OpenIM Helm Charts](https://github.com/openimsdk/helm-charts)), pay attention to configuring relevant monitoring settings.

###  2.3. <a name='NodeExporter'></a>Node Exporter

Node Exporter is used to collect hardware and operating system-level metrics for Kubernetes nodes, such as CPU, memory, disk usage, and more. To integrate Node Exporter in Kubernetes:

+ Deploy Node Exporter using the appropriate Helm chart. You can find information and guides on [Prometheus Community](https://prometheus.io/docs/guides/node-exporter/).
+ Ensure Node Exporter's data is collected by Prometheus instances within your cluster.



##  3. <a name='SettingUpandConfiguringAlertManagerUsingEnvironmentVariablesandmakeinit'></a>Setting Up and Configuring AlertManager Using Environment Variables and `make init`

###  3.1. <a name='Introduction-1'></a>Introduction

AlertManager, a component of the Prometheus monitoring system, handles alerts sent by client applications such as the Prometheus server. It takes care of deduplicating, grouping, and routing them to the correct receiver. This document outlines how to set up and configure AlertManager using environment variables and the `make init` command. We will focus on configuring key fields like the sender's email, SMTP settings, and SMTP authentication password.

###  3.2. <a name='Prerequisites-1'></a>Prerequisites

+ Basic knowledge of terminal and command-line operations.
+ AlertManager installed on your system.
+ Access to an SMTP server for sending emails.

###  3.3. <a name='ConfigurationSteps'></a>Configuration Steps

####  3.3.1. <a name='ExportingEnvironmentVariables'></a>Exporting Environment Variables

Before initializing AlertManager, you need to set environment variables. These variables are used to configure the AlertManager settings without altering the code. Use the `export` command in your terminal. Here are some key variables you might set:

+ `export ALERTMANAGER_RESOLVE_TIMEOUT='5m'`
+ `export ALERTMANAGER_SMTP_FROM='alert@example.com'`
+ `export ALERTMANAGER_SMTP_SMARTHOST='smtp.example.com:465'`
+ `export ALERTMANAGER_SMTP_AUTH_USERNAME='alert@example.com'`
+ `export ALERTMANAGER_SMTP_AUTH_PASSWORD='your_password'`
+ `export ALERTMANAGER_SMTP_REQUIRE_TLS='false'`

####  3.3.2. <a name='InitializingAlertManager'></a>Initializing AlertManager

After setting the necessary environment variables, you can initialize AlertManager by running the `make init` command. This command typically runs a script that prepares AlertManager with the provided configuration.

####  3.3.3. <a name='KeyConfigurationFields'></a>Key Configuration Fields

##### a. Sender's Email (`ALERTMANAGER_SMTP_FROM`)

This variable sets the email address that will appear as the sender in the notifications sent by AlertManager.

##### b. SMTP Configuration

+ **SMTP Server (`ALERTMANAGER_SMTP_SMARTHOST`):** Specifies the address and port of the SMTP server used for sending emails.
+ **SMTP Authentication Username (`ALERTMANAGER_SMTP_AUTH_USERNAME`):** The username for authenticating with the SMTP server.
+ **SMTP Authentication Password (`ALERTMANAGER_SMTP_AUTH_PASSWORD`):** The password for SMTP server authentication. It's crucial to keep this value secure.

####  3.3.4. <a name='ConfiguringSMTPAuthenticationPassword'></a>Configuring SMTP Authentication Password

The SMTP authentication password can be set using the `ALERTMANAGER_SMTP_AUTH_PASSWORD` environment variable. It's recommended to use a secure method to set this variable to avoid exposing sensitive information. For instance, you might read the password from a secure file or a secret management tool.

####  3.3.5. <a name='UsefulLinksforCommonEmailServers'></a>Useful Links for Common Email Servers

For specific configurations related to common email servers, you may refer to their respective documentation:

+ Gmail SMTP Settings:
  + [Gmail SMTP Configuration](https://support.google.com/mail/answer/7126229?hl=en)
+ Microsoft Outlook SMTP Settings:
  + [Outlook Email Settings](https://support.microsoft.com/en-us/office/pop-imap-and-smtp-settings-8361e398-8af4-4e97-b147-6c6c4ac95353)
+ Yahoo Mail SMTP Settings:
  + [Yahoo SMTP Configuration](https://help.yahoo.com/kb/SLN4724.html)

###  3.4. <a name='Conclusion'></a>Conclusion

Setting up and configuring AlertManager with environment variables provides a flexible and secure way to manage alert settings. By following the above steps, you can easily configure AlertManager for your monitoring needs. Always ensure to secure sensitive information, especially when dealing with SMTP authentication credentials.