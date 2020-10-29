# KubeLinter

KubeLinter is a static analysis tool that checks Kubernetes YAML files and Helm charts to ensure the applications represented in them adhere to best practices.

KubeLinter accepts YAML files as input and runs a series of checks on them.
If it finds any issues, it reports them and returns a non-zero exit
code.

KubeLinter is:
- **Configurable**: includes multiple built-in checks which you can enable or
  disable.
- **Extensible**: you can define and configure custom checks.

## Installation

If you have [Go](https://golang.org/) installed, run the following command:

```bash
GO111MODULE=on go get -u golang.stackrox.io/kube-linter/cmd/kube-linter
```

Otherwise, download the latest binary from
[Releases](https://github.com/stackrox/kube-linter/releases) and add it to your
PATH.

### Building from source

> **NOTE**: Before you build, make sure that you have [installed
> Go](https://golang.org/doc/install).

To build KubeLinter from source, follow these instructions:

1. Clone the KubeLinter repository.
   ```bash
   git clone git@github.com:stackrox/kube-linter.git
   ```
1. Run the `make build` command. This command compiles the source code and
   creates `kube-linter` binary files for each platform in the `.gobin` folder.

## Usage

1. Consider the following sample pod specification file `pod.yaml`:
   ```yaml
   apiVersion: v1
   kind: Pod
   metadata:
     name: security-context-demo
   spec:
     securityContext:
       runAsUser: 1000
       runAsGroup: 3000
       fsGroup: 2000
     volumes:
     - name: sec-ctx-vol
       emptyDir: {}
     containers:
     - name: sec-ctx-demo
       image: busybox
       resources:
         requests:
           memory: "64Mi"
           cpu: "250m"
       command: [ "sh", "-c", "sleep 1h" ]
       volumeMounts:
       - name: sec-ctx-vol
         mountPath: /data/demo
       securityContext:
         allowPrivilegeEscalation: false
   ```
1. To lint this file with KubeLinter, run the following command:
   ```bash
   kube-linter lint pod.yaml
   ```
1. KubeLinter runs the default checks and reports errors.
   ```
   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container "sec-ctx-demo" does not have a read-only root file system (check: no-read-only-root-fs, remediation: Set readOnlyRootFilesystem to true in your container's securityContext.)

   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container    "sec-ctx-demo" has cpu limit 0 (check: unset-cpu-requirements, remediation: Set    your container's CPU requests and limits depending on its requirements. See    https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/   #requests-and-limits for more details.)
   
   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container    "sec-ctx-demo" has memory limit 0 (check: unset-memory-requirements, remediation:    Set your container's memory requests and limits depending on its requirements.    See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/   #requests-and-limits for more details.)
   
   Error: found 3 lint errors
   ```

For more details about using and configuring KubeLinter, see the
[documentation](./docs) page.

# WARNING: Alpha release

KubeLinter is at an early stage of development. There may be breaking changes in
the future to the command usage, flags, and configuration file formats. However,
we encourage you to use KubeLinter to test your environment YAML files, see what
breaks, and [contribute](./CONTRIBUTING.md).

# LICENSE 

KubeLinter is licensed under the [Apache License 2.0](./LICENSE).
