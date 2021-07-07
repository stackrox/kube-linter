# KubeLinter templates

KubeLinter supports the following templates:

## Access to Resources

**Key**: `access-to-resources`

**Description**: Flag cluster role bindings and role bindings that grant access to the specified resource kinds and verbs

**Supported Objects**: Role,ClusterRole,ClusterRoleBinding,RoleBinding

**Parameters**:

```json
[
  {
    "name": "flagRolesNotFound",
    "type": "boolean",
    "description": "Set to true to flag the roles that are referenced in bindings but not found in the context",
    "required": false
  },
  {
    "name": "resources",
    "type": "array",
    "description": "An array of regular expressions specifying resources. e.g. ^secrets$ for secrets and ^*$ for any resources",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": false,
    "arrayElemType": "string"
  },
  {
    "name": "verbs",
    "type": "array",
    "description": "An array of regular expressions specifying verbs. e.g. ^create$ for create and ^*$ for any k8s verbs",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": false,
    "arrayElemType": "string"
  }
]
```

## Anti affinity not specified

**Key**: `anti-affinity`

**Description**: Flag objects with multiple replicas but inter-pod anti affinity not specified in the pod template spec

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "minReplicas",
    "type": "integer",
    "description": "The minimum number of replicas a deployment must have before anti-affinity is enforced on it",
    "required": false
  },
  {
    "name": "topologyKey",
    "type": "string",
    "description": "The topology key that the anti-affinity term should use. If not specified, it defaults to \"kubernetes.io/hostname\".",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## cluster-admin Role Binding

**Key**: `cluster-admin-role-binding`

**Description**: Flag bindings of cluster-admin role to service accounts, users, or groups

**Supported Objects**: ClusterRoleBinding

**Parameters**:

```json
[]
```

## CPU Requirements

**Key**: `cpu-requirements`

**Description**: Flag containers with CPU requirements in the given range

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "requirementsType",
    "type": "string",
    "description": "The type of requirement. Use any to apply to both requests and limits.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "lowerBoundMillis",
    "type": "integer",
    "description": "The lower bound of the requirement (inclusive), specified as a number of milli-cores. If not specified, it is treated as a lower bound of zero.",
    "required": false
  },
  {
    "name": "upperBoundMillis",
    "type": "integer",
    "description": "The upper bound of the requirement (inclusive), specified as a number of milli-cores. If not specified, it is treated as \"no upper bound\".",
    "required": false
  }
]
```

## Dangling Services

**Key**: `dangling-service`

**Description**: Flag services which do not match any application

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Deprecated Service Account Field

**Key**: `deprecated-service-account-field`

**Description**: Flag uses of the deprecated serviceAccount field, which should be migrated to serviceAccountName

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Disallowed API Objects

**Key**: `disallowed-api-obj`

**Description**: Flag disallowed API object kinds

**Supported Objects**: Any

**Parameters**:

```json
[
  {
    "name": "group",
    "type": "string",
    "description": "The disallowed object group.",
    "required": false,
    "examples": [
      "apps"
    ],
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "version",
    "type": "string",
    "description": "The disallowed object API version.",
    "required": false,
    "examples": [
      "v1",
      "v1beta1"
    ],
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "kind",
    "type": "string",
    "description": "The disallowed kind.",
    "required": false,
    "examples": [
      "Deployment",
      "DaemonSet"
    ],
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Environment Variables

**Key**: `env-var`

**Description**: Flag environment variables that match the provided patterns

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "name",
    "type": "string",
    "description": "The name of the environment variable.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "value",
    "type": "string",
    "description": "The value of the environment variable.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Forbidden Service Types

**Key**: `forbidden-service-types`

**Description**: Flag forbidden services

**Supported Objects**: Service

**Parameters**:

```json
[
  {
    "name": "forbiddenServiceTypes",
    "type": "array",
    "description": "An array of service types that should not be used",
    "required": false,
    "regexAllowed": false,
    "negationAllowed": false,
    "arrayElemType": "string"
  }
]
```

## Host IPC

**Key**: `host-ipc`

**Description**: Flag Pod sharing host's IPC namespace

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Host Mounts

**Key**: `host-mounts`

**Description**: Flag volume mounts of sensitive system directories

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "dirs",
    "type": "array",
    "description": "An array of regular expressions specifying system directories to be mounted on containers. e.g. ^/usr$ for /usr",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": false,
    "arrayElemType": "string"
  }
]
```

## Host Network

**Key**: `host-network`

**Description**: Flag Pod sharing host's network namespace

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Host PID

**Key**: `host-pid`

**Description**: Flag Pod sharing host's process namespace

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Latest Tag

**Key**: `latest-tag`

**Description**: Flag applications running containers with floating container image tag, "latest"

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "blockList",
    "type": "array",
    "description": "list of regular expressions for blocked or bad container image tags",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true,
    "arrayElemType": "string"
  }
]
```

## Liveness Probe Not Specified

**Key**: `liveness-probe`

**Description**: Flag containers that don't specify a liveness probe

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Memory Requirements

**Key**: `memory-requirements`

