The following table enumerates built-in checks:

| Name | Enabled by default | Description | Template | Parameters |
| ---- | ------------------ | ----------- | -------- | ---------- |
 | dangling-service | Yes | Alert on services that don't have any matching deployments | dangling-service | `{}` |
 | env-var-secret | Yes | Alert on objects using a secret in an environment variable | env-var | `{"name":".*secret.*"}` |
 | no-extensions-v1beta | Yes | Alert on objects using deprecated API versions under extensions v1beta | disallowed-api-obj | `{"group":"extensions","version":"v1beta.+"}` |
 | no-read-only-root-fs | Yes | Alert on containers not running with a read-only root filesystem | read-only-root-fs | `{}` |
 | privileged-container | Yes | Alert on deployments with containers running in privileged mode | privileged | `{}` |
 | required-label-owner | No | Alert on objects without the 'owner' label | required-label | `{"key":"owner"}` |
 | run-as-non-root | Yes | Alert on containers not set to runAsNonRoot | run-as-non-root | `{}` |
 | unset-cpu-requirements | Yes | Alert on containers without CPU requests and limits set | cpu-requirements | `{"lowerBoundMillis":0,"requirementsType":"any","upperBoundMillis":0}` |
 | unset-memory-requirements | Yes | Alert on containers without memory requests and limits set | memory-requirements | `{"lowerBoundMB":0,"requirementsType":"any","upperBoundMB":0}` |
