{{/*
Copyright VMware, Inc.
SPDX-License-Identifier: APACHE-2.0
*/}}

{{/* vim: set filetype=mustache: */}}

{{/*
Expand the name of the chart.
*/}}
{{- define "kafka.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified zookeeper name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "kafka.zookeeper.fullname" -}}
{{- if .Values.zookeeper.fullnameOverride -}}
{{- .Values.zookeeper.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default "zookeeper" .Values.zookeeper.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{/*
 Create the name of the service account to use
 */}}
{{- define "kafka.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "common.names.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return the proper Kafka image name
*/}}
{{- define "kafka.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper image name (for the init container auto-discovery image)
*/}}
{{- define "kafka.externalAccess.autoDiscovery.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.externalAccess.autoDiscovery.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper image name (for the init container volume-permissions image)
*/}}
{{- define "kafka.volumePermissions.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.volumePermissions.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper Kafka exporter image name
*/}}
{{- define "kafka.metrics.kafka.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.metrics.kafka.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper JMX exporter image name
*/}}
{{- define "kafka.metrics.jmx.image" -}}
{{ include "common.images.image" (dict "imageRoot" .Values.metrics.jmx.image "global" .Values.global) }}
{{- end -}}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "kafka.imagePullSecrets" -}}
{{ include "common.images.pullSecrets" (dict "images" (list .Values.image .Values.externalAccess.autoDiscovery.image .Values.volumePermissions.image .Values.metrics.kafka.image .Values.metrics.jmx.image) "global" .Values.global) }}
{{- end -}}

{{/*
Create a default fully qualified Kafka exporter name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
*/}}
{{- define "kafka.metrics.kafka.fullname" -}}
  {{- printf "%s-exporter" (include "common.names.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end -}}

{{/*
 Create the name of the service account to use for Kafka exporter pods
 */}}
{{- define "kafka.metrics.kafka.serviceAccountName" -}}
{{- if .Values.metrics.kafka.serviceAccount.create -}}
    {{ default (include "kafka.metrics.kafka.fullname" .) .Values.metrics.kafka.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.metrics.kafka.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return true if encryption via TLS for client connections should be configured
*/}}
{{- define "kafka.sslEnabled" -}}
{{- $res := "" -}}
{{- $listeners := list .Values.listeners.client .Values.listeners.interbroker -}}
{{- range $i := .Values.listeners.extraListeners -}}
{{- $listeners = append $listeners $i -}}
{{- end -}}
{{- if and .Values.externalAccess.enabled -}}
{{- $listeners = append $listeners .Values.listeners.external -}}
{{- end -}}
{{- if and .Values.kraft.enabled -}}
{{- $listeners = append $listeners .Values.listeners.controller -}}
{{- end -}}
{{- range $listener := $listeners -}}
{{- if regexFind "SSL" (upper $listener.protocol) -}}
{{- $res = "true" -}}
{{- end -}}
{{- end -}}
{{- if $res -}}
{{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if SASL connections should be configured
*/}}
{{- define "kafka.saslEnabled" -}}
{{- $res := "" -}}
{{- if (include "kafka.client.saslEnabled" .) -}}
{{- $res = "true" -}}
{{- else -}}
{{- $listeners := list .Values.listeners.interbroker -}}
{{- if and .Values.kraft.enabled -}}
{{- $listeners = append $listeners .Values.listeners.controller -}}
{{- end -}}
{{- range $listener := $listeners -}}
{{- if regexFind "SASL" (upper $listener.protocol) -}}
{{- $res = "true" -}}
{{- end -}}
{{- end -}}
{{- end -}}
{{- if $res -}}
{{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if SASL connections should be configured
*/}}
{{- define "kafka.client.saslEnabled" -}}
{{- $res := "" -}}
{{- $listeners := list .Values.listeners.client -}}
{{- range $i := .Values.listeners.extraListeners -}}
{{- $listeners = append $listeners $i -}}
{{- end -}}
{{- if and .Values.externalAccess.enabled -}}
{{- $listeners = append $listeners .Values.listeners.external -}}
{{- end -}}
{{- range $listener := $listeners -}}
{{- if regexFind "SASL" (upper $listener.protocol) -}}
{{- $res = "true" -}}
{{- end -}}
{{- end -}}
{{- if $res -}}
{{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka SASL credentials secret
*/}}
{{- define "kafka.saslSecretName" -}}
{{- if .Values.sasl.existingSecret -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.sasl.existingSecret "context" $) -}}
{{- else -}}
    {{- printf "%s-user-passwords" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a SASL credentials secret object should be created
*/}}
{{- define "kafka.createSaslSecret" -}}
{{- $secretName := .Values.sasl.existingSecret -}}
{{- if and (or (include "kafka.saslEnabled" .) (or .Values.zookeeper.auth.client.enabled .Values.sasl.zookeeper.user)) (empty $secretName) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.tlsSecretName" -}}
{{- if .Values.tls.existingSecret -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.tls.existingSecret "context" $) -}}
{{- else -}}
    {{- printf "%s-tls" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.createTlsSecret" -}}
{{- if and (include "kafka.sslEnabled" .) (empty .Values.tls.existingSecret) .Values.tls.autoGenerated -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka TLS credentials secret
*/}}
{{- define "kafka.tlsPasswordsSecretName" -}}
{{- if .Values.tls.passwordsSecret -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.tls.passwordsSecret "context" $) -}}
{{- else -}}
    {{- printf "%s-tls-passwords" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.createTlsPasswordsSecret" -}}
{{- $secretName := .Values.tls.passwordsSecret -}}
{{- if and (include "kafka.sslEnabled" .) (or (empty $secretName) .Values.tls.autoGenerated ) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka TLS credentials secret
*/}}
{{- define "kafka.zookeeper.tlsPasswordsSecretName" -}}
{{- if .Values.tls.zookeeper.passwordsSecret -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.tls.zookeeper.passwordsSecret "context" $) -}}
{{- else -}}
    {{- printf "%s-zookeeper-tls-passwords" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a TLS credentials secret object should be created
*/}}
{{- define "kafka.zookeeper.createTlsPasswordsSecret" -}}
{{- $secretName := .Values.tls.zookeeper.passwordsSecret -}}
{{- if and .Values.tls.zookeeper.enabled (or (empty $secretName) .Values.tls.zookeeper.keystorePassword .Values.tls.zookeeper.truststorePassword ) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Returns the secret name for the Kafka Provisioning client
*/}}
{{- define "kafka.client.passwordsSecretName" -}}
{{- if .Values.provisioning.auth.tls.passwordsSecret -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.provisioning.auth.tls.passwordsSecret "context" $) -}}
{{- else -}}
    {{- printf "%s-client-secret" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account to use for the Kafka Provisioning client
*/}}
{{- define "kafka.provisioning.serviceAccountName" -}}
{{- if .Values.provisioning.serviceAccount.create -}}
    {{ default (include "common.names.fullname" .) .Values.provisioning.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.provisioning.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka controller-eligible configuration configmap
*/}}
{{- define "kafka.controller.configmapName" -}}
{{- if .Values.controller.existingConfigmap -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.controller.existingConfigmap "context" $) -}}
{{- else if .Values.existingConfigmap -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.existingConfigmap "context" $) -}}
{{- else -}}
    {{- printf "%s-controller-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka controller-eligible secret configuration
*/}}
{{- define "kafka.controller.secretConfigName" -}}
{{- if .Values.controller.existingSecretConfig -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.controller.existingSecretConfig "context" $) -}}
{{- else if .Values.existingSecretConfig -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.existingSecretConfig "context" $) -}}
{{- else -}}
    {{- printf "%s-controller-secret-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka controller-eligible secret configuration values 
*/}}
{{- define "kafka.controller.secretConfig" -}}
{{- if .Values.secretConfig }}
{{- include "common.tplvalues.render" ( dict "value" .Values.secretConfig "context" $ ) }}
{{- end }}
{{- if .Values.controller.secretConfig }}
{{- include "common.tplvalues.render" ( dict "value" .Values.controller.secretConfig "context" $ ) }}
{{- end }}
{{- end -}}

{{/*
Return true if a configmap object should be created for controller-eligible pods
*/}}
{{- define "kafka.controller.createConfigmap" -}}
{{- if and (not .Values.controller.existingConfigmap) (not .Values.existingConfigmap) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a secret object with config should be created for controller-eligible pods
*/}}
{{- define "kafka.controller.createSecretConfig" -}}
{{- if and (or .Values.controller.secretConfig .Values.secretConfig) (and (not .Values.controller.existingSecretConfig) (not .Values.existingSecretConfig)) }}
    {{- true -}}
{{- end -}}
{{- end -}}
{{/*
Return true if a secret object with config exists for controller-eligible pods
*/}}
{{- define "kafka.controller.secretConfigExists" -}}
{{- if or .Values.controller.secretConfig .Values.secretConfig .Values.controller.existingSecretConfig .Values.existingSecretConfig }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka broker configuration configmap
*/}}
{{- define "kafka.broker.configmapName" -}}
{{- if .Values.broker.existingConfigmap -}}
    {{- printf "%s" (tpl .Values.broker.existingConfigmap $) -}}
{{- else if .Values.existingConfigmap -}}
    {{- printf "%s" (tpl .Values.existingConfigmap $) -}}
{{- else -}}
    {{- printf "%s-broker-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka broker secret configuration
*/}}
{{- define "kafka.broker.secretConfigName" -}}
{{- if .Values.broker.existingSecretConfig -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.broker.existingSecretConfig "context" $) -}}
{{- else if .Values.existingSecretConfig -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.existingSecretConfig "context" $) -}}
{{- else -}}
    {{- printf "%s-broker-secret-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka broker secret configuration values 
*/}}
{{- define "kafka.broker.secretConfig" -}}
{{- if .Values.secretConfig }}
{{- include "common.tplvalues.render" ( dict "value" .Values.secretConfig "context" $ ) }}
{{- end }}
{{- if .Values.broker.secretConfig }}
{{- include "common.tplvalues.render" ( dict "value" .Values.broker.secretConfig "context" $ ) }}
{{- end }}
{{- end -}}

