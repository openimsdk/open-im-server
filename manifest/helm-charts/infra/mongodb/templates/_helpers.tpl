{{/*
Copyright VMware, Inc.
SPDX-License-Identifier: APACHE-2.0
*/}}

{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "mongodb.name" -}}
{{- include "common.names.name" . -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "mongodb.fullname" -}}
{{- include "common.names.fullname" . -}}
{{- end -}}

{{/*
Create a default mongo service name which can be overridden.
*/}}
{{- define "mongodb.service.nameOverride" -}}
    {{- if and .Values.service .Values.service.nameOverride -}}
        {{- print .Values.service.nameOverride -}}
    {{- else -}}
        {{- if eq .Values.architecture "replicaset" -}}
            {{- printf "%s-headless" (include "mongodb.fullname" .) -}}
        {{- else -}}
            {{- printf "%s" (include "mongodb.fullname" .) -}}
        {{- end -}}
    {{- end -}}
{{- end -}}

{{/*
Create a default mongo arbiter service name which can be overridden.
*/}}
{{- define "mongodb.arbiter.service.nameOverride" -}}
    {{- if and .Values.arbiter.service .Values.arbiter.service.nameOverride -}}
        {{- print .Values.arbiter.service.nameOverride -}}
    {{- else -}}
        {{- printf "%s-arbiter-headless" (include "mongodb.fullname" .) -}}
    {{- end }}
{{- end }}

{{/*
Return the proper MongoDB&reg; image name
*/}}
{{- define "mongodb.image" -}}
{{- include "common.images.image" (dict "imageRoot" .Values.image "global" .Values.global) -}}
{{- end -}}

{{/*
Return the proper image name (for the metrics image)
*/}}
{{- define "mongodb.metrics.image" -}}
{{- include "common.images.image" (dict "imageRoot" .Values.metrics.image "global" .Values.global) -}}
{{- end -}}

{{/*
Return the proper image name (for the init container volume-permissions image)
*/}}
{{- define "mongodb.volumePermissions.image" -}}
{{- include "common.images.image" (dict "imageRoot" .Values.volumePermissions.image "global" .Values.global) -}}
{{- end -}}

{{/*
Return the proper image name (for the init container auto-discovery image)
*/}}
{{- define "mongodb.externalAccess.autoDiscovery.image" -}}
{{- include "common.images.image" (dict "imageRoot" .Values.externalAccess.autoDiscovery.image "global" .Values.global) -}}
{{- end -}}

{{/*
Return the proper image name (for the TLS Certs image)
*/}}
{{- define "mongodb.tls.image" -}}
{{- include "common.images.image" (dict "imageRoot" .Values.tls.image "global" .Values.global) -}}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "mongodb.imagePullSecrets" -}}
{{- include "common.images.pullSecrets" (dict "images" (list .Values.image .Values.metrics.image .Values.volumePermissions.image .Values.tls.image) "global" .Values.global) -}}
{{- end -}}

{{/*
Allow the release namespace to be overridden for multi-namespace deployments in combined charts.
*/}}
{{- define "mongodb.namespace" -}}
    {{- if and .Values.global .Values.global.namespaceOverride -}}
        {{- print .Values.global.namespaceOverride -}}
    {{- else -}}
        {{- print .Release.Namespace -}}
    {{- end }}
{{- end -}}
{{- define "mongodb.serviceMonitor.namespace" -}}
    {{- if .Values.metrics.serviceMonitor.namespace -}}
        {{- print .Values.metrics.serviceMonitor.namespace -}}
    {{- else -}}
        {{- include "mongodb.namespace" . -}}
    {{- end }}
{{- end -}}
{{- define "mongodb.prometheusRule.namespace" -}}
    {{- if .Values.metrics.prometheusRule.namespace -}}
        {{- print .Values.metrics.prometheusRule.namespace -}}
    {{- else -}}
        {{- include "mongodb.namespace" . -}}
    {{- end }}
{{- end -}}

