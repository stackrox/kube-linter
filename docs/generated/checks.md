# KubeLinter checks

KubeLinter includes the following built-in checks:

## cluster-admin-role-binding

**Enabled by default**: No

**Description**: CIS Benchmark 5.1.1 Ensure that the cluster-admin role is only used where required

**Remediation**: Create and assign a separate role that has access to specific resources/actions needed for the service account.

**Template**: [cluster-admin-role-binding](generated/templates.md#cluster-admin-role-binding)

**Parameters**:

```json
{}
```

## dangling-service

**Enabled by default**: Yes

**Description**: Alert on services that don't have any matching deployments

**Remediation**: Make sure your service's selector correctly matches the labels on one of your deployments.

**Template**: [dangling-service](generated/templates.md#dangling-services)

**Parameters**:

```json
{}
```

## default-service-account

**Enabled by default**: No

**Description**: Alert on pods that use the default service account

**Remediation**: Create a dedicated service account for your pod. See https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/ for more details.

**Template**: [service-account](generated/templates.md#service-account)

**Parameters**:

```json
{"serviceAccount":"^(|default)$"}
```

## deprecated-service-account-field

**Enabled by default**: Yes

**Description**: Alert on deployments that use the deprecated serviceAccount field

**Remediation**: Use the serviceAccountName field instead of the serviceAccount field.

**Template**: [deprecated-service-account-field](generated/templates.md#deprecated-service-account-field)

**Parameters**:

```json
{}
```

## docker-sock

**Enabled by default**: No

**Description**: Alert on deployments with docker.sock mounted on containers

**Remediation**: Ensure the Docker socket is not mounted inside any containers.

**Template**: [docker-sock](generated/templates.md#docker-sock)

**Parameters**:

```json
{}
```

## drop-net-raw-capability

**Enabled by default**: Yes

**Description**: Alert on containers not dropping NET_RAW capability

**Remediation**: NET_RAW grants an application within the container the ability to craft raw packets, use raw sockets, and it also allows an application to bind to any address. Please specify to drop this capability in the containers under containers security contexts.

**Template**: [verify-container-capabilities](generated/templates.md#verify-container-capabilities)

**Parameters**:

```json
{"forbiddenCapabilities":["NET_RAW"]}
```

## env-var-secret

**Enabled by default**: Yes

**Description**: Alert on objects using a secret in an environment variable

**Remediation**: Don't use raw secrets in an environment variable. Instead, either mount the secret as a file or use a secretKeyRef. See https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets for more details.

**Template**: [env-var](generated/templates.md#environment-variables)

**Parameters**:

```json
{"name":"(?i).*secret.*","value":".+"}
```

## forbidden-service-types

**Enabled by default**: No

**Description**: Alert on services for forbidden types

**Remediation**: Ensure containers are not exposed through a forbidden service type such as NodePort or LoadBalancer.

**Template**: [forbidden-service-types](generated/templates.md#forbidden-service-types)

**Parameters**:

```json
{"forbiddenServiceTypes":["NodePort","LoadBalancer"]}
```

## host-ipc

**Enabled by default**: No

**Description**: Alert on pods/deployment-likes with sharing host's IPC namespace

**Remediation**: Ensure the host's IPC namespace is not shared.

**Template**: [host-ipc](generated/templates.md#host-ipc)

**Parameters**:

```json
{}
```

## host-network

**Enabled by default**: No

**Description**: Alert on pods/deployment-likes with sharing host's network namespace

**Remediation**: Ensure the host's network namespace is not shared.

**Template**: [host-network](generated/templates.md#host-network)

**Parameters**:

```json
{}
```

## host-pid

**Enabled by default**: No

**Description**: Alert on pods/deployment-likes with sharing host's process namespace

**Remediation**: Ensure the host's process namespace is not shared.

**Template**: [host-pid](generated/templates.md#host-pid)

**Parameters**:

```json
{}
```

## mismatching-selector

**Enabled by default**: Yes

**Description**: Alert on deployments where the selector doesn't match the pod template labels

**Remediation**: Make sure your deployment's selector correctly matches the labels in its pod template.

**Template**: [mismatching-selector](generated/templates.md#mismatching-selector)

**Parameters**:

```json
{}
```

## no-anti-affinity

**Enabled by default**: Yes

**Description**: Alert on deployments with multiple replicas that don't specify inter pod anti-affinity to ensure that the orchestrator attempts to schedule replicas on different nodes

**Remediation**: Specify anti-affinity in your pod spec to ensure that the orchestrator attempts to schedule replicas on different nodes. You can do this by using podAntiAffinity, specifying a labelSelector that matches pods of this deployment, and setting the topologyKey to kubernetes.io/hostname. See https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#inter-pod-affinity-and-anti-affinity for more details.

**Template**: [anti-affinity](generated/templates.md#anti-affinity-not-specified)

**Parameters**:

```json
{"minReplicas":2}
```

## no-extensions-v1beta

**Enabled by default**: Yes

**Description**: Alert on objects using deprecated API versions under extensions v1beta

**Remediation**: Migrate to using the apps/v1 API versions for these objects. See https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/ for more details.

**Template**: [disallowed-api-obj](generated/templates.md#disallowed-api-objects)

**Parameters**:

```json
{"group":"extensions","version":"v1beta.+"}
```

## no-liveness-probe

**Enabled by default**: No

**Description**: Alert on containers which don't specify a liveness probe

**Remediation**: Specify a liveness probe in your container. See https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ for more details.

**Template**: [liveness-probe](generated/templates.md#liveness-probe-not-specified)

**Parameters**:

```json
{}
```

## no-read-only-root-fs

**Enabled by default**: Yes

**Description**: Alert on containers not running with a read-only root filesystem

**Remediation**: Set readOnlyRootFilesystem to true in your container's securityContext.

**Template**: [read-only-root-fs](generated/templates.md#read-only-root-filesystems)

**Parameters**:

```json
{}
```

## no-readiness-probe

**Enabled by default**: No

**Description**: Alert on containers which don't specify a readiness probe

**Remediation**: Specify a readiness probe in your container. See https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ for more details.

**Template**: [readiness-probe](generated/templates.md#readiness-probe-not-specified)

**Parameters**:

```json
{}
```

## non-existent-service-account

**Enabled by default**: Yes

**Description**: Alert on pods referencing a service account that isn't found

**Remediation**: Make sure to create the service account, or to refer to an existing service account.

**Template**: [non-existent-service-account](generated/templates.md#non-existent-service-account)

**Parameters**:

```json
{}
```

## privilege-escalation-container

**Enabled by default**: No

**Description**: Alert on containers with allowing privilege escalation

**Remediation**: Ensure containers do not allow privilege escalation.

**Template**: [privilege-escalation-container](generated/templates.md#privilege-escalation-on-containers)

**Parameters**:

```json
{}
```

## privileged-container

**Enabled by default**: Yes

**Description**: Alert on deployments with containers running in privileged mode

**Remediation**: Don't run your container as privileged unless required.

**Template**: [privileged](generated/templates.md#privileged-containers)

**Parameters**:

```json
{}
```

## privileged-ports

**Enabled by default**: No

**Description**: Alert on deployments with privileged ports mapped in containers

**Remediation**: Ensure privileged ports [0, 1024] are not mapped within containers.

**Template**: [privileged-ports](generated/templates.md#privileged-ports)

**Parameters**:

```json
{}
```

## required-annotation-email

**Enabled by default**: No

**Description**: Alert on objects without an 'email' annotation with a valid email

**Remediation**: Add an email annotation to your object with the contact information of the object's owner.

**Template**: [required-annotation](generated/templates.md#required-annotation)

**Parameters**:

```json
{"key":"email","value":"[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+"}
```

## required-label-owner

**Enabled by default**: No

**Description**: Alert on objects without the 'owner' label

**Remediation**: Add an email annotation to your object with information about the object's owner.

**Template**: [required-label](generated/templates.md#required-label)

**Parameters**:

```json
{"key":"owner"}
```

## run-as-non-root

**Enabled by default**: Yes

**Description**: Alert on containers not set to runAsNonRoot

**Remediation**: Set runAsUser to a non-zero number, and runAsNonRoot to true, in your pod or container securityContext. See https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ for more details.

**Template**: [run-as-non-root](generated/templates.md#run-as-non-root-user)

**Parameters**:

```json
{}
```

## ssh-port

**Enabled by default**: Yes

**Description**: Alert on deployments exposing port 22, commonly reserved for SSH access

**Remediation**: Ensure that non-SSH services are not using port 22. Ensure that any actual SSH servers have been vetted.

**Template**: [ports](generated/templates.md#ports)

**Parameters**:

```json
{"port":22,"protocol":"TCP"}
```

## unsafe-proc-mount

**Enabled by default**: No

**Description**: Alert on deployments with unsafe /proc mount on containers

**Remediation**: Ensure container does not exposes unsafe parts of /proc.

**Template**: [unsafe-proc-mount](generated/templates.md#unsafe-proc-mount)

**Parameters**:

```json
{}
```

## unsafe-sysctls

**Enabled by default**: No

**Description**: Alert on deployments with unsafe allocation of sysctls

**Remediation**: Ensure containers does not allow unsafe allocation of CPU resources.

**Template**: [unsafe-sysctls](generated/templates.md#unsafe-sysctls)

**Parameters**:

```json
{"unsafeSysCtls":["kernel.msg","kernel.sem","kernel.shm","fs.mqueue.","net."]}
```

## unset-cpu-requirements

**Enabled by default**: Yes

**Description**: Alert on containers without CPU requests and limits set

**Remediation**: Set your container's CPU requests and limits depending on its requirements. See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for more details.

**Template**: [cpu-requirements](generated/templates.md#cpu-requirements)

**Parameters**:

```json
{"lowerBoundMillis":0,"requirementsType":"any","upperBoundMillis":0}
```

## unset-memory-requirements

**Enabled by default**: Yes

**Description**: Alert on containers without memory requests and limits set

**Remediation**: Set your container's memory requests and limits depending on its requirements. See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for more details.

**Template**: [memory-requirements](generated/templates.md#memory-requirements)

**Parameters**:

```json
{"lowerBoundMB":0,"requirementsType":"any","upperBoundMB":0}
```

## volume-mounts

**Enabled by default**: No

**Description**: Alert on deployments with sensitive host system directories mounted on containers

**Remediation**: Ensure sensitive host system directories are not mounted on containers.

**Template**: [volume-mounts](generated/templates.md#volume-mounts)

**Parameters**:

```json
{"sensitiveSysDirs":["/","/boot","/dev","/etc","/lib","/proc","/sys","/usr"]}
```

## writable-host-mount

**Enabled by default**: No

**Description**: Alert on containers that mount a host path as writable

**Remediation**: If you need to access files on the host, mount them as readOnly.

**Template**: [writable-host-mount](generated/templates.md#writable-host-mounts)

**Parameters**:

```json
{}
```

