# Deployment and Design of OpenIM's Management Backend and Monitoring

<!-- vscode-markdown-toc -->
* 1. [Source Code & Docker](#SourceCodeDocker)
	* 1.1. [Deployment](#Deployment)
	* 1.2. [Configuration](#Configuration)
* 2. [Kubernetes](#Kubernetes)
	* 2.1. [Middleware Monitoring](#MiddlewareMonitoring)
	* 2.2. [Custom OpenIM Metrics](#CustomOpenIMMetrics)
	* 2.3. [Node Exporter](#NodeExporter)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

OpenIM offers various flexible deployment options to suit different environments and requirements. Here is a simplified and optimized description of these deployment options:

1. Source Code Deployment:
   + **Regular Source Code Deployment**: Deployment using the `nohup` method. This is a basic deployment method suitable for development and testing environments. For details, refer to the [Regular Source Code Deployment Guide](https://docs.openim.io/guides/gettingStarted/imSourceCodeDeployment).
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

1. **Log in to Grafana**: First, open your web browser and access the Grafana URL. If you haven't changed the port, the address is typically [http://localhost:3000](http://localhost:3000/).

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

1. **OpenIM Metrics (`prometheus-dashboard.yaml`)**:
   + **Configuration File Path**: Located at `config/prometheus-dashboard.yaml`.
   + **Enabling Monitoring**: Set the environment variable `export PROMETHEUS_ENABLE=true` to enable Prometheus monitoring.
   + **More Information**: Refer to the [OpenIM Configuration Guide](https://docs.openim.io/configurations/prometheus-integration).
2. **Node Exporter**:
   + **Container Deployment**: Deploy the `quay.io/prometheus/node-exporter` container for node monitoring.
   + **Get Dashboard**: Access the [Node Exporter Full Feature Dashboard](https://grafana.com/grafana/dashboards/1860-node-exporter-full/) and import it using YAML file download or ID import.
   + **Deployment Guide**: Refer to the [Node Exporter Deployment Documentation](https://prometheus.io/docs/guides/node-exporter/).
3. **Middleware Metrics**: Each middleware requires specific steps and configurations to enable monitoring. Here is a list of common middleware and links to their respective setup guides:
   + MySQL:
     + **Configuration**: Ensure MySQL has performance monitoring enabled.
     + **Link**: Refer to the [MySQL Monitoring Configuration Guide](https://grafana.com/docs/grafana/latest/datasources/mysql/).
   + Redis:
     + **Configuration**: Configure Redis to allow monitoring data export.
     + **Link**: Refer to the [Redis Monitoring Guide](https://grafana.com/docs/grafana/latest/datasources/redis/).
   + MongoDB:
     + **Configuration**: Set up monitoring metrics for MongoDB.
     + **Link**: Refer to the [MongoDB Monitoring Guide](https://grafana.com/grafana/plugins/grafana-mongodb-datasource/).
   + Kafka:
     + **Configuration**: Integrate Kafka with Prometheus monitoring.
     + **Link**: Refer to the [Kafka Monitoring Guide](https://grafana.com/grafana/plugins/grafana-kafka-datasource/).
   + Zookeeper:
     + **Configuration**: Ensure Zookeeper can be monitored by Prometheus.
     + **Link**: Refer to the [Zookeeper Monitoring Configuration](https://grafana.com/docs/grafana/latest/datasources/zookeeper/).



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