{{/*
Returns the proper service account name depending if an explicit service account name is set
in the values file. If the name is not set it will default to either mongodb.fullname if serviceAccount.create
is true or default otherwise.
*/}}
{{- define "mongodb.serviceAccountName" -}}
    {{- if .Values.serviceAccount.create -}}
        {{- default (include "mongodb.fullname" .) (print .Values.serviceAccount.name) -}}
    {{- else -}}
        {{- default "default" (print .Values.serviceAccount.name) -}}
    {{- end -}}
{{- end -}}

{{/*
Return the list of custom users to create during the initialization (string format)
*/}}
{{- define "mongodb.customUsers" -}}
    {{- $customUsers := list -}}
    {{- if .Values.auth.username -}}
        {{- $customUsers = append $customUsers .Values.auth.username }}
    {{- end }}
    {{- range .Values.auth.usernames }}
        {{- $customUsers = append $customUsers . }}
    {{- end }}
    {{- printf "%s" (default "" (join "," $customUsers)) -}}
{{- end -}}

{{/*
Return the list of passwords for the custom users (string format)
*/}}
{{- define "mongodb.customPasswords" -}}
    {{- $customPasswords := list -}}
    {{- if .Values.auth.password -}}
        {{- $customPasswords = append $customPasswords .Values.auth.password }}
    {{- end }}
    {{- range .Values.auth.passwords }}
        {{- $customPasswords = append $customPasswords . }}
    {{- end }}
    {{- printf "%s" (default "" (join "," $customPasswords)) -}}
{{- end -}}

{{/*
Return the list of custom databases to create during the initialization (string format)
*/}}
{{- define "mongodb.customDatabases" -}}
    {{- $customDatabases := list -}}
    {{- if .Values.auth.database -}}
        {{- $customDatabases = append $customDatabases .Values.auth.database }}
    {{- end }}
    {{- range .Values.auth.databases }}
        {{- $customDatabases = append $customDatabases . }}
    {{- end }}
    {{- printf "%s" (default "" (join "," $customDatabases)) -}}
{{- end -}}

