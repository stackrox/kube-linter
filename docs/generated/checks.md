# KubeLinter checks

KubeLinter includes the following built-in checks:

## access-to-create-pods

**Enabled by default**: No

**Description**: Indicates when a subject (Group/User/ServiceAccount) has create access to Pods. CIS Benchmark 5.1.4: The ability to create pods in a cluster opens up possibilities for privilege escalation and should be restricted, where possible.

**Remediation**: Where possible, remove create access to pod objects in the cluster.

**Template**: [access-to-resources](templates.md#access-to-resources)

**Parameters**:

```yaml
resources:
- ^pods$
- ^deployments$
- ^statefulsets$
- ^replicasets$
- ^cronjob$
- ^jobs$
- ^daemonsets$
verbs:
- ^create$
```
## access-to-secrets

**Enabled by default**: No

**Description**: Indicates when a subject (Group/User/ServiceAccount) has access to Secrets. CIS Benchmark 5.1.2: Access to secrets should be restricted to the smallest possible group of users to reduce the risk of privilege escalation.

**Remediation**: Where possible, remove get, list and watch access to secret objects in the cluster.

**Template**: [access-to-resources](templates.md#access-to-resources)

**Parameters**:

```yaml
resources:
- ^secrets$
verbs:
- ^get$
- ^list$
- ^delete$
- ^create$
- ^watch$
- ^*$
```
## cluster-admin-role-binding

**Enabled by default**: No

**Description**: CIS Benchmark 5.1.1 Ensure that the cluster-admin role is only used where required

**Remediation**: Create and assign a separate role that has access to specific resources/actions needed for the service account.

**Template**: [cluster-admin-role-binding](templates.md#cluster-admin-role-binding)
## dangling-horizontalpodautoscaler

**Enabled by default**: No

**Description**: Indicates when HorizontalPodAutoscalers target a missing resource.

**Remediation**: Confirm that your HorizontalPodAutoscaler's scaleTargetRef correctly matches one of your deployments.

**Template**: [dangling-horizontalpodautoscaler](templates.md#dangling-horizontalpodautoscalers)
## dangling-ingress

**Enabled by default**: No

**Description**: Indicates when ingress do not have any associated services.

**Remediation**: Confirm that your ingress's backend correctly matches the name and port on one of your services.

**Template**: [dangling-ingress](templates.md#dangling-ingress)
## dangling-networkpolicy

**Enabled by default**: No

**Description**: Indicates when networkpolicies do not have any associated deployments.

**Remediation**: Confirm that your networkPolicy's podselector correctly matches the labels on one of your deployments.

**Template**: [dangling-networkpolicy](templates.md#dangling-networkpolicies)
## dangling-networkpolicypeer-podselector

**Enabled by default**: No

**Description**: Indicates when NetworkPolicyPeer in Egress/Ingress rules -in the Spec of NetworkPolicy- do not have any associated deployments. Applied on peer specified with podSelectors only.

**Remediation**: Confirm that your NetworkPolicy's Ingress/Egress peer's podselector correctly matches the labels on one of your deployments.

**Template**: [dangling-networkpolicypeer-podselector](templates.md#dangling-networkpolicypeer-podselector)
## dangling-service

**Enabled by default**: Yes

**Description**: Indicates when services do not have any associated deployments.

**Remediation**: Confirm that your service's selector correctly matches the labels on one of your deployments.

**Template**: [dangling-service](templates.md#dangling-services)
## dangling-servicemonitor

**Enabled by default**: No

**Description**: Indicates when a service monitor's selectors don't match any service. ServiceMonitors are a custom resource only used by the Prometheus operator (https://prometheus-operator.dev/docs/operator/design/#servicemonitor).

**Remediation**: Check selectors and your services.

**Template**: [dangling-servicemonitor](templates.md#dangling-service-monitor)
## default-service-account

**Enabled by default**: No

**Description**: Indicates when pods use the default service account.

**Remediation**: Create a dedicated service account for your pod. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/ for details.

**Template**: [service-account](templates.md#service-account)

**Parameters**:

```yaml
serviceAccount: ^(|default)$
```
## deprecated-service-account-field

**Enabled by default**: Yes

**Description**: Indicates when deployments use the deprecated serviceAccount field.

**Remediation**: Use the serviceAccountName field instead. If you must specify serviceAccount, ensure values for serviceAccount and serviceAccountName match.

**Template**: [deprecated-service-account-field](templates.md#deprecated-service-account-field)
## dnsconfig-options

**Enabled by default**: No

**Description**: Alert on deployments that have no specified dnsConfig options

**Remediation**: Specify dnsconfig options in your Pod specification to ensure the expected DNS setting on the Pod. Refer to https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-dns-config for details.

**Template**: [dnsconfig-options](templates.md#dnsconfig-options)

**Parameters**:

```yaml
Key: ndots
Value: "2"
```
## docker-sock

**Enabled by default**: Yes

**Description**: Alert on deployments with docker.sock mounted in containers. 

**Remediation**: Ensure the Docker socket is not mounted inside any containers by removing the associated  Volume and VolumeMount in deployment yaml specification. If the Docker socket is mounted inside a container it could allow processes running within  the container to execute Docker commands which would effectively allow for full control of the host.

**Template**: [host-mounts](templates.md#host-mounts)

**Parameters**:

```yaml
dirs:
- docker.sock$
```
## drop-net-raw-capability

**Enabled by default**: Yes

**Description**: Indicates when containers do not drop NET_RAW capability

**Remediation**: NET_RAW makes it so that an application within the container is able to craft raw packets, use raw sockets, and bind to any address. Remove this capability in the containers under containers security contexts.

**Template**: [verify-container-capabilities](templates.md#verify-container-capabilities)

**Parameters**:

```yaml
forbiddenCapabilities:
- NET_RAW
```
## duplicate-env-var

**Enabled by default**: Yes

**Description**: Check that duplicate named env vars aren't passed to a deployment like.

**Remediation**: Confirm that your DeploymentLike doesn't have duplicate env vars names.

**Template**: [duplicate-env-var](templates.md#duplicate-environment-variables)
## env-var-secret

**Enabled by default**: Yes

**Description**: Indicates when objects use a secret in an environment variable.

**Remediation**: Do not use raw secrets in environment variables. Instead, either mount the secret as a file or use a secretKeyRef. Refer to https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets for details.

**Template**: [env-var](templates.md#environment-variables)

**Parameters**:

```yaml
name: (?i).*secret.*
value: .+
```
## exposed-services

**Enabled by default**: No

**Description**: Alert on services for forbidden types

**Remediation**: Ensure containers are not exposed through a forbidden service type such as NodePort or LoadBalancer.

**Template**: [forbidden-service-types](templates.md#forbidden-service-types)

**Parameters**:

```yaml
forbiddenServiceTypes:
- NodePort
- LoadBalancer
```
## host-ipc

**Enabled by default**: Yes

**Description**: Alert on pods/deployment-likes with sharing host's IPC namespace

**Remediation**: Ensure the host's IPC namespace is not shared.

**Template**: [host-ipc](templates.md#host-ipc)
## host-network

**Enabled by default**: Yes

**Description**: Alert on pods/deployment-likes with sharing host's network namespace

**Remediation**: Ensure the host's network namespace is not shared.

**Template**: [host-network](templates.md#host-network)
## host-pid

**Enabled by default**: Yes

**Description**: Alert on pods/deployment-likes with sharing host's process namespace

**Remediation**: Ensure the host's process namespace is not shared.

**Template**: [host-pid](templates.md#host-pid)
## hpa-minimum-three-replicas

**Enabled by default**: No

**Description**: Indicates when a HorizontalPodAutoscaler specifies less than three minReplicas

**Remediation**: Increase the number of replicas in the HorizontalPodAutoscaler to at least three to increase fault tolerance.

**Template**: [hpa-minimum-replicas](templates.md#horizontalpodautoscaler-minimum-replicas)

**Parameters**:

```yaml
minReplicas: 3
```
## invalid-target-ports

**Enabled by default**: Yes

**Description**: Indicates when deployments or services are using port names that are violating specifications.

**Remediation**: Ensure that port naming is in conjunction with the specification. For more information, please look at the Kubernetes Service specification on this page: https://kubernetes.io/docs/reference/_print/#ServiceSpec. And additional information about IANA Service naming can be found on the following page: https://www.rfc-editor.org/rfc/rfc6335.html#section-5.1.

**Template**: [target-port](templates.md#target-port)
## latest-tag

**Enabled by default**: Yes

**Description**: Indicates when a deployment-like object is running a container with an invalid container image

**Remediation**: Use a container image with a specific tag other than latest.

**Template**: [latest-tag](templates.md#latest-tag)

**Parameters**:

```yaml
BlockList:
- .*:(latest)$
- ^[^:]*$
- (.*/[^:]+)$
```
## liveness-port

**Enabled by default**: Yes

**Description**: Indicates when containers have a liveness probe to a not exposed port.

**Remediation**: Check which ports you've exposed and ensure they match what you have specified in the liveness probe.

**Template**: [liveness-port](templates.md#liveness-port-exposed)
## minimum-three-replicas

**Enabled by default**: No

**Description**: Indicates when a deployment uses less than three replicas

**Remediation**: Increase the number of replicas in the deployment to at least three to increase the fault tolerance of the deployment.

**Template**: [minimum-replicas](templates.md#minimum-replicas)

**Parameters**:

```yaml
minReplicas: 3
```
## mismatching-selector

**Enabled by default**: Yes

**Description**: Indicates when deployment selectors fail to match the pod template labels.

**Remediation**: Confirm that your deployment selector correctly matches the labels in its pod template.

**Template**: [mismatching-selector](templates.md#mismatching-selector)
## no-anti-affinity

**Enabled by default**: Yes

**Description**: Indicates when deployments with multiple replicas fail to specify inter-pod anti-affinity, to ensure that the orchestrator attempts to schedule replicas on different nodes.

**Remediation**: Specify anti-affinity in your pod specification to ensure that the orchestrator attempts to schedule replicas on different nodes. Using podAntiAffinity, specify a labelSelector that matches pods for the deployment, and set the topologyKey to kubernetes.io/hostname. Refer to https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#inter-pod-affinity-and-anti-affinity for details.

**Template**: [anti-affinity](templates.md#anti-affinity-not-specified)

**Parameters**:

```yaml
minReplicas: 2
```
## no-extensions-v1beta

**Enabled by default**: Yes

**Description**: Indicates when objects use deprecated API versions under extensions/v1beta.

**Remediation**: Migrate using the apps/v1 API versions for the objects. Refer to https://kubernetes.io/blog/2019/07/18/api-deprecations-in-1-16/ for details.

**Template**: [disallowed-api-obj](templates.md#disallowed-api-objects)

**Parameters**:

```yaml
group: extensions
version: v1beta.+
```
## no-liveness-probe

**Enabled by default**: No

**Description**: Indicates when containers fail to specify a liveness probe.

**Remediation**: Specify a liveness probe in your container. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ for details.

**Template**: [liveness-probe](templates.md#liveness-probe-not-specified)
## no-node-affinity

**Enabled by default**: No

**Description**: Alert on deployments that have no node affinity defined

**Remediation**: Specify node-affinity in your pod specification to ensure that the orchestrator attempts to schedule replicas on specified nodes. Refer to https://kubernetes.io/docs/concepts/scheduling-eviction/assign-pod-node/#node-affinity for details.

**Template**: [no-node-affinity](templates.md#node-affinity)
## no-read-only-root-fs

**Enabled by default**: Yes

**Description**: Indicates when containers are running without a read-only root filesystem.

**Remediation**: Set readOnlyRootFilesystem to true in the container securityContext.

**Template**: [read-only-root-fs](templates.md#read-only-root-filesystems)
## no-readiness-probe

**Enabled by default**: No

**Description**: Indicates when containers fail to specify a readiness probe.

**Remediation**: Specify a readiness probe in your container. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/ for details.

**Template**: [readiness-probe](templates.md#readiness-probe-not-specified)
## no-rolling-update-strategy

**Enabled by default**: No

**Description**: Indicates when a deployment doesn't use a rolling update strategy

**Remediation**: Use a rolling update strategy to avoid service disruption during an update. A rolling update strategy allows for pods to be systematicaly replaced in a controlled fashion to ensure no service disruption.

**Template**: [update-configuration](templates.md#update-configuration)

**Parameters**:

```yaml
strategyTypeRegex: ^(RollingUpdate|Rolling)$
```
## non-existent-service-account

**Enabled by default**: Yes

**Description**: Indicates when pods reference a service account that is not found.

**Remediation**: Create the missing service account, or refer to an existing service account.

**Template**: [non-existent-service-account](templates.md#non-existent-service-account)
## non-isolated-pod

**Enabled by default**: No

**Description**: Alert on deployment-like objects that are not selected by any NetworkPolicy.

**Remediation**: Ensure pod does not accept unsafe traffic by isolating it with a NetworkPolicy. See https://cloud.redhat.com/blog/guide-to-kubernetes-ingress-network-policies for more details.

**Template**: [non-isolated-pod](templates.md#non-isolated-pods)
## pdb-max-unavailable

**Enabled by default**: Yes

**Description**: Indicates when a PodDisruptionBudget has a maxUnavailable value that will always prevent disruptions of pods created by related deployment-like objects.

**Remediation**: Change the PodDisruptionBudget to have maxUnavailable set to a value greater than 0. Refer to https://kubernetes.io/docs/tasks/run-application/configure-pdb/ for more information.

**Template**: [pdb-max-unavailable](templates.md#no-pod-disruptions-allowed---maxunavailable)
## pdb-min-available

**Enabled by default**: Yes

**Description**: Indicates when a PodDisruptionBudget sets a minAvailable value that will always prevent disruptions of pods created by related deployment-like objects.

**Remediation**: Change the PodDisruptionBudget to have minAvailable set to a number lower than the number of replicas in the related deployment-like objects. Refer to https://kubernetes.io/docs/tasks/run-application/configure-pdb/ for more information.

**Template**: [pdb-min-available](templates.md#no-pod-disruptions-allowed---minavailable)
## pdb-unhealthy-pod-eviction-policy

**Enabled by default**: Yes

**Description**: Indicates when a PodDisruptionBudget does not explicitly set the unhealthyPodEvictionPolicy field.

**Remediation**: Set unhealthyPodEvictionPolicy to AlwaysAllow. Refer to https://kubernetes.io/docs/tasks/run-application/configure-pdb/#unhealthy-pod-eviction-policy for more information.

**Template**: [pdb-unhealthy-pod-eviction-policy](templates.md#.spec.unhealthypodevictionpolicy-in-pdb-is-set-to-default)
## privilege-escalation-container

**Enabled by default**: Yes

**Description**: Alert on containers of allowing privilege escalation that could gain more privileges than its parent process.

**Remediation**: Ensure containers do not allow privilege escalation by setting allowPrivilegeEscalation=false, privileged=false and removing CAP_SYS_ADMIN capability. See https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ for more details.

**Template**: [privilege-escalation-container](templates.md#privilege-escalation-on-containers)
## privileged-container

**Enabled by default**: Yes

**Description**: Indicates when deployments have containers running in privileged mode.

**Remediation**: Do not run your container as privileged unless it is required.

**Template**: [privileged](templates.md#privileged-containers)
## privileged-ports

**Enabled by default**: No

**Description**: Alert on deployments with privileged ports mapped in containers

**Remediation**: Ensure privileged ports [0, 1024] are not mapped within containers.

**Template**: [privileged-ports](templates.md#privileged-ports)
## read-secret-from-env-var

**Enabled by default**: No

**Description**: Indicates when a deployment reads secret from environment variables. CIS Benchmark 5.4.1: "Prefer using secrets as files over secrets as environment variables. "

**Remediation**: If possible, rewrite application code to read secrets from mounted secret files, rather than from environment variables. Refer to https://kubernetes.io/docs/concepts/configuration/secret/#using-secrets for details.

**Template**: [read-secret-from-env-var](templates.md#read-secret-from-environment-variables)
## readiness-port

**Enabled by default**: Yes

**Description**: Indicates when containers have a readiness probe to a not exposed port.

**Remediation**: Check which ports you've exposed and ensure they match what you have specified in the readiness probe.

**Template**: [readiness-port](templates.md#readiness-port-not-exposed)
## required-annotation-email

**Enabled by default**: No

**Description**: Indicates when objects do not have an email annotation with a valid email address.

**Remediation**: Add an email annotation to your object with the email address of the object's owner.

**Template**: [required-annotation](templates.md#required-annotation)

**Parameters**:

```yaml
key: email
value: '[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+'
```
## required-label-owner

**Enabled by default**: No

**Description**: Indicates when objects do not have an email annotation with an owner label.

**Remediation**: Add an email annotation to your object with the name of the object's owner.

**Template**: [required-label](templates.md#required-label)

**Parameters**:

```yaml
key: owner
```
## run-as-non-root

**Enabled by default**: Yes

**Description**: Indicates when containers are not set to runAsNonRoot.

**Remediation**: Set runAsUser to a non-zero number and runAsNonRoot to true in your pod or container securityContext. Refer to https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ for details.

**Template**: [run-as-non-root](templates.md#run-as-non-root-user)
## scc-deny-privileged-container

**Enabled by default**: No

**Description**: Indicates when allowPrivilegedContainer SecurityContextConstraints set to true

**Remediation**: SecurityContextConstraints has AllowPrivilegedContainer set to "true". Using this option is dangerous, please consider using allowedCapabilities instead. Refer to https://docs.openshift.com/container-platform/4.12/authentication/managing-security-context-constraints.html#scc-settings_configuring-internal-oauth for details.

**Template**: [scc-deny-privileged-container](templates.md#securitycontextconstraints-allowprivilegedcontainer)

**Parameters**:

```yaml
AllowPrivilegedContainer: true
```
## sensitive-host-mounts

**Enabled by default**: Yes

**Description**: Alert on deployments with sensitive host system directories mounted in containers

**Remediation**: Ensure sensitive host system directories are not mounted in containers by removing those Volumes and VolumeMounts.

**Template**: [host-mounts](templates.md#host-mounts)

**Parameters**:

```yaml
dirs:
- ^/$
- ^/boot$
- ^/dev$
- ^/etc$
- ^/lib$
- ^/proc$
- ^/sys$
- ^/usr$
```
## ssh-port

**Enabled by default**: Yes

**Description**: Indicates when deployments expose port 22, which is commonly reserved for SSH access.

**Remediation**: Ensure that non-SSH services are not using port 22. Confirm that any actual SSH servers have been vetted.

**Template**: [ports](templates.md#ports)

**Parameters**:

```yaml
port: 22
protocol: TCP
```
## startup-port

**Enabled by default**: Yes

**Description**: Indicates when containers have a liveness probe to a not exposed port.

**Remediation**: Check which ports you've exposed and ensure they match what you have specified in the liveness probe.

**Template**: [startup-port](templates.md#startup-port-exposed)
## unsafe-proc-mount

**Enabled by default**: No

**Description**: Alert on deployments with unsafe /proc mount (procMount=Unmasked) that will bypass the default masking behavior of the container runtime

**Remediation**: Ensure container does not unsafely exposes parts of /proc by setting procMount=Default.  Unmasked ProcMount bypasses the default masking behavior of the container runtime. See https://kubernetes.io/docs/concepts/security/pod-security-standards/ for more details.

**Template**: [unsafe-proc-mount](templates.md#unsafe-proc-mount)
## unsafe-sysctls

**Enabled by default**: Yes

**Description**: Alert on deployments specifying unsafe sysctls that may lead to severe problems like wrong behavior of containers

**Remediation**: Ensure container does not allow unsafe allocation of system resources by removing unsafe sysctls configurations. For more details see https://kubernetes.io/docs/tasks/administer-cluster/sysctl-cluster/ https://docs.docker.com/engine/reference/commandline/run/#configure-namespaced-kernel-parameters-sysctls-at-runtime.

**Template**: [unsafe-sysctls](templates.md#unsafe-sysctls)

**Parameters**:

```yaml
unsafeSysCtls:
- kernel.msg
- kernel.sem
- kernel.shm
- fs.mqueue.
- net.
```
## unset-cpu-requirements

**Enabled by default**: Yes

**Description**: Indicates when containers do not have CPU requests and limits set.

**Remediation**: Set CPU requests for your container based on its requirements. Refer to https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for details.

**Template**: [cpu-requirements](templates.md#cpu-requirements)

**Parameters**:

```yaml
lowerBoundMillis: 0
requirementsType: request
upperBoundMillis: 0
```
## unset-memory-requirements

**Enabled by default**: Yes

**Description**: Indicates when containers do not have memory requests and limits set.

**Remediation**: Set memory limits for your container based on its requirements. Refer to https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for details.

**Template**: [memory-requirements](templates.md#memory-requirements)

**Parameters**:

```yaml
lowerBoundMB: 0
requirementsType: limit
upperBoundMB: 0
```
## use-namespace

**Enabled by default**: No

**Description**: Indicates when a resource is deployed to the default namespace.   CIS Benchmark 5.7.1: Create administrative boundaries between resources using namespaces. CIS Benchmark 5.7.4: The default namespace should not be used.

**Remediation**: Create namespaces for objects in your deployment.

**Template**: [use-namespace](templates.md#use-namespaces-for-administrative-boundaries-between-resources)
## wildcard-in-rules

**Enabled by default**: No

**Description**: Indicate when a wildcard is used in Role or ClusterRole rules. CIS Benchmark 5.1.3 Use of wildcards is not optimal from a security perspective as it may allow for inadvertent access to be granted when new resources are added to the Kubernetes API either as CRDs or in later versions of the product.

**Remediation**: Where possible replace any use of wildcards in clusterroles and roles with specific objects or actions.

**Template**: [wildcard-in-rules](templates.md#wildcard-use-in-role-and-clusterrole-rules)
## writable-host-mount

**Enabled by default**: No

**Description**: Indicates when containers mount a host path as writable.

**Remediation**: Set containers to mount host paths as readOnly, if you need to access files on the host.

**Template**: [writable-host-mount](templates.md#writable-host-mounts)