{{/*
Return true if a configmap object should be created for broker pods
*/}}
{{- define "kafka.broker.createConfigmap" -}}
{{- if and (not .Values.broker.existingConfigmap) (not .Values.existingConfigmap) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a secret object with config should be created for broker pods
*/}}
{{- define "kafka.broker.createSecretConfig" -}}
{{- if and (or .Values.broker.secretConfig .Values.secretConfig) (and (not .Values.broker.existingSecretConfig) (not .Values.existingSecretConfig)) }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a secret object with config exists for broker pods
*/}}
{{- define "kafka.broker.secretConfigExists" -}}
{{- if or .Values.broker.secretConfig .Values.secretConfig .Values.broker.existingSecretConfig .Values.existingSecretConfig }}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka log4j ConfigMap name.
*/}}
{{- define "kafka.log4j.configMapName" -}}
{{- if .Values.existingLog4jConfigMap -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.existingLog4jConfigMap "context" $) -}}
{{- else -}}
    {{- printf "%s-log4j-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return the SASL mechanism to use for the Kafka exporter to access Kafka
The exporter uses a different nomenclature so we need to do this hack
*/}}
{{- define "kafka.metrics.kafka.saslMechanism" -}}
{{- $saslMechanisms := .Values.sasl.enabledMechanisms }}
{{- if contains "SCRAM-SHA-512" (upper $saslMechanisms) }}
    {{- print "scram-sha512" -}}
{{- else if contains "SCRAM-SHA-256" (upper $saslMechanisms) }}
    {{- print "scram-sha256" -}}
{{- else if contains "PLAIN" (upper $saslMechanisms) }}
    {{- print "plain" -}}
{{- end -}}
{{- end -}}

{{/*
Return the Kafka configuration configmap
*/}}
{{- define "kafka.metrics.jmx.configmapName" -}}
{{- if .Values.metrics.jmx.existingConfigmap -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.metrics.jmx.existingConfigmap "context" $) -}}
{{- else -}}
    {{ printf "%s-jmx-configuration" (include "common.names.fullname" .) -}}
{{- end -}}
{{- end -}}

{{/*
Return true if a configmap object should be created
*/}}
{{- define "kafka.metrics.jmx.createConfigmap" -}}
{{- if and .Values.metrics.jmx.enabled .Values.metrics.jmx.config (not .Values.metrics.jmx.existingConfigmap) -}}
    {{- true -}}
{{- end -}}
{{- end -}}

{{/*
Returns the Kafka listeners settings based on the listeners.* object
*/}}
{{- define "kafka.listeners" -}}
{{- if .context.Values.listeners.overrideListeners -}}
  {{- printf "%s" .context.Values.listeners.overrideListeners -}}
{{- else -}}
  {{- $listeners := list .context.Values.listeners.client .context.Values.listeners.interbroker -}}
  {{- if and .context.Values.externalAccess.enabled -}}
  {{- $listeners = append $listeners .context.Values.listeners.external -}}
  {{- end -}}
  {{- if and .context.Values.kraft.enabled .isController -}}
  {{- if and .context.Values.controller.controllerOnly -}}
  {{- $listeners = list .context.Values.listeners.controller -}}
  {{- else -}}
  {{- $listeners = append $listeners .context.Values.listeners.controller -}}
  {{- end -}}
  {{- end -}}
  {{- $res := list -}}
  {{- range $listener := $listeners -}}
  {{- $res = append $res (printf "%s://:%d" (upper $listener.name) (int $listener.containerPort)) -}}
  {{- end -}}
  {{- printf "%s" (join "," $res) -}}
{{- end -}}
{{- end -}}

