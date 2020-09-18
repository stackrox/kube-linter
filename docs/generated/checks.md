The following table enumerates built-in checks:

| Name | Enabled by default | Description | Template | Parameters |
| ---- | ------------------ | ----------- | -------- | ---------- |
 | env-var-secret | Yes | Alert on objects using a secret in an environment variable | env-var |- `name`: `.*secret.*` <br />|
 | no-read-only-root-fs | Yes | Alert on containers not running with a read-only root filesystem | read-only-root-fs | none |
 | privileged-container | Yes | Alert on deployments with containers running in privileged mode | privileged | none |
 | required-label-owner | No | Alert on objects without the 'owner' label | required-label |- `key`: `owner` <br />|
 | run-as-non-root | Yes | Alert on containers not set to runAsNonRoot | run-as-non-root | none |
