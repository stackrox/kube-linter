This page lists supported check templates.

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