{{/*
Returns the list of advertised listeners, although the advertised address will be replaced during each node init time
*/}}
{{- define "kafka.advertisedListeners" -}}
{{- if .Values.listeners.advertisedListeners -}}
  {{- printf "%s" .Values.listeners.advertisedListeners -}}
{{- else -}}
  {{- $listeners := list .Values.listeners.client .Values.listeners.interbroker -}}
  {{- range $i := .Values.listeners.extraListeners -}}
  {{- $listeners = append $listeners $i -}}
  {{- end -}}
  {{- $res := list -}}
  {{- range $listener := $listeners -}}
  {{- $res = append $res (printf "%s://advertised-address-placeholder:%d" (upper $listener.name) (int $listener.containerPort)) -}}
  {{- end -}}
  {{- printf "%s" (join "," $res) -}}
{{- end -}}
{{- end -}}

{{/*
Returns the value listener.security.protocol.map based on the values of 'listeners.*.protocol'
*/}}
{{- define "kafka.securityProtocolMap" -}}
{{- if .Values.listeners.securityProtocolMap -}}
  {{- printf "%s" .Values.listeners.securityProtocolMap -}}
{{- else -}}
  {{- $listeners := list .Values.listeners.client .Values.listeners.interbroker -}}
  {{- range $i := .Values.listeners.extraListeners -}}
  {{- $listeners = append $listeners $i -}}
  {{- end -}}
  {{- if .Values.kraft.enabled -}}
  {{- $listeners = append $listeners .Values.listeners.controller -}}
  {{- end -}}
  {{- if and .Values.externalAccess.enabled -}}
  {{- $listeners = append $listeners .Values.listeners.external -}}
  {{- end -}}
  {{- $res := list -}}
  {{- range $listener := $listeners -}}
  {{- $res = append $res (printf "%s:%s" (upper $listener.name) (upper $listener.protocol)) -}}
  {{- end -}}
  {{ printf "%s" (join "," $res)}}
{{- end -}}
{{- end -}}

{{/*
Returns the containerPorts for listeneres.extraListeners
*/}}
{{- define "kafka.extraListeners.containerPorts" -}}
{{- range $listener := .Values.listeners.extraListeners -}}
- name: {{ lower $listener.name}}
  containerPort: {{ $listener.containerPort }}
{{- end -}}
{{- end -}}

{{/*
Returns the zookeeper.connect setting value
*/}}
{{- define "kafka.zookeeperConnect" -}}
{{- if .Values.zookeeper.enabled -}}
{{- printf "%s:%s%s" (include "kafka.zookeeper.fullname" .) (ternary "3181" "2181" .Values.tls.zookeeper.enabled) (tpl .Values.zookeeperChrootPath .) -}}
{{- else -}}
{{- printf "%s%s" (join "," .Values.externalZookeeper.servers) (tpl .Values.zookeeperChrootPath .) -}}
{{- end -}}
{{- end -}}

{{/*
Returns the controller quorum voters based on the number of controller-eligible nodes
*/}}
{{- define "kafka.kraft.controllerQuorumVoters" -}}
{{- if .Values.kraft.controllerQuorumVoters -}}
    {{- include "common.tplvalues.render" (dict "value" .Values.kraft.controllerQuorumVoters "context" $) -}}
{{- else -}}
  {{- $controllerVoters := list -}}
  {{- $fullname := include "common.names.fullname" . -}}
  {{- $releaseNamespace := include "common.names.namespace" . -}}
  {{- range $i := until (int .Values.controller.replicaCount) -}}
  {{- $nodeId := add (int $i) (int $.Values.controller.minId) -}}
  {{- $nodeAddress := printf "%s-controller-%d.%s-controller-headless.%s.svc.%s:%d" $fullname (int $i) $fullname $releaseNamespace $.Values.clusterDomain (int $.Values.listeners.controller.containerPort) -}}
  {{- $controllerVoters = append $controllerVoters (printf "%d@%s" $nodeId $nodeAddress ) -}}
  {{- end -}}
  {{- join "," $controllerVoters -}}
{{- end -}}
{{- end -}}

{{/*
Section of the server.properties configmap shared by both controller-eligible and broker nodes
*/}}
{{- define "kafka.commonConfig" -}}
log.dir={{ printf "%s/data" .Values.controller.persistence.mountPath }}
{{- if or (include "kafka.saslEnabled" .) }}
sasl.enabled.mechanisms={{ upper .Values.sasl.enabledMechanisms }}
{{- end }}
# Interbroker configuration
inter.broker.listener.name={{ .Values.listeners.interbroker.name }}
{{- if regexFind "SASL" (upper .Values.listeners.interbroker.protocol) }}
sasl.mechanism.inter.broker.protocol={{ upper .Values.sasl.interBrokerMechanism }}
{{- end }}
{{- if (include "kafka.sslEnabled" .) }}
# TLS configuration
ssl.keystore.type=JKS
ssl.truststore.type=JKS
ssl.keystore.location=/opt/bitnami/kafka/config/certs/kafka.keystore.jks
ssl.truststore.location=/opt/bitnami/kafka/config/certs/kafka.truststore.jks
#ssl.keystore.password=
#ssl.truststore.password=
#ssl.key.password=
ssl.client.auth={{ .Values.tls.sslClientAuth }}
ssl.endpoint.identification.algorithm={{ .Values.tls.endpointIdentificationAlgorithm }}
{{- end }}
{{- if (include "kafka.saslEnabled" .) }}
# Listeners SASL JAAS configuration
{{- $listeners := list .Values.listeners.client .Values.listeners.interbroker }}
{{- range $i := .Values.listeners.extraListeners }}
{{- $listeners = append $listeners $i }}
{{- end }}
{{- if .Values.externalAccess.enabled }}
{{- $listeners = append $listeners .Values.listeners.external }}
{{- end }}
{{- range $listener := $listeners }}
{{- if and $listener.sslClientAuth (regexFind "SSL" (upper $listener.protocol)) }}
listener.name.{{lower $listener.name}}.ssl.client.auth={{ $listener.sslClientAuth }}
{{- end }}
{{- if regexFind "SASL" (upper $listener.protocol) }}
{{- range $mechanism := ( splitList "," $.Values.sasl.enabledMechanisms )}}
  {{- $securityModule := ternary "org.apache.kafka.common.security.plain.PlainLoginModule required" "org.apache.kafka.common.security.scram.ScramLoginModule required" (eq "PLAIN" (upper $mechanism)) }}
  {{- $saslJaasConfig := list $securityModule }}
  {{- if eq $listener.name $.Values.listeners.interbroker.name }}
  {{- $saslJaasConfig = append $saslJaasConfig (printf "username=\"%s\"" $.Values.sasl.interbroker.user) }}
  {{- $saslJaasConfig = append $saslJaasConfig (print "password=\"interbroker-password-placeholder\"") }}
  {{- end }}
  {{- if eq (upper $mechanism) "PLAIN" }}
  {{- if eq $listener.name $.Values.listeners.interbroker.name }}
  {{- $saslJaasConfig = append $saslJaasConfig (printf "user_%s=\"interbroker-password-placeholder\"" $.Values.sasl.interbroker.user) }}
  {{- end }}
  {{- range $i, $user := $.Values.sasl.client.users }}
  {{- $saslJaasConfig = append $saslJaasConfig (printf "user_%s=\"password-placeholder-%d\"" $user (int $i)) }}
  {{- end }}
  {{- end }}