**Description**: Flag containers with memory requirements in the given range

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "requirementsType",
    "type": "string",
    "description": "The type of requirement. Use any to apply to both requests and limits.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "lowerBoundMB",
    "type": "integer",
    "description": "The lower bound of the requirement (inclusive), specified as a number of MB.",
    "required": false
  },
  {
    "name": "upperBoundMB",
    "type": "integer",
    "description": "The upper bound of the requirement (inclusive), specified as a number of MB. If not specified, it is treated as \"no upper bound\".",
    "required": false
  }
]
```

## Minimum replicas

**Key**: `minimum-replicas`

**Description**: Flag applications running fewer than the specified number of replicas

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "minReplicas",
    "type": "integer",
    "description": "The minimum number of replicas a deployment should have",
    "required": false
  }
]
```

## Mismatching Selector

**Key**: `mismatching-selector`

**Description**: Flag deployments where the selector doesn't match the labels in the pod template spec

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Non-Existent Service Account

**Key**: `non-existent-service-account`

**Description**: Flag cases where a pod references a non-existent service account

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Ports

**Key**: `ports`

**Description**: Flag containers exposing ports under protocols that match the supplied parameters

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "port",
    "type": "integer",
    "description": "The port",
    "required": false
  },
  {
    "name": "protocol",
    "type": "string",
    "description": "The protocol",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Privilege Escalation on Containers

**Key**: `privilege-escalation-container`

**Description**: Flag containers of allowing privilege escalation

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Privileged Containers

**Key**: `privileged`

**Description**: Flag privileged containers

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Privileged Ports

**Key**: `privileged-ports`

**Description**: Flag privileged ports

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Read-only Root Filesystems

**Key**: `read-only-root-fs`

**Description**: Flag containers without read-only root file systems

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Read Secret From Environment Variables

**Key**: `read-secret-from-env-var`

**Description**: Flag environment variables that use SecretKeyRef

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Readiness Probe Not Specified

**Key**: `readiness-probe`

**Description**: Flag containers that don't specify a readiness probe

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Required Annotation

**Key**: `required-annotation`

**Description**: Flag objects not carrying at least one annotation matching the provided patterns

**Supported Objects**: Any

**Parameters**:

```json
[
  {
    "name": "key",
    "type": "string",
    "description": "Key of the required label.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "value",
    "type": "string",
    "description": "Value of the required label.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Required Label

**Key**: `required-label`

**Description**: Flag objects not carrying at least one label matching the provided patterns

**Supported Objects**: Any

**Parameters**:

```json
[
  {
    "name": "key",
    "type": "string",
    "description": "Key of the required label.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "value",
    "type": "string",
    "description": "Value of the required label.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Run as non-root user

**Key**: `run-as-non-root`

**Description**: Flag containers set to run as a root user

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Service Account

**Key**: `service-account`

**Description**: Flag containers which use a matching service account

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "serviceAccount",
    "type": "string",
    "description": "A regex specifying the required service account to match.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Unsafe Proc Mount

**Key**: `unsafe-proc-mount`

**Description**: Flag containers of unsafe proc mount

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

## Unsafe Sysctls

**Key**: `unsafe-sysctls`

**Description**: Flag unsafe sysctls

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "unsafeSysCtls",
    "type": "array",
    "description": "An array of unsafe system controls",
    "required": false,
    "regexAllowed": false,
    "negationAllowed": false,
    "arrayElemType": "string"
  }
]
```

## Update configuration

**Key**: `update-configuration`

**Description**: Flag configurations that do not meet the specified update configuration

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "strategyTypeRegex",
    "type": "string",
    "description": "A regular expression the defines the type of update strategy allowed.",
    "required": true,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "maxPodsUnavailable",
    "type": "string",
    "description": "The maximum value that be set in a RollingUpdate configuration for the MaxUnavailable.  This can be an integer or a percent.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "minPodsUnavailable",
    "type": "string",
    "description": "The minimum value that be set in a RollingUpdate configuration for the MaxUnavailable.  This can be an integer or a percent.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "maxSurge",
    "type": "string",
    "description": "The maximum value that be set in a RollingUpdate configuration for the MaxSurge.  This can be an integer or a percent.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  },
  {
    "name": "minSurge",
    "type": "string",
    "description": "The minimum value that be set in a RollingUpdate configuration for the MaxSurge.  This can be an integer or a percent.",
    "required": false,
    "regexAllowed": true,
    "negationAllowed": true
  }
]
```

## Use Namespaces for Administrative Boundaries between Resources

**Key**: `use-namespace`

**Description**: Flag resources with no namespace specified or using default namespace

**Supported Objects**: DeploymentLike,Service

**Parameters**:

```json
[]
```

## Verify container capabilities

**Key**: `verify-container-capabilities`

**Description**: Flag containers that do not match capabilities requirements

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[
  {
    "name": "forbiddenCapabilities",
    "type": "array",
    "description": "List of capabilities that needs to be removed from containers.",
    "required": false,
    "regexAllowed": false,
    "negationAllowed": false,
    "arrayElemType": "string"
  },
  {
    "name": "exceptions",
    "type": "array",
    "description": "List of capabilities that are exceptions to the above list. This should only be filled when the above contains \"all\", and is used to forgive capabilities in ADD list.",
    "required": false,
    "regexAllowed": false,
    "negationAllowed": false,
    "arrayElemType": "string"
  }
]
```

## Wildcard Use in Role and ClusterRole Rules

**Key**: `wildcard-in-rules`

**Description**: Flag Roles and ClusterRoles that use wildcard * in rules

**Supported Objects**: Role,ClusterRole

**Parameters**:

```json
[]
```

## Writable Host Mounts

**Key**: `writable-host-mount`

**Description**: Flag containers that have mounted a directory on the host as writable

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

