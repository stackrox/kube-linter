# KubeLinter templates

KubeLinter supports the following templates:

## Access to Resources

**Key**: `access-to-resources`

**Description**: Flag cluster role bindings and role bindings that grant access to the specified resource kinds and verbs

**Supported Objects**: Role,ClusterRole,ClusterRoleBinding,RoleBinding


**Parameters**:

```yaml
- description: Set to true to flag the roles that are referenced in bindings but not
    found in the context
  name: flagRolesNotFound
  required: false
  type: boolean
- arrayElemType: string
  description: An array of regular expressions specifying resources. e.g. ^secrets$
    for secrets and ^*$ for any resources
  name: resources
  negationAllowed: false
  regexAllowed: true
  required: false
  type: array
- arrayElemType: string
  description: An array of regular expressions specifying verbs. e.g. ^create$ for
    create and ^*$ for any k8s verbs
  name: verbs
  negationAllowed: false
  regexAllowed: true
  required: false
  type: array
```

## Anti affinity not specified

**Key**: `anti-affinity`

**Description**: Flag objects with multiple replicas but inter-pod anti affinity not specified in the pod template spec

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: The minimum number of replicas a deployment must have before anti-affinity
    is enforced on it
  name: minReplicas
  required: false
  type: integer
- description: The topology key that the anti-affinity term should use. If not specified,
    it defaults to "kubernetes.io/hostname".
  name: topologyKey
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## cluster-admin Role Binding

**Key**: `cluster-admin-role-binding`

**Description**: Flag bindings of cluster-admin role to service accounts, users, or groups

**Supported Objects**: ClusterRoleBinding


## CPU Requirements

**Key**: `cpu-requirements`

**Description**: Flag containers with CPU requirements in the given range

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: The type of requirement. Use any to apply to both requests and limits.
  name: requirementsType
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: The lower bound of the requirement (inclusive), specified as a number
    of milli-cores. If not specified, it is treated as a lower bound of zero.
  name: lowerBoundMillis
  required: false
  type: integer
- description: The upper bound of the requirement (inclusive), specified as a number
    of milli-cores. If not specified, it is treated as "no upper bound".
  name: upperBoundMillis
  required: false
  type: integer
```

## Dangling HorizontalPodAutoscalers

**Key**: `dangling-horizontalpodautoscaler`

**Description**: Flag HorizontalPodAutoscalers that target a resource that does not exist

**Supported Objects**: HorizontalPodAutoscaler


## Dangling Ingress

**Key**: `dangling-ingress`

**Description**: Flag ingress which do not match any service and port

**Supported Objects**: Ingress


## Dangling NetworkPolicies

**Key**: `dangling-networkpolicy`

**Description**: Flag NetworkPolicies which do not match any application

**Supported Objects**: DeploymentLike


## Dangling NetworkPolicyPeer PodSelector

**Key**: `dangling-networkpolicypeer-podselector`

**Description**: Flag NetworkPolicyPeer in Ingress/Egress rules which their podselector do not match any application. Applied to peers consisting with podSelectors only.

**Supported Objects**: DeploymentLike


## Dangling Services

**Key**: `dangling-service`

**Description**: Flag services which do not match any application

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- arrayElemType: string
  description: A list of labels that will not cause the check to fail. For example,
    a label that is known to be populated at runtime by Kubernetes.
  name: ignoredLabels
  negationAllowed: true
  regexAllowed: true
  required: false
  type: array
```

## Dangling Service Monitor

**Key**: `dangling-servicemonitor`

**Description**: Flag service monitors which do not match any service

**Supported Objects**: ServiceMonitor


## Deprecated Service Account Field

**Key**: `deprecated-service-account-field`

**Description**: Flag uses of the deprecated serviceAccount field, which should be migrated to serviceAccountName

**Supported Objects**: DeploymentLike


## Disallowed API Objects

**Key**: `disallowed-api-obj`