{{/*
Return the configmap with the MongoDB&reg; configuration
*/}}
{{- define "mongodb.configmapName" -}}
{{- if .Values.existingConfigmap -}}
    {{- printf "%s" (tpl .Values.existingConfigmap $) -}}
{{- else -}}
    {{- printf "%s" (include "mongodb.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a configmap object should be created for MongoDB&reg;
*/}}
{{- define "mongodb.createConfigmap" -}}
{{- if and .Values.configuration (not .Values.existingConfigmap) }}
    {{- true -}}
{{- else -}}
{{- end -}}
{{- end -}}

{{/*
Return the secret with MongoDB&reg; credentials
*/}}
{{- define "mongodb.secretName" -}}
    {{- if .Values.auth.existingSecret -}}
        {{- printf "%s" (tpl .Values.auth.existingSecret $) -}}
    {{- else -}}
        {{- printf "%s" (include "mongodb.fullname" .) -}}
    {{- end -}}
{{- end -}}

{{/*
Return true if a secret object should be created for MongoDB&reg;
*/}}
{{- define "mongodb.createSecret" -}}
{{- if and .Values.auth.enabled (not .Values.auth.existingSecret) }}
    {{- true -}}
{{- else -}}
{{- end -}}
{{- end -}}

{{/*
Get the initialization scripts ConfigMap name.
*/}}
{{- define "mongodb.initdbScriptsCM" -}}
{{- if .Values.initdbScriptsConfigMap -}}
{{- printf "%s" .Values.initdbScriptsConfigMap -}}
{{- else -}}
{{- printf "%s-init-scripts" (include "mongodb.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if the Arbiter should be deployed
*/}}
{{- define "mongodb.arbiter.enabled" -}}
{{- if and (eq .Values.architecture "replicaset") .Values.arbiter.enabled }}
    {{- true -}}
{{- else -}}
{{- end -}}
{{- end -}}

{{/*
Return the configmap with the MongoDB&reg; configuration for the Arbiter
*/}}
{{- define "mongodb.arbiter.configmapName" -}}
{{- if .Values.arbiter.existingConfigmap -}}
    {{- printf "%s" (tpl .Values.arbiter.existingConfigmap $) -}}
{{- else -}}
    {{- printf "%s-arbiter" (include "mongodb.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a configmap object should be created for MongoDB&reg; Arbiter
*/}}
{{- define "mongodb.arbiter.createConfigmap" -}}
{{- if and (eq .Values.architecture "replicaset") .Values.arbiter.enabled .Values.arbiter.configuration (not .Values.arbiter.existingConfigmap) }}
    {{- true -}}
{{- else -}}
{{- end -}}
{{- end -}}

{{/*
Return true if the Hidden should be deployed
*/}}
{{- define "mongodb.hidden.enabled" -}}
{{- if and (eq .Values.architecture "replicaset") .Values.hidden.enabled }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the configmap with the MongoDB&reg; configuration for the Hidden
*/}}
{{- define "mongodb.hidden.configmapName" -}}
{{- if .Values.hidden.existingConfigmap -}}
    {{- printf "%s" (tpl .Values.hidden.existingConfigmap $) -}}
{{- else -}}
    {{- printf "%s-hidden" (include "mongodb.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a configmap object should be created for MongoDB&reg; Hidden
*/}}
{{- define "mongodb.hidden.createConfigmap" -}}
{{- if and  (include "mongodb.hidden.enabled" .) .Values.hidden.enabled .Values.hidden.configuration (not .Values.hidden.existingConfigmap) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Compile all warnings into a single message, and call fail.
*/}}
{{- define "mongodb.validateValues" -}}
{{- $messages := list -}}
{{- $messages := append $messages (include "mongodb.validateValues.pspAndRBAC" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.architecture" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.customUsersDBs" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.customUsersDBsLength" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.externalAccessServiceType" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.loadBalancerIPsListLength" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.nodePortListLength" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.externalAccessAutoDiscoveryRBAC" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.replicaset.existingSecrets" .) -}}
{{- $messages := append $messages (include "mongodb.validateValues.hidden.existingSecrets" .) -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message | fail -}}
{{- end -}}
{{- end -}}

{{/* Validate RBAC is created when using PSP */}}
{{- define "mongodb.validateValues.pspAndRBAC" -}}
{{- if and (.Values.podSecurityPolicy.create) (not .Values.rbac.create) -}}
mongodb: podSecurityPolicy.create, rbac.create
    Both podSecurityPolicy.create and rbac.create must be true, if you want
    to create podSecurityPolicy
{{- end -}}
{{- end -}}

{{/* Validate values of MongoDB&reg; - must provide a valid architecture */}}
{{- define "mongodb.validateValues.architecture" -}}
{{- if and (ne .Values.architecture "standalone") (ne .Values.architecture "replicaset") -}}
mongodb: architecture
    Invalid architecture selected. Valid values are "standalone" and
    "replicaset". Please set a valid architecture (--set mongodb.architecture="xxxx")
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - both auth.usernames and auth.databases are necessary
to create a custom user and database during 1st initialization
*/}}
{{- define "mongodb.validateValues.customUsersDBs" -}}
{{- $customUsers := include "mongodb.customUsers" . -}}
{{- $customDatabases := include "mongodb.customDatabases" . -}}
{{- if or (and (empty $customUsers) (not (empty $customDatabases))) (and (not (empty $customUsers)) (empty $customDatabases)) }}
mongodb: auth.usernames, auth.databases
    Both auth.usernames and auth.databases must be provided to create
    custom users and databases during 1st initialization.
    Please set both of them (--set auth.usernames[0]="xxxx",auth.databases[0]="yyyy")
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - both auth.usernames and auth.databases arrays should have the same length
to create a custom user and database during 1st initialization
*/}}
{{- define "mongodb.validateValues.customUsersDBsLength" -}}
{{- if ne (len .Values.auth.usernames) (len .Values.auth.databases) }}
mongodb: auth.usernames, auth.databases
    Both auth.usernames and auth.databases arrays should have the same length
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - service type for external access
*/}}
{{- define "mongodb.validateValues.externalAccessServiceType" -}}
{{- if and (eq .Values.architecture "replicaset") (not (eq .Values.externalAccess.service.type "NodePort")) (not (eq .Values.externalAccess.service.type "LoadBalancer")) (not (eq .Values.externalAccess.service.type "ClusterIP")) -}}
mongodb: externalAccess.service.type
    Available service type for external access are NodePort, LoadBalancer or ClusterIP.
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - number of replicas must be the same than LoadBalancer IPs list
*/}}
{{- define "mongodb.validateValues.loadBalancerIPsListLength" -}}
{{- $replicaCount := int .Values.replicaCount }}
{{- $loadBalancerListLength := len .Values.externalAccess.service.loadBalancerIPs }}
{{- if and (eq .Values.architecture "replicaset") .Values.externalAccess.enabled (not .Values.externalAccess.autoDiscovery.enabled ) (eq .Values.externalAccess.service.type "LoadBalancer") (not (eq $replicaCount $loadBalancerListLength )) -}}
mongodb: .Values.externalAccess.service.loadBalancerIPs
    Number of replicas and loadBalancerIPs array length must be the same.
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - number of replicas must be the same than NodePort list
*/}}
{{- define "mongodb.validateValues.nodePortListLength" -}}
{{- $replicaCount := int .Values.replicaCount }}
{{- $nodePortListLength := len .Values.externalAccess.service.nodePorts }}
{{- if and (eq .Values.architecture "replicaset") .Values.externalAccess.enabled (eq .Values.externalAccess.service.type "NodePort") (not (eq $replicaCount $nodePortListLength )) -}}
mongodb: .Values.externalAccess.service.nodePorts
    Number of replicas and nodePorts array length must be the same.
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - RBAC should be enabled when autoDiscovery is enabled
*/}}
{{- define "mongodb.validateValues.externalAccessAutoDiscoveryRBAC" -}}
{{- if and (eq .Values.architecture "replicaset") .Values.externalAccess.enabled .Values.externalAccess.autoDiscovery.enabled (not .Values.rbac.create ) }}
mongodb: rbac.create
    By specifying "externalAccess.enabled=true" and "externalAccess.autoDiscovery.enabled=true"
    an initContainer will be used to autodetect the external IPs/ports by querying the
    K8s API. Please note this initContainer requires specific RBAC resources. You can create them
    by specifying "--set rbac.create=true".
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - Number of replicaset secrets must be the same than number of replicaset nodes.
*/}}
{{- define "mongodb.validateValues.replicaset.existingSecrets" -}}
{{- if and .Values.tls.enabled (eq .Values.architecture "replicaset") (not (empty .Values.tls.replicaset.existingSecrets)) }}
{{- $nbSecrets := len .Values.tls.replicaset.existingSecrets -}}
{{- if not (eq $nbSecrets (int .Values.replicaCount)) }}
mongodb: tls.replicaset.existingSecrets
    tls.replicaset.existingSecrets Number of secrets and number of replicaset nodes must be the same.
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; - Number of hidden secrets must be the same than number of hidden nodes.
*/}}
{{- define "mongodb.validateValues.hidden.existingSecrets" -}}
{{- if and .Values.tls.enabled (include "mongodb.hidden.enabled" .) (not (empty .Values.tls.hidden.existingSecrets)) }}
{{- $nbSecrets := len .Values.tls.hidden.existingSecrets -}}
{{- if not (eq $nbSecrets (int .Values.hidden.replicaCount)) }}
mongodb: tls.hidden.existingSecrets
    tls.hidden.existingSecrets Number of secrets and number of hidden nodes must be the same.
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Validate values of MongoDB&reg; exporter URI string - auth.enabled and/or tls.enabled must be enabled or it defaults
*/}}
{{- define "mongodb.mongodb_exporter.uri" -}}
    {{- $uriTlsArgs := ternary "tls=true&tlsCertificateKeyFile=/certs/mongodb.pem&tlsCAFile=/certs/mongodb-ca-cert" "" .Values.tls.enabled -}}
    {{- if .Values.metrics.username }}
        {{- $uriAuth := ternary "$(echo $MONGODB_METRICS_USERNAME | sed -r \"s/@/%40/g;s/:/%3A/g\"):$(echo $MONGODB_METRICS_PASSWORD | sed -r \"s/@/%40/g;s/:/%3A/g\")@" "" .Values.auth.enabled -}}
        {{- printf "mongodb://%slocalhost:%d/admin?%s" $uriAuth (int .Values.containerPorts.mongodb) $uriTlsArgs -}}
    {{- else -}}
        {{- $uriAuth := ternary "$MONGODB_ROOT_USER:$(echo $MONGODB_ROOT_PASSWORD | sed -r \"s/@/%40/g;s/:/%3A/g\")@" "" .Values.auth.enabled -}}
        {{- printf "mongodb://%slocalhost:%d/admin?%s" $uriAuth (int .Values.containerPorts.mongodb) $uriTlsArgs -}}
    {{- end -}}
{{- end -}}

