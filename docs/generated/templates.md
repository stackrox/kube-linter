This page lists supported check templates.

## CPU Requirements

**Key**: `cpu-requirements`

**Description**: Flag containers with CPU requirements in the given range

**Supported Objects**: DeploymentLike

**Parameters**:
```
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
```
[]

```

## Deprecated Service Account Field

**Key**: `deprecated-service-account-field`

**Description**: Flag uses of the deprecated serviceAccount field, which should be migrated to serviceAccountName

**Supported Objects**: DeploymentLike

**Parameters**:
```
[]

```

## Disallowed API Objects

**Key**: `disallowed-api-obj`

**Description**: Flag disallowed API object kinds

**Supported Objects**: Any

**Parameters**:
```
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
```
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
```
[]

```

## Memory Requirements

**Key**: `memory-requirements`

**Description**: Flag containers with memory requirements in the given range

**Supported Objects**: DeploymentLike

**Parameters**:
```
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

## Non-Existent Service Account

**Key**: `non-existent-service-account`

**Description**: Flag cases where a pod references a non-existent service account

**Supported Objects**: DeploymentLike

**Parameters**:
```
[]

```

## Privileged Containers

**Key**: `privileged`

**Description**: Flag privileged containers

**Supported Objects**: DeploymentLike

**Parameters**:
```
[]

```

## Read-only Root Filesystems

**Key**: `read-only-root-fs`

**Description**: Flag containers without read-only root file systems

**Supported Objects**: DeploymentLike

**Parameters**:
```
[]

```

## Readiness Probe Not Specified

**Key**: `readiness-probe`

**Description**: Flag containers that don't specify a readiness probe

**Supported Objects**: DeploymentLike

**Parameters**:
```
[]

```

## Required Annotation

**Key**: `required-annotation`

**Description**: Flag objects not carrying at least one annotation matching the provided patterns

**Supported Objects**: Any

**Parameters**:
```
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
```
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
```
[]

```

## Service Account

**Key**: `service-account`

**Description**: Flag containers which use a matching service account

**Supported Objects**: DeploymentLike

**Parameters**:
```
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

## Writable Host Mounts

**Key**: `writable-host-mount`

**Description**: Flag containers that have mounted a directory on the host as writable

**Supported Objects**: DeploymentLike

**Parameters**:
```
[]

```

