The following table enumerates supported check templates:

| Name | Description | Supported Objects | Parameters |
| ---- | ----------- | ----------------- | ---------- |
 | env-var | Flag environment variables that match the provided patterns | DeploymentLike |- `name` (required): A regex for the env var name <br />- `value`: A regex for the env var value <br />|
 | privileged | Flag privileged containers | DeploymentLike | none |
 | read-only-root-fs | Flag containers without read-only root file systems | DeploymentLike | none |
 | required-label | Flag objects not carrying at least one label matching the provided patterns | Any |- `key` (required): A regex for the key of the required label <br />- `value`: A regex for the value of the required label <br />|
 | run-as-non-root | Flag containers without runAsUser specified | DeploymentLike | none |