{{/*
Return the appropriate apiGroup for PodSecurityPolicy.
*/}}
{{- define "podSecurityPolicy.apiGroup" -}}
{{- if semverCompare ">=1.14-0" .Capabilities.KubeVersion.GitVersion -}}
{{- print "policy" -}}
{{- else -}}
{{- print "extensions" -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS secret object should be created
*/}}
{{- define "mongodb.createTlsSecret" -}}
{{- if and .Values.tls.enabled (not .Values.tls.existingSecret) (include "mongodb.autoGenerateCerts" .) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the secret containing MongoDB&reg; TLS certificates
*/}}
{{- define "mongodb.tlsSecretName" -}}
{{- $secretName := .Values.tls.existingSecret -}}
{{- if $secretName -}}
    {{- printf "%s" (tpl $secretName $) -}}
{{- else -}}
    {{- printf "%s-ca" (include "mongodb.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if certificates must be auto generated
*/}}
{{- define "mongodb.autoGenerateCerts" -}}
{{- $standalone := (eq .Values.architecture "standalone") | ternary (not .Values.tls.standalone.existingSecret) true -}}
{{- $replicaset := (eq .Values.architecture "replicaset") | ternary (empty .Values.tls.replicaset.existingSecrets) true -}}
{{- $arbiter := (eq (include "mongodb.arbiter.enabled" .) "true") | ternary (not .Values.tls.arbiter.existingSecret) true -}}
{{- $hidden := (eq (include "mongodb.hidden.enabled" .) "true") | ternary (empty .Values.tls.hidden.existingSecrets) true -}}
{{- if and $standalone $replicaset $arbiter $hidden -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Generate argument list for mongodb-exporter
reference: https://github.com/percona/mongodb_exporter/blob/main/REFERENCE.md
*/}}
{{- define "mongodb.exporterArgs" -}}
{{- with .Values.metrics.collector -}}
{{- ternary " --collect-all" "" .all -}}
{{- ternary " --collector.diagnosticdata" "" .diagnosticdata -}}
{{- ternary " --collector.replicasetstatus" "" .replicasetstatus -}}
{{- ternary " --collector.dbstats" "" .dbstats -}}
{{- ternary " --collector.topmetrics" "" .topmetrics -}}
{{- ternary " --collector.indexstats" "" .indexstats -}}
{{- ternary " --collector.collstats" "" .collstats -}}
{{- if .collstatsColls -}}
{{- " --mongodb.collstats-colls=" -}}
{{- join "," .collstatsColls -}}
{{- end -}}
{{- if .indexstatsColls -}}
{{- " --mongodb.indexstats-colls=" -}}
{{- join "," .indexstatsColls -}}
{{- end -}}
{{- $limitArg := print " --collector.collstats-limit=" .collstatsLimit -}}
{{- ne (print .collstatsLimit) "0" | ternary $limitArg "" -}}
{{- end -}}
{{- ternary " --compatible-mode" "" .Values.metrics.compatibleMode -}}
{{- end -}}