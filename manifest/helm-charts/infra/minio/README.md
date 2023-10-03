<!--- app-name: Bitnami Object Storage based on MinIO&reg; -->

# Bitnami Object Storage based on MinIO(R)

MinIO(R) is an object storage server, compatible with Amazon S3 cloud storage service, mainly used for storing unstructured data (such as photos, videos, log files, etc.).

[Overview of Bitnami Object Storage based on MinIO&reg;](https://min.io/)

Disclaimer: All software products, projects and company names are trademark(TM) or registered(R) trademarks of their respective holders, and use of them does not imply any affiliation or endorsement. This software is licensed to you subject to one or more open source licenses and VMware provides the software on an AS-IS basis. MinIO(R) is a registered trademark of the MinIO Inc. in the US and other countries. Bitnami is not affiliated, associated, authorized, endorsed by, or in any way officially connected with MinIO Inc. MinIO(R) is licensed under GNU AGPL v3.0.

## TL;DR

```console
helm install my-release oci://registry-1.docker.io/bitnamicharts/minio
```

## Introduction

This chart bootstraps a [MinIO&reg;](https://github.com/bitnami/containers/tree/main/bitnami/minio) deployment on a [Kubernetes](https://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

Bitnami charts can be used with [Kubeapps](https://kubeapps.dev/) for deployment and management of Helm Charts in clusters.

Looking to use Bitnami Object Storage based on MinIOreg; in production? Try [VMware Application Catalog](https://bitnami.com/enterprise), the enterprise edition of Bitnami Application Catalog.

## Prerequisites

- Kubernetes 1.19+
- Helm 3.2.0+
- PV provisioner support in the underlying infrastructure

## Installing the Chart

To install the chart with the release name `my-release`:

```console
helm install my-release oci://registry-1.docker.io/bitnamicharts/minio
```

These commands deploy MinIO&reg; on the Kubernetes cluster in the default configuration. The [Parameters](#parameters) section lists the parameters that can be configured during installation.

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

| Name                | Description                                                                                  | Value           |
| ------------------- | -------------------------------------------------------------------------------------------- | --------------- |
| `nameOverride`      | String to partially override common.names.fullname template (will maintain the release name) | `""`            |
| `fullnameOverride`  | String to fully override common.names.fullname template                                      | `""`            |
| `commonLabels`      | Labels to add to all deployed objects                                                        | `{}`            |
| `commonAnnotations` | Annotations to add to all deployed objects                                                   | `{}`            |
| `kubeVersion`       | Force target Kubernetes version (using Helm capabilities if not set)                         | `""`            |
| `clusterDomain`     | Default Kubernetes cluster domain                                                            | `cluster.local` |
| `extraDeploy`       | Array of extra objects to deploy with the release                                            | `[]`            |

### MinIO&reg; parameters

| Name                       | Description                                                                                                                                                                                               | Value                    |
| -------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| `image.registry`           | MinIO&reg; image registry                                                                                                                                                                                 | `docker.io`              |
| `image.repository`         | MinIO&reg; image repository                                                                                                                                                                               | `bitnami/minio`          |
| `image.tag`                | MinIO&reg; image tag (immutable tags are recommended)                                                                                                                                                     | `2023.9.20-debian-11-r0` |
| `image.digest`             | MinIO&reg; image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                                | `""`                     |
| `image.pullPolicy`         | Image pull policy                                                                                                                                                                                         | `IfNotPresent`           |
| `image.pullSecrets`        | Specify docker-registry secret names as an array                                                                                                                                                          | `[]`                     |
| `image.debug`              | Specify if debug logs should be enabled                                                                                                                                                                   | `false`                  |
| `clientImage.registry`     | MinIO&reg; Client image registry                                                                                                                                                                          | `docker.io`              |
| `clientImage.repository`   | MinIO&reg; Client image repository                                                                                                                                                                        | `bitnami/minio-client`   |
| `clientImage.tag`          | MinIO&reg; Client image tag (immutable tags are recommended)                                                                                                                                              | `2023.9.20-debian-11-r0` |
| `clientImage.digest`       | MinIO&reg; Client image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag                                                                                         | `""`                     |
| `mode`                     | MinIO&reg; server mode (`standalone` or `distributed`)                                                                                                                                                    | `standalone`             |
| `auth.rootUser`            | MinIO&reg; root username                                                                                                                                                                                  | `admin`                  |
| `auth.rootPassword`        | Password for MinIO&reg; root user                                                                                                                                                                         | `""`                     |
| `auth.existingSecret`      | Use existing secret for credentials details (`auth.rootUser` and `auth.rootPassword` will be ignored and picked up from this secret). The secret has to contain the keys `root-user` and `root-password`) | `""`                     |
| `auth.forcePassword`       | Force users to specify required passwords                                                                                                                                                                 | `false`                  |
| `auth.useCredentialsFiles` | Mount credentials as a files instead of using an environment variable                                                                                                                                     | `false`                  |
| `auth.forceNewKeys`        | Force root credentials (user and password) to be reconfigured every time they change in the secrets                                                                                                       | `false`                  |
| `defaultBuckets`           | Comma, semi-colon or space separated list of buckets to create at initialization (only in standalone mode)                                                                                                | `""`                     |
| `disableWebUI`             | Disable MinIO&reg; Web UI                                                                                                                                                                                 | `false`                  |
| `tls.enabled`              | Enable tls in front of the container                                                                                                                                                                      | `false`                  |
| `tls.autoGenerated`        | Generate automatically self-signed TLS certificates                                                                                                                                                       | `false`                  |
| `tls.existingSecret`       | Name of an existing secret holding the certificate information                                                                                                                                            | `""`                     |
| `tls.mountPath`            | The mount path where the secret will be located                                                                                                                                                           | `""`                     |
| `extraEnvVars`             | Extra environment variables to be set on MinIO&reg; container                                                                                                                                             | `[]`                     |
| `extraEnvVarsCM`           | ConfigMap with extra environment variables                                                                                                                                                                | `""`                     |
| `extraEnvVarsSecret`       | Secret with extra environment variables                                                                                                                                                                   | `""`                     |
| `command`                  | Default container command (useful when using custom images). Use array form                                                                                                                               | `[]`                     |
| `args`                     | Default container args (useful when using custom images). Use array form                                                                                                                                  | `[]`                     |

### MinIO&reg; deployment/statefulset parameters

| Name                                                 | Description                                                                                                                                                                                   | Value           |
| ---------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------- |
| `schedulerName`                                      | Specifies the schedulerName, if it's nil uses kube-scheduler                                                                                                                                  | `""`            |
| `terminationGracePeriodSeconds`                      | In seconds, time the given to the MinIO pod needs to terminate gracefully                                                                                                                     | `""`            |
| `deployment.updateStrategy.type`                     | Deployment strategy type                                                                                                                                                                      | `Recreate`      |
| `statefulset.updateStrategy.type`                    | StatefulSet strategy type                                                                                                                                                                     | `RollingUpdate` |
| `statefulset.podManagementPolicy`                    | StatefulSet controller supports relax its ordering guarantees while preserving its uniqueness and identity guarantees. There are two valid pod management policies: OrderedReady and Parallel | `Parallel`      |
| `statefulset.replicaCount`                           | Number of pods per zone (only for MinIO&reg; distributed mode). Should be even and `>= 4`                                                                                                     | `4`             |
| `statefulset.zones`                                  | Number of zones (only for MinIO&reg; distributed mode)                                                                                                                                        | `1`             |
| `statefulset.drivesPerNode`                          | Number of drives attached to every node (only for MinIO&reg; distributed mode)                                                                                                                | `1`             |
| `provisioning.enabled`                               | Enable MinIO&reg; provisioning Job                                                                                                                                                            | `false`         |
| `provisioning.schedulerName`                         | Name of the k8s scheduler (other than default) for MinIO&reg; provisioning                                                                                                                    | `""`            |
| `provisioning.podLabels`                             | Extra labels for provisioning pods                                                                                                                                                            | `{}`            |
| `provisioning.podAnnotations`                        | Provisioning Pod annotations.                                                                                                                                                                 | `{}`            |
| `provisioning.command`                               | Default provisioning container command (useful when using custom images). Use array form                                                                                                      | `[]`            |
| `provisioning.args`                                  | Default provisioning container args (useful when using custom images). Use array form                                                                                                         | `[]`            |
| `provisioning.extraCommands`                         | Optionally specify extra list of additional commands for MinIO&reg; provisioning pod                                                                                                          | `[]`            |
| `provisioning.extraVolumes`                          | Optionally specify extra list of additional volumes for MinIO&reg; provisioning pod                                                                                                           | `[]`            |
| `provisioning.extraVolumeMounts`                     | Optionally specify extra list of additional volumeMounts for MinIO&reg; provisioning container                                                                                                | `[]`            |
| `provisioning.resources.limits`                      | The resources limits for the container                                                                                                                                                        | `{}`            |
| `provisioning.resources.requests`                    | The requested resources for the container                                                                                                                                                     | `{}`            |
| `provisioning.policies`                              | MinIO&reg; policies provisioning                                                                                                                                                              | `[]`            |
| `provisioning.users`                                 | MinIO&reg; users provisioning. Can be used in addition to provisioning.usersExistingSecrets.                                                                                                  | `[]`            |
| `provisioning.usersExistingSecrets`                  | Array if existing secrets containing MinIO&reg; users to be provisioned. Can be used in addition to provisioning.users.                                                                       | `[]`            |
| `provisioning.groups`                                | MinIO&reg; groups provisioning                                                                                                                                                                | `[]`            |
| `provisioning.buckets`                               | MinIO&reg; buckets, versioning, lifecycle, quota and tags provisioning                                                                                                                        | `[]`            |
| `provisioning.config`                                | MinIO&reg; config provisioning                                                                                                                                                                | `[]`            |
| `provisioning.podSecurityContext.enabled`            | Enable pod Security Context                                                                                                                                                                   | `true`          |
| `provisioning.podSecurityContext.fsGroup`            | Group ID for the container                                                                                                                                                                    | `1001`          |
| `provisioning.containerSecurityContext.enabled`      | Enable container Security Context                                                                                                                                                             | `true`          |
| `provisioning.containerSecurityContext.runAsUser`    | User ID for the container                                                                                                                                                                     | `1001`          |
| `provisioning.containerSecurityContext.runAsNonRoot` | Avoid running as root User                                                                                                                                                                    | `true`          |
| `provisioning.cleanupAfterFinished.enabled`          | Enables Cleanup for Finished Jobs                                                                                                                                                             | `false`         |
| `provisioning.cleanupAfterFinished.seconds`          | Sets the value of ttlSecondsAfterFinished                                                                                                                                                     | `600`           |
| `hostAliases`                                        | MinIO&reg; pod host aliases                                                                                                                                                                   | `[]`            |
| `containerPorts.api`                                 | MinIO&reg; container port to open for MinIO&reg; API                                                                                                                                          | `9000`          |
| `containerPorts.console`                             | MinIO&reg; container port to open for MinIO&reg; Console                                                                                                                                      | `9001`          |
| `podSecurityContext.enabled`                         | Enable pod Security Context                                                                                                                                                                   | `true`          |
| `podSecurityContext.fsGroup`                         | Group ID for the container                                                                                                                                                                    | `1001`          |
| `containerSecurityContext.enabled`                   | Enable container Security Context                                                                                                                                                             | `true`          |
| `containerSecurityContext.runAsUser`                 | User ID for the container                                                                                                                                                                     | `1001`          |
| `containerSecurityContext.runAsNonRoot`              | Avoid running as root User                                                                                                                                                                    | `true`          |
| `podLabels`                                          | Extra labels for MinIO&reg; pods                                                                                                                                                              | `{}`            |
| `podAnnotations`                                     | Annotations for MinIO&reg; pods                                                                                                                                                               | `{}`            |
| `podAffinityPreset`                                  | Pod affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                           | `""`            |
| `podAntiAffinityPreset`                              | Pod anti-affinity preset. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                      | `soft`          |
| `nodeAffinityPreset.type`                            | Node affinity preset type. Ignored if `affinity` is set. Allowed values: `soft` or `hard`                                                                                                     | `""`            |
| `nodeAffinityPreset.key`                             | Node label key to match. Ignored if `affinity` is set.                                                                                                                                        | `""`            |
| `nodeAffinityPreset.values`                          | Node label values to match. Ignored if `affinity` is set.                                                                                                                                     | `[]`            |
| `affinity`                                           | Affinity for pod assignment. Evaluated as a template.                                                                                                                                         | `{}`            |
| `nodeSelector`                                       | Node labels for pod assignment. Evaluated as a template.                                                                                                                                      | `{}`            |
| `tolerations`                                        | Tolerations for pod assignment. Evaluated as a template.                                                                                                                                      | `[]`            |
| `topologySpreadConstraints`                          | Topology Spread Constraints for MinIO&reg; pods assignment spread across your cluster among failure-domains                                                                                   | `[]`            |
| `priorityClassName`                                  | MinIO&reg; pods' priorityClassName                                                                                                                                                            | `""`            |
| `resources.limits`                                   | The resources limits for the MinIO&reg; container                                                                                                                                             | `{}`            |
| `resources.requests`                                 | The requested resources for the MinIO&reg; container                                                                                                                                          | `{}`            |
| `livenessProbe.enabled`                              | Enable livenessProbe                                                                                                                                                                          | `true`          |
| `livenessProbe.initialDelaySeconds`                  | Initial delay seconds for livenessProbe                                                                                                                                                       | `5`             |
| `livenessProbe.periodSeconds`                        | Period seconds for livenessProbe                                                                                                                                                              | `5`             |
| `livenessProbe.timeoutSeconds`                       | Timeout seconds for livenessProbe                                                                                                                                                             | `5`             |
| `livenessProbe.failureThreshold`                     | Failure threshold for livenessProbe                                                                                                                                                           | `5`             |
| `livenessProbe.successThreshold`                     | Success threshold for livenessProbe                                                                                                                                                           | `1`             |
| `readinessProbe.enabled`                             | Enable readinessProbe                                                                                                                                                                         | `true`          |
| `readinessProbe.initialDelaySeconds`                 | Initial delay seconds for readinessProbe                                                                                                                                                      | `5`             |
| `readinessProbe.periodSeconds`                       | Period seconds for readinessProbe                                                                                                                                                             | `5`             |
| `readinessProbe.timeoutSeconds`                      | Timeout seconds for readinessProbe                                                                                                                                                            | `1`             |
| `readinessProbe.failureThreshold`                    | Failure threshold for readinessProbe                                                                                                                                                          | `5`             |
| `readinessProbe.successThreshold`                    | Success threshold for readinessProbe                                                                                                                                                          | `1`             |
| `startupProbe.enabled`                               | Enable startupProbe                                                                                                                                                                           | `false`         |
| `startupProbe.initialDelaySeconds`                   | Initial delay seconds for startupProbe                                                                                                                                                        | `0`             |
| `startupProbe.periodSeconds`                         | Period seconds for startupProbe                                                                                                                                                               | `10`            |
| `startupProbe.timeoutSeconds`                        | Timeout seconds for startupProbe                                                                                                                                                              | `5`             |
| `startupProbe.failureThreshold`                      | Failure threshold for startupProbe                                                                                                                                                            | `60`            |
| `startupProbe.successThreshold`                      | Success threshold for startupProbe                                                                                                                                                            | `1`             |
| `customLivenessProbe`                                | Override default liveness probe                                                                                                                                                               | `{}`            |
| `customReadinessProbe`                               | Override default readiness probe                                                                                                                                                              | `{}`            |
| `customStartupProbe`                                 | Override default startup probe                                                                                                                                                                | `{}`            |
| `lifecycleHooks`                                     | for the MinIO&reg container(s) to automate configuration before or after startup                                                                                                              | `{}`            |
| `extraVolumes`                                       | Optionally specify extra list of additional volumes for MinIO&reg; pods                                                                                                                       | `[]`            |
| `extraVolumeMounts`                                  | Optionally specify extra list of additional volumeMounts for MinIO&reg; container(s)                                                                                                          | `[]`            |
| `initContainers`                                     | Add additional init containers to the MinIO&reg; pods                                                                                                                                         | `[]`            |
| `sidecars`                                           | Add additional sidecar containers to the MinIO&reg; pods                                                                                                                                      | `[]`            |

### Traffic exposure parameters

| Name                               | Description                                                                                                                      | Value                    |
| ---------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- | ------------------------ |
| `service.type`                     | MinIO&reg; service type                                                                                                          | `ClusterIP`              |
| `service.ports.api`                | MinIO&reg; API service port                                                                                                      | `9000`                   |
| `service.ports.console`            | MinIO&reg; Console service port                                                                                                  | `9001`                   |
| `service.nodePorts.api`            | Specify the MinIO&reg API nodePort value for the LoadBalancer and NodePort service types                                         | `""`                     |
| `service.nodePorts.console`        | Specify the MinIO&reg Console nodePort value for the LoadBalancer and NodePort service types                                     | `""`                     |
| `service.clusterIP`                | Service Cluster IP                                                                                                               | `""`                     |
| `service.loadBalancerIP`           | loadBalancerIP if service type is `LoadBalancer` (optional, cloud specific)                                                      | `""`                     |
| `service.loadBalancerSourceRanges` | Addresses that are allowed when service is LoadBalancer                                                                          | `[]`                     |
| `service.externalTrafficPolicy`    | Enable client source IP preservation                                                                                             | `Cluster`                |
| `service.extraPorts`               | Extra ports to expose in the service (normally used with the `sidecar` value)                                                    | `[]`                     |
| `service.annotations`              | Annotations for MinIO&reg; service                                                                                               | `{}`                     |
| `service.headless.annotations`     | Annotations for the headless service.                                                                                            | `{}`                     |
| `ingress.enabled`                  | Enable ingress controller resource for MinIO Console                                                                             | `false`                  |
| `ingress.apiVersion`               | Force Ingress API version (automatically detected if not set)                                                                    | `""`                     |
| `ingress.ingressClassName`         | IngressClass that will be be used to implement the Ingress (Kubernetes 1.18+)                                                    | `""`                     |
| `ingress.hostname`                 | Default host for the ingress resource                                                                                            | `minio.local`            |
| `ingress.path`                     | The Path to MinIO&reg;. You may need to set this to '/*' in order to use this with ALB ingress controllers.                      | `/`                      |
| `ingress.pathType`                 | Ingress path type                                                                                                                | `ImplementationSpecific` |
| `ingress.servicePort`              | Service port to be used                                                                                                          | `minio-console`          |
| `ingress.annotations`              | Additional annotations for the Ingress resource. To enable certificate autogeneration, place here your cert-manager annotations. | `{}`                     |
| `ingress.tls`                      | Enable TLS configuration for the hostname defined at `ingress.hostname` parameter                                                | `false`                  |
| `ingress.selfSigned`               | Create a TLS secret for this ingress record using self-signed certificates generated by Helm                                     | `false`                  |
| `ingress.extraHosts`               | The list of additional hostnames to be covered with this ingress record.                                                         | `[]`                     |
| `ingress.extraPaths`               | Any additional paths that may need to be added to the ingress under the main host                                                | `[]`                     |
| `ingress.extraTls`                 | The tls configuration for additional hostnames to be covered with this ingress record.                                           | `[]`                     |
| `ingress.secrets`                  | If you're providing your own certificates, please use this to add the certificates as secrets                                    | `[]`                     |
| `ingress.extraRules`               | Additional rules to be covered with this ingress record                                                                          | `[]`                     |
| `apiIngress.enabled`               | Enable ingress controller resource for MinIO API                                                                                 | `false`                  |
| `apiIngress.apiVersion`            | Force Ingress API version (automatically detected if not set)                                                                    | `""`                     |
| `apiIngress.ingressClassName`      | IngressClass that will be be used to implement the Ingress (Kubernetes 1.18+)                                                    | `""`                     |
| `apiIngress.hostname`              | Default host for the ingress resource                                                                                            | `minio.local`            |
| `apiIngress.path`                  | The Path to MinIO&reg;. You may need to set this to '/*' in order to use this with ALB ingress controllers.                      | `/`                      |
| `apiIngress.pathType`              | Ingress path type                                                                                                                | `ImplementationSpecific` |
| `apiIngress.servicePort`           | Service port to be used                                                                                                          | `minio-api`              |
| `apiIngress.annotations`           | Additional annotations for the Ingress resource. To enable certificate autogeneration, place here your cert-manager annotations. | `{}`                     |
| `apiIngress.tls`                   | Enable TLS configuration for the hostname defined at `apiIngress.hostname` parameter                                             | `false`                  |
| `apiIngress.selfSigned`            | Create a TLS secret for this ingress record using self-signed certificates generated by Helm                                     | `false`                  |
| `apiIngress.extraHosts`            | The list of additional hostnames to be covered with this ingress record.                                                         | `[]`                     |
| `apiIngress.extraPaths`            | Any additional paths that may need to be added to the ingress under the main host                                                | `[]`                     |
| `apiIngress.extraTls`              | The tls configuration for additional hostnames to be covered with this ingress record.                                           | `[]`                     |
| `apiIngress.secrets`               | If you're providing your own certificates, please use this to add the certificates as secrets                                    | `[]`                     |
| `apiIngress.extraRules`            | Additional rules to be covered with this ingress record                                                                          | `[]`                     |
| `networkPolicy.enabled`            | Enable the default NetworkPolicy policy                                                                                          | `false`                  |
| `networkPolicy.allowExternal`      | Don't require client label for connections                                                                                       | `true`                   |
| `networkPolicy.extraFromClauses`   | Allows to add extra 'from' clauses to the NetworkPolicy                                                                          | `[]`                     |

### Persistence parameters

| Name                        | Description                                                          | Value                 |
| --------------------------- | -------------------------------------------------------------------- | --------------------- |
| `persistence.enabled`       | Enable MinIO&reg; data persistence using PVC. If false, use emptyDir | `true`                |
| `persistence.storageClass`  | PVC Storage Class for MinIO&reg; data volume                         | `""`                  |
| `persistence.mountPath`     | Data volume mount path                                               | `/bitnami/minio/data` |
| `persistence.accessModes`   | PVC Access Modes for MinIO&reg; data volume                          | `["ReadWriteOnce"]`   |
| `persistence.size`          | PVC Storage Request for MinIO&reg; data volume                       | `8Gi`                 |
| `persistence.annotations`   | Annotations for the PVC                                              | `{}`                  |
| `persistence.existingClaim` | Name of an existing PVC to use (only in `standalone` mode)           | `""`                  |

### Volume Permissions parameters

| Name                                                   | Description                                                                                                                       | Value              |
| ------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------- | ------------------ |
| `volumePermissions.enabled`                            | Enable init container that changes the owner and group of the persistent volume(s) mountpoint to `runAsUser:fsGroup`              | `false`            |
| `volumePermissions.image.registry`                     | Init container volume-permissions image registry                                                                                  | `docker.io`        |
| `volumePermissions.image.repository`                   | Init container volume-permissions image repository                                                                                | `bitnami/os-shell` |
| `volumePermissions.image.tag`                          | Init container volume-permissions image tag (immutable tags are recommended)                                                      | `11-debian-11-r72` |
| `volumePermissions.image.digest`                       | Init container volume-permissions image digest in the way sha256:aa.... Please note this parameter, if set, will override the tag | `""`               |
| `volumePermissions.image.pullPolicy`                   | Init container volume-permissions image pull policy                                                                               | `IfNotPresent`     |
| `volumePermissions.image.pullSecrets`                  | Specify docker-registry secret names as an array                                                                                  | `[]`               |
| `volumePermissions.resources.limits`                   | Init container volume-permissions resource limits                                                                                 | `{}`               |
| `volumePermissions.resources.requests`                 | Init container volume-permissions resource requests                                                                               | `{}`               |
| `volumePermissions.containerSecurityContext.runAsUser` | User ID for the init container                                                                                                    | `0`                |

### RBAC parameters

| Name                                          | Description                                                 | Value  |
| --------------------------------------------- | ----------------------------------------------------------- | ------ |
| `serviceAccount.create`                       | Enable the creation of a ServiceAccount for MinIO&reg; pods | `true` |
| `serviceAccount.name`                         | Name of the created ServiceAccount                          | `""`   |
| `serviceAccount.automountServiceAccountToken` | Enable/disable auto mounting of the service account token   | `true` |
| `serviceAccount.annotations`                  | Custom annotations for MinIO&reg; ServiceAccount            | `{}`   |

### Other parameters

| Name                 | Description                                                                       | Value   |
| -------------------- | --------------------------------------------------------------------------------- | ------- |
| `pdb.create`         | Enable/disable a Pod Disruption Budget creation                                   | `false` |
| `pdb.minAvailable`   | Minimum number/percentage of pods that must still be available after the eviction | `1`     |
| `pdb.maxUnavailable` | Maximum number/percentage of pods that may be made unavailable after the eviction | `""`    |

### Metrics parameters

| Name                                       | Description                                                                                                                   | Value                       |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| `metrics.prometheusAuthType`               | Authentication mode for Prometheus (`jwt` or `public`)                                                                        | `public`                    |
| `metrics.serviceMonitor.enabled`           | If the operator is installed in your cluster, set to true to create a Service Monitor Entry                                   | `false`                     |
| `metrics.serviceMonitor.namespace`         | Namespace which Prometheus is running in                                                                                      | `""`                        |
| `metrics.serviceMonitor.labels`            | Extra labels for the ServiceMonitor                                                                                           | `{}`                        |
| `metrics.serviceMonitor.jobLabel`          | The name of the label on the target service to use as the job name in Prometheus                                              | `""`                        |
| `metrics.serviceMonitor.path`              | HTTP path to scrape for metrics                                                                                               | `/minio/v2/metrics/cluster` |
| `metrics.serviceMonitor.interval`          | Interval at which metrics should be scraped                                                                                   | `30s`                       |
| `metrics.serviceMonitor.scrapeTimeout`     | Specify the timeout after which the scrape is ended                                                                           | `""`                        |
| `metrics.serviceMonitor.metricRelabelings` | MetricRelabelConfigs to apply to samples before ingestion                                                                     | `[]`                        |
| `metrics.serviceMonitor.relabelings`       | Metrics relabelings to add to the scrape endpoint, applied before scraping                                                    | `[]`                        |
| `metrics.serviceMonitor.honorLabels`       | Specify honorLabels parameter to add the scrape endpoint                                                                      | `false`                     |
| `metrics.serviceMonitor.selector`          | Prometheus instance selector labels                                                                                           | `{}`                        |
| `metrics.serviceMonitor.apiVersion`        | ApiVersion for the serviceMonitor Resource (defaults to "monitoring.coreos.com/v1")                                           | `""`                        |
| `metrics.serviceMonitor.tlsConfig`         | Additional TLS configuration for metrics endpoint with "https" scheme                                                         | `{}`                        |
| `metrics.prometheusRule.enabled`           | Create a Prometheus Operator PrometheusRule (also requires `metrics.enabled` to be `true` and `metrics.prometheusRule.rules`) | `false`                     |
| `metrics.prometheusRule.namespace`         | Namespace for the PrometheusRule Resource (defaults to the Release Namespace)                                                 | `""`                        |
| `metrics.prometheusRule.additionalLabels`  | Additional labels that can be used so PrometheusRule will be discovered by Prometheus                                         | `{}`                        |
| `metrics.prometheusRule.rules`             | Prometheus Rule definitions                                                                                                   | `[]`                        |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example,

```console
helm install my-release \
  --set auth.rootUser=minio-admin \
  --set auth.rootPassword=minio-secret-password \
    oci://registry-1.docker.io/bitnamicharts/minio
```

The above command sets the MinIO&reg; Server root user and password to `minio-admin` and `minio-secret-password`, respectively.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```console
helm install my-release -f values.yaml oci://registry-1.docker.io/bitnamicharts/minio
```

> **Tip**: You can use the default [values.yaml](values.yaml)

## Configuration and installation details

### [Rolling VS Immutable tags](https://docs.bitnami.com/containers/how-to/understand-rolling-tags-containers/)

It is strongly recommended to use immutable tags in a production environment. This ensures your deployment does not change automatically if the same tag is updated with a different image.

Bitnami will release a new chart updating its containers if a new version of the main container, significant changes, or critical vulnerabilities exist.

### Distributed mode

By default, this chart provisions a MinIO&reg; server in standalone mode. You can start MinIO&reg; server in [distributed mode](https://docs.minio.io/docs/distributed-minio-quickstart-guide) with the following parameter: `mode=distributed`

This chart bootstrap MinIO&reg; server in distributed mode with 4 nodes by default. You can change the number of nodes using the `statefulset.replicaCount` parameter. For instance, you can deploy the chart with 8 nodes using the following parameters:

```console
mode=distributed
statefulset.replicaCount=8
```

You can also bootstrap MinIO&reg; server in distributed mode in several zones, and using multiple drives per node. For instance, you can deploy the chart with 2 nodes per zone on 2 zones, using 2 drives per node:

```console
mode=distributed
statefulset.replicaCount=2
statefulset.zones=2
statefulset.drivesPerNode=2
```

> Note: The total number of drives should be greater than 4 to guarantee erasure coding. Please set a combination of nodes, and drives per node that match this condition.

### Prometheus exporter

MinIO&reg; exports Prometheus metrics at `/minio/v2/metrics/cluster`. To allow Prometheus collecting your MinIO&reg; metrics, modify the `values.yaml` adding the corresponding annotations:

```diff
- podAnnotations: {}
+ podAnnotations:
+   prometheus.io/scrape: "true"
+   prometheus.io/path: "/minio/v2/metrics/cluster"
+   prometheus.io/port: "9000"
```

> Find more information about MinIO&reg; metrics at <https://docs.min.io/docs/how-to-monitor-minio-using-prometheus.html>

## Persistence

The [Bitnami Object Storage based on MinIO(&reg;)](https://github.com/bitnami/containers/tree/main/bitnami/minio) image stores data at the `/data` path of the container.

The chart mounts a [Persistent Volume](https://kubernetes.io/docs/concepts/storage/persistent-volumes/) at this location. The volume is created using dynamic volume provisioning.

### Adjust permissions of persistent volume mountpoint

As the image run as non-root by default, it is necessary to adjust the ownership of the persistent volume so that the container can write data into it.

By default, the chart is configured to use Kubernetes Security Context to automatically change the ownership of the volume. However, this feature does not work in all Kubernetes distributions.
As an alternative, this chart supports using an initContainer to change the ownership of the volume before mounting it in the final destination.

You can enable this initContainer by setting `volumePermissions.enabled` to `true`.

### Ingress

This chart provides support for Ingress resources. If you have an ingress controller installed on your cluster, such as [nginx-ingress-controller](https://github.com/bitnami/charts/tree/main/bitnami/nginx-ingress-controller) or [contour](https://github.com/bitnami/charts/tree/main/bitnami/contour) you can utilize the ingress controller to serve your application.

To enable Ingress integration, set `ingress.enabled` to `true`. The `ingress.hostname` property can be used to set the host name. The `ingress.tls` parameter can be used to add the TLS configuration for this host. It is also possible to have more than one host, with a separate TLS configuration for each host. [Learn more about configuring and using Ingress](https://docs.bitnami.com/kubernetes/infrastructure/minio/configuration/configure-ingress/).

### TLS secrets

The chart also facilitates the creation of TLS secrets for use with the Ingress controller, with different options for certificate management. [Learn more about TLS secrets](https://docs.bitnami.com/kubernetes/infrastructure/minio/administration/enable-tls-ingress/).

### Adding extra environment variables

In case you want to add extra environment variables (useful for advanced operations like custom init scripts), you can use the `extraEnvVars` property.

```yaml
extraEnvVars:
  - name: MINIO_LOG_LEVEL
    value: DEBUG
```

Alternatively, you can use a ConfigMap or a Secret with the environment variables. To do so, use the `extraEnvVarsCM` or the `extraEnvVarsSecret` values.

### Sidecars and Init Containers

If you have a need for additional containers to run within the same pod as the MinIO&reg; app (e.g. an additional metrics or logging exporter), you can do so via the `sidecars` config parameter. Simply define your container according to the Kubernetes container spec.

```yaml
sidecars:
  - name: your-image-name
    image: your-image
    imagePullPolicy: Always
    ports:
      - name: portname
       containerPort: 1234
```

Similarly, you can add extra init containers using the `initContainers` parameter.

```yaml
initContainers:
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

There are cases where you may want to deploy extra objects, such a ConfigMap containing your app's configuration or some extra deployment with a micro service used by your app. For covering this case, the chart allows adding the full specification of other objects using the `extraDeploy` parameter.

## Troubleshooting

Find more information about how to deal with common errors related to Bitnami's Helm charts in [this troubleshooting guide](https://docs.bitnami.com/general/how-to/troubleshoot-helm-chart-issues).

## Upgrading

### To 12.0.0

This version updates MinIO&reg; to major version 2023. All gateway features have been removed from Minio since upstream completely dropped this feature. The related options have been removed in version 12.1.0.

### To 11.0.0

This version deprecates the usage of `MINIO_ACCESS_KEY` and `MINIO_SECRET_KEY` environment variables in MINIO&reg; container in favor of `MINIO_ROOT_USER` and `MINIO_ROOT_PASSWORD`.

If you were already using the new variables, no issues are expected during upgrade.

### To 9.0.0

This version updates MinIO&reg; authentication parameters so they're aligned with the [current terminology](https://docs.min.io/minio/baremetal/security/minio-identity-management/user-management.html#minio-users-root). As a result the following parameters have been affected:

- `accessKey.password` has been renamed to `auth.rootUser`.
- `secretKey.password` has been renamed to `auth.rootPassword`.
- `accessKey.forcePassword` and `secretKey.forcePassword` have been unified into `auth.forcePassword`.
- `existingSecret`, `useCredentialsFile` and `forceNewKeys` have been renamed to `auth.existingSecret`, `auth.useCredentialsFiles` and `forceNewKeys`, respectively.

### To 8.0.0

This version updates MinIO&reg; after some major changes, affecting its Web UI. MinIO&reg; has replaced its MinIO&reg; Browser with the MinIO&reg; Console, and Web UI has been moved to a separated port. As a result the following variables have been affected:

- `service.port` has been slit into `service.ports.api` (default: 9000) and `service.ports.console` (default: 9001).
- `containerPort` has been slit into `containerPorts.api` (default: 9000) and `containerPort.console` (default: 9001).
- `service.nodePort`has been slit into `service.nodePorts.api` and `service.nodePorts.console`.
- Service port `minio` has been replaced with `minio-api` and `minio-console` with target ports minio-api and minio-console respectively.
- Liveness, readiness and startup probes now use port `minio-console` instead of `minio`.

Please note that Web UI, previously running on port 9000 will now use port 9001 leaving port 9000 for the MinIO&reg; Server API.

### To 7.0.0

This version introduces pod and container securityContext support. The previous configuration of `securityContext` has moved to `podSecurityContext` and `containerSecurityContext`. Apart from this case, no issues are expected to appear when upgrading.

### To 5.0.0

This version standardizes the way of defining Ingress rules. When configuring a single hostname for the Ingress rule, set the `ingress.hostname` value. When defining more than one, set the `ingress.extraHosts` array. Apart from this case, no issues are expected to appear when upgrading.

### To 4.1.0

This version introduces `bitnami/common`, a [library chart](https://helm.sh/docs/topics/library_charts/#helm) as a dependency. More documentation about this new utility could be found [here](https://github.com/bitnami/charts/tree/main/bitnami/common#bitnami-common-library-chart). Please, make sure that you have updated the chart dependencies before executing any upgrade.

### To 4.0.0

[On November 13, 2020, Helm v2 support was formally finished](https://github.com/helm/charts#status-of-the-project), this major version is the result of the required changes applied to the Helm Chart to be able to incorporate the different features added in Helm v3 and to be consistent with the Helm project itself regarding the Helm v2 EOL.

#### What changes were introduced in this major version?

- Previous versions of this Helm Chart use `apiVersion: v1` (installable by both Helm 2 and 3), this Helm Chart was updated to `apiVersion: v2` (installable by Helm 3 only). [Here](https://helm.sh/docs/topics/charts/#the-apiversion-field) you can find more information about the `apiVersion` field.
- The different fields present in the *Chart.yaml* file has been ordered alphabetically in a homogeneous way for all the Bitnami Helm Charts

#### Considerations when upgrading to this version

- If you want to upgrade to this version from a previous one installed with Helm v3, you shouldn't face any issues
- If you want to upgrade to this version using Helm v2, this scenario is not supported as this version doesn't support Helm v2 anymore
- If you installed the previous version with Helm v2 and wants to upgrade to this version with Helm v3, please refer to the [official Helm documentation](https://helm.sh/docs/topics/v2_v3_migration/#migration-use-cases) about migrating from Helm v2 to v3

#### Useful links

- <https://docs.bitnami.com/tutorials/resolve-helm2-helm3-post-migration-issues/>
- <https://helm.sh/docs/topics/v2_v3_migration/>
- <https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/>

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