**Description**: Flag disallowed API object kinds

**Supported Objects**: Any


**Parameters**:

```yaml
- description: The disallowed object group.
  examples:
  - apps
  name: group
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
- description: The disallowed object API version.
  examples:
  - v1
  - v1beta1
  name: version
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
- description: The disallowed kind.
  examples:
  - Deployment
  - DaemonSet
  name: kind
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## DnsConfig Options

**Key**: `dnsconfig-options`

**Description**: Flag objects that don't have specified DNSConfig Options

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: Key of the dnsConfig option.
  name: key
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
- description: Value of the dnsConfig option.
  name: value
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Duplicate Environment Variables

**Key**: `duplicate-env-var`

**Description**: Flag Duplicate Env Variables names

**Supported Objects**: DeploymentLike


## Environment Variables

**Key**: `env-var`

**Description**: Flag environment variables that match the provided patterns

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: The name of the environment variable.
  name: name
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: The value of the environment variable.
  name: value
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Forbidden Annotation

**Key**: `forbidden-annotation`

**Description**: Flag objects carrying at least one annotation matching the provided patterns

**Supported Objects**: Any


**Parameters**:

```yaml
- description: Key of the forbidden annotation.
  name: key
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: Value of the forbidden annotation.
  name: value
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Forbidden Service Types

**Key**: `forbidden-service-types`

**Description**: Flag forbidden services

**Supported Objects**: Service


**Parameters**:

```yaml
- arrayElemType: string
  description: An array of service types that should not be used
  name: forbiddenServiceTypes
  negationAllowed: false
  regexAllowed: false
  required: false
  type: array
```

## Host IPC

**Key**: `host-ipc`

**Description**: Flag Pod sharing host's IPC namespace

**Supported Objects**: DeploymentLike


## Host Mounts

**Key**: `host-mounts`

**Description**: Flag volume mounts of sensitive system directories

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- arrayElemType: string
  description: An array of regular expressions specifying system directories to be
    mounted on containers. e.g. ^/usr$ for /usr
  name: dirs
  negationAllowed: false
  regexAllowed: true
  required: false
  type: array
```

## Host Network

**Key**: `host-network`

**Description**: Flag Pod sharing host's network namespace

**Supported Objects**: DeploymentLike


## Host PID

**Key**: `host-pid`

**Description**: Flag Pod sharing host's process namespace

**Supported Objects**: DeploymentLike


## HorizontalPodAutoscaler Minimum replicas

**Key**: `hpa-minimum-replicas`

**Description**: Flag applications running fewer than the specified number of replicas

**Supported Objects**: HorizontalPodAutoscaler


**Parameters**:

```yaml
- description: The minimum number of replicas a HorizontalPodAutoscaler should have
  name: minReplicas
  required: false
  type: integer
```

## Image Pull Policy

**Key**: `image-pull-policy`

**Description**: Flag containers with forbidden image pull policy

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- arrayElemType: string
  description: list of forbidden image pull policy
  name: forbiddenPolicies
  negationAllowed: false
  regexAllowed: false
  required: false
  type: array
```

## Latest Tag

**Key**: `latest-tag`

**Description**: Flag applications running container images that do not satisfies "allowList" & "blockList" parameters criteria.

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- arrayElemType: string
  description: list of regular expressions specifying pattern(s) for container images
    that will be blocked. */
  name: blockList
  negationAllowed: true
  regexAllowed: true
  required: false
  type: array
- arrayElemType: string
  description: list of regular expressions specifying pattern(s) for container images
    that will be allowed.
  name: allowList
  negationAllowed: true
  regexAllowed: true
  required: false
  type: array
