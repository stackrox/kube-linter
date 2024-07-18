# KubeLinter

KubeLinter analyzes Kubernetes YAML files and Helm charts and checks them
against various best practices, with a focus on production readiness and
security.

KubeLinter runs sensible default checks designed to give you useful information
about your Kubernetes YAML files and Helm charts. Use it to check early and
often for security misconfigurations and DevOps best practices. Some common
issues that KubeLinter identifies are running containers as a non-root user,
enforcing least privilege, and storing sensitive information only in secrets.

KubeLinter is configurable, so you can enable and disable checks and create your
custom checks, depending on the policies you want to follow within your
organization. When a lint check fails, KubeLinter also reports recommendations
for resolving any potential issues and returns a non-zero exit code.

> [!WARNING]
> KubeLinter is at an early stage of development. There may be breaking changes
> in the future to the command usage, flags, and configuration file formats.
> However, we encourage you to use KubeLinter to test your environment YAML
> files, see what breaks, and
> [contribute](https://github.com/stackrox/kube-linter/blob/main/CONTRIBUTING.md)
> to its development.

## Installing KubeLinter

### Using Go

To install by using [Go](https://golang.org/), run the following command:

```bash
go install golang.stackrox.io/kube-linter/cmd/kube-linter@latest
```
Otherwise, download the latest binary from
[Releases](https://github.com/stackrox/kube-linter/releases) and add it to your
PATH.

### Using Homebrew

To install by using [Homebrew](https://brew.sh/) on macOS, Linux, and [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/about),
run the following command:

```bash
brew install kube-linter
```

### Using nix-shell
To install by using [nix](https://nixos.org/) on macOS, Linux, and [Windows Subsystem for Linux (WSL)](https://docs.microsoft.com/en-us/windows/wsl/about),
run the following command:
```bash
nix-shell -p kube-linter
```

### Using Docker

1. Get the latest KubeLinter Docker image:
   ```bash
   docker pull stackrox/kube-linter:latest
   ```
   > [!NOTE] While we provide the `:latest` tag for convenience and ease of experimentation,
   > we recommend using a tag corresponding to a specific release
   > when incorporating KubeLinter into your workflows to avoid unexpected breakages.
   > See the [Releases](https://github.com/stackrox/kube-linter/releases) page to view
   > available tags.
1. Add path to a directory containing your `yaml` files:
   `docker run` command:
   ```bash
   docker run -v /path/to/files/you/want/to/lint:/dir -v /path/to/config.yaml:/etc/config.yaml stackrox/kube-linter lint /dir --config /etc/config.yaml
   ```

### KubeLinter Github Action

You can also run KubeLinter as a GitHub Action. To use the KubeLinter Github Action, create a `kubelint.yml` file (or choose custom `*.yml` file name) in the `.github/workflows/` directory and use `stackrox/kube-linter-action@v1`.
```yaml
- name: Scan yamls
  id: kube-lint-scan
  uses: stackrox/kube-linter-action@v1
  with:
    directory: yamls
    config: .kube-linter/config.yaml
```
The KubeLinter Github Action accepts the following inputs:

|Parameter|Description|
|:--|:--|
|`directory`|(Mandatory) A directory path that contains the Kubernetes YAML files or `Chart.yaml` file.|
|`config`|(Optional) A path to your custom [KubeLinter configuration file](configuring-kubelinter.md).

## Development

### Building from source

> [!NOTE] Before you build, make sure that you have [installed Go](https://golang.org/doc/install).

To build KubeLinter from source:

1. Clone the KubeLinter repository:
   ```bash
   git clone git@github.com:stackrox/kube-linter.git
   ```
1. Compile the source code:
   ```bash
   make build
   ```
   This command compiles the source code and creates a `kube-linter` binary file
   for your platform in the `.gobin` folder.
1. Verify that the compiled binary is working:
   ```bash
   .gobin/kube-linter version
   ```
1. (Optional) Add the generated binary to your path. Run the following command and
   add the output to your shell profile (`~/.bash_profile`,
   `~/.bashrc` or `~/.zshenv`):
   ```bash
   echo export PATH='"${PATH}:'"$(pwd)/.gobin"'"'
   ```

## Verifying KubeLinter images

KubeLinter images are signed by [cosign](https://github.com/sigstore/cosign).
We recommend verifying the image before using it.

Once you've installed cosign, you can use the [KubeLinter public key](https://github.com/stackrox/kube-linter/blob/main/kubelinter-cosign.pub) to verify the KubeLinter image with:

```shell
cat kubelinter-cosign.pub
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEl0HCkCRzYv0qH5QiazoXeXe2qwFX
DmAszeH26g1s3OSsG/focPWkN88wEKQ5eiE95v+Z2snUQPl/mjPdvqpyjA==
-----END PUBLIC KEY-----


cosign verify --key kubelinter-cosign $IMAGE_NAME
```

KubeLinter also provides [cosign keyless signatures](https://github.com/sigstore/cosign/blob/623d50f9b77ee85886a166daac648455e65003ec/KEYLESS.md).

You can verify the KubeLinter image with:
```shell
# NOTE: Keyless signatures are NOT PRODUCTION ready.

COSIGN_EXPERIMENTAL=1 cosign verify $IMAGE_NAME
```

## Usage

<!-- tabs:start -->

#### ** Kubernetes **

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
   > [!NOTE] This sample file has two production readiness issues and one
   > security issue.
   >
   > **Security issue**
   > - The container in this pod is not running as a read-only file system,
   >   allowing it to write to the root filesystem.
   >
   > **Production readiness issue**
   > - The configuration doesn't specify the container's CPU limits,
   >   allowing it to consume excessive CPU.
   > - The configuration doesn't specify the container's memory limits,
   >   allowing it to consume excessive memory.

1. To lint this file with KubeLinter, run the following command:
   ```bash
   kube-linter lint pod.yaml
   ```
1. KubeLinter runs the default checks and reports errors.
   ```
   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container "sec-ctx-demo" does not have a read-only root file system (check: no-read-only-root-fs, remediation: Set readOnlyRootFilesystem to true in your container's securityContext.)

   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container "sec-ctx-demo" has cpu limit 0 (check: unset-cpu-requirements, remediation: Set your container's CPU requests and limits depending on its requirements. See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for more details.)

   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container "sec-ctx-demo" has memory limit 0 (check: unset-memory-requirements, remediation: Set your container's memory requests and limits depending on its requirements.    See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for more details.)

   Error: found 3 lint errors
   ```


#### ** Helm **

To run KubeLinter on Helm charts, provide a path to the directory which contains
the `Chart.yaml` file. For example, consider running KubeLinter on a sample Helm
chart:

1. Create a new Helm chart:
   ```bash
   helm create helm-chart-sample
   ```
1. To lint this Helm chart with KubeLinter, run the following command:
   ```bash
   kube-linter lint helm-chart-sample/
   ```
1. KubeLinter runs the default checks and reports errors.
   ```
   helm-chart-sample/helm-chart-sample/templates/tests/test-connection.yaml: (object: <no namespace>/test-release-helm-chart-sample-test-connection /v1, Kind=Pod) container "wget" does not have a read-only root file system (check: no-read-only-root-fs, remediation: Set readOnlyRootFilesystem to true in your container's securityContext.)

   helm-chart-sample/helm-chart-sample/templates/tests/test-connection.yaml: (object: <no namespace>/test-release-helm-chart-sample-test-connection /v1, Kind=Pod) container "wget" is not set to runAsNonRoot (check: run-as-non-root, remediation: Set runAsUser to a non-zero number, and runAsNonRoot to true, in your pod or container securityContext. See https://kubernetes.io/docs/tasks/configure-pod-container/security-context/ for more details.)

   helm-chart-sample/helm-chart-sample/templates/tests/test-connection.yaml: (object: <no namespace>/test-release-helm-chart-sample-test-connection /v1, Kind=Pod) container "wget" has cpu request 0 (check: unset-cpu-requirements, remediation: Set your container's CPU requests and limits depending on its requirements. See https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for more details.)

   ...

   Error: found 12 lint errors
   ```


<!-- tabs:end -->


For more details about using and configuring KubeLinter, see the
[Using KubeLinter](using-kubelinter.md) topic.

# Community

To engage with the KubeLinter community, including maintainers and other
  users, join [KubeLinter on Slack <span class="iconify" data-icon="logos:slack-icon"></span>](https://kube-linter.slack.com/join/shared_invite/zt-icv44kde-gfpmAtrT6toeqYYd7JOVTA#/).

To contribute, see the [contributing guide](https://github.com/stackrox/kube-linter/blob/main/CONTRIBUTING.md).

> [!ATTENTION]
> Our [code of conduct](https://github.com/stackrox/kube-linter/blob/main/CODE_OF_CONDUCT.md) governs all participation in the KubeLinter community.

# License

KubeLinter is licensed under the [Apache License 2.0](https://github.com/stackrox/kube-linter/blob/main/LICENSE).
