This page lists supported check templates.

## Disallowed API Objects

**Key**: `disallowed-api-obj`

**Description**: Flag disallowed API object kinds

**Supported Objects**: Any

Parameters:
``` 
[
	{
		"description": "The disallowed object group.",
		"examples": [
			"apps"
		],
		"name": "group",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": false,
		"type": "string"
	},
	{
		"description": "The disallowed object API version.",
		"examples": [
			"v1",
			"v1beta1"
		],
		"name": "version",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": false,
		"type": "string"
	},
	{
		"description": "The disallowed kind.",
		"examples": [
			"Deployment",
			"DaemonSet"
		],
		"name": "kind",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": false,
		"type": "string"
	}
]

``` 

## Environment Variables

**Key**: `env-var`

**Description**: Flag environment variables that match the provided patterns

**Supported Objects**: DeploymentLike

Parameters:
``` 
[
	{
		"description": "The name of the environment variable.",
		"name": "name",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": true,
		"type": "string"
	},
	{
		"description": "The value of the environment variable.",
		"name": "value",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": false,
		"type": "string"
	}
]

``` 

## Privileged Containers

**Key**: `privileged`

**Description**: Flag privileged containers

**Supported Objects**: DeploymentLike

Parameters:
``` 
[]

``` 

## Read-only Root Filesystems

**Key**: `read-only-root-fs`

**Description**: Flag containers without read-only root file systems

**Supported Objects**: DeploymentLike

Parameters:
``` 
[]

``` 

## Required Label

**Key**: `required-label`

**Description**: Flag objects not carrying at least one label matching the provided patterns

**Supported Objects**: Any

Parameters:
``` 
[
	{
		"description": "Key of the required label.",
		"name": "key",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": true,
		"type": "string"
	},
	{
		"description": "Value of the required label.",
		"name": "value",
		"negationAllowed": true,
		"regexAllowed": true,
		"required": false,
		"type": "string"
	}
]

``` 

## Run as non-root user

**Key**: `run-as-non-root`

**Description**: Flag containers set to run as a root user

**Supported Objects**: DeploymentLike

Parameters:
``` 
[]

``` 

