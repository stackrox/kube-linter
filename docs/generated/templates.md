# KubeLinter templates

KubeLinter supports the following templates:

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

## Privileged Containers

**Key**: `privileged`

**Description**: Flag privileged containers

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

## Writable Host Mounts

**Key**: `writable-host-mount`

**Description**: Flag containers that have mounted a directory on the host as writable

**Supported Objects**: DeploymentLike

**Parameters**:

```json
[]
```

