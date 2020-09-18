The following table enumerates built-in checks:


| Name | Enabled by default | Description | Template | Parameters |
 --- | --- | --- | --- | --- | 
|env-var-secret|No|Alert on objects using a secret in an environment variable|env-var|- name: .*secret.* <br />|
|privileged-container|Yes|Alert on deployments with containers running in privileged mode|privileged|none|
|required-label-owner|No|Alert on objects without the 'owner' label|required-label|- key: owner <br />|
