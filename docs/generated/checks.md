# KubeLinter checks

KubeLinter includes the following built-in checks:

## dangling-service

**Enabled by default**: Yes

**Description**: Indicates when services do not have any associated deployments.

**Remediation**: Confirm that your service's selector correctly matches the labels on one of your deployments.

**Template**: [dangling-service](generated/templates.md#dangling-services)

**Parameters**:

```json
{}
```

## default-service-account

**Enabled by default**: No

**Description**: Indicates when pods use the default service account.

**Remediation**: Create a dedicated service account for your pod. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/ for details.

**Template**: [service-account](generated/templates.md#service-account)

**Parameters**:

```json
{"serviceAccount":"^(|default)$"}
```

## deprecated-service-account-field

**Enabled by default**: Yes

**Description**: Indicates when deployments use the deprecated serviceAccount field.

**Remediation**: Use the serviceAccountName field instead.

**Template**: [deprecated-service-account-field](generated/templates.md#deprecated-service-account-field)

**Parameters**:

```json
{}
```

## drop-net-raw-capability

**Enabled by default**: Yes

**Description**: Indicates when containers do not drop NET_RAW capability

**Remediation**: NET_RAW makes it so that an application within the container is able to craft raw packets, use raw sockets, and bind to any address. Remove this capability in the containers under containers security contexts.

**Template**: [verify-container-capabilities](generated/templates.md#verify-container-capabilities)

**Parameters**:

```json
{"forbiddenCapabilities":["NET_RAW"]}
```

## env-var-secret

**Enabled by default**: Yes

**Description**: Indicates when objects use a secret in an environment variable.

**Remediation**: Do not use raw secrets in environment variables. Instead, either mount the secret as a file or use a secretKeyRef. Refer to https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets for details.

**Template**: [env-var](generated/templates.md#environment-variables)

**Parameters**:

```json
{"name":"(?i).*secret.*","value":".+"}
```

## mismatching-selector

**Enabled by default**: Yes

**Description**: Indicates when deployment selectors fail to match the pod template labels.

**Remediation**: Confirm that your deployment selector correctly matches the labels in its pod template.

**Template**: [mismatching-selector](generated/templates.md#mismatching-selector)

**Parameters**:

```json
{}
```

## no-anti-affinity

**Enabled by default**: Yes

**Description**: Indicates when deployments with multiple replicas fail to specify inter-pod anti-affinity, to ensure that the orchestrator attempts to schedule replicas on different nodes.

**Remediation**: Specify anti-affinity in your pod specification to ensure that the orchestrator attempts to schedule replicas on different nodes. Using podAntiAffinity, specify a labelSelector that matches pods for the deployment, and set the topologyKey to kubernetes.io/hostname. Refer to https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#inter-pod-affinity-and-anti-affinity for details.

**Template**: [anti-affinity](generated/templates.md#anti-affinity-not-specified)

**Parameters**:

```json
{"minReplicas":2}
```

## no-extensions-v1beta

**Enabled by default**: Yes

**Description**: Indicates when objects use deprecated API versions under extensions/v1beta.

**Remediation**: Migrate using the apps/v1 API versions for the objects. Refer to https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/ for details.

**Template**: [disallowed-api-obj](generated/templates.md#disallowed-api-objects)

**Parameters**:

```json
{"group":"extensions","version":"v1beta.+"}
```

## no-liveness-probe

**Enabled by default**: No

**Description**: Indicates when containers fail to specify a liveness probe.

**Remediation**: Specify a liveness probe in your container. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ for details.

**Template**: [liveness-probe](generated/templates.md#liveness-probe-not-specified)

**Parameters**:

```json
{}
```

## no-read-only-root-fs

**Enabled by default**: Yes

**Description**: Indicates when containers are running without a read-only root filesystem.

**Remediation**: Set readOnlyRootFilesystem to true in the container securityContext.

**Template**: [read-only-root-fs](generated/templates.md#read-only-root-filesystems)

**Parameters**:

```json
{}
```

## no-readiness-probe

**Enabled by default**: No

**Description**: Indicates when containers fail to specify a readiness probe.

**Remediation**: Specify a readiness probe in your container. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ for details.

**Template**: [readiness-probe](generated/templates.md#readiness-probe-not-specified)

**Parameters**:

```json
{}
```

## non-existent-service-account

**Enabled by default**: Yes

**Description**: Indicates when pods reference a service account that is not found.

**Remediation**: Create the missing service account, or refer to an existing service account.

**Template**: [non-existent-service-account](generated/templates.md#non-existent-service-account)

**Parameters**:

```json
{}
```

## privileged-container

**Enabled by default**: Yes

**Description**: Indicates when deployments have containers running in privileged mode.

**Remediation**: Do not run your container as privileged unless it is required.

**Template**: [privileged](generated/templates.md#privileged-containers)

**Parameters**:

```json
{}
```

## required-annotation-email

**Enabled by default**: No

**Description**: Indicates when objects do not have an email annotation with a valid email address.

**Remediation**: Add an email annotation to your object with the email address of the object's owner.

**Template**: [required-annotation](generated/templates.md#required-annotation)

**Parameters**:

```json
{"key":"email","value":"[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+"}
```

## required-label-owner

**Enabled by default**: No

**Description**: Indicates when objects do not have an email annotation with an owner label.

**Remediation**: Add an email annotation to your object with the name of the object's owner.

**Template**: [required-label](generated/templates.md#required-label)

**Parameters**:

```json
{"key":"owner"}
```

## run-as-non-root

**Enabled by default**: Yes

**Description**: Indicates when containers are not set to runAsNonRoot.

**Remediation**: Set runAsUser to a non-zero number and runAsNonRoot to true in your pod or container securityContext. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ for details.

**Template**: [run-as-non-root](generated/templates.md#run-as-non-root-user)

**Parameters**:

```json
{}
```

## ssh-port

**Enabled by default**: Yes

**Description**: Indicates when deployments expose port 22, which is commonly reserved for SSH access.

**Remediation**: Ensure that non-SSH services are not using port 22. Confirm that any actual SSH servers have been vetted.

**Template**: [ports](generated/templates.md#ports)

**Parameters**:

```json
{"port":22,"protocol":"TCP"}
```

## unset-cpu-requirements

**Enabled by default**: Yes

**Description**: Indicates when containers do not have CPU requests and limits set.

**Remediation**: Set CPU requests and limits for your container based on its requirements. Refer to https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for details.

**Template**: [cpu-requirements](generated/templates.md#cpu-requirements)

**Parameters**:

```json
{"lowerBoundMillis":0,"requirementsType":"any","upperBoundMillis":0}
```

## unset-memory-requirements

**Enabled by default**: Yes

**Description**: Indicates when containers do not have memory requests and limits set.

**Remediation**: Set memory requests and limits for your container based on its requirements. Refer to https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for details.

**Template**: [memory-requirements](generated/templates.md#memory-requirements)

**Parameters**:

```json
{"lowerBoundMB":0,"requirementsType":"any","upperBoundMB":0}
```

## writable-host-mount

**Enabled by default**: No

**Description**: Indicates when containers mount a host path as writable.

**Remediation**: Set containers to mount host paths as readOnly, if you need to access files on the host.

**Template**: [writable-host-mount](generated/templates.md#writable-host-mounts)

**Parameters**:

```json
{}
```