listener.name.{{lower $listener.name}}.{{lower $mechanism}}.sasl.jaas.config={{ join " " $saslJaasConfig }};
{{- end }}
{{- end }}
{{- end }}
# End of SASL JAAS configuration
{{- end }}
{{- end -}}

{{/*
Zookeeper connection section of the server.properties
*/}}
{{- define "kafka.zookeeperConfig" -}}
zookeeper.connect={{ include "kafka.zookeeperConnect" . }}
#broker.id=
{{- if .Values.sasl.zookeeper.user }}
sasl.jaas.config=org.apache.kafka.common.security.plain.PlainLoginModule required \
    username="{{ .Values.sasl.zookeeper.user }}" \
    password="zookeeper-password-placeholder";
{{- end }}
{{- if and .Values.tls.zookeeper.enabled .Values.tls.zookeeper.existingSecret }}
zookeeper.clientCnxnSocket=org.apache.zookeeper.ClientCnxnSocketNetty
zookeeper.ssl.client.enable=true
zookeeper.ssl.keystore.location=/opt/bitnami/kafka/config/certs/zookeeper.keystore.jks
zookeeper.ssl.truststore.location=/opt/bitnami/kafka/config/certs/zookeeper.truststore.jks
zookeeper.ssl.hostnameVerification={{ .Values.tls.zookeeper.verifyHostname }}
#zookeeper.ssl.keystore.password=
#zookeeper.ssl.truststore.password=
{{- end }}
{{- end -}}

{{/*
Kraft section of the server.properties
*/}}
{{- define "kafka.kraftConfig" -}}
#node.id=
controller.listener.names={{ .Values.listeners.controller.name }}
controller.quorum.voters={{ include "kafka.kraft.controllerQuorumVoters" . }}
{{- $listener := $.Values.listeners.controller }}
{{- if and $listener.sslClientAuth (regexFind "SSL" (upper $listener.protocol)) }}
# Kraft Controller listener SSL settings
listener.name.{{lower $listener.name}}.ssl.client.auth={{ $listener.sslClientAuth }}
{{- end }}
{{- if regexFind "SASL" (upper $listener.protocol) }}
  {{- $mechanism := $.Values.sasl.controllerMechanism }}
  {{- $securityModule := ternary "org.apache.kafka.common.security.plain.PlainLoginModule required" "org.apache.kafka.common.security.scram.ScramLoginModule required" (eq "PLAIN" (upper $mechanism)) }}
  {{- $saslJaasConfig := list $securityModule }}
  {{- $saslJaasConfig = append $saslJaasConfig (printf "username=\"%s\"" $.Values.sasl.controller.user) }}
  {{- $saslJaasConfig = append $saslJaasConfig (print "password=\"controller-password-placeholder\"") }}
  {{- if eq (upper $mechanism) "PLAIN" }}
  {{- $saslJaasConfig = append $saslJaasConfig (printf "user_%s=\"controller-password-placeholder\"" $.Values.sasl.controller.user) }}
  {{- end }}
# Kraft Controller listener SASL settings
sasl.mechanism.controller.protocol={{ upper $mechanism }}
listener.name.{{lower $listener.name}}.sasl.enabled.mechanisms={{ upper $mechanism }}
listener.name.{{lower $listener.name}}.{{lower $mechanism }}.sasl.jaas.config={{ join " " $saslJaasConfig }};
{{- end }}
{{- end -}}