```

## Liveness Port Exposed

**Key**: `liveness-port`

**Description**: Flag containers with an liveness probe to not exposed port.

**Supported Objects**: DeploymentLike


## Liveness Probe Not Specified

**Key**: `liveness-probe`

**Description**: Flag containers that don't specify a liveness probe

**Supported Objects**: DeploymentLike


## Memory Requirements

**Key**: `memory-requirements`

**Description**: Flag containers with memory requirements in the given range

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: The type of requirement. Use any to apply to both requests and limits.
  name: requirementsType
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: The lower bound of the requirement (inclusive), specified as a number
    of MB.
  name: lowerBoundMB
  required: false
  type: integer
- description: The upper bound of the requirement (inclusive), specified as a number
    of MB. If not specified, it is treated as "no upper bound".
  name: upperBoundMB
  required: false
  type: integer
```

## Minimum replicas

**Key**: `minimum-replicas`

**Description**: Flag applications running fewer than the specified number of replicas

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: The minimum number of replicas a deployment should have
  name: minReplicas
  required: false
  type: integer
```

## Mismatching Selector

**Key**: `mismatching-selector`

**Description**: Flag deployments where the selector doesn't match the labels in the pod template spec

**Supported Objects**: DeploymentLike


## Node Affinity

**Key**: `no-node-affinity`

**Description**: Flag objects that don't have node affinity rules set

**Supported Objects**: DeploymentLike


## Non-Existent Service Account

**Key**: `non-existent-service-account`

**Description**: Flag cases where a pod references a non-existent service account

**Supported Objects**: DeploymentLike


## Non Isolated Pods

**Key**: `non-isolated-pod`

**Description**: Flag Pod that is not selected by any networkPolicy

**Supported Objects**: NetworkPolicy


## No pod disruptions allowed - maxUnavailable

**Key**: `pdb-max-unavailable`

**Description**: Flag PodDisruptionBudgets whose maxUnavailable value will always prevent pod disruptions.

**Supported Objects**: PodDisruptionBudget


## No pod disruptions allowed - minAvailable

**Key**: `pdb-min-available`

**Description**: Flag PodDisruptionBudgets whose minAvailable value will always prevent pod disruptions.

**Supported Objects**: PodDisruptionBudget


## .spec.unhealthyPodEvictionPolicy in PDB is set to default

**Key**: `pdb-unhealthy-pod-eviction-policy`

**Description**: Flag PodDisruptionBudget objects that do not explicitly set unhealthyPodEvictionPolicy.

**Supported Objects**: PodDisruptionBudget


## Ports

**Key**: `ports`

**Description**: Flag containers exposing ports under protocols that match the supplied parameters

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: The port
  name: port
  required: false
  type: integer
- description: The protocol
  name: protocol
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Privilege Escalation on Containers

**Key**: `privilege-escalation-container`

**Description**: Flag containers of allowing privilege escalation

**Supported Objects**: DeploymentLike


## Privileged Containers

**Key**: `privileged`

**Description**: Flag privileged containers

**Supported Objects**: DeploymentLike


## Privileged Ports

**Key**: `privileged-ports`

**Description**: Flag privileged ports

**Supported Objects**: DeploymentLike


## Read-only Root Filesystems

**Key**: `read-only-root-fs`

**Description**: Flag containers without read-only root file systems

**Supported Objects**: DeploymentLike


## Read Secret From Environment Variables

**Key**: `read-secret-from-env-var`

**Description**: Flag environment variables that use SecretKeyRef

**Supported Objects**: DeploymentLike


## Readiness Port Not Exposed

**Key**: `readiness-port`

**Description**: Flag containers with an Readiness probe to not exposed port.

**Supported Objects**: DeploymentLike


## Readiness Probe Not Specified

**Key**: `readiness-probe`

**Description**: Flag containers that don't specify a readiness probe

**Supported Objects**: DeploymentLike


## Required Annotation

**Key**: `required-annotation`

**Description**: Flag objects not carrying at least one annotation matching the provided patterns

**Supported Objects**: Any


**Parameters**:

```yaml
- description: Key of the required label.
  name: key
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: Value of the required label.
  name: value
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Required Label

**Key**: `required-label`

