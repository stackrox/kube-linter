# KubeLinter checks

KubeLinter includes the following built-in checks:

## access-to-create-pods

**Enabled by default**: No

**Description**: Indicates when a subject (Group/User/ServiceAccount) has create access to Pods. CIS Benchmark 5.1.4: The ability to create pods in a cluster opens up possibilities for privilege escalation and should be restricted, where possible.

**Remediation**: Where possible, remove create access to pod objects in the cluster.

**Template**: [access-to-resources](generated/templates.md#access-to-resources)

**Parameters**:

```json
{"resources":["^pods$","^deployments$","^statefulsets$","^replicasets$","^cronjob$","^jobs$","^daemonsets$"],"verbs":["^create$"]}
```

## access-to-secrets

**Enabled by default**: No

**Description**: Indicates when a subject (Group/User/ServiceAccount) has access to Secrets. CIS Benchmark 5.1.2: Access to secrets should be restricted to the smallest possible group of users to reduce the risk of privilege escalation.

**Remediation**: Where possible, remove get, list and watch access to secret objects in the cluster.

**Template**: [access-to-resources](generated/templates.md#access-to-resources)

**Parameters**:

```json
{"resources":["^secrets$"],"verbs":["^get$","^list$","^delete$","^create$","^watch$","^*$"]}
```

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

## docker-sock

**Enabled by default**: Yes

**Description**: Alert on deployments with docker.sock mounted in containers. 

**Remediation**: Ensure the Docker socket is not mounted inside any containers by removing the associated  Volume and VolumeMount in deployment yaml specification. If the Docker socket is mounted inside a container it could allow processes running within  the container to execute Docker commands which would effectively allow for full control of the host.

**Template**: [host-mounts](generated/templates.md#host-mounts)

**Parameters**:

```json
{"dirs":["docker.sock$"]}
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

## exposed-services

**Enabled by default**: No

**Description**: Alert on services for forbidden types

**Remediation**: Ensure containers are not exposed through a forbidden service type such as NodePort or LoadBalancer.

**Template**: [forbidden-service-types](generated/templates.md#forbidden-service-types)

**Parameters**:

```json
{"forbiddenServiceTypes":["NodePort","LoadBalancer"]}
```

## host-ipc

**Enabled by default**: Yes

**Description**: Alert on pods/deployment-likes with sharing host's IPC namespace

**Remediation**: Ensure the host's IPC namespace is not shared.

**Template**: [host-ipc](generated/templates.md#host-ipc)

**Parameters**:

```json
{}
```

## host-network

**Enabled by default**: Yes

**Description**: Alert on pods/deployment-likes with sharing host's network namespace

**Remediation**: Ensure the host's network namespace is not shared.

**Template**: [host-network](generated/templates.md#host-network)

**Parameters**:

```json
{}
```

## host-pid

**Enabled by default**: Yes

**Description**: Alert on pods/deployment-likes with sharing host's process namespace

**Remediation**: Ensure the host's process namespace is not shared.

**Template**: [host-pid](generated/templates.md#host-pid)

**Parameters**:

```json
{}
```

## latest-tag

**Enabled by default**: Yes

**Description**: Indicates when a deployment-like object is running a container with a floating image tag, "latest"

**Remediation**: Use a container image with a proper image tag, outside the set blocked tag regex ".*:(latest)$".

**Template**: [latest-tag](generated/templates.md#latest-tag)

**Parameters**:

```json
{"BlockList":[".*:(latest)$"]}
```

## minimum-three-replicas

**Enabled by default**: No

**Description**: Indicates when a deployment uses less than three replicas

**Remediation**: Increase be number of replicas in the deployment to at least three to increase the fault tolerancy of the deployment.

**Template**: [minimum-replicas](generated/templates.md#minimum-replicas)

**Parameters**:

```json
{"minReplicas":3}
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

## no-rolling-update-strategy

**Enabled by default**: No

**Description**: Indicates when a deployment doesn't use a rolling update strategy

**Remediation**: Use a rolling update strategy to avoid service disruption during an update. A rolling update strategy allows for pods to be systematicaly replaced in a controlled fashion to ensure no service disruption.

**Template**: [update-configuration](generated/templates.md#update-configuration)

**Parameters**:

```json
{"strategyTypeRegex":"^(RollingUpdate|Rolling)$"}
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

## privilege-escalation-container

**Enabled by default**: Yes

**Description**: Alert on containers of allowing privilege escalation that could gain more privileges than its parent process.

**Remediation**: Ensure containers do not allow privilege escalation by setting allowPrivilegeEscalation=false." See https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ for more details.

**Template**: [privilege-escalation-container](generated/templates.md#privilege-escalation-on-containers)

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

## privileged-ports

**Enabled by default**: No

**Description**: Alert on deployments with privileged ports mapped in containers

**Remediation**: Ensure privileged ports [0, 1024] are not mapped within containers.

**Template**: [privileged-ports](generated/templates.md#privileged-ports)

**Parameters**:

```json
{}
```

## read-secret-from-env-var

**Enabled by default**: No

**Description**: Indicates when a deployment reads secret from environment variables. CIS Benchmark 5.4.1: "Prefer using secrets as files over secrets as environment variables. "

**Remediation**: If possible, rewrite application code to read secrets from mounted secret files, rather than from environment variables. Refer to https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets for details.

**Template**: [read-secret-from-env-var](generated/templates.md#read-secret-from-environment-variables)

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

## sensitive-host-mounts

**Enabled by default**: Yes

**Description**: Alert on deployments with sensitive host system directories mounted in containers

**Remediation**: Ensure sensitive host system directories are not mounted in containers by removing those Volumes and VolumeMounts.

**Template**: [host-mounts](generated/templates.md#host-mounts)

**Parameters**:

```json
{"dirs":["^/$","^/boot$","^/dev$","^/etc$","^/lib$","^/proc$","^/sys$","^/usr$"]}
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

## unsafe-proc-mount

**Enabled by default**: No

**Description**: Alert on deployments with unsafe /proc mount (procMount=Unmasked) that will bypass the default masking behavior of the container runtime

**Remediation**: Ensure container does not unsafely exposes parts of /proc by setting procMount=Default.  Unmasked ProcMount bypasses the default masking behavior of the container runtime. See https://kubernetes.io/docs/concepts/security/pod-security-standards/ for more details.

**Template**: [unsafe-proc-mount](generated/templates.md#unsafe-proc-mount)

**Parameters**:

```json
{}
```

## unsafe-sysctls

**Enabled by default**: Yes

**Description**: Alert on deployments specifying unsafe sysctls that may lead to severe problems like wrong behavior of containers

**Remediation**: Ensure container does not allow unsafe allocation of system resources by removing unsafe sysctls configurations. For more details see https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster/ https://docs.docker.com/engine/reference/commandline/run/#configure-namespaced-kernel-parameters-sysctls-at-runtime.

**Template**: [unsafe-sysctls](generated/templates.md#unsafe-sysctls)

**Parameters**:

```json
{"unsafeSysCtls":["kernel.msg","kernel.sem","kernel.shm","fs.mqueue.","net."]}
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

## use-namespace

**Enabled by default**: No

**Description**: Indicates when a resource is deployed to the default namespace.   CIS Benchmark 5.7.1: Create administrative boundaries between resources using namespaces. CIS Benchmark 5.7.4: The default namespace should not be used.

**Remediation**: Create namespaces for objects in your deployment.

**Template**: [use-namespace](generated/templates.md#use-namespaces-for-administrative-boundaries-between-resources)

**Parameters**:

```json
{}
```

## wildcard-in-rules

**Enabled by default**: No

**Description**: Indicate when a wildcard is used in Role or ClusterRole rules. CIS Benchmark 5.1.3 Use of wildcards is not optimal from a security perspective as it may allow for inadvertent access to be granted when new resources are added to the Kubernetes API either as CRDs or in later versions of the product.

**Remediation**: Where possible replace any use of wildcards in clusterroles and roles with specific objects or actions.

**Template**: [wildcard-in-rules](generated/templates.md#wildcard-use-in-role-and-clusterrole-rules)

**Parameters**:

```json
{}
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