{{/*
Init container definition for Kafka initialization
*/}}
{{- define "kafka.prepareKafkaInitContainer" -}}
{{- $role := .role -}}
{{- $roleSettings := index .context.Values .role -}}
- name: kafka-init
  image: {{ include "kafka.image" .context }}
  imagePullPolicy: {{ .context.Values.image.pullPolicy }}
  {{- if $roleSettings.containerSecurityContext.enabled }}
  securityContext: {{- omit $roleSettings.containerSecurityContext "enabled" | toYaml | nindent 4 }}
  {{- end }}
  command:
    - /bin/bash
  args:
    - -ec
    - |
      /scripts/kafka-init.sh
  env:
    - name: BITNAMI_DEBUG
      value: {{ ternary "true" "false" (or .context.Values.image.debug .context.Values.diagnosticMode.enabled) | quote }}
    - name: MY_POD_NAME
      valueFrom:
        fieldRef:
            fieldPath: metadata.name
    - name: KAFKA_VOLUME_DIR
      value: {{ $roleSettings.persistence.mountPath | quote }}
    - name: KAFKA_MIN_ID
      value: {{ $roleSettings.minId | quote }}
    {{- if or (and (eq .role "broker") .context.Values.externalAccess.enabled) (and (eq .role "controller") .context.Values.externalAccess.enabled (or .context.Values.externalAccess.controller.forceExpose (not .context.Values.controller.controllerOnly))) }}
    {{- $externalAccess := index .context.Values.externalAccess .role }}
    - name: EXTERNAL_ACCESS_ENABLED
      value: "true"
    {{- if eq $externalAccess.service.type "LoadBalancer" }}
    {{- if not .context.Values.externalAccess.autoDiscovery.enabled }}
    - name: EXTERNAL_ACCESS_HOSTS_LIST
      value: {{ join "," (default $externalAccess.service.loadBalancerIPs $externalAccess.service.loadBalancerNames) | quote }}
    {{- end }}
    - name: EXTERNAL_ACCESS_PORT
      value: {{ $externalAccess.service.ports.external | quote }}
    {{- else if eq $externalAccess.service.type "NodePort" }}
    {{- if $externalAccess.service.domain }}
    - name: EXTERNAL_ACCESS_HOST
      value: {{ $externalAccess.service.domain | quote }}
    {{- else if and $externalAccess.service.usePodIPs .context.Values.externalAccess.autoDiscovery.enabled }}
    - name: MY_POD_IP
      valueFrom:
        fieldRef:
          fieldPath: status.podIP
    - name: EXTERNAL_ACCESS_HOST
      value: "$(MY_POD_IP)"
    {{- else if or $externalAccess.service.useHostIPs .context.Values.externalAccess.autoDiscovery.enabled }}
    - name: HOST_IP
      valueFrom:
        fieldRef:
          fieldPath: status.hostIP
    - name: EXTERNAL_ACCESS_HOST
      value: "$(HOST_IP)"
    {{- else if and $externalAccess.service.externalIPs (not .context.Values.externalAccess.autoDiscovery.enabled) }}
    - name: EXTERNAL_ACCESS_HOSTS_LIST
      value: {{ join "," $externalAccess.service.externalIPs }}
    {{- else }}
    - name: EXTERNAL_ACCESS_HOST_USE_PUBLIC_IP
      value: "true"
    {{- end }}
    {{- if not .context.Values.externalAccess.autoDiscovery.enabled }}
    {{- if and $externalAccess.service.externalIPs (empty $externalAccess.service.nodePorts)}}
    - name: EXTERNAL_ACCESS_PORT
      value: {{ $externalAccess.service.ports.external | quote }}
    {{- else }}
    - name: EXTERNAL_ACCESS_PORTS_LIST
      value: {{ join "," $externalAccess.service.nodePorts | quote }}
    {{- end }}
    {{- end }}
    {{- else if eq $externalAccess.service.type "ClusterIP" }}
    - name: EXTERNAL_ACCESS_HOST
      value: {{ $externalAccess.service.domain | quote }}
    - name: EXTERNAL_ACCESS_PORT
      value: {{ $externalAccess.service.ports.external | quote}}
    - name: EXTERNAL_ACCESS_PORT_AUTOINCREMENT
      value: "true"
    {{- end }}
    {{- end }}
    {{- if and (include "kafka.client.saslEnabled" .context ) .context.Values.sasl.client.users }}
    - name: KAFKA_CLIENT_USERS
      value: {{ join "," .context.Values.sasl.client.users | quote }}
    - name: KAFKA_CLIENT_PASSWORDS
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.saslSecretName" .context }}
          key: client-passwords
    {{- end }}
    {{- if regexFind "SASL" (upper .context.Values.listeners.interbroker.protocol) }}
    - name: KAFKA_INTER_BROKER_USER
      value: {{ .context.Values.sasl.interbroker.user | quote }}
    - name: KAFKA_INTER_BROKER_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.saslSecretName" .context }}
          key: inter-broker-password
    {{- end }}
    {{- if and .context.Values.kraft.enabled (regexFind "SASL" (upper .context.Values.listeners.controller.protocol)) }}
    - name: KAFKA_CONTROLLER_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.saslSecretName" .context }}
          key: controller-password
    {{- end }}
    {{- if (include "kafka.sslEnabled" .context )  }}
    - name: KAFKA_TLS_TYPE
      value: {{ ternary "PEM" "JKS" (or .context.Values.tls.autoGenerated (eq (upper .context.Values.tls.type) "PEM")) }}
    - name: KAFKA_TLS_KEYSTORE_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.tlsPasswordsSecretName" .context }}
          key: {{ .context.Values.tls.passwordsSecretKeystoreKey | quote }}
    - name: KAFKA_TLS_TRUSTSTORE_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.tlsPasswordsSecretName" .context }}
          key: {{ .context.Values.tls.passwordsSecretTruststoreKey | quote }}
    {{- if and (not .context.Values.tls.autoGenerated) (or .context.Values.tls.keyPassword (and .context.Values.tls.passwordsSecret .context.Values.tls.passwordsSecretPemPasswordKey)) }}
    - name: KAFKA_TLS_PEM_KEY_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.tlsPasswordsSecretName" .context }}
          key: {{ default "key-password" .context.Values.tls.passwordsSecretPemPasswordKey | quote }}
    {{- end }}
    {{- end }}
    {{- if or .context.Values.zookeeper.enabled .context.Values.externalZookeeper.servers }}
    {{- if .context.Values.sasl.zookeeper.user }}
    - name: KAFKA_ZOOKEEPER_USER
      value: {{ .context.Values.sasl.zookeeper.user | quote }}
    - name: KAFKA_ZOOKEEPER_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.saslSecretName" .context }}
          key: zookeeper-password
    {{- end }}
    {{- if .context.Values.tls.zookeeper.enabled }}
    {{- if and .context.Values.tls.zookeeper.passwordsSecretKeystoreKey (or .context.Values.tls.zookeeper.passwordsSecret .context.Values.tls.zookeeper.keystorePassword) }}
    - name: KAFKA_ZOOKEEPER_TLS_KEYSTORE_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.zookeeper.tlsPasswordsSecretName" .context }}
          key: {{ .context.Values.tls.zookeeper.passwordsSecretKeystoreKey | quote }}
    {{- end }}
    {{- if and .context.Values.tls.zookeeper.passwordsSecretTruststoreKey (or .context.Values.tls.zookeeper.passwordsSecret .context.Values.tls.zookeeper.truststorePassword) }}
    - name: KAFKA_ZOOKEEPER_TLS_TRUSTSTORE_PASSWORD
      valueFrom:
        secretKeyRef:
          name: {{ include "kafka.zookeeper.tlsPasswordsSecretName" .context }}
          key: {{ .context.Values.tls.zookeeper.passwordsSecretTruststoreKey | quote }}
    {{- end }}
    {{- end }}
    {{- end }}
  volumeMounts:
    - name: data
      mountPath: /bitnami/kafka
    - name: kafka-config
      mountPath: /config
    - name: kafka-configmaps
      mountPath: /configmaps
    - name: kafka-secret-config
      mountPath: /secret-config
    - name: scripts
      mountPath: /scripts
    - name: tmp
      mountPath: /tmp
    {{- if and .context.Values.externalAccess.enabled .context.Values.externalAccess.autoDiscovery.enabled }}
    - name: kafka-autodiscovery-shared
      mountPath: /shared
    {{- end }}
    {{- if or (include "kafka.sslEnabled" .context) .context.Values.tls.zookeeper.enabled }}
    - name: kafka-shared-certs
      mountPath: /certs
    {{- if and (include "kafka.sslEnabled" .context) (or .context.Values.tls.existingSecret .context.Values.tls.autoGenerated) }}
    - name: kafka-certs
      mountPath: /mounted-certs
      readOnly: true
    {{- end }}
    {{- if and .context.Values.tls.zookeeper.enabled .context.Values.tls.zookeeper.existingSecret }}
    - name: kafka-zookeeper-cert
      mountPath: /zookeeper-certs
      readOnly: true
    {{- end }}
    {{- end }}
{{- end -}}