**Description**: Flag objects not carrying at least one label matching the provided patterns

**Supported Objects**: Any


**Parameters**:

```yaml
- description: Key of the required label.
  name: key
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: Value of the required label.
  name: value
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Run as non-root user

**Key**: `run-as-non-root`

**Description**: Flag containers set to run as a root user

**Supported Objects**: DeploymentLike


## SecurityContextConstraints allowPrivilegedContainer

**Key**: `scc-deny-privileged-container`

**Description**: Flag SCC with allowPrivilegedContainer set to true

**Supported Objects**: SecurityContextConstraints


**Parameters**:

```yaml
- description: allowPrivilegedContainer value
  name: allowPrivilegedContainer
  required: false
  type: boolean
```

## Service Account

**Key**: `service-account`

**Description**: Flag containers which use a matching service account

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: A regex specifying the required service account to match.
  name: serviceAccount
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
```

## Startup Port Exposed

**Key**: `startup-port`

**Description**: Flag containers with an Startup probe to not exposed port.

**Supported Objects**: DeploymentLike


## Target Port

**Key**: `target-port`

**Description**: Flag containers and services using not allowed port names or numbers

**Supported Objects**: DeploymentLike,Service


## Unsafe Proc Mount

**Key**: `unsafe-proc-mount`

**Description**: Flag containers of unsafe proc mount

**Supported Objects**: DeploymentLike


## Unsafe Sysctls

**Key**: `unsafe-sysctls`

**Description**: Flag unsafe sysctls

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- arrayElemType: string
  description: An array of unsafe system controls
  name: unsafeSysCtls
  negationAllowed: false
  regexAllowed: false
  required: false
  type: array
```

## Update configuration

**Key**: `update-configuration`

**Description**: Flag configurations that do not meet the specified update configuration

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- description: A regular expression the defines the type of update strategy allowed.
  name: strategyTypeRegex
  negationAllowed: true
  regexAllowed: true
  required: true
  type: string
- description: The maximum value that be set in a RollingUpdate configuration for
    the MaxUnavailable.  This can be an integer or a percent.
  name: maxPodsUnavailable
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
- description: The minimum value that be set in a RollingUpdate configuration for
    the MaxUnavailable.  This can be an integer or a percent.
  name: minPodsUnavailable
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
- description: The maximum value that be set in a RollingUpdate configuration for
    the MaxSurge.  This can be an integer or a percent.
  name: maxSurge
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
- description: The minimum value that be set in a RollingUpdate configuration for
    the MaxSurge.  This can be an integer or a percent.
  name: minSurge
  negationAllowed: true
  regexAllowed: true
  required: false
  type: string
```

## Use Namespaces for Administrative Boundaries between Resources

**Key**: `use-namespace`

**Description**: Flag resources with no namespace specified or using default namespace

**Supported Objects**: DeploymentLike,Service


## Verify container capabilities

**Key**: `verify-container-capabilities`

**Description**: Flag containers that do not match capabilities requirements

**Supported Objects**: DeploymentLike


**Parameters**:

```yaml
- arrayElemType: string
  description: List of capabilities that needs to be removed from containers.
  name: forbiddenCapabilities
  negationAllowed: false
  regexAllowed: false
  required: false
  type: array
- arrayElemType: string
  description: List of capabilities that are exceptions to the above list. This should
    only be filled when the above contains "all", and is used to forgive capabilities
    in ADD list.
  name: exceptions
  negationAllowed: false
  regexAllowed: false
  required: false
  type: array
```

## Wildcard Use in Role and ClusterRole Rules

**Key**: `wildcard-in-rules`

**Description**: Flag Roles and ClusterRoles that use wildcard * in rules

**Supported Objects**: Role,ClusterRole


## Writable Host Mounts

**Key**: `writable-host-mount`

**Description**: Flag containers that have mounted a directory on the host as writable

**Supported Objects**: DeploymentLike


