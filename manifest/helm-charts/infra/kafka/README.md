<!--- app-name: Apache Kafka -->

# Apache Kafka packaged by Bitnami

Apache Kafka is a distributed streaming platform designed to build real-time pipelines and can be used as a message broker or as a replacement for a log aggregation solution for big data applications.

[Overview of Apache Kafka](http://kafka.apache.org/)

Trademarks: This software listing is packaged by Bitnami. The respective trademarks mentioned in the offering are owned by the respective companies, and use of them does not imply any affiliation or endorsement.

## TL;DR

```console
helm install my-release oci://registry-1.docker.io/bitnamicharts/kafka
```

## Introduction

This chart bootstraps a [Kafka](https://github.com/bitnami/containers/tree/main/bitnami/kafka) deployment on a [Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

Bitnami charts can be used with [Kubeapps](https://kubeapps.dev/) for deployment and management of Helm Charts in clusters.

Looking to use Apache Kafka in production? Try [VMware Application Catalog](https://bitnami.com/enterprise), the enterprise edition of Bitnami Application Catalog.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release oci://registry-1.docker.io/bitnamicharts/kafka
```

These commands deploy Kafka on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```console
helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Parameters

### Global parameters

| Name                      | Description                                     | Value |
| ------------------------- | ----------------------------------------------- | ----- |
| `global.imageRegistry`    | Global Docker image registry                    | `""`  |
| `global.imagePullSecrets` | Global Docker registry secret names as an array | `[]`  |
| `global.storageClass`     | Global StorageClass for Persistent Volume(s)    | `""`  |

### Common parameters

| Name                      | Description                                                                             | Value           |
| ------------------------- | --------------------------------------------------------------------------------------- | --------------- |
| `kubeVersion`             | Override Kubernetes version                                                             | `""`            |
| `nameOverride`            | String to partially override common.names.fullname                                      | `""`            |
| `fullnameOverride`        | String to fully override common.names.fullname                                          | `""`            |
| `clusterDomain`           | Default Kubernetes cluster domain                                                       | `cluster.local` |
| `commonLabels`            | Labels to add to all deployed objects                                                   | `{}`            |
| `commonAnnotations`       | Annotations to add to all deployed objects                                              | `{}`            |
| `extraDeploy`             | Array of extra objects to deploy with the release                                       | `[]`            |
| `serviceBindings.enabled` | Create secret for service binding (Experimental)                                        | `false`         |
| `diagnosticMode.enabled`  | Enable diagnostic mode (all probes will be disabled and the command will be overridden) | `false`         |
| `diagnosticMode.command`  | Command to override all containers in the statefulset                                   | `["sleep"]`     |
| `diagnosticMode.args`     | Args to override all containers in the statefulset                                      | `["infinity"]`  |

### Kafka parameters

| Name                                  | Description                                                                                                                                                                                                | Value                 |
| ------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `image.registry`                      | Kafka image registry                                                                                                                                                                                       | `docker.io`           |
| `image.repository`                    | Kafka image repository                                                                                                                                                                                     | `bitnami/kafka`       |
| `image.tag`                           | Kafka image tag (immutable tags are recommended)                                                                                                                                                           | `3.5.1-debian-11-r44` |
| `image.digest`                        | Kafka image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                                      | `""`                  |
| `image.pullPolicy`                    | Kafka image pull policy                                                                                                                                                                                    | `IfNotPresent`        |
| `image.pullSecrets`                   | Specify docker-registry secret names as an array                                                                                                                                                           | `[]`                  |
| `image.debug`                         | Specify if debug values should be set                                                                                                                                                                      | `false`               |
| `extraInit`                           | Additional content for the kafka init script, rendered as a template.                                                                                                                                      | `""`                  |
| `config`                              | Configuration file for Kafka, rendered as a template. Auto-generated based on chart values when not specified.                                                                                             | `""`                  |
| `existingConfigmap`                   | ConfigMap with Kafka Configuration                                                                                                                                                                         | `""`                  |
| `extraConfig`                         | Additional configuration to be appended at the end of the generated Kafka configuration file.                                                                                                              | `""`                  |
| `secretConfig`                        | Additional configuration to be appended at the end of the generated Kafka configuration file.                                                                                                              | `""`                  |
| `existingSecretConfig`                | Secret with additonal configuration that will be appended to the end of the generated Kafka configuration file                                                                                             | `""`                  |
| `log4j`                               | An optional log4j.properties file to overwrite the default of the Kafka brokers                                                                                                                            | `""`                  |
| `existingLog4jConfigMap`              | The name of an existing ConfigMap containing a log4j.properties file                                                                                                                                       | `""`                  |
| `heapOpts`                            | Kafka Java Heap size                                                                                                                                                                                       | `-Xmx1024m -Xms1024m` |
| `interBrokerProtocolVersion`          | Override the setting 'inter.broker.protocol.version' during the ZK migration.                                                                                                                              | `""`                  |
| `listeners.client.name`               | Name for the Kafka client listener                                                                                                                                                                         | `CLIENT`              |
| `listeners.client.containerPort`      | Port for the Kafka client listener                                                                                                                                                                         | `9092`                |
| `listeners.client.protocol`           | Security protocol for the Kafka client listener. Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                                    | `SASL_PLAINTEXT`      |
| `listeners.client.sslClientAuth`      | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.authType for this listener. Allowed values are 'none', 'requested' and 'required'      | `""`                  |
| `listeners.controller.name`           | Name for the Kafka controller listener                                                                                                                                                                     | `CONTROLLER`          |
| `listeners.controller.containerPort`  | Port for the Kafka controller listener                                                                                                                                                                     | `9093`                |
| `listeners.controller.protocol`       | Security protocol for the Kafka controller listener. Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                                | `SASL_PLAINTEXT`      |
| `listeners.controller.sslClientAuth`  | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.authType for this listener. Allowed values are 'none', 'requested' and 'required'      | `""`                  |
| `listeners.interbroker.name`          | Name for the Kafka inter-broker listener                                                                                                                                                                   | `INTERNAL`            |
| `listeners.interbroker.containerPort` | Port for the Kafka inter-broker listener                                                                                                                                                                   | `9094`                |
| `listeners.interbroker.protocol`      | Security protocol for the Kafka inter-broker listener. Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                              | `SASL_PLAINTEXT`      |
| `listeners.interbroker.sslClientAuth` | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.authType for this listener. Allowed values are 'none', 'requested' and 'required'      | `""`                  |
| `listeners.external.containerPort`    | Port for the Kafka external listener                                                                                                                                                                       | `9095`                |
| `listeners.external.protocol`         | Security protocol for the Kafka external listener. . Allowed values are 'PLAINTEXT', 'SASL_PLAINTEXT', 'SASL_SSL' and 'SSL'                                                                                | `SASL_PLAINTEXT`      |
| `listeners.external.name`             | Name for the Kafka external listener                                                                                                                                                                       | `EXTERNAL`            |
| `listeners.external.sslClientAuth`    | Optional. If SASL_SSL is enabled, configure mTLS TLS authentication type. If SSL protocol is enabled, overrides tls.sslClientAuth for this listener. Allowed values are 'none', 'requested' and 'required' | `""`                  |
| `listeners.extraListeners`            | Array of listener objects to be appended to already existing listeners                                                                                                                                     | `[]`                  |
| `listeners.overrideListeners`         | Overrides the Kafka 'listeners' configuration setting.                                                                                                                                                     | `""`                  |
| `listeners.advertisedListeners`       | Overrides the Kafka 'advertised.listener' configuration setting.                                                                                                                                           | `""`                  |
| `listeners.securityProtocolMap`       | Overrides the Kafka 'security.protocol.map' configuration setting.                                                                                                                                         | `""`                  |

### Kafka SASL parameters

| Name                        | Description                                                                                                                                                     | Value                               |
| --------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------- | ----------------------------------- |
| `sasl.enabledMechanisms`    | Comma-separated list of allowed SASL mechanisms when SASL listeners are configured. Allowed types: `PLAIN`, `SCRAM-SHA-256`, `SCRAM-SHA-512`                    | `PLAIN,SCRAM-SHA-256,SCRAM-SHA-512` |
| `sasl.interBrokerMechanism` | SASL mechanism for inter broker communication.                                                                                                                  | `PLAIN`                             |
| `sasl.controllerMechanism`  | SASL mechanism for controller communications.                                                                                                                   | `PLAIN`                             |
| `sasl.interbroker.user`     | Username for inter-broker communications when SASL is enabled                                                                                                   | `inter_broker_user`                 |
| `sasl.interbroker.password` | Password for inter-broker communications when SASL is enabled. If not set and SASL is enabled for the controller listener, a random password will be generated. | `""`                                |
| `sasl.controller.user`      | Username for controller communications when SASL is enabled                                                                                                     | `controller_user`                   |
| `sasl.controller.password`  | Password for controller communications when SASL is enabled. If not set and SASL is enabled for the inter-broker listener, a random password will be generated. | `""`                                |
| `sasl.client.users`         | Comma-separated list of usernames for client communications when SASL is enabled                                                                                | `["user1"]`                         |
| `sasl.client.passwords`     | Comma-separated list of passwords for client communications when SASL is enabled, must match the number of client.users                                         | `""`                                |
| `sasl.zookeeper.user`       | Username for zookeeper communications when SASL is enabled.                                                                                                     | `""`                                |
| `sasl.zookeeper.password`   | Password for zookeeper communications when SASL is enabled.                                                                                                     | `""`                                |
| `sasl.existingSecret`       | Name of the existing secret containing credentials for clientUsers, interBrokerUser, controllerUser and zookeeperUser                                           | `""`                                |

### Kafka TLS parameters

| Name                                         | Description                                                                                                                             | Value                      |
| -------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- |
| `tls.type`                                   | Format to use for TLS certificates. Allowed types: `JKS` and `PEM`                                                                      | `JKS`                      |
| `tls.pemChainIncluded`                       | Flag to denote that the Certificate Authority (CA) certificates are bundled with the endpoint cert.                                     | `false`                    |
| `tls.existingSecret`                         | Name of the existing secret containing the TLS certificates for the Kafka nodes.                                                        | `""`                       |
| `tls.autoGenerated`                          | Generate automatically self-signed TLS certificates for Kafka brokers. Currently only supported if `tls.type` is `PEM`                  | `false`                    |
| `tls.passwordsSecret`                        | Name of the secret containing the password to access the JKS files or PEM key when they are password-protected. (`key`: `password`)     | `""`                       |
| `tls.passwordsSecretKeystoreKey`             | The secret key from the tls.passwordsSecret containing the password for the Keystore.                                                   | `keystore-password`        |
| `tls.passwordsSecretTruststoreKey`           | The secret key from the tls.passwordsSecret containing the password for the Truststore.                                                 | `truststore-password`      |
| `tls.passwordsSecretPemPasswordKey`          | The secret key from the tls.passwordsSecret containing the password for the PEM key inside 'tls.passwordsSecret'.                       | `""`                       |
| `tls.keystorePassword`                       | Password to access the JKS keystore when it is password-protected. Ignored when 'tls.passwordsSecret' is provided.                      | `""`                       |
| `tls.truststorePassword`                     | Password to access the JKS truststore when it is password-protected. Ignored when 'tls.passwordsSecret' is provided.                    | `""`                       |
| `tls.keyPassword`                            | Password to access the PEM key when it is password-protected.                                                                           | `""`                       |
| `tls.jksTruststoreSecret`                    | Name of the existing secret containing your truststore if truststore not existing or different from the one in the `tls.existingSecret` | `""`                       |
| `tls.jksTruststoreKey`                       | The secret key from the `tls.existingSecret` or `tls.jksTruststoreSecret` containing the truststore                                     | `""`                       |
| `tls.endpointIdentificationAlgorithm`        | The endpoint identification algorithm to validate server hostname using server certificate                                              | `https`                    |
| `tls.sslClientAuth`                          | Sets the default value for the ssl.client.auth Kafka setting.                                                                           | `required`                 |
| `tls.zookeeper.enabled`                      | Enable TLS for Zookeeper client connections.                                                                                            | `false`                    |
| `tls.zookeeper.verifyHostname`               | Hostname validation.                                                                                                                    | `true`                     |
| `tls.zookeeper.existingSecret`               | Name of the existing secret containing the TLS certificates for ZooKeeper client communications.                                        | `""`                       |
| `tls.zookeeper.existingSecretKeystoreKey`    | The secret key from the  tls.zookeeper.existingSecret containing the Keystore.                                                          | `zookeeper.keystore.jks`   |
| `tls.zookeeper.existingSecretTruststoreKey`  | The secret key from the tls.zookeeper.existingSecret containing the Truststore.                                                         | `zookeeper.truststore.jks` |
| `tls.zookeeper.passwordsSecret`              | Existing secret containing Keystore and Truststore passwords.                                                                           | `""`                       |
| `tls.zookeeper.passwordsSecretKeystoreKey`   | The secret key from the tls.zookeeper.passwordsSecret containing the password for the Keystore.                                         | `keystore-password`        |
| `tls.zookeeper.passwordsSecretTruststoreKey` | The secret key from the tls.zookeeper.passwordsSecret containing the password for the Truststore.                                       | `truststore-password`      |
| `tls.zookeeper.keystorePassword`             | Password to access the JKS keystore when it is password-protected. Ignored when 'tls.passwordsSecret' is provided.                      | `""`                       |
| `tls.zookeeper.truststorePassword`           | Password to access the JKS truststore when it is password-protected. Ignored when 'tls.passwordsSecret' is provided.                    | `""`                       |
| `extraEnvVars`                               | Extra environment variables to add to Kafka pods                                                                                        | `[]`                       |
| `extraEnvVarsCM`                             | ConfigMap with extra environment variables                                                                                              | `""`                       |
| `extraEnvVarsSecret`                         | Secret with extra environment variables                                                                                                 | `""`                       |
| `extraVolumes`                               | Optionally specify extra list of additional volumes for the Kafka pod(s)                                                                | `[]`                       |
| `extraVolumeMounts`                          | Optionally specify extra list of additional volumeMounts for the Kafka container(s)                                                     | `[]`                       |
| `sidecars`                                   | Add additional sidecar containers to the Kafka pod(s)                                                                                   | `[]`                       |
| `initContainers`                             | Add additional Add init containers to the Kafka pod(s)                                                                                  | `[]`                       |

### Controller-eligible statefulset parameters

| Name                                                           | Description                                                                                                                                                                                   | Value                     |
| -------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- |
| `controller.replicaCount`                                      | Number of Kafka controller-eligible nodes                                                                                                                                                     | `3`                       |
| `controller.controllerOnly`                                    | If set to true, controller nodes will be deployed as dedicated controllers, instead of controller+broker processes.                                                                           | `false`                   |
| `controller.minId`                                             | Minimal node.id values for controller-eligible nodes. Do not change after first initialization.                                                                                               | `0`                       |
| `controller.zookeeperMigrationMode`                            | Set to true to deploy cluster controller quorum                                                                                                                                               | `false`                   |
| `controller.config`                                            | Configuration file for Kafka controller-eligible nodes, rendered as a template. Auto-generated based on chart values when not specified.                                                      | `""`                      |
| `controller.existingConfigmap`                                 | ConfigMap with Kafka Configuration for controller-eligible nodes.                                                                                                                             | `""`                      |
| `controller.extraConfig`                                       | Additional configuration to be appended at the end of the generated Kafka controller-eligible nodes configuration file.                                                                       | `""`                      |
| `controller.secretConfig`                                      | Additional configuration to be appended at the end of the generated Kafka controller-eligible nodes configuration file.                                                                       | `""`                      |
| `controller.existingSecretConfig`                              | Secret with additonal configuration that will be appended to the end of the generated Kafka controller-eligible nodes configuration file                                                      | `""`                      |
| `controller.heapOpts`                                          | Kafka Java Heap size for controller-eligible nodes                                                                                                                                            | `-Xmx1024m -Xms1024m`     |
| `controller.command`                                           | Override Kafka container command                                                                                                                                                              | `[]`                      |
| `controller.args`                                              | Override Kafka container arguments                                                                                                                                                            | `[]`                      |
| `controller.extraEnvVars`                                      | Extra environment variables to add to Kafka pods                                                                                                                                              | `[]`                      |
| `controller.extraEnvVarsCM`                                    | ConfigMap with extra environment variables                                                                                                                                                    | `""`                      |
| `controller.extraEnvVarsSecret`                                | Secret with extra environment variables                                                                                                                                                       | `""`                      |
| `controller.extraContainerPorts`                               | Kafka controller-eligible extra containerPorts.                                                                                                                                               | `[]`                      |
| `controller.livenessProbe.enabled`                             | Enable livenessProbe on Kafka containers                                                                                                                                                      | `true`                    |
| `controller.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                                                                                       | `10`                      |
| `controller.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                                                                                              | `10`                      |
| `controller.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                                                                                             | `5`                       |
| `controller.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                                                                                           | `3`                       |
| `controller.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                                                                                           | `1`                       |
| `controller.readinessProbe.enabled`                            | Enable readinessProbe on Kafka containers                                                                                                                                                     | `true`                    |
| `controller.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                                                                                                      | `5`                       |
| `controller.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                                                                                             | `10`                      |
| `controller.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                                                                                            | `5`                       |
| `controller.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                                                                                          | `6`                       |
| `controller.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                                                                                          | `1`                       |
| `controller.startupProbe.enabled`                              | Enable startupProbe on Kafka containers                                                                                                                                                       | `false`                   |
| `controller.startupProbe.initialDelaySeconds`                  | Initial delay seconds for startupProbe                                                                                                                                                        | `30`                      |
| `controller.startupProbe.periodSeconds`                        | Period seconds for startupProbe                                                                                                                                                               | `10`                      |
| `controller.startupProbe.timeoutSeconds`                       | Timeout seconds for startupProbe                                                                                                                                                              | `1`                       |
| `controller.startupProbe.failureThreshold`                     | Failure threshold for startupProbe                                                                                                                                                            | `15`                      |
| `controller.startupProbe.successThreshold`                     | Success threshold for startupProbe                                                                                                                                                            | `1`                       |
| `controller.customLivenessProbe`                               | Custom livenessProbe that overrides the default one                                                                                                                                           | `{}`                      |
| `controller.customReadinessProbe`                              | Custom readinessProbe that overrides the default one                                                                                                                                          | `{}`                      |
| `controller.customStartupProbe`                                | Custom startupProbe that overrides the default one                                                                                                                                            | `{}`                      |
| `controller.lifecycleHooks`                                    | lifecycleHooks for the Kafka container to automate configuration before or after startup                                                                                                      | `{}`                      |
| `controller.resources.limits`                                  | The resources limits for the container                                                                                                                                                        | `{}`                      |
| `controller.resources.requests`                                | The requested resources for the container                                                                                                                                                     | `{}`                      |
| `controller.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                                                                                          | `true`                    |
| `controller.podSecurityContext.fsGroup`                        | Set Kafka pod's Security Context fsGroup                                                                                                                                                      | `1001`                    |
| `controller.podSecurityContext.seccompProfile.type`            | Set Kafka pods's Security Context seccomp profile                                                                                                                                             | `RuntimeDefault`          |
| `controller.containerSecurityContext.enabled`                  | Enable Kafka containers' Security Context                                                                                                                                                     | `true`                    |
| `controller.containerSecurityContext.runAsUser`                | Set Kafka containers' Security Context runAsUser                                                                                                                                              | `1001`                    |
| `controller.containerSecurityContext.runAsNonRoot`             | Set Kafka containers' Security Context runAsNonRoot                                                                                                                                           | `true`                    |
| `controller.containerSecurityContext.allowPrivilegeEscalation` | Force the child process to be run as non-privileged                                                                                                                                           | `false`                   |
| `controller.containerSecurityContext.readOnlyRootFilesystem`   | Allows the pod to mount the RootFS as ReadOnly only                                                                                                                                           | `true`                    |
| `controller.containerSecurityContext.capabilities.drop`        | Set Kafka containers' server Security Context capabilities to be dropped                                                                                                                      | `["ALL"]`                 |
| `controller.hostAliases`                                       | Kafka pods host aliases                                                                                                                                                                       | `[]`                      |
| `controller.hostNetwork`                                       | Specify if host network should be enabled for Kafka pods                                                                                                                                      | `false`                   |
| `controller.hostIPC`                                           | Specify if host IPC should be enabled for Kafka pods                                                                                                                                          | `false`                   |
| `controller.podLabels`                                         | Extra labels for Kafka pods                                                                                                                                                                   | `{}`                      |
| `controller.podAnnotations`                                    | Extra annotations for Kafka pods                                                                                                                                                              | `{}`                      |
| `controller.podAffinityPreset`                                 | Pod affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                           | `""`                      |
| `controller.podAntiAffinityPreset`                             | Pod anti-affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                      | `soft`                    |
| `controller.nodeAffinityPreset.type`                           | Node affinity preset type. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                     | `""`                      |
| `controller.nodeAffinityPreset.key`                            | Node label key to match Ignored if `affinity` is set.                                                                                                                                         | `""`                      |
| `controller.nodeAffinityPreset.values`                         | Node label values to match. Ignored if `affinity` is set.                                                                                                                                     | `[]`                      |
| `controller.affinity`                                          | Affinity for pod assignment                                                                                                                                                                   | `{}`                      |
| `controller.nodeSelector`                                      | Node labels for pod assignment                                                                                                                                                                | `{}`                      |
| `controller.tolerations`                                       | Tolerations for pod assignment                                                                                                                                                                | `[]`                      |
| `controller.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment spread across your cluster among failure-domains. Evaluated as a template                                                                      | `[]`                      |
| `controller.terminationGracePeriodSeconds`                     | Seconds the pod needs to gracefully terminate                                                                                                                                                 | `""`                      |
| `controller.podManagementPolicy`                               | StatefulSet controller supports relax its ordering guarantees while preserving its uniqueness and identity guarantees. There are two valid pod management policies: OrderedReady and Parallel | `Parallel`                |
| `controller.priorityClassName`                                 | Name of the existing priority class to be used by kafka pods                                                                                                                                  | `""`                      |
| `controller.runtimeClassName`                                  | Name of the runtime class to be used by pod(s)                                                                                                                                                | `""`                      |
| `controller.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                                                                                         | `true`                    |
| `controller.schedulerName`                                     | Name of the k8s scheduler (other than default)                                                                                                                                                | `""`                      |
| `controller.updateStrategy.type`                               | Kafka statefulset strategy type                                                                                                                                                               | `RollingUpdate`           |
| `controller.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka pod(s)                                                                                                                      | `[]`                      |
| `controller.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka container(s)                                                                                                           | `[]`                      |
| `controller.sidecars`                                          | Add additional sidecar containers to the Kafka pod(s)                                                                                                                                         | `[]`                      |
| `controller.initContainers`                                    | Add additional Add init containers to the Kafka pod(s)                                                                                                                                        | `[]`                      |
| `controller.pdb.create`                                        | Deploy a pdb object for the Kafka pod                                                                                                                                                         | `false`                   |
| `controller.pdb.minAvailable`                                  | Maximum number/percentage of unavailable Kafka replicas                                                                                                                                       | `""`                      |
| `controller.pdb.maxUnavailable`                                | Maximum number/percentage of unavailable Kafka replicas                                                                                                                                       | `1`                       |
| `controller.persistence.enabled`                               | Enable Kafka data persistence using PVC, note that ZooKeeper persistence is unaffected                                                                                                        | `true`                    |
| `controller.persistence.existingClaim`                         | A manually managed Persistent Volume and Claim                                                                                                                                                | `""`                      |
| `controller.persistence.storageClass`                          | PVC Storage Class for Kafka data volume                                                                                                                                                       | `""`                      |
| `controller.persistence.accessModes`                           | Persistent Volume Access Modes                                                                                                                                                                | `["ReadWriteOnce"]`       |
| `controller.persistence.size`                                  | PVC Storage Request for Kafka data volume                                                                                                                                                     | `8Gi`                     |
| `controller.persistence.annotations`                           | Annotations for the PVC                                                                                                                                                                       | `{}`                      |
| `controller.persistence.labels`                                | Labels for the PVC                                                                                                                                                                            | `{}`                      |
| `controller.persistence.selector`                              | Selector to match an existing Persistent Volume for Kafka data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                                            | `{}`                      |
| `controller.persistence.mountPath`                             | Mount path of the Kafka data volume                                                                                                                                                           | `/bitnami/kafka`          |
| `controller.logPersistence.enabled`                            | Enable Kafka logs persistence using PVC, note that ZooKeeper persistence is unaffected                                                                                                        | `false`                   |
| `controller.logPersistence.existingClaim`                      | A manually managed Persistent Volume and Claim                                                                                                                                                | `""`                      |
| `controller.logPersistence.storageClass`                       | PVC Storage Class for Kafka logs volume                                                                                                                                                       | `""`                      |
| `controller.logPersistence.accessModes`                        | Persistent Volume Access Modes                                                                                                                                                                | `["ReadWriteOnce"]`       |
| `controller.logPersistence.size`                               | PVC Storage Request for Kafka logs volume                                                                                                                                                     | `8Gi`                     |
| `controller.logPersistence.annotations`                        | Annotations for the PVC                                                                                                                                                                       | `{}`                      |
| `controller.logPersistence.selector`                           | Selector to match an existing Persistent Volume for Kafka log data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                                        | `{}`                      |
| `controller.logPersistence.mountPath`                          | Mount path of the Kafka logs volume                                                                                                                                                           | `/opt/bitnami/kafka/logs` |

### Broker-only statefulset parameters

| Name                                                       | Description                                                                                                                                                                                   | Value                     |
| ---------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------- |
| `broker.replicaCount`                                      | Number of Kafka broker-only nodes                                                                                                                                                             | `0`                       |
| `broker.minId`                                             | Minimal node.id values for broker-only nodes. Do not change after first initialization.                                                                                                       | `100`                     |
| `broker.zookeeperMigrationMode`                            | Set to true to deploy cluster controller quorum                                                                                                                                               | `false`                   |
| `broker.config`                                            | Configuration file for Kafka broker-only nodes, rendered as a template. Auto-generated based on chart values when not specified.                                                              | `""`                      |
| `broker.existingConfigmap`                                 | ConfigMap with Kafka Configuration for broker-only nodes.                                                                                                                                     | `""`                      |
| `broker.extraConfig`                                       | Additional configuration to be appended at the end of the generated Kafka broker-only nodes configuration file.                                                                               | `""`                      |
| `broker.secretConfig`                                      | Additional configuration to be appended at the end of the generated Kafka broker-only nodes configuration file.                                                                               | `""`                      |
| `broker.existingSecretConfig`                              | Secret with additonal configuration that will be appended to the end of the generated Kafka broker-only nodes configuration file                                                              | `""`                      |
| `broker.heapOpts`                                          | Kafka Java Heap size for broker-only nodes                                                                                                                                                    | `-Xmx1024m -Xms1024m`     |
| `broker.command`                                           | Override Kafka container command                                                                                                                                                              | `[]`                      |
| `broker.args`                                              | Override Kafka container arguments                                                                                                                                                            | `[]`                      |
| `broker.extraEnvVars`                                      | Extra environment variables to add to Kafka pods                                                                                                                                              | `[]`                      |
| `broker.extraEnvVarsCM`                                    | ConfigMap with extra environment variables                                                                                                                                                    | `""`                      |
| `broker.extraEnvVarsSecret`                                | Secret with extra environment variables                                                                                                                                                       | `""`                      |
| `broker.extraContainerPorts`                               | Kafka broker-only extra containerPorts.                                                                                                                                                       | `[]`                      |
| `broker.livenessProbe.enabled`                             | Enable livenessProbe on Kafka containers                                                                                                                                                      | `true`                    |
| `broker.livenessProbe.initialDelaySeconds`                 | Initial delay seconds for livenessProbe                                                                                                                                                       | `10`                      |
| `broker.livenessProbe.periodSeconds`                       | Period seconds for livenessProbe                                                                                                                                                              | `10`                      |
| `broker.livenessProbe.timeoutSeconds`                      | Timeout seconds for livenessProbe                                                                                                                                                             | `5`                       |
| `broker.livenessProbe.failureThreshold`                    | Failure threshold for livenessProbe                                                                                                                                                           | `3`                       |
| `broker.livenessProbe.successThreshold`                    | Success threshold for livenessProbe                                                                                                                                                           | `1`                       |
| `broker.readinessProbe.enabled`                            | Enable readinessProbe on Kafka containers                                                                                                                                                     | `true`                    |
| `broker.readinessProbe.initialDelaySeconds`                | Initial delay seconds for readinessProbe                                                                                                                                                      | `5`                       |
| `broker.readinessProbe.periodSeconds`                      | Period seconds for readinessProbe                                                                                                                                                             | `10`                      |
| `broker.readinessProbe.timeoutSeconds`                     | Timeout seconds for readinessProbe                                                                                                                                                            | `5`                       |
| `broker.readinessProbe.failureThreshold`                   | Failure threshold for readinessProbe                                                                                                                                                          | `6`                       |
| `broker.readinessProbe.successThreshold`                   | Success threshold for readinessProbe                                                                                                                                                          | `1`                       |
| `broker.startupProbe.enabled`                              | Enable startupProbe on Kafka containers                                                                                                                                                       | `false`                   |
| `broker.startupProbe.initialDelaySeconds`                  | Initial delay seconds for startupProbe                                                                                                                                                        | `30`                      |
| `broker.startupProbe.periodSeconds`                        | Period seconds for startupProbe                                                                                                                                                               | `10`                      |
| `broker.startupProbe.timeoutSeconds`                       | Timeout seconds for startupProbe                                                                                                                                                              | `1`                       |
| `broker.startupProbe.failureThreshold`                     | Failure threshold for startupProbe                                                                                                                                                            | `15`                      |
| `broker.startupProbe.successThreshold`                     | Success threshold for startupProbe                                                                                                                                                            | `1`                       |
| `broker.customLivenessProbe`                               | Custom livenessProbe that overrides the default one                                                                                                                                           | `{}`                      |
| `broker.customReadinessProbe`                              | Custom readinessProbe that overrides the default one                                                                                                                                          | `{}`                      |
| `broker.customStartupProbe`                                | Custom startupProbe that overrides the default one                                                                                                                                            | `{}`                      |
| `broker.lifecycleHooks`                                    | lifecycleHooks for the Kafka container to automate configuration before or after startup                                                                                                      | `{}`                      |
| `broker.resources.limits`                                  | The resources limits for the container                                                                                                                                                        | `{}`                      |
| `broker.resources.requests`                                | The requested resources for the container                                                                                                                                                     | `{}`                      |
| `broker.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                                                                                          | `true`                    |
| `broker.podSecurityContext.fsGroup`                        | Set Kafka pod's Security Context fsGroup                                                                                                                                                      | `1001`                    |
| `broker.podSecurityContext.seccompProfile.type`            | Set Kafka pod's Security Context seccomp profile                                                                                                                                              | `RuntimeDefault`          |
| `broker.containerSecurityContext.enabled`                  | Enable Kafka containers' Security Context                                                                                                                                                     | `true`                    |
| `broker.containerSecurityContext.runAsUser`                | Set Kafka containers' Security Context runAsUser                                                                                                                                              | `1001`                    |
| `broker.containerSecurityContext.runAsNonRoot`             | Set Kafka containers' Security Context runAsNonRoot                                                                                                                                           | `true`                    |
| `broker.containerSecurityContext.allowPrivilegeEscalation` | Force the child process to be run as non-privileged                                                                                                                                           | `false`                   |
| `broker.containerSecurityContext.readOnlyRootFilesystem`   | Allows the pod to mount the RootFS as ReadOnly only                                                                                                                                           | `true`                    |
| `broker.containerSecurityContext.capabilities.drop`        | Set Kafka containers' server Security Context capabilities to be dropped                                                                                                                      | `["ALL"]`                 |
| `broker.hostAliases`                                       | Kafka pods host aliases                                                                                                                                                                       | `[]`                      |
| `broker.hostNetwork`                                       | Specify if host network should be enabled for Kafka pods                                                                                                                                      | `false`                   |
| `broker.hostIPC`                                           | Specify if host IPC should be enabled for Kafka pods                                                                                                                                          | `false`                   |
| `broker.podLabels`                                         | Extra labels for Kafka pods                                                                                                                                                                   | `{}`                      |
| `broker.podAnnotations`                                    | Extra annotations for Kafka pods                                                                                                                                                              | `{}`                      |
| `broker.podAffinityPreset`                                 | Pod affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                           | `""`                      |
| `broker.podAntiAffinityPreset`                             | Pod anti-affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                      | `soft`                    |
| `broker.nodeAffinityPreset.type`                           | Node affinity preset type. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                     | `""`                      |
| `broker.nodeAffinityPreset.key`                            | Node label key to match Ignored if `affinity` is set.                                                                                                                                         | `""`                      |
| `broker.nodeAffinityPreset.values`                         | Node label values to match. Ignored if `affinity` is set.                                                                                                                                     | `[]`                      |
| `broker.affinity`                                          | Affinity for pod assignment                                                                                                                                                                   | `{}`                      |
| `broker.nodeSelector`                                      | Node labels for pod assignment                                                                                                                                                                | `{}`                      |
| `broker.tolerations`                                       | Tolerations for pod assignment                                                                                                                                                                | `[]`                      |
| `broker.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment spread across your cluster among failure-domains. Evaluated as a template                                                                      | `[]`                      |
| `broker.terminationGracePeriodSeconds`                     | Seconds the pod needs to gracefully terminate                                                                                                                                                 | `""`                      |
| `broker.podManagementPolicy`                               | StatefulSet controller supports relax its ordering guarantees while preserving its uniqueness and identity guarantees. There are two valid pod management policies: OrderedReady and Parallel | `Parallel`                |
| `broker.priorityClassName`                                 | Name of the existing priority class to be used by kafka pods                                                                                                                                  | `""`                      |
| `broker.runtimeClassName`                                  | Name of the runtime class to be used by pod(s)                                                                                                                                                | `""`                      |
| `broker.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                                                                                         | `true`                    |
| `broker.schedulerName`                                     | Name of the k8s scheduler (other than default)                                                                                                                                                | `""`                      |
| `broker.updateStrategy.type`                               | Kafka statefulset strategy type                                                                                                                                                               | `RollingUpdate`           |
| `broker.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka pod(s)                                                                                                                      | `[]`                      |
| `broker.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka container(s)                                                                                                           | `[]`                      |
| `broker.sidecars`                                          | Add additional sidecar containers to the Kafka pod(s)                                                                                                                                         | `[]`                      |
| `broker.initContainers`                                    | Add additional Add init containers to the Kafka pod(s)                                                                                                                                        | `[]`                      |
| `broker.pdb.create`                                        | Deploy a pdb object for the Kafka pod                                                                                                                                                         | `false`                   |
| `broker.pdb.minAvailable`                                  | Maximum number/percentage of unavailable Kafka replicas                                                                                                                                       | `""`                      |
| `broker.pdb.maxUnavailable`                                | Maximum number/percentage of unavailable Kafka replicas                                                                                                                                       | `1`                       |
| `broker.persistence.enabled`                               | Enable Kafka data persistence using PVC, note that ZooKeeper persistence is unaffected                                                                                                        | `true`                    |
| `broker.persistence.existingClaim`                         | A manually managed Persistent Volume and Claim                                                                                                                                                | `""`                      |
| `broker.persistence.storageClass`                          | PVC Storage Class for Kafka data volume                                                                                                                                                       | `""`                      |
| `broker.persistence.accessModes`                           | Persistent Volume Access Modes                                                                                                                                                                | `["ReadWriteOnce"]`       |
| `broker.persistence.size`                                  | PVC Storage Request for Kafka data volume                                                                                                                                                     | `8Gi`                     |
| `broker.persistence.annotations`                           | Annotations for the PVC                                                                                                                                                                       | `{}`                      |
| `broker.persistence.labels`                                | Labels for the PVC                                                                                                                                                                            | `{}`                      |
| `broker.persistence.selector`                              | Selector to match an existing Persistent Volume for Kafka data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                                            | `{}`                      |
| `broker.persistence.mountPath`                             | Mount path of the Kafka data volume                                                                                                                                                           | `/bitnami/kafka`          |
| `broker.logPersistence.enabled`                            | Enable Kafka logs persistence using PVC, note that ZooKeeper persistence is unaffected                                                                                                        | `false`                   |
| `broker.logPersistence.existingClaim`                      | A manually managed Persistent Volume and Claim                                                                                                                                                | `""`                      |
| `broker.logPersistence.storageClass`                       | PVC Storage Class for Kafka logs volume                                                                                                                                                       | `""`                      |
| `broker.logPersistence.accessModes`                        | Persistent Volume Access Modes                                                                                                                                                                | `["ReadWriteOnce"]`       |
| `broker.logPersistence.size`                               | PVC Storage Request for Kafka logs volume                                                                                                                                                     | `8Gi`                     |
| `broker.logPersistence.annotations`                        | Annotations for the PVC                                                                                                                                                                       | `{}`                      |
| `broker.logPersistence.selector`                           | Selector to match an existing Persistent Volume for Kafka log data PVC. If set, the PVC can't have a PV dynamically provisioned for it                                                        | `{}`                      |
| `broker.logPersistence.mountPath`                          | Mount path of the Kafka logs volume                                                                                                                                                           | `/opt/bitnami/kafka/logs` |

### Traffic Exposure parameters

| Name                                                         | Description                                                                                                                               | Value                   |
| ------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------- | ----------------------- |
| `service.type`                                               | Kubernetes Service type                                                                                                                   | `ClusterIP`             |
| `service.ports.client`                                       | Kafka svc port for client connections                                                                                                     | `9092`                  |
| `service.ports.controller`                                   | Kafka svc port for controller connections. It is used if "kraft.enabled: true"                                                            | `9093`                  |
| `service.ports.interbroker`                                  | Kafka svc port for inter-broker connections                                                                                               | `9094`                  |
| `service.ports.external`                                     | Kafka svc port for external connections                                                                                                   | `9095`                  |
| `service.extraPorts`                                         | Extra ports to expose in the Kafka service (normally used with the `sidecar` value)                                                       | `[]`                    |
| `service.nodePorts.client`                                   | Node port for the Kafka client connections                                                                                                | `""`                    |
| `service.nodePorts.external`                                 | Node port for the Kafka external connections                                                                                              | `""`                    |
| `service.sessionAffinity`                                    | Control where client requests go, to the same pod or round-robin                                                                          | `None`                  |
| `service.sessionAffinityConfig`                              | Additional settings for the sessionAffinity                                                                                               | `{}`                    |
| `service.clusterIP`                                          | Kafka service Cluster IP                                                                                                                  | `""`                    |
| `service.loadBalancerIP`                                     | Kafka service Load Balancer IP                                                                                                            | `""`                    |
| `service.loadBalancerSourceRanges`                           | Kafka service Load Balancer sources                                                                                                       | `[]`                    |
| `service.externalTrafficPolicy`                              | Kafka service external traffic policy                                                                                                     | `Cluster`               |
| `service.annotations`                                        | Additional custom annotations for Kafka service                                                                                           | `{}`                    |
| `service.headless.controller.annotations`                    | Annotations for the controller-eligible headless service.                                                                                 | `{}`                    |
| `service.headless.controller.labels`                         | Labels for the controller-eligible headless service.                                                                                      | `{}`                    |
| `service.headless.broker.annotations`                        | Annotations for the broker-only headless service.                                                                                         | `{}`                    |
| `service.headless.broker.labels`                             | Labels for the broker-only headless service.                                                                                              | `{}`                    |
| `externalAccess.enabled`                                     | Enable Kubernetes external cluster access to Kafka brokers                                                                                | `false`                 |
| `externalAccess.autoDiscovery.enabled`                       | Enable using an init container to auto-detect external IPs/ports by querying the K8s API                                                  | `false`                 |
| `externalAccess.autoDiscovery.image.registry`                | Init container auto-discovery image registry                                                                                              | `docker.io`             |
| `externalAccess.autoDiscovery.image.repository`              | Init container auto-discovery image repository                                                                                            | `bitnami/kubectl`       |
| `externalAccess.autoDiscovery.image.tag`                     | Init container auto-discovery image tag (immutable tags are recommended)                                                                  | `1.25.13-debian-11-r11` |
| `externalAccess.autoDiscovery.image.digest`                  | Kubectl image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                   | `""`                    |
| `externalAccess.autoDiscovery.image.pullPolicy`              | Init container auto-discovery image pull policy                                                                                           | `IfNotPresent`          |
| `externalAccess.autoDiscovery.image.pullSecrets`             | Init container auto-discovery image pull secrets                                                                                          | `[]`                    |
| `externalAccess.autoDiscovery.resources.limits`              | The resources limits for the auto-discovery init container                                                                                | `{}`                    |
| `externalAccess.autoDiscovery.resources.requests`            | The requested resources for the auto-discovery init container                                                                             | `{}`                    |
| `externalAccess.controller.forceExpose`                      | If set to true, force exposing controller-eligible nodes although they are configured as controller-only nodes                            | `false`                 |
| `externalAccess.controller.service.type`                     | Kubernetes Service type for external access. It can be NodePort, LoadBalancer or ClusterIP                                                | `LoadBalancer`          |
| `externalAccess.controller.service.ports.external`           | Kafka port used for external access when service type is LoadBalancer                                                                     | `9094`                  |
| `externalAccess.controller.service.loadBalancerIPs`          | Array of load balancer IPs for each Kafka broker. Length must be the same as replicaCount                                                 | `[]`                    |
| `externalAccess.controller.service.loadBalancerNames`        | Array of load balancer Names for each Kafka broker. Length must be the same as replicaCount                                               | `[]`                    |
| `externalAccess.controller.service.loadBalancerAnnotations`  | Array of load balancer annotations for each Kafka broker. Length must be the same as replicaCount                                         | `[]`                    |
| `externalAccess.controller.service.loadBalancerSourceRanges` | Address(es) that are allowed when service is LoadBalancer                                                                                 | `[]`                    |
| `externalAccess.controller.service.nodePorts`                | Array of node ports used for each Kafka broker. Length must be the same as replicaCount                                                   | `[]`                    |
| `externalAccess.controller.service.externalIPs`              | Use distinct service host IPs to configure Kafka external listener when service type is NodePort. Length must be the same as replicaCount | `[]`                    |
| `externalAccess.controller.service.useHostIPs`               | Use service host IPs to configure Kafka external listener when service type is NodePort                                                   | `false`                 |
| `externalAccess.controller.service.usePodIPs`                | using the MY_POD_IP address for external access.                                                                                          | `false`                 |
| `externalAccess.controller.service.domain`                   | Domain or external ip used to configure Kafka external listener when service type is NodePort or ClusterIP                                | `""`                    |
| `externalAccess.controller.service.publishNotReadyAddresses` | Indicates that any agent which deals with endpoints for this Service should disregard any indications of ready/not-ready                  | `false`                 |
| `externalAccess.controller.service.labels`                   | Service labels for external access                                                                                                        | `{}`                    |
| `externalAccess.controller.service.annotations`              | Service annotations for external access                                                                                                   | `{}`                    |
| `externalAccess.controller.service.extraPorts`               | Extra ports to expose in the Kafka external service                                                                                       | `[]`                    |
| `externalAccess.broker.service.type`                         | Kubernetes Service type for external access. It can be NodePort, LoadBalancer or ClusterIP                                                | `LoadBalancer`          |
| `externalAccess.broker.service.ports.external`               | Kafka port used for external access when service type is LoadBalancer                                                                     | `9094`                  |
| `externalAccess.broker.service.loadBalancerIPs`              | Array of load balancer IPs for each Kafka broker. Length must be the same as replicaCount                                                 | `[]`                    |
| `externalAccess.broker.service.loadBalancerNames`            | Array of load balancer Names for each Kafka broker. Length must be the same as replicaCount                                               | `[]`                    |
| `externalAccess.broker.service.loadBalancerAnnotations`      | Array of load balancer annotations for each Kafka broker. Length must be the same as replicaCount                                         | `[]`                    |
| `externalAccess.broker.service.loadBalancerSourceRanges`     | Address(es) that are allowed when service is LoadBalancer                                                                                 | `[]`                    |
| `externalAccess.broker.service.nodePorts`                    | Array of node ports used for each Kafka broker. Length must be the same as replicaCount                                                   | `[]`                    |
| `externalAccess.broker.service.externalIPs`                  | Use distinct service host IPs to configure Kafka external listener when service type is NodePort. Length must be the same as replicaCount | `[]`                    |
| `externalAccess.broker.service.useHostIPs`                   | Use service host IPs to configure Kafka external listener when service type is NodePort                                                   | `false`                 |
| `externalAccess.broker.service.usePodIPs`                    | using the MY_POD_IP address for external access.                                                                                          | `false`                 |
| `externalAccess.broker.service.domain`                       | Domain or external ip used to configure Kafka external listener when service type is NodePort or ClusterIP                                | `""`                    |
| `externalAccess.broker.service.publishNotReadyAddresses`     | Indicates that any agent which deals with endpoints for this Service should disregard any indications of ready/not-ready                  | `false`                 |
| `externalAccess.broker.service.labels`                       | Service labels for external access                                                                                                        | `{}`                    |
| `externalAccess.broker.service.annotations`                  | Service annotations for external access                                                                                                   | `{}`                    |
| `externalAccess.broker.service.extraPorts`                   | Extra ports to expose in the Kafka external service                                                                                       | `[]`                    |
| `networkPolicy.enabled`                                      | Specifies whether a NetworkPolicy should be created                                                                                       | `false`                 |
| `networkPolicy.allowExternal`                                | Don't require client label for connections                                                                                                | `true`                  |
| `networkPolicy.explicitNamespacesSelector`                   | A Kubernetes LabelSelector to explicitly select namespaces from which traffic could be allowed                                            | `{}`                    |
| `networkPolicy.externalAccess.from`                          | customize the from section for External Access on tcp-external port                                                                       | `[]`                    |
| `networkPolicy.egressRules.customRules`                      | Custom network policy rule                                                                                                                | `{}`                    |

### Volume Permissions parameters

| Name                                                   | Description                                                                                                                       | Value              |
| ------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------- | ------------------ |
| `volumePermissions.enabled`                            | Enable init container that changes the owner and group of the persistent volume                                                   | `false`            |
| `volumePermissions.image.registry`                     | Init container volume-permissions image registry                                                                                  | `docker.io`        |
| `volumePermissions.image.repository`                   | Init container volume-permissions image repository                                                                                | `bitnami/os-shell` |
| `volumePermissions.image.tag`                          | Init container volume-permissions image tag (immutable tags are recommended)                                                      | `11-debian-11-r60` |
| `volumePermissions.image.digest`                       | Init container volume-permissions image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag | `""`               |
| `volumePermissions.image.pullPolicy`                   | Init container volume-permissions image pull policy                                                                               | `IfNotPresent`     |
| `volumePermissions.image.pullSecrets`                  | Init container volume-permissions image pull secrets                                                                              | `[]`               |
| `volumePermissions.resources.limits`                   | Init container volume-permissions resource limits                                                                                 | `{}`               |
| `volumePermissions.resources.requests`                 | Init container volume-permissions resource requests                                                                               | `{}`               |
| `volumePermissions.containerSecurityContext.runAsUser` | User ID for the init container                                                                                                    | `0`                |

### Other Parameters

| Name                                          | Description                                                                                    | Value   |
| --------------------------------------------- | ---------------------------------------------------------------------------------------------- | ------- |
| `serviceAccount.create`                       | Enable creation of ServiceAccount for Kafka pods                                               | `true`  |
| `serviceAccount.name`                         | The name of the service account to use. If not set and `create` is `true`, a name is generated | `""`    |
| `serviceAccount.automountServiceAccountToken` | Allows auto mount of ServiceAccountToken on the serviceAccount created                         | `true`  |
| `serviceAccount.annotations`                  | Additional custom annotations for the ServiceAccount                                           | `{}`    |
| `rbac.create`                                 | Whether to create & use RBAC resources or not                                                  | `false` |

### Metrics parameters

| Name                                                              | Description                                                                                                                      | Value                                                                                   |
| ----------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------- |
| `metrics.kafka.enabled`                                           | Whether or not to create a standalone Kafka exporter to expose Kafka metrics                                                     | `false`                                                                                 |
| `metrics.kafka.image.registry`                                    | Kafka exporter image registry                                                                                                    | `docker.io`                                                                             |
| `metrics.kafka.image.repository`                                  | Kafka exporter image repository                                                                                                  | `bitnami/kafka-exporter`                                                                |
| `metrics.kafka.image.tag`                                         | Kafka exporter image tag (immutable tags are recommended)                                                                        | `1.7.0-debian-11-r102`                                                                  |
| `metrics.kafka.image.digest`                                      | Kafka exporter image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                   | `""`                                                                                    |
| `metrics.kafka.image.pullPolicy`                                  | Kafka exporter image pull policy                                                                                                 | `IfNotPresent`                                                                          |
| `metrics.kafka.image.pullSecrets`                                 | Specify docker-registry secret names as an array                                                                                 | `[]`                                                                                    |
| `metrics.kafka.certificatesSecret`                                | Name of the existing secret containing the optional certificate and key files                                                    | `""`                                                                                    |
| `metrics.kafka.tlsCert`                                           | The secret key from the certificatesSecret if 'client-cert' key different from the default (cert-file)                           | `cert-file`                                                                             |
| `metrics.kafka.tlsKey`                                            | The secret key from the certificatesSecret if 'client-key' key different from the default (key-file)                             | `key-file`                                                                              |
| `metrics.kafka.tlsCaSecret`                                       | Name of the existing secret containing the optional ca certificate for Kafka exporter client authentication                      | `""`                                                                                    |
| `metrics.kafka.tlsCaCert`                                         | The secret key from the certificatesSecret or tlsCaSecret if 'ca-cert' key different from the default (ca-file)                  | `ca-file`                                                                               |
| `metrics.kafka.extraFlags`                                        | Extra flags to be passed to Kafka exporter                                                                                       | `{}`                                                                                    |
| `metrics.kafka.command`                                           | Override Kafka exporter container command                                                                                        | `[]`                                                                                    |
| `metrics.kafka.args`                                              | Override Kafka exporter container arguments                                                                                      | `[]`                                                                                    |
| `metrics.kafka.containerPorts.metrics`                            | Kafka exporter metrics container port                                                                                            | `9308`                                                                                  |
| `metrics.kafka.resources.limits`                                  | The resources limits for the container                                                                                           | `{}`                                                                                    |
| `metrics.kafka.resources.requests`                                | The requested resources for the container                                                                                        | `{}`                                                                                    |
| `metrics.kafka.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                             | `true`                                                                                  |
| `metrics.kafka.podSecurityContext.fsGroup`                        | Set Kafka exporter pod's Security Context fsGroup                                                                                | `1001`                                                                                  |
| `metrics.kafka.podSecurityContext.seccompProfile.type`            | Set Kafka exporter pod's Security Context seccomp profile                                                                        | `RuntimeDefault`                                                                        |
| `metrics.kafka.containerSecurityContext.enabled`                  | Enable Kafka exporter containers' Security Context                                                                               | `true`                                                                                  |
| `metrics.kafka.containerSecurityContext.runAsUser`                | Set Kafka exporter containers' Security Context runAsUser                                                                        | `1001`                                                                                  |
| `metrics.kafka.containerSecurityContext.runAsNonRoot`             | Set Kafka exporter containers' Security Context runAsNonRoot                                                                     | `true`                                                                                  |
| `metrics.kafka.containerSecurityContext.allowPrivilegeEscalation` | Set Kafka exporter containers' Security Context allowPrivilegeEscalation                                                         | `false`                                                                                 |
| `metrics.kafka.containerSecurityContext.readOnlyRootFilesystem`   | Set Kafka exporter containers' Security Context readOnlyRootFilesystem                                                           | `true`                                                                                  |
| `metrics.kafka.containerSecurityContext.capabilities.drop`        | Set Kafka exporter containers' Security Context capabilities to be dropped                                                       | `["ALL"]`                                                                               |
| `metrics.kafka.hostAliases`                                       | Kafka exporter pods host aliases                                                                                                 | `[]`                                                                                    |
| `metrics.kafka.podLabels`                                         | Extra labels for Kafka exporter pods                                                                                             | `{}`                                                                                    |
| `metrics.kafka.podAnnotations`                                    | Extra annotations for Kafka exporter pods                                                                                        | `{}`                                                                                    |
| `metrics.kafka.podAffinityPreset`                                 | Pod affinity preset. Ignored if `metrics.kafka.affinity` is set. Allowed values: `soft` or `hard`                                | `""`                                                                                    |
| `metrics.kafka.podAntiAffinityPreset`                             | Pod anti-affinity preset. Ignored if `metrics.kafka.affinity` is set. Allowed values: `soft` or `hard`                           | `soft`                                                                                  |
| `metrics.kafka.nodeAffinityPreset.type`                           | Node affinity preset type. Ignored if `metrics.kafka.affinity` is set. Allowed values: `soft` or `hard`                          | `""`                                                                                    |
| `metrics.kafka.nodeAffinityPreset.key`                            | Node label key to match Ignored if `metrics.kafka.affinity` is set.                                                              | `""`                                                                                    |
| `metrics.kafka.nodeAffinityPreset.values`                         | Node label values to match. Ignored if `metrics.kafka.affinity` is set.                                                          | `[]`                                                                                    |
| `metrics.kafka.affinity`                                          | Affinity for pod assignment                                                                                                      | `{}`                                                                                    |
| `metrics.kafka.nodeSelector`                                      | Node labels for pod assignment                                                                                                   | `{}`                                                                                    |
| `metrics.kafka.tolerations`                                       | Tolerations for pod assignment                                                                                                   | `[]`                                                                                    |
| `metrics.kafka.schedulerName`                                     | Name of the k8s scheduler (other than default) for Kafka exporter                                                                | `""`                                                                                    |
| `metrics.kafka.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                            | `true`                                                                                  |
| `metrics.kafka.priorityClassName`                                 | Kafka exporter pods' priorityClassName                                                                                           | `""`                                                                                    |
| `metrics.kafka.topologySpreadConstraints`                         | Topology Spread Constraints for pod assignment                                                                                   | `[]`                                                                                    |
| `metrics.kafka.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka exporter pod(s)                                                | `[]`                                                                                    |
| `metrics.kafka.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka exporter container(s)                                     | `[]`                                                                                    |
| `metrics.kafka.sidecars`                                          | Add additional sidecar containers to the Kafka exporter pod(s)                                                                   | `[]`                                                                                    |
| `metrics.kafka.initContainers`                                    | Add init containers to the Kafka exporter pods                                                                                   | `[]`                                                                                    |
| `metrics.kafka.service.ports.metrics`                             | Kafka exporter metrics service port                                                                                              | `9308`                                                                                  |
| `metrics.kafka.service.clusterIP`                                 | Static clusterIP or None for headless services                                                                                   | `""`                                                                                    |
| `metrics.kafka.service.sessionAffinity`                           | Control where client requests go, to the same pod or round-robin                                                                 | `None`                                                                                  |
| `metrics.kafka.service.annotations`                               | Annotations for the Kafka exporter service                                                                                       | `{}`                                                                                    |
| `metrics.kafka.serviceAccount.create`                             | Enable creation of ServiceAccount for Kafka exporter pods                                                                        | `true`                                                                                  |
| `metrics.kafka.serviceAccount.name`                               | The name of the service account to use. If not set and `create` is `true`, a name is generated                                   | `""`                                                                                    |
| `metrics.kafka.serviceAccount.automountServiceAccountToken`       | Allows auto mount of ServiceAccountToken on the serviceAccount created                                                           | `true`                                                                                  |
| `metrics.jmx.enabled`                                             | Whether or not to expose JMX metrics to Prometheus                                                                               | `false`                                                                                 |
| `metrics.jmx.kafkaJmxPort`                                        | JMX port where the exporter will collect metrics, exposed in the Kafka container.                                                | `5555`                                                                                  |
| `metrics.jmx.image.registry`                                      | JMX exporter image registry                                                                                                      | `docker.io`                                                                             |
| `metrics.jmx.image.repository`                                    | JMX exporter image repository                                                                                                    | `bitnami/jmx-exporter`                                                                  |
| `metrics.jmx.image.tag`                                           | JMX exporter image tag (immutable tags are recommended)                                                                          | `0.19.0-debian-11-r66`                                                                  |
| `metrics.jmx.image.digest`                                        | JMX exporter image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                     | `""`                                                                                    |
| `metrics.jmx.image.pullPolicy`                                    | JMX exporter image pull policy                                                                                                   | `IfNotPresent`                                                                          |
| `metrics.jmx.image.pullSecrets`                                   | Specify docker-registry secret names as an array                                                                                 | `[]`                                                                                    |
| `metrics.jmx.containerSecurityContext.enabled`                    | Enable Prometheus JMX exporter containers' Security Context                                                                      | `true`                                                                                  |
| `metrics.jmx.containerSecurityContext.runAsUser`                  | Set Prometheus JMX exporter containers' Security Context runAsUser                                                               | `1001`                                                                                  |
| `metrics.jmx.containerSecurityContext.runAsNonRoot`               | Set Prometheus JMX exporter containers' Security Context runAsNonRoot                                                            | `true`                                                                                  |
| `metrics.jmx.containerSecurityContext.allowPrivilegeEscalation`   | Set Prometheus JMX exporter containers' Security Context allowPrivilegeEscalation                                                | `false`                                                                                 |
| `metrics.jmx.containerSecurityContext.readOnlyRootFilesystem`     | Set Prometheus JMX exporter containers' Security Context readOnlyRootFilesystem                                                  | `true`                                                                                  |
| `metrics.jmx.containerSecurityContext.capabilities.drop`          | Set Prometheus JMX exporter containers' Security Context capabilities to be dropped                                              | `["ALL"]`                                                                               |
| `metrics.jmx.containerPorts.metrics`                              | Prometheus JMX exporter metrics container port                                                                                   | `5556`                                                                                  |
| `metrics.jmx.resources.limits`                                    | The resources limits for the JMX exporter container                                                                              | `{}`                                                                                    |
| `metrics.jmx.resources.requests`                                  | The requested resources for the JMX exporter container                                                                           | `{}`                                                                                    |
| `metrics.jmx.service.ports.metrics`                               | Prometheus JMX exporter metrics service port                                                                                     | `5556`                                                                                  |
| `metrics.jmx.service.clusterIP`                                   | Static clusterIP or None for headless services                                                                                   | `""`                                                                                    |
| `metrics.jmx.service.sessionAffinity`                             | Control where client requests go, to the same pod or round-robin                                                                 | `None`                                                                                  |
| `metrics.jmx.service.annotations`                                 | Annotations for the Prometheus JMX exporter service                                                                              | `{}`                                                                                    |
| `metrics.jmx.whitelistObjectNames`                                | Allows setting which JMX objects you want to expose to via JMX stats to JMX exporter                                             | `["kafka.controller:*","kafka.server:*","java.lang:*","kafka.network:*","kafka.log:*"]` |
| `metrics.jmx.config`                                              | Configuration file for JMX exporter                                                                                              | `""`                                                                                    |
| `metrics.jmx.existingConfigmap`                                   | Name of existing ConfigMap with JMX exporter configuration                                                                       | `""`                                                                                    |
| `metrics.jmx.extraRules`                                          | Add extra rules to JMX exporter configuration                                                                                    | `""`                                                                                    |
| `metrics.serviceMonitor.enabled`                                  | if `true`, creates a Prometheus Operator ServiceMonitor (requires `metrics.kafka.enabled` or `metrics.jmx.enabled` to be `true`) | `false`                                                                                 |
| `metrics.serviceMonitor.namespace`                                | Namespace in which Prometheus is running                                                                                         | `""`                                                                                    |
| `metrics.serviceMonitor.interval`                                 | Interval at which metrics should be scraped                                                                                      | `""`                                                                                    |
| `metrics.serviceMonitor.scrapeTimeout`                            | Timeout after which the scrape is ended                                                                                          | `""`                                                                                    |
| `metrics.serviceMonitor.labels`                                   | Additional labels that can be used so ServiceMonitor will be discovered by Prometheus                                            | `{}`                                                                                    |
| `metrics.serviceMonitor.selector`                                 | Prometheus instance selector labels                                                                                              | `{}`                                                                                    |
| `metrics.serviceMonitor.relabelings`                              | RelabelConfigs to apply to samples before scraping                                                                               | `[]`                                                                                    |
| `metrics.serviceMonitor.metricRelabelings`                        | MetricRelabelConfigs to apply to samples before ingestion                                                                        | `[]`                                                                                    |
| `metrics.serviceMonitor.honorLabels`                              | Specify honorLabels parameter to add the scrape endpoint                                                                         | `false`                                                                                 |
| `metrics.serviceMonitor.jobLabel`                                 | The name of the label on the target service to use as the job name in prometheus.                                                | `""`                                                                                    |
| `metrics.prometheusRule.enabled`                                  | if `true`, creates a Prometheus Operator PrometheusRule (requires `metrics.kafka.enabled` or `metrics.jmx.enabled` to be `true`) | `false`                                                                                 |
| `metrics.prometheusRule.namespace`                                | Namespace in which Prometheus is running                                                                                         | `""`                                                                                    |
| `metrics.prometheusRule.labels`                                   | Additional labels that can be used so PrometheusRule will be discovered by Prometheus                                            | `{}`                                                                                    |
| `metrics.prometheusRule.groups`                                   | Prometheus Rule Groups for Kafka                                                                                                 | `[]`                                                                                    |

### Kafka provisioning parameters

| Name                                                             | Description                                                                                                                   | Value                 |
| ---------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------- | --------------------- |
| `provisioning.enabled`                                           | Enable kafka provisioning Job                                                                                                 | `false`               |
| `provisioning.numPartitions`                                     | Default number of partitions for topics when unspecified                                                                      | `1`                   |
| `provisioning.replicationFactor`                                 | Default replication factor for topics when unspecified                                                                        | `1`                   |
| `provisioning.topics`                                            | Kafka topics to provision                                                                                                     | `[]`                  |
| `provisioning.nodeSelector`                                      | Node labels for pod assignment                                                                                                | `{}`                  |
| `provisioning.tolerations`                                       | Tolerations for pod assignment                                                                                                | `[]`                  |
| `provisioning.extraProvisioningCommands`                         | Extra commands to run to provision cluster resources                                                                          | `[]`                  |
| `provisioning.parallel`                                          | Number of provisioning commands to run at the same time                                                                       | `1`                   |
| `provisioning.preScript`                                         | Extra bash script to run before topic provisioning. $CLIENT_CONF is path to properties file with most needed configurations   | `""`                  |
| `provisioning.postScript`                                        | Extra bash script to run after topic provisioning. $CLIENT_CONF is path to properties file with most needed configurations    | `""`                  |
| `provisioning.auth.tls.type`                                     | Format to use for TLS certificates. Allowed types: `JKS` and `PEM`.                                                           | `jks`                 |
| `provisioning.auth.tls.certificatesSecret`                       | Existing secret containing the TLS certificates for the Kafka provisioning Job.                                               | `""`                  |
| `provisioning.auth.tls.cert`                                     | The secret key from the certificatesSecret if 'cert' key different from the default (tls.crt)                                 | `tls.crt`             |
| `provisioning.auth.tls.key`                                      | The secret key from the certificatesSecret if 'key' key different from the default (tls.key)                                  | `tls.key`             |
| `provisioning.auth.tls.caCert`                                   | The secret key from the certificatesSecret if 'caCert' key different from the default (ca.crt)                                | `ca.crt`              |
| `provisioning.auth.tls.keystore`                                 | The secret key from the certificatesSecret if 'keystore' key different from the default (keystore.jks)                        | `keystore.jks`        |
| `provisioning.auth.tls.truststore`                               | The secret key from the certificatesSecret if 'truststore' key different from the default (truststore.jks)                    | `truststore.jks`      |
| `provisioning.auth.tls.passwordsSecret`                          | Name of the secret containing passwords to access the JKS files or PEM key when they are password-protected.                  | `""`                  |
| `provisioning.auth.tls.keyPasswordSecretKey`                     | The secret key from the passwordsSecret if 'keyPasswordSecretKey' key different from the default (key-password)               | `key-password`        |
| `provisioning.auth.tls.keystorePasswordSecretKey`                | The secret key from the passwordsSecret if 'keystorePasswordSecretKey' key different from the default (keystore-password)     | `keystore-password`   |
| `provisioning.auth.tls.truststorePasswordSecretKey`              | The secret key from the passwordsSecret if 'truststorePasswordSecretKey' key different from the default (truststore-password) | `truststore-password` |
| `provisioning.auth.tls.keyPassword`                              | Password to access the password-protected PEM key if necessary. Ignored if 'passwordsSecret' is provided.                     | `""`                  |
| `provisioning.auth.tls.keystorePassword`                         | Password to access the JKS keystore. Ignored if 'passwordsSecret' is provided.                                                | `""`                  |
| `provisioning.auth.tls.truststorePassword`                       | Password to access the JKS truststore. Ignored if 'passwordsSecret' is provided.                                              | `""`                  |
| `provisioning.command`                                           | Override provisioning container command                                                                                       | `[]`                  |
| `provisioning.args`                                              | Override provisioning container arguments                                                                                     | `[]`                  |
| `provisioning.extraEnvVars`                                      | Extra environment variables to add to the provisioning pod                                                                    | `[]`                  |
| `provisioning.extraEnvVarsCM`                                    | ConfigMap with extra environment variables                                                                                    | `""`                  |
| `provisioning.extraEnvVarsSecret`                                | Secret with extra environment variables                                                                                       | `""`                  |
| `provisioning.podAnnotations`                                    | Extra annotations for Kafka provisioning pods                                                                                 | `{}`                  |
| `provisioning.podLabels`                                         | Extra labels for Kafka provisioning pods                                                                                      | `{}`                  |
| `provisioning.serviceAccount.create`                             | Enable creation of ServiceAccount for Kafka provisioning pods                                                                 | `false`               |
| `provisioning.serviceAccount.name`                               | The name of the service account to use. If not set and `create` is `true`, a name is generated                                | `""`                  |
| `provisioning.serviceAccount.automountServiceAccountToken`       | Allows auto mount of ServiceAccountToken on the serviceAccount created                                                        | `true`                |
| `provisioning.resources.limits`                                  | The resources limits for the Kafka provisioning container                                                                     | `{}`                  |
| `provisioning.resources.requests`                                | The requested resources for the Kafka provisioning container                                                                  | `{}`                  |
| `provisioning.podSecurityContext.enabled`                        | Enable security context for the pods                                                                                          | `true`                |
| `provisioning.podSecurityContext.fsGroup`                        | Set Kafka provisioning pod's Security Context fsGroup                                                                         | `1001`                |
| `provisioning.podSecurityContext.seccompProfile.type`            | Set Kafka provisioning pod's Security Context seccomp profile                                                                 | `RuntimeDefault`      |
| `provisioning.containerSecurityContext.enabled`                  | Enable Kafka provisioning containers' Security Context                                                                        | `true`                |
| `provisioning.containerSecurityContext.runAsUser`                | Set Kafka provisioning containers' Security Context runAsUser                                                                 | `1001`                |
| `provisioning.containerSecurityContext.runAsNonRoot`             | Set Kafka provisioning containers' Security Context runAsNonRoot                                                              | `true`                |
| `provisioning.containerSecurityContext.allowPrivilegeEscalation` | Set Kafka provisioning containers' Security Context allowPrivilegeEscalation                                                  | `false`               |
| `provisioning.containerSecurityContext.readOnlyRootFilesystem`   | Set Kafka provisioning containers' Security Context readOnlyRootFilesystem                                                    | `true`                |
| `provisioning.containerSecurityContext.capabilities.drop`        | Set Kafka provisioning containers' Security Context capabilities to be dropped                                                | `["ALL"]`             |
| `provisioning.schedulerName`                                     | Name of the k8s scheduler (other than default) for kafka provisioning                                                         | `""`                  |
| `provisioning.enableServiceLinks`                                | Whether information about services should be injected into pod's environment variable                                         | `true`                |
| `provisioning.extraVolumes`                                      | Optionally specify extra list of additional volumes for the Kafka provisioning pod(s)                                         | `[]`                  |
| `provisioning.extraVolumeMounts`                                 | Optionally specify extra list of additional volumeMounts for the Kafka provisioning container(s)                              | `[]`                  |
| `provisioning.sidecars`                                          | Add additional sidecar containers to the Kafka provisioning pod(s)                                                            | `[]`                  |
| `provisioning.initContainers`                                    | Add additional Add init containers to the Kafka provisioning pod(s)                                                           | `[]`                  |
| `provisioning.waitForKafka`                                      | If true use an init container to wait until kafka is ready before starting provisioning                                       | `true`                |

### KRaft chart parameters

| Name                           | Description                                                                                                                                                   | Value  |
| ------------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------ |
| `kraft.enabled`                | Switch to enable or disable the KRaft mode for Kafka                                                                                                          | `true` |
| `kraft.clusterId`              | Kafka Kraft cluster ID. If not set, a random cluster ID will be generated the first time Kraft is initialized.                                                | `""`   |
| `kraft.controllerQuorumVoters` | Override the Kafka controller quorum voters of the Kafka Kraft cluster. If not set, it will be automatically configured to use all controller-elegible nodes. | `""`   |

### ZooKeeper chart parameters

| Name                                    | Description                                                                                                                                                             | Value               |
| --------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------- |
| `zookeeperChrootPath`                   | Path which puts data under some path in the global ZooKeeper namespace                                                                                                  | `""`                |
| `zookeeper.enabled`                     | Switch to enable or disable the ZooKeeper helm chart. Must be false if you use KRaft mode.                                                                              | `false`             |
| `zookeeper.replicaCount`                | Number of ZooKeeper nodes                                                                                                                                               | `1`                 |
| `zookeeper.auth.client.enabled`         | Enable ZooKeeper auth                                                                                                                                                   | `false`             |
| `zookeeper.auth.client.clientUser`      | User that will use ZooKeeper client (zkCli.sh) to authenticate. Must exist in the serverUsers comma-separated list.                                                     | `""`                |
| `zookeeper.auth.client.clientPassword`  | Password that will use ZooKeeper client (zkCli.sh) to authenticate. Must exist in the serverPasswords comma-separated list.                                             | `""`                |
| `zookeeper.auth.client.serverUsers`     | Comma, semicolon or whitespace separated list of user to be created. Specify them as a string, for example: "user1,user2,admin"                                         | `""`                |
| `zookeeper.auth.client.serverPasswords` | Comma, semicolon or whitespace separated list of passwords to assign to users when created. Specify them as a string, for example: "pass4user1, pass4user2, pass4admin" | `""`                |
| `zookeeper.persistence.enabled`         | Enable persistence on ZooKeeper using PVC(s)                                                                                                                            | `true`              |
| `zookeeper.persistence.storageClass`    | Persistent Volume storage class                                                                                                                                         | `""`                |
| `zookeeper.persistence.accessModes`     | Persistent Volume access modes                                                                                                                                          | `["ReadWriteOnce"]` |
| `zookeeper.persistence.size`            | Persistent Volume size                                                                                                                                                  | `8Gi`               |
| `externalZookeeper.servers`             | List of external zookeeper servers to use. Typically used in combination with 'zookeeperChrootPath'. Must be empty if you use KRaft mode.                               | `[]`                |

```console
helm install my-release \
  --set replicaCount=3 \
  oci://registry-1.docker.io/bitnamicharts/kafka
```

The above command deploys Kafka with 3 brokers (replicas).

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```console
helm install my-release -f values.yaml oci://registry-1.docker.io/bitnamicharts/kafka
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Configuration and installation details

### [Rolling VS Immutable tags](https://docs.bitnami.com/containers/how-to/understand-rolling-tags-containers/)

It is strongly recommended to use immutable tags in a production environment. This ensures your deployment does not change automatically if the same tag is updated with a different image.

Bitnami will release a new chart updating its containers if a new version of the main container, significant changes, or critical vulnerabilities exist.

### Setting custom parameters

Any environment variable beginning with `KAFKA_CFG_` will be mapped to its corresponding Kafka key. For example, use `KAFKA_CFG_BACKGROUND_THREADS` in order to set `background.threads`. In order to pass custom environment variables use the `extraEnvVars` property.

Using `extraEnvVars` with `KAFKA_CFG_` is the preferred and simplest way to add custom Kafka parameters not otherwise specified in this chart. Alternatively, you can provide a *full* Kafka configuration using `config` or `existingConfigmap`.
Setting either `config` or `existingConfigmap` will cause the chart to disregard `KAFKA_CFG_` settings, which are used by many other Kafka-related chart values described above, as well as dynamically generated parameters such as `zookeeper.connect`. This can cause unexpected behavior.

### Listeners configuration

This chart allows you to automatically configure Kafka with 3 listeners:

- One for inter-broker communications.
- A second one for communications with clients within the K8s cluster.
- (optional) a third listener for communications with clients outside the K8s cluster. Check [this section](#accessing-kafka-brokers-from-outside-the-cluster) for more information.

For more complex configurations, set the `listeners`, `advertisedListeners` and `listenerSecurityProtocolMap` parameters as needed.

### Enable security for Kafka and Zookeeper

You can configure different authentication protocols for each listener you configure in Kafka. For instance, you can use `sasl_tls` authentication for client communications, while using `tls` for inter-broker communications. This table shows the available protocols and the security they provide:

| Method    | Authentication               | Encryption via TLS |
|-----------|------------------------------|--------------------|
| plaintext | None                         | No                 |
| tls       | None                         | Yes                |
| mtls      | Yes (two-way authentication) | Yes                |
| sasl      | Yes (via SASL)               | No                 |
| sasl_tls  | Yes (via SASL)               | Yes                |

Learn more about how to configure Kafka to use the different authentication protocols in the [chart documentation](https://docs.bitnami.com/kubernetes/infrastructure/kafka/administration/enable-security/).

If you enabled SASL authentication on any listener, you can set the SASL credentials using the parameters below:

- `auth.sasl.jaas.clientUsers`/`auth.sasl.jaas.clientPasswords`: when enabling SASL authentication for communications with clients.
- `auth.sasl.jaas.interBrokerUser`/`auth.sasl.jaas.interBrokerPassword`:  when enabling SASL authentication for inter-broker communications.
- `auth.jaas.zookeeperUser`/`auth.jaas.zookeeperPassword`: In the case that the Zookeeper chart is deployed with SASL authentication enabled.

In order to configure TLS authentication/encryption, you **can** create a secret per Kafka broker you have in the cluster containing the Java Key Stores (JKS) files: the truststore (`kafka.truststore.jks`) and the keystore (`kafka.keystore.jks`). Then, you need pass the secret names with the `tls.existingSecret` parameter when deploying the chart.

> **Note**: If the JKS files are password protected (recommended), you will need to provide the password to get access to the keystores. To do so, use the `tls.password` parameter to provide your password.

For instance, to configure TLS authentication on a Kafka cluster with 2 Kafka brokers use the commands below to create the secrets:

```console
kubectl create secret generic kafka-jks-0 --from-file=kafka.truststore.jks=./kafka.truststore.jks --from-file=kafka.keystore.jks=./kafka-0.keystore.jks
kubectl create secret generic kafka-jks-1 --from-file=kafka.truststore.jks=./kafka.truststore.jks --from-file=kafka.keystore.jks=./kafka-1.keystore.jks
```

> **Note**: the command above assumes you already created the truststore and keystores files. This [script](https://raw.githubusercontent.com/confluentinc/confluent-platform-security-tools/master/kafka-generate-ssl.sh) can help you with the JKS files generation.

If, for some reason (like using Cert-Manager) you can not use the default JKS secret scheme, you can use the additional parameters:

- `tls.jksTruststoreSecret` to define additional secret, where the `kafka.truststore.jks` is being kept. The truststore password **must** be the same as in `tls.password`
- `tls.jksTruststore` to overwrite the default value of the truststore key (`kafka.truststore.jks`).

> **Note**: If you are using cert-manager, particularly when an ACME issuer is used, the `ca.crt` field is not put in the `Secret` that cert-manager creates. To handle this, the `tls.pemChainIncluded` property can be set to `true` and the initContainer created by this Chart will attempt to extract the intermediate certs from the `tls.crt` field of the secret (which is a PEM chain)
> **Note**: The truststore/keystore from above **must** be protected with the same password as in `tls.password`

You can deploy the chart with authentication using the following parameters:

```console
replicaCount=2
listeners.client.client.protocol=SASL
listeners.client.interbroker.protocol=TLS
tls.existingSecret=kafka-jks
tls.password=jksPassword
sasl.client.users[0]=brokerUser
sasl.client.passwords[0]=brokerPassword
sasl.zookeeper.user=zookeeperUser
sasl.zookeeper.password=zookeeperPassword
zookeeper.auth.enabled=true
zookeeper.auth.serverUsers=zookeeperUser
zookeeper.auth.serverPasswords=zookeeperPassword
zookeeper.auth.clientUser=zookeeperUser
zookeeper.auth.clientPassword=zookeeperPassword
```

You can deploy the chart with AclAuthorizer using the following parameters:

```console
replicaCount=2
listeners.client.protocol=SASL
listeners.interbroker.protocol=SASL_TLS
tls.existingSecret=kafka-jks-0
tls.password=jksPassword
sasl.client.users[0]=brokerUser
sasl.client.passwords[0]=brokerPassword
sasl.zookeeper.user=zookeeperUser
sasl.zookeeper.password=zookeeperPassword
zookeeper.auth.enabled=true
zookeeper.auth.serverUsers=zookeeperUser
zookeeper.auth.serverPasswords=zookeeperPassword
zookeeper.auth.clientUser=zookeeperUser
zookeeper.auth.clientPassword=zookeeperPassword
authorizerClassName=kafka.security.authorizer.AclAuthorizer
allowEveryoneIfNoAclFound=false
superUsers=User:admin
```

If you are using Kafka ACLs, you might encounter in kafka-authorizer.log the following event: `[...] Principal = User:ANONYMOUS is Allowed Operation [...]`.

By setting the following parameter: `listeners.client.protocol=SSL` and `listener.client.sslClientAuth=required`, Kafka will require the clients to authenticate to Kafka brokers via certificate.

As result, we will be able to see in kafka-authorizer.log the events specific Subject: `[...] Principal = User:CN=kafka,OU=...,O=...,L=...,C=..,ST=... is [...]`.

If you also enable exposing metrics using the Kafka exporter, and you are using `SSL` or `SASL_SSL` security protocols protocols, you need to mount the CA certificated used to sign the brokers certificates in the exporter so it can validate the Kafka brokers. To do so, create a secret containing the CA, and set the `metrics.certificatesSecret` parameter. As an alternative, you can skip TLS validation using extra flags:

```console
metrics.kafka.extraFlags={tls.insecure-skip-tls-verify: ""}
```

### Accessing Kafka brokers from outside the cluster

In order to access Kafka Brokers from outside the cluster, an additional listener and advertised listener must be configured. Additionally, a specific service per kafka pod will be created.

There are three ways of configuring external access. Using LoadBalancer services, using NodePort services or using ClusterIP services.

#### Using LoadBalancer services

You have two alternatives to use LoadBalancer services:

- Option A) Use random load balancer IPs using an **initContainer** that waits for the IPs to be ready and discover them automatically.

```console
externalAccess.enabled=true
externalAccess.service.broker.type=LoadBalancer
externalAccess.service.controller.type=LoadBalancer
externalAccess.service.broker.ports.external=9094
externalAccess.service.controller.containerPorts.external=9094
externalAccess.autoDiscovery.enabled=true
serviceAccount.create=true
rbac.create=true
```

Note: This option requires creating RBAC rules on clusters where RBAC policies are enabled.

- Option B) Manually specify the load balancer IPs:

```console
externalAccess.enabled=true
externalAccess.service.controller.type=LoadBalancer
externalAccess.service.controller.containerPorts.external=9094
externalAccess.service.controller.loadBalancerIPs[0]='external-ip-1'
externalAccess.service.controller.loadBalancerIPs[1]='external-ip-2'
externalAccess.service.broker.type=LoadBalancer
externalAccess.service.broker.ports.external=9094
externalAccess.service.broker.loadBalancerIPs[0]='external-ip-3'
externalAccess.service.broker.loadBalancerIPs[1]='external-ip-4'
```

Note: You need to know in advance the load balancer IPs so each Kafka broker advertised listener is configured with it.

Following the aforementioned steps will also allow to connect the brokers from the outside using the cluster's default service (when `service.type` is `LoadBalancer` or `NodePort`). Use the property `service.externalPort` to specify the port used for external connections.

#### Using NodePort services

You have two alternatives to use NodePort services:

- Option A) Use random node ports using an **initContainer** that discover them automatically.

  ```console
  externalAccess.enabled=true
  externalAccess.controller.service.type=NodePort
  externalAccess.broker.service.type=NodePort
  externalAccess.autoDiscovery.enabled=true
  serviceAccount.create=true
  rbac.create=true
  ```

  Note: This option requires creating RBAC rules on clusters where RBAC policies are enabled.

- Option B) Manually specify the node ports:

  ```console
  externalAccess.enabled=true
  externalAccess.controller.service.type=NodePort
  externalAccess.controller.service.nodePorts[0]='node-port-1'
  externalAccess.controller.service.nodePorts[1]='node-port-2'
  ```

  Note: You need to know in advance the node ports that will be exposed so each Kafka broker advertised listener is configured with it.

  The pod will try to get the external ip of the node using `curl -s https://ipinfo.io/ip` unless `externalAccess.service.domain` or `externalAccess.service.useHostIPs` is provided.

- Option C) Manually specify distinct external IPs (using controller+broker nodes)

  ```console
  externalAccess.enabled=true
  externalAccess.controller.service.type=NodePort
  externalAccess.controller.service.externalIPs[0]='172.16.0.20'
  externalAccess.controller.service.externalIPs[1]='172.16.0.21'
  externalAccess.controller.service.externalIPs[2]='172.16.0.22'
  ```

  Note: You need to know in advance the available IP of your cluster that will be exposed so each Kafka broker advertised listener is configured with it.

#### Using ClusterIP services

Note: This option requires that an ingress is deployed within your cluster

```console
externalAccess.enabled=true
externalAccess.controller.service.type=ClusterIP
externalAccess.controller.service.ports.external=9094
externalAccess.controller.service.domain='ingress-ip'
externalAccess.broker.service.type=ClusterIP
externalAccess.broker.service.ports.external=9094
externalAccess.broker.service.domain='ingress-ip'
```

Note: the deployed ingress must contain the following block:

```console
tcp:
  9094: "{{ include "common.names.namespace" . }}/{{ include "common.names.fullname" . }}-0-external:9094"
  9095: "{{ include "common.names.namespace" . }}/{{ include "common.names.fullname" . }}-1-external:9094"
  9096: "{{ include "common.names.namespace" . }}/{{ include "common.names.fullname" . }}-2-external:9094"
```

#### Name resolution with External-DNS

You can use the following values to generate External-DNS annotations which automatically creates DNS records for each ReplicaSet pod:

```yaml
externalAccess:
  service:
    annotations:
      external-dns.alpha.kubernetes.io/hostname: "{{ .targetPod }}.example.com"
```

### Sidecars

If you have a need for additional containers to run within the same pod as Kafka (e.g. an additional metrics or logging exporter), you can do so via the `sidecars` config parameter. Simply define your container according to the Kubernetes container spec.

```yaml
sidecars:
  - name: your-image-name
    image: your-image
    imagePullPolicy: Always
    ports:
      - name: portname
       containerPort: 1234
```

### Setting Pod's affinity

This chart allows you to set your custom affinity using the `affinity` parameter. Find more information about Pod's affinity in the [kubernetes documentation](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity).

As an alternative, you can use of the preset configurations for pod affinity, pod anti-affinity, and node affinity available at the [bitnami/common](https://github.com/bitnami/charts/tree/main/bitnami/common#affinities) chart. To do so, set the `podAffinityPreset`, `podAntiAffinityPreset`, or `nodeAffinityPreset` parameters.

### Deploying extra resources

There are cases where you may want to deploy extra objects, such as Kafka Connect. For covering this case, the chart allows adding the full specification of other objects using the `extraDeploy` parameter. The following example would create a deployment including a Kafka Connect deployment so you can connect Kafka with MongoDB&reg;:

```yaml
## Extra objects to deploy (value evaluated as a template)
##
extraDeploy:
  - |
    apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: {{ include "common.names.fullname" . }}-connect
      labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
    spec:
      replicas: 1
      selector:
        matchLabels: {{- include "common.labels.matchLabels" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 6 }}
          app.kubernetes.io/component: connector
      template:
        metadata:
          labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 8 }}
            app.kubernetes.io/component: connector
        spec:
          containers:
            - name: connect
              image: KAFKA-CONNECT-IMAGE
              imagePullPolicy: IfNotPresent
              ports:
                - name: connector
                  containerPort: 8083
              volumeMounts:
                - name: configuration
                  mountPath: /bitnami/kafka/config
          volumes:
            - name: configuration
              configMap:
                name: {{ include "common.names.fullname" . }}-connect
  - |
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: {{ include "common.names.fullname" . }}-connect
      labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
    data:
      connect-standalone.properties: |-
        bootstrap.servers = {{ include "common.names.fullname" . }}-0.{{ include "common.names.fullname" . }}-headless.{{ include "common.names.namespace" . }}.svc.{{ .Values.clusterDomain }}:{{ .Values.service.port }}
        ...
      mongodb.properties: |-
        connection.uri=mongodb://root:password@mongodb-hostname:27017
        ...
  - |
    apiVersion: v1
    kind: Service
    metadata:
      name: {{ include "common.names.fullname" . }}-connect
      labels: {{- include "common.labels.standard" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
    spec:
      ports:
        - protocol: TCP
          port: 8083
          targetPort: connector
      selector: {{- include "common.labels.matchLabels" ( dict "customLabels" .Values.commonLabels "context" $ ) | nindent 4 }}
        app.kubernetes.io/component: connector
```

You can create the Kafka Connect image using the Dockerfile below:

```Dockerfile
FROM bitnami/kafka:latest
# Download MongoDB&reg; Connector for Apache Kafka https://www.confluent.io/hub/mongodb/kafka-connect-mongodb
RUN mkdir -p /opt/bitnami/kafka/plugins && \
    cd /opt/bitnami/kafka/plugins && \
    curl --remote-name --location --silent https://search.maven.org/remotecontent?filepath=org/mongodb/kafka/mongo-kafka-connect/1.2.0/mongo-kafka-connect-1.2.0-all.jar
CMD /opt/bitnami/kafka/bin/connect-standalone.sh /opt/bitnami/kafka/config/connect-standalone.properties /opt/bitnami/kafka/config/mongo.properties
```

## Persistence

The [Bitnami Kafka](https://github.com/bitnami/containers/tree/main/bitnami/kafka) image stores the Kafka data at the `/bitnami/kafka` path of the container. Persistent Volume Claims are used to keep the data across deployments. This is known to work in GCE, AWS, and minikube.

### Adjust permissions of persistent volume mountpoint

As the image run as non-root by default, it is necessary to adjust the ownership of the persistent volume so that the container can write data into it.

By default, the chart is configured to use Kubernetes Security Context to automatically change the ownership of the volume. However, this feature does not work in all Kubernetes distributions.
As an alternative, this chart supports using an initContainer to change the ownership of the volume before mounting it in the final destination.

You can enable this initContainer by setting `volumePermissions.enabled` to `true`.

## Troubleshooting

Find more information about how to deal with common errors related to Bitnami's Helm charts in [this troubleshooting guide](https://docs.bitnami.com/general/how-to/troubleshoot-helm-chart-issues).

## Migrating from Zookeeper (Early access)

This guide is an adaptation from upstream documentation: [Migrate from ZooKeeper to KRaft](https://docs.confluent.io/platform/current/installation/migrate-zk-kraft.html)

1. Retrieve the cluster ID from Zookeeper:

    ```console
    $ kubectl exec -it <your-zookeeper-pod> -- zkCli.sh get /cluster/id
    /opt/bitnami/java/bin/java
    Connecting to localhost:2181

    WATCHER::

    WatchedEvent state:SyncConnected type:None path:null
    {"version":"1","id":"TEr3HVPvTqSWixWRHngP5g"}
    ```

2. Deploy at least one Kraft controller-only in your deployment and enable `zookeeperMigrationMode=true`. The Kraft controllers will migrate the data from your Kafka ZkBroker to Kraft mode.

    To do so add the following values to your Zookeeper deployment when upgrading:

    ```yaml
    controller:
      replicaCount: 1
      controllerOnly: true
      zookeeperMigrationMode: true
      # If needed, set controllers minID to avoid conflict with your ZK brokers' ids.
      # minID: 0
    broker:
      zookeeperMigrationMode: true
    kraft:
      enabled: true
      clusterId: "<your_cluster_id>"
    ```

3. Wait until until all brokers are ready. You should see the following log in the broker logs:

    ```console
    INFO [KafkaServer id=100] Finished catching up on KRaft metadata log, requesting that the KRaft controller unfence this broker (kafka.server.KafkaServer)
    INFO [BrokerLifecycleManager id=100 isZkBroker=true] The broker has been unfenced. Transitioning from RECOVERY to RUNNING. (kafka.server.BrokerLifecycleManager)
    ```

    In the controllers, the following message should show up:

    ```console
    Transitioning ZK migration state from PRE_MIGRATION to MIGRATION (org.apache.kafka.controller.FeatureControlManager)
    ```

4. Once all brokers have been successfully migrated, set `broker.zookeeperMigrationMode=false` to fully migrate them.

    ```yaml
    broker:
      zookeeperMigrationMode: false
    ```

5. To conclude the migration, switch off migration mode on controllers and stop Zookeeper:

    ```yaml
    controller:
      zookeeperMigrationMode: false
    zookeeper:
      enabled: false
    ```

    After migration is complete, you should see the following message in your controllers:

    ```console
    [2023-07-13 13:07:45,226] INFO [QuorumController id=1] Transitioning ZK migration state from MIGRATION to POST_MIGRATION (org.apache.kafka.controller.FeatureControlManager)
    ```

6. (**Optional**) If you would like to switch to a non-dedicated cluster, set `controller.controllerOnly=false`. This will cause controller-only nodes to switch to controller+broker nodes.

    At that point, you could manually decommission broker-only nodes by reassigning its partitions to controller-eligible nodes.

    For more information about decommissioning kafka broker check the [Kafka documentation](https://www.confluent.io/blog/remove-kafka-brokers-from-any-cluster-the-easy-way/).

## Upgrading

### To 25.0.0

This major updates the Zookeeper subchart to it newest major, 12.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-1200).

### To 24.0.0

This major version is a refactor of the Kafka chart and its architecture, to better adapt to Kraft features introduced in version 22.0.0.

The changes introduced in this version are:

- New architecture. The chart now has two statefulsets, one for controller-eligible nodes (controller or controller+broker) and another one for broker-only nodes. Please take a look at the subsections [Upgrading from Kraft mode](#upgrading-from-kraft-mode) and [Upgrading from Zookeeper mode](#upgrading-from-zookeeper-mode) for more information about how to upgrade this chart depending on which mode you were using.

  The new architecture is designed to support two main features:
  - Deployment of dedicated nodes
  - Support for Zookeeper to Kraft migration

- Adds compatibility with `securityContext.readOnlyRootFs=true`, which is now the execution default.
  - The Kafka configuration is now mounted as a ConfigMap instead of generated at runtime.
  - Due to the implementation of readOnlyRootFs support, the following settings have been removed and will now rely on Kafka defaults. To override them, please use `extraConfig` to extend your Kafka configuration instead.
    - `deleteTopicEnable`
    - `autoCreateTopicsEnable`
    - `logFlushIntervalMessages`
    - `logFlushIntervalMs`
    - `logRetentionBytes`
    - `logRetentionCheckIntervalMs`
    - `logRetentionHours`
    - `logSegmentBytes`
    - `logsDirs`
    - `maxMessageBytes`
    - `defaultReplicationFactor`
    - `offsetsTopicReplicationFactor`
    - `transactionStateLogReplicationFactor`
    - `transactionStateLogMinIsr`
    - `numIoThreads`
    - `numNetworkThreads`
    - `numPartitions`
    - `numRecoveryThreadsPerDataDir`
    - `socketReceiveBufferBytes`
    - `socketRequestMaxBytes`
    - `socketSendBufferBytes`
    - `zookeeperConnectionTimeoutMs`
    - `authorizerClassName`
    - `allowEveryoneIfNoAclFound`
    - `superUsers`
- All listeners are configured with protocol 'SASL_PLAINTEXT' by default.
- Support for SCRAM authentication in KRaft mode
- All statefulset settings have been moved from values' root to `controller.*` and `broker.*`.
- Refactor of listeners configuration:
  - Settings `listeners`, `advertisedListeners` and `listenerSecurityProtocolMap` have been replaced with `listeners.*` object, which includes default listeners and each listener can be configured individually and extended using `listeners.extraListeners`.
  - Values `interBrokerListenerName`, `allowPlaintextListener` have been removed.
- Refactor of SASL, SSL and ACL settings:
  - Authentication nomenclature `plaintext,tls,mtls,sasl,sasl_tls` has been removed. Listeners are now configured using Kafka nomenclature `PLAINTEXT,SASL_PLAINTEXT,SASL_SSL,SSL` in `listeners.*.protocol`.
  - mTLS is configured by default for SSL protocol listeners, while it can now also be configured for SASL_SSL listeners if `listener.*.sslClientAuth` is set.
  - All SASL settings are now grouped under `sasl.*`.
    - `auth.sasl.mechanisms` -> `sasl.enabledMechanisms`
    - `auth.interBrokerMechanism` -> `sasl.interBrokerMechanism`
    - `auth.sasl.jaas.clientUSers` -> `sasl.client.users`
    - `auth.sasl.jaas.clientPasswords` -> `sasl.client.passwords`
    - `auth.sasl.jaas.interBrokerUser` -> `sasl.interbroker.user`
    - `auth.sasl.jaas.interBrokerPassword` -> `sasl.interbroker.password`
    - `auth.sasl.jaas.zookeeperUser` -> `sasl.zookeeper.user`
    - `auth.sasl.jaas.zookeeperPassword` -> `sasl.zookeeper.password`
    - `auth.sasl.jaas.existingSecret` -> `sasl.existingSecret`
  - Added support for Controller listener protocols other than PLAINTEXT.
  - TLS settings have been moved from `auth.tls.*` to `tls.*`.
  - Zookeeper TLS settings have been moved from `auth.zookeeper*` to `tls.zookeeper.*`
- Refactor externalAccess to support the new architecture:
  - `externalAccess.service.*` have been renamed to `externalAccess.controller.service.*` and `externalAccess.controller.service.*`.
  - Controller pods will not configure externalAccess unless:
    - `controller.controllerOnly=false` (default), meaning the pods are running as 'controller+broker' nodes.
    - `externalAccess.controller.service.forceExpose=true`, for use cases where controller-only nodes want to be exposed externally.

#### Upgrading from Kraft mode

If upgrading from Kraft mode, existing PVCs from Kafka containers should be reattached to 'controller' pods.

#### Upgrading from Zookeeper mode

If upgrading from Zookeeper mode, make sure you set 'controller.replicaCount=0' and reattach the existing PVCs to 'broker' pods.
This will allow you to perform a migration to Kraft mode in the future by following the 'Migrating from Zookeeper' section of this documentation.

#### Retaining PersistentVolumes

When upgrading the Kafka chart, you may want to retain your existing data. To do so, we recommend following this guide:

**NOTE**: This guide requires the binaries 'kubectl' and 'jq'.

```console
# Env variables
REPLICA=0
OLD_PVC="data-<your_release_name>-kafka-${REPLICA}"
NEW_PVC="data-<your_release_name>-kafka-<your_kafka_role>-${REPLICA}"
PV_NAME=$(kubectl get pvc $OLD_PVC -o jsonpath="{.spec.volumeName}")
NEW_PVC_MANIFEST_FILE="$NEW_PVC.yaml"

# Modify PV reclaim policy
kubectl patch pv $PV_NAME -p '{"spec":{"persistentVolumeReclaimPolicy":"Retain"}}'
# Manually check field 'RECLAIM POLICY'
kubectl get pv $PV_NAME

# Create new PVC manifest
kubectl get pvc $OLD_PVC -o json | jq "
  .metadata.name = \"$NEW_PVC\"
  | with_entries(
      select([.key] |
        inside([\"metadata\", \"spec\", \"apiVersion\", \"kind\"]))
    )
  | del(
      .metadata.annotations, .metadata.creationTimestamp,
      .metadata.finalizers, .metadata.resourceVersion,
      .metadata.selfLink, .metadata.uid
    )
  " > $NEW_PVC_MANIFEST_FILE
# Check manifest
cat $NEW_PVC_MANIFEST_FILE

# Delete your old Statefulset and PVC
kubectl delete sts "<your_release_name>-kafka"
kubectl delete pvc $OLD_PVC
# Make PV available again and create the new PVC
kubectl patch pv $PV_NAME -p '{"spec":{"claimRef": null}}'
kubectl apply -f $NEW_PVC_MANIFEST_FILE
```

Repeat this process for each replica you had in your Kafka cluster. Once completed, upgrade the cluster and the new Statefulset should reuse the existing PVCs.

### To 23.0.0

This major updates Kafka to its newest version, 3.5.x. For more information, please refer to [kafka upgrade notes](https://kafka.apache.org/35/documentation.html#upgrade).

### To 22.0.0

This major updates the Kafka's configuration to use Kraft by default. You can learn more about this configuration [here](https://developer.confluent.io/learn/kraft). Apart from seting the `kraft.enabled` parameter to `true`, we also made the following changes:

- Renamed `minBrokerId` parameter to `minId` to set the minimum ID to use when configuring the node.id or broker.id parameter depending on the Kafka's configuration. This parameter sets the `KAFKA_CFG_NODE_ID` env var in the container.
- Updated the `containerPorts` and `service.ports` parameters to include the new controller port.

### To 21.0.0

This major updates Kafka to its newest version, 3.4.x. For more information, please refer to [kafka upgrade notes](https://kafka.apache.org/34/documentation.html#upgrade).

### To 20.0.0

This major updates the Zookeeper subchart to it newest major, 11.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-1100).

### To 19.0.0

This major updates Kafka to its newest version, 3.3.x. For more information, please refer to [kafka upgrade notes](https://kafka.apache.org/33/documentation.html#upgrade).

### To 18.0.0

This major updates the Zookeeper subchart to it newest major, 10.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-1000).

### To 16.0.0

This major updates the Zookeeper subchart to it newest major, 9.0.0. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-900).

### To 15.0.0

This major release bumps Kafka major version to `3.x` series.
It also renames several values in this chart and adds missing features, in order to be inline with the rest of assets in the Bitnami charts repository. Some affected values are:

- `service.port`, `service.internalPort` and `service.externalPort` have been regrouped under the `service.ports` map.
- `metrics.kafka.service.port` has been regrouped under the `metrics.kafka.service.ports` map.
- `metrics.jmx.service.port` has been regrouped under the `metrics.jmx.service.ports` map.
- `updateStrategy` (string) and `rollingUpdatePartition` are regrouped under the `updateStrategy` map.
- Several parameters marked as deprecated `14.x.x` are not supported anymore.

Additionally updates the ZooKeeper subchart to it newest major, `8.0.0`, which contains similar changes.

### To 14.0.0

In this version, the `image` block is defined once and is used in the different templates, while in the previous version, the `image` block was duplicated for the main container and the provisioning one

```yaml
image:
  registry: docker.io
  repository: bitnami/kafka
  tag: 2.8.0
```

VS

```yaml
image:
  registry: docker.io
  repository: bitnami/kafka
  tag: 2.8.0
...
provisioning:
  image:
    registry: docker.io
    repository: bitnami/kafka
    tag: 2.8.0
```

See [PR#7114](https://github.com/bitnami/charts/pull/7114) for more info about the implemented changes

### To 13.0.0

This major updates the Zookeeper subchart to it newest major, 7.0.0, which renames all TLS-related settings. For more information on this subchart's major, please refer to [zookeeper upgrade notes](https://github.com/bitnami/charts/tree/main/bitnami/zookeeper#to-700).

### To 12.2.0

This version also introduces `bitnami/common`, a [library chart](https://helm.sh/docs/topics/library_charts/#helm) as a dependency. More documentation about this new utility could be found [here](https://github.com/bitnami/charts/tree/main/bitnami/common#bitnami-common-library-chart). Please, make sure that you have updated the chart dependencies before executing any upgrade.

### To 12.0.0

[On November 13, 2020, Helm v2 support was formally finished](https://github.com/helm/charts#status-of-the-project), this major version is the result of the required changes applied to the Helm Chart to be able to incorporate the different features added in Helm v3 and to be consistent with the Helm project itself regarding the Helm v2 EOL.

#### What changes were introduced in this major version?

- Previous versions of this Helm Chart use `apiVersion: v1` (installable by both Helm 2 and 3), this Helm Chart was updated to `apiVersion: v2` (installable by Helm 3 only). [Here](https://helm.sh/docs/topics/charts/#the-apiversion-field) you can find more information about the `apiVersion` field.
- Move dependency information from the *requirements.yaml* to the *Chart.yaml*
- After running `helm dependency update`, a *Chart.lock* file is generated containing the same structure used in the previous *requirements.lock*
- The different fields present in the *Chart.yaml* file has been ordered alphabetically in a homogeneous way for all the Bitnami Helm Charts

#### Considerations when upgrading to this version

- If you want to upgrade to this version from a previous one installed with Helm v3, you shouldn't face any issues
- If you want to upgrade to this version using Helm v2, this scenario is not supported as this version doesn't support Helm v2 anymore
- If you installed the previous version with Helm v2 and wants to upgrade to this version with Helm v3, please refer to the [official Helm documentation](https://helm.sh/docs/topics/v2_v3_migration/#migration-use-cases) about migrating from Helm v2 to v3

#### Useful links

- <https://docs.bitnami.com/tutorials/resolve-helm2-helm3-post-migration-issues/>
- <https://helm.sh/docs/topics/v2_v3_migration/>
- <https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/>

### To 11.8.0

External access to brokers can now be achieved through the cluster's Kafka service.

- `service.nodePort` -> deprecated  in favor of `service.nodePorts.client` and `service.nodePorts.external`

### To 11.7.0

The way to configure the users and passwords changed. Now it is allowed to create multiple users during the installation by providing the list of users and passwords.

- `auth.jaas.clientUser` (string) -> deprecated  in favor of `auth.jaas.clientUsers` (array).
- `auth.jaas.clientPassword` (string) -> deprecated  in favor of `auth.jaas.clientPasswords` (array).

### To 11.0.0

The way to configure listeners and athentication on Kafka is totally refactored allowing users to configure different authentication protocols on different listeners. Please check the [Listeners Configuration](#listeners-configuration) section for more information.

Backwards compatibility is not guaranteed you adapt your values.yaml to the new format. Here you can find some parameters that were renamed or disappeared in favor of new ones on this major version:

- `auth.enabled` -> deprecated in favor of `auth.clientProtocol` and `auth.interBrokerProtocol` parameters.
- `auth.ssl` -> deprecated in favor of `auth.clientProtocol` and `auth.interBrokerProtocol` parameters.
- `auth.certificatesSecret` -> renamed to `auth.jksSecret`.
- `auth.certificatesPassword` -> renamed to `auth.jksPassword`.
- `sslEndpointIdentificationAlgorithm` -> renamedo to `auth.tlsEndpointIdentificationAlgorithm`.
- `auth.interBrokerUser` -> renamed to `auth.jaas.interBrokerUser`
- `auth.interBrokerPassword` -> renamed to `auth.jaas.interBrokerPassword`
- `auth.zookeeperUser` -> renamed to `auth.jaas.zookeeperUser`
- `auth.zookeeperPassword` -> renamed to `auth.jaas.zookeeperPassword`
- `auth.existingSecret` -> renamed to `auth.jaas.existingSecret`
- `service.sslPort` -> deprecated in favor of `service.internalPort`
- `service.nodePorts.kafka` and `service.nodePorts.ssl` -> deprecated in favor of `service.nodePort`
- `metrics.kafka.extraFlag` -> new parameter
- `metrics.kafka.certificatesSecret` -> new parameter

### To 10.0.0

If you are setting the `config` or `log4j` parameter, backwards compatibility is not guaranteed, because the `KAFKA_MOUNTED_CONFDIR` has moved from `/opt/bitnami/kafka/conf` to `/bitnami/kafka/config`. In order to continue using these parameters, you must also upgrade your image to `docker.io/bitnami/kafka:2.4.1-debian-10-r38` or later.

### To 9.0.0

Backwards compatibility is not guaranteed you adapt your values.yaml to the new format. Here you can find some parameters that were renamed on this major version:

```diff
- securityContext.enabled
- securityContext.fsGroup
- securityContext.fsGroup
+ podSecurityContext
- externalAccess.service.loadBalancerIP
+ externalAccess.service.loadBalancerIPs
- externalAccess.service.nodePort
+ externalAccess.service.nodePorts
- metrics.jmx.configMap.enabled
- metrics.jmx.configMap.overrideConfig
+ metrics.jmx.config
- metrics.jmx.configMap.overrideName
+ metrics.jmx.existingConfigmap
```

Ports names were prefixed with the protocol to comply with Istio (see <https://istio.io/docs/ops/deployment/requirements/>).

### To 8.0.0

There is not backwards compatibility since the brokerID changes to the POD_NAME. For more information see [this PR](https://github.com/bitnami/charts/pull/2028).

### To 7.0.0

Backwards compatibility is not guaranteed when Kafka metrics are enabled, unless you modify the labels used on the exporter deployments.
Use the workaround below to upgrade from versions previous to 7.0.0. The following example assumes that the release name is kafka:

```console
helm upgrade kafka oci://registry-1.docker.io/bitnamicharts/kafka --version 6.1.8 --set metrics.kafka.enabled=false
helm upgrade kafka oci://registry-1.docker.io/bitnamicharts/kafka --version 7.0.0 --set metrics.kafka.enabled=true
```

### To 2.0.0

Backwards compatibility is not guaranteed unless you modify the labels used on the chart's deployments.
Use the workaround below to upgrade from versions previous to 2.0.0. The following example assumes that the release name is kafka:

```console
kubectl delete statefulset kafka-kafka --cascade=false
kubectl delete statefulset kafka-zookeeper --cascade=false
```

### To 1.0.0

Backwards compatibility is not guaranteed unless you modify the labels used on the chart's deployments.
Use the workaround below to upgrade from versions previous to 1.0.0. The following example assumes that the release name is kafka:

```console
kubectl delete statefulset kafka-kafka --cascade=false
kubectl delete statefulset kafka-zookeeper --cascade=false
```

## License

Copyright &copy; 2023 VMware, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.