{{/*
Init container definition for waiting for Kubernetes autodiscovery
*/}}
{{- define "kafka.autoDiscoveryInitContainer" -}}
{{- $externalAccessService := index .context.Values.externalAccess .role }}
- name: auto-discovery
  image: {{ include "kafka.externalAccess.autoDiscovery.image" .context }}
  imagePullPolicy: {{ .context.Values.externalAccess.autoDiscovery.image.pullPolicy | quote }}
  command:
    - /scripts/auto-discovery.sh
  env:
    - name: MY_POD_NAME
      valueFrom:
        fieldRef:
          fieldPath: metadata.name
    - name: AUTODISCOVERY_SERVICE_TYPE
      value: {{ $externalAccessService.service.type | quote }}
  {{- if .context.Values.externalAccess.autoDiscovery.resources }}
  resources: {{- toYaml .context.Values.externalAccess.autoDiscovery.resources | nindent 12 }}
  {{- end }}
  volumeMounts:
    - name: scripts
      mountPath: /scripts/auto-discovery.sh
      subPath: auto-discovery.sh
    - name: kafka-autodiscovery-shared
      mountPath: /shared
{{- end -}}

{{/*
Check if there are rolling tags in the images
*/}}
{{- define "kafka.checkRollingTags" -}}
{{- include "common.warnings.rollingTag" .Values.image }}
{{- include "common.warnings.rollingTag" .Values.externalAccess.autoDiscovery.image }}
{{- include "common.warnings.rollingTag" .Values.metrics.kafka.image }}
{{- include "common.warnings.rollingTag" .Values.metrics.jmx.image }}
{{- include "common.warnings.rollingTag" .Values.volumePermissions.image }}
{{- end -}}

{{/*
Compile all warnings into a single message, and call fail.
*/}}
{{- define "kafka.validateValues" -}}
{{- $messages := list -}}
{{- $messages := append $messages (include "kafka.validateValues.listener.protocols" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.controller.nodePortListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.broker.nodePortListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.controller.externalIPListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.broker.externalIPListLength" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.domainSpecified" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceType" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessAutoDiscoveryRBAC" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessAutoDiscoveryIPsOrNames" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceList" (dict "element" "loadBalancerIPs" "context" .)) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceList" (dict "element" "loadBalancerNames" "context" .)) -}}
{{- $messages := append $messages (include "kafka.validateValues.externalAccessServiceList" (dict "element" "loadBalancerAnnotations" "context" . )) -}}
{{- $messages := append $messages (include "kafka.validateValues.saslMechanisms" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.tlsSecret" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.provisioning.tlsPasswords" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.kraftMode" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.kraftMissingControllers" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.zookeeperMissingBrokers" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.zookeeperNoControllers" .) -}}
{{- $messages := append $messages (include "kafka.validateValues.modeEmpty" .) -}}
{{- $messages := without $messages "" -}}
{{- $message := join "\n" $messages -}}

{{- if $message -}}
{{-   printf "\nVALUES VALIDATION:\n%s" $message | fail -}}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - Authentication protocols for Kafka */}}
{{- define "kafka.validateValues.listener.protocols" -}}
{{- $authProtocols := list "PLAINTEXT" "SASL_PLAINTEXT" "SASL_SSL" "SSL" -}}
{{- if not .Values.listeners.securityProtocolMap -}}
{{- $listeners := list .Values.listeners.client .Values.listeners.interbroker -}}
{{- if .Values.kraft.enabled -}}
{{- $listeners = append $listeners .Values.listeners.controller -}}
{{- end -}}
{{- if and .Values.externalAccess.enabled -}}
{{- $listeners = append $listeners .Values.listeners.external -}}
{{- end -}}
{{- $error := false -}}
{{- range $listener := $listeners -}}
{{- if not (has (upper $listener.protocol) $authProtocols) -}}
{{- $error := true -}}
{{- end -}}
{{- end -}}
{{- if $error -}}
kafka: listeners.*.protocol
    Available authentication protocols are "PLAINTEXT" "SASL_PLAINTEXT" "SSL" "SASL_SSL"
{{- end -}}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of controller-eligible replicas must be the same as NodePort list in controller-eligible external service */}}
{{- define "kafka.validateValues.controller.nodePortListLength" -}}
{{- $replicaCount := int .Values.controller.replicaCount -}}
{{- $nodePortListLength := len .Values.externalAccess.controller.service.nodePorts -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.controller.service.nodePorts -}}
{{- $nodePortListLengthEqualsReplicaCount := eq $nodePortListLength $replicaCount -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.controller.service.externalIPs -}}
{{- if and .Values.externalAccess.enabled (not .Values.externalAccess.autoDiscovery.enabled) (eq .Values.externalAccess.controller.service.type "NodePort") (or (and (not $nodePortListIsEmpty) (not $nodePortListLengthEqualsReplicaCount)) (and $nodePortListIsEmpty $externalIPListIsEmpty)) -}}
kafka: .Values.externalAccess.controller.service.nodePorts
    Number of controller-eligible replicas and externalAccess.controller.service.nodePorts array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length nodePorts = {{ $nodePortListLength }} - {{ $externalIPListIsEmpty }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of broker replicas must be the same as NodePort list in broker external service */}}
{{- define "kafka.validateValues.broker.nodePortListLength" -}}
{{- $replicaCount := int .Values.broker.replicaCount -}}
{{- $nodePortListLength := len .Values.externalAccess.broker.service.nodePorts -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.broker.service.nodePorts -}}
{{- $nodePortListLengthEqualsReplicaCount := eq $nodePortListLength $replicaCount -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.broker.service.externalIPs -}}
{{- if and .Values.externalAccess.enabled (not .Values.externalAccess.autoDiscovery.enabled) (eq .Values.externalAccess.broker.service.type "NodePort") (or (and (not $nodePortListIsEmpty) (not $nodePortListLengthEqualsReplicaCount)) (and $nodePortListIsEmpty $externalIPListIsEmpty)) -}}
kafka: .Values.externalAccess.broker.service.nodePorts
    Number of broker replicas and externalAccess.broker.service.nodePorts array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length nodePorts = {{ $nodePortListLength }} - {{ $externalIPListIsEmpty }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of replicas must be the same as externalIPs list */}}
{{- define "kafka.validateValues.controller.externalIPListLength" -}}
{{- $replicaCount := int .Values.controller.replicaCount -}}
{{- $externalIPListLength := len .Values.externalAccess.controller.service.externalIPs -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.controller.service.externalIPs -}}
{{- $externalIPListEqualsReplicaCount := eq $externalIPListLength $replicaCount -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.controller.service.nodePorts -}}
{{- if and .Values.externalAccess.enabled (or .Values.externalAccess.controller.forceExpose (not .Values.controller.controllerOnly)) (not .Values.externalAccess.autoDiscovery.enabled) (eq .Values.externalAccess.controller.service.type "NodePort") (or (and (not $externalIPListIsEmpty) (not $externalIPListEqualsReplicaCount)) (and $externalIPListIsEmpty $nodePortListIsEmpty)) -}}
kafka: .Values.externalAccess.controller.service.externalIPs
    Number of controller-eligible replicas and externalAccess.controller.service.externalIPs array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length externalIPs = {{ $externalIPListLength }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of replicas must be the same as externalIPs list */}}
{{- define "kafka.validateValues.broker.externalIPListLength" -}}
{{- $replicaCount := int .Values.broker.replicaCount -}}
{{- $externalIPListLength := len .Values.externalAccess.broker.service.externalIPs -}}
{{- $externalIPListIsEmpty := empty .Values.externalAccess.broker.service.externalIPs -}}
{{- $externalIPListEqualsReplicaCount := eq $externalIPListLength $replicaCount -}}
{{- $nodePortListIsEmpty := empty .Values.externalAccess.broker.service.nodePorts -}}
{{- if and .Values.externalAccess.enabled (not .Values.externalAccess.autoDiscovery.enabled) (eq .Values.externalAccess.broker.service.type "NodePort") (or (and (not $externalIPListIsEmpty) (not $externalIPListEqualsReplicaCount)) (and $externalIPListIsEmpty $nodePortListIsEmpty)) -}}
kafka: .Values.externalAccess.broker.service.externalIPs
    Number of broker replicas and externalAccess.broker.service.externalIPs array length must be the same. Currently: replicaCount = {{ $replicaCount }} and length externalIPs = {{ $externalIPListLength }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - domain must be defined if external service type ClusterIP */}}
{{- define "kafka.validateValues.domainSpecified" -}}
{{- if and (eq .Values.externalAccess.controller.service.type "ClusterIP") (empty .Values.externalAccess.controller.service.domain) -}}
kafka: .Values.externalAccess.controller.service.domain
    Domain must be specified if service type ClusterIP is set for external service
{{- end -}}
{{- if and (eq .Values.externalAccess.broker.service.type "ClusterIP") (empty .Values.externalAccess.broker.service.domain) -}}
kafka: .Values.externalAccess.broker.service.domain
    Domain must be specified if service type ClusterIP is set for external service
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - service type for external access */}}
{{- define "kafka.validateValues.externalAccessServiceType" -}}
{{- if and (not (eq .Values.externalAccess.controller.service.type "NodePort")) (not (eq .Values.externalAccess.controller.service.type "LoadBalancer")) (not (eq .Values.externalAccess.controller.service.type "ClusterIP")) -}}
kafka: externalAccess.controller.service.type
    Available service type for external access are NodePort, LoadBalancer or ClusterIP.
{{- end -}}
{{- if and (not (eq .Values.externalAccess.broker.service.type "NodePort")) (not (eq .Values.externalAccess.broker.service.type "LoadBalancer")) (not (eq .Values.externalAccess.broker.service.type "ClusterIP")) -}}
kafka: externalAccess.broker.service.type
    Available service type for external access are NodePort, LoadBalancer or ClusterIP.
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - RBAC should be enabled when autoDiscovery is enabled */}}
{{- define "kafka.validateValues.externalAccessAutoDiscoveryRBAC" -}}
{{- if and .Values.externalAccess.enabled .Values.externalAccess.autoDiscovery.enabled (not .Values.rbac.create ) }}
kafka: rbac.create
    By specifying "externalAccess.enabled=true" and "externalAccess.autoDiscovery.enabled=true"
    an initContainer will be used to auto-detect the external IPs/ports by querying the
    K8s API. Please note this initContainer requires specific RBAC resources. You can create them
    by specifying "--set rbac.create=true".
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - LoadBalancerIPs or LoadBalancerNames should be set when autoDiscovery is disabled */}}
{{- define "kafka.validateValues.externalAccessAutoDiscoveryIPsOrNames" -}}
{{- $loadBalancerNameListLength := len .Values.externalAccess.controller.service.loadBalancerNames -}}
{{- $loadBalancerIPListLength := len .Values.externalAccess.controller.service.loadBalancerIPs -}}
{{- if and .Values.externalAccess.enabled (or .Values.externalAccess.controller.forceExpose (not .Values.controller.controllerOnly)) (eq .Values.externalAccess.controller.service.type "LoadBalancer") (not .Values.externalAccess.autoDiscovery.enabled) (eq $loadBalancerNameListLength 0) (eq $loadBalancerIPListLength 0) }}
kafka: externalAccess.controller.service.loadBalancerNames or externalAccess.controller.service.loadBalancerIPs
    By specifying "externalAccess.enabled=true", "externalAccess.autoDiscovery.enabled=false" and
    "externalAccess.controller.service.type=LoadBalancer" at least one of externalAccess.controller.service.loadBalancerNames
    or externalAccess.controller.service.loadBalancerIPs  must be set and the length of those arrays must be equal
    to the number of replicas.
{{- end -}}
{{- $loadBalancerNameListLength := len .Values.externalAccess.broker.service.loadBalancerNames -}}
{{- $loadBalancerIPListLength := len .Values.externalAccess.broker.service.loadBalancerIPs -}}
{{- $replicaCount := int .Values.broker.replicaCount }}
{{- if and .Values.externalAccess.enabled (gt 0 $replicaCount) (eq .Values.externalAccess.broker.service.type "LoadBalancer") (not .Values.externalAccess.autoDiscovery.enabled) (eq $loadBalancerNameListLength 0) (eq $loadBalancerIPListLength 0) }}
kafka: externalAccess.broker.service.loadBalancerNames or externalAccess.broker.service.loadBalancerIPs
    By specifying "externalAccess.enabled=true", "externalAccess.autoDiscovery.enabled=false" and
    "externalAccess.broker.service.type=LoadBalancer" at least one of externalAccess.broker.service.loadBalancerNames
    or externalAccess.broker.service.loadBalancerIPs  must be set and the length of those arrays must be equal
    to the number of replicas.
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - number of replicas must be the same as loadBalancerIPs list */}}
{{- define "kafka.validateValues.externalAccessServiceList" -}}
{{- $replicaCount := int .context.Values.controller.replicaCount }}
{{- $listLength := len (get .context.Values.externalAccess.controller.service .element) -}}
{{- if and .context.Values.externalAccess.enabled (or .context.Values.externalAccess.controller.forceExpose (not .context.Values.controller.controllerOnly)) (not .context.Values.externalAccess.autoDiscovery.enabled) (eq .context.Values.externalAccess.controller.service.type "LoadBalancer") (gt $listLength 0) (not (eq $replicaCount $listLength)) }}
kafka: externalAccess.service.{{ .element }}
    Number of replicas and {{ .element }} array length must be the same. Currently: replicaCount = {{ $replicaCount }} and {{ .element }} = {{ $listLength }}
{{- end -}}
{{- $replicaCount := int .context.Values.broker.replicaCount }}
{{- $listLength := len (get .context.Values.externalAccess.broker.service .element) -}}
{{- if and .context.Values.externalAccess.enabled (gt 0 $replicaCount) (not .context.Values.externalAccess.autoDiscovery.enabled) (eq .context.Values.externalAccess.broker.service.type "LoadBalancer") (gt $listLength 0) (not (eq $replicaCount $listLength)) }}
kafka: externalAccess.service.{{ .element }}
    Number of replicas and {{ .element }} array length must be the same. Currently: replicaCount = {{ $replicaCount }} and {{ .element }} = {{ $listLength }}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - SASL mechanisms must be provided when using SASL */}}
{{- define "kafka.validateValues.saslMechanisms" -}}
{{- if and (include "kafka.saslEnabled" .) (not .Values.sasl.enabledMechanisms) }}
kafka: sasl.enabledMechanisms
    The SASL mechanisms are required when listeners use SASL security protocol.
{{- end }}
{{- if not (contains .Values.sasl.interBrokerMechanism .Values.sasl.enabledMechanisms) }}
kafka: sasl.enabledMechanisms
    sasl.interBrokerMechanism must be provided and it should be one of the specified mechanisms at sasl.enabledMechanisms
{{- end -}}
{{- if and .Values.kraft.enabled (not (contains .Values.sasl.controllerMechanism .Values.sasl.enabledMechanisms)) }}
kafka: sasl.enabledMechanisms
    sasl.controllerMechanism must be provided and it should be one of the specified mechanisms at sasl.enabledMechanisms
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka - Secrets containing TLS certs must be provided when TLS authentication is enabled */}}
{{- define "kafka.validateValues.tlsSecret" -}}
{{- if and (include "kafka.sslEnabled" .) (eq (upper .Values.tls.type) "JKS") (empty .Values.tls.existingSecret) (not .Values.tls.autoGenerated) }}
kafka: tls.existingSecret
    A secret containing the Kafka JKS keystores and truststore is required
    when TLS encryption in enabled and TLS format is "JKS"
{{- else if and (include "kafka.sslEnabled" .) (eq (upper .Values.tls.type) "PEM") (empty .Values.tls.existingSecret) (not .Values.tls.autoGenerated) }}
kafka: tls.existingSecret
    A secret containing the Kafka TLS certificates and keys is required
    when TLS encryption in enabled and TLS format is "PEM"
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka provisioning - keyPasswordSecretKey, keystorePasswordSecretKey or truststorePasswordSecretKey must not be used without passwordsSecret */}}
{{- define "kafka.validateValues.provisioning.tlsPasswords" -}}
{{- if and (regexFind "SSL" (upper .Values.listeners.client.protocol)) .Values.provisioning.enabled (not .Values.provisioning.auth.tls.passwordsSecret) }}
{{- if or .Values.provisioning.auth.tls.keyPasswordSecretKey .Values.provisioning.auth.tls.keystorePasswordSecretKey .Values.provisioning.auth.tls.truststorePasswordSecretKey }}
kafka: tls.keyPasswordSecretKey,tls.keystorePasswordSecretKey,tls.truststorePasswordSecretKey
    tls.keyPasswordSecretKey,tls.keystorePasswordSecretKey,tls.truststorePasswordSecretKey
    must not be used without passwordsSecret setted.
{{- end -}}
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka Kraft mode. It cannot be used with Zookeeper unless migration is enabled */}}
{{- define "kafka.validateValues.kraftMode" -}}
{{- if and .Values.kraft.enabled (or .Values.zookeeper.enabled .Values.externalZookeeper.servers) (and (not .Values.controller.zookeeperMigrationMode ) (not .Values.broker.zookeeperMigrationMode )) }}
kafka: Simultaneous KRaft and Zookeeper modes
    Both Zookeeper and KRaft modes have been configured simultaneously, but migration mode has not been enabled.
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka Kraft mode. At least 1 controller is configured or controller.quorum.voters is set  */}}
{{- define "kafka.validateValues.kraftMissingControllers" -}}
{{- if and .Values.kraft.enabled (le (int .Values.controller.replicaCount) 0) (not .Values.kraft.controllerQuorumVoters) }}
kafka: Kraft mode - Missing controller-eligible nodes
    Kraft mode has been enabled, but no controller-eligible nodes have been configured
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka Zookeper mode. At least 1 broker is configured  */}}
{{- define "kafka.validateValues.zookeeperMissingBrokers" -}}
{{- if and (or .Values.zookeeper.enabled .Values.externalZookeeper.servers) (le (int .Values.broker.replicaCount) 0)}}
kafka: Zookeeper mode - No Kafka brokers configured
    Zookeper mode has been enabled, but no Kafka brokers nodes have been configured
{{- end -}}
{{- end -}}

{{/* Validate values of Kafka Zookeper mode. Controller nodes not enabled in Zookeeper mode unless migration enabled */}}
{{- define "kafka.validateValues.zookeeperNoControllers" -}}
{{- if and (or .Values.zookeeper.enabled .Values.externalZookeeper.servers) (gt (int .Values.controller.replicaCount) 0) (and (not .Values.controller.zookeeperMigrationMode ) (not .Values.broker.zookeeperMigrationMode )) }}
kafka: Zookeeper mode - Controller nodes not supported
    Controller replicas have been enabled in Zookeeper mode, set controller.replicaCount to zero or enable migration mode to migrate to Kraft mode
{{- end -}}
{{- end -}}

{{/* Validate either KRaft or Zookeeper mode are enabled */}}
{{- define "kafka.validateValues.modeEmpty" -}}
{{- if and (not .Values.kraft.enabled) (not (or .Values.zookeeper.enabled .Values.externalZookeeper.servers)) }}
kafka: Missing KRaft or Zookeeper mode settings
    The Kafka chart has been deployed but neither KRaft or Zookeeper modes have been enabled.
    Please configure 'kraft.enabled', 'zookeeper.enabled' or `externalZookeeper.servers` before proceeding.
{{- end -}}
{{- end -}}
