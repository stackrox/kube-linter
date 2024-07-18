<p align="center"><img src="images/logo/KubeLinter-horizontal.svg" width="360"></p>
<p align="center"><b>Static analysis for Kubernetes</b></p>

[![Go Report Card](https://goreportcard.com/badge/github.com/stackrox/kube-linter)](https://goreportcard.com/report/github.com/stackrox/kube-linter)

# What is KubeLinter?

KubeLinter analyzes Kubernetes YAML files and Helm charts, and checks them against a variety of best practices, with a focus on production readiness and security.

KubeLinter runs sensible default checks, designed to give you useful information about your Kubernetes YAML files and Helm charts. This is to help teams check early and often for security misconfigurations and DevOps best practices. Some common examples of these include running containers as a non-root user, enforcing least privilege, and storing sensitive information only in secrets.

KubeLinter is configurable, so you can enable and disable checks, as well as create your own custom checks, depending on the policies you want to follow within your organization.

When a lint check fails, KubeLinter reports recommendations for how to resolve any potential issues and returns a non-zero exit code.

## Documentation
Visit https://docs.kubelinter.io for detailed documentation on installing, using and configuring KubeLinter.

## Installing KubeLinter

Kube-linter binaries could be found here: https://github.com/stackrox/kube-linter/releases/latest

### Using Go

To install using [Go](https://golang.org/), run the following command:

```bash
go install golang.stackrox.io/kube-linter/cmd/kube-linter@latest
```
Otherwise, download the latest binary from [Releases](https://github.com/stackrox/kube-linter/releases) and add it to your
PATH.

### Using Homebrew for macOS or LinuxBrew for Linux

To install using Homebrew or LinuxBrew, run the following command:

```bash
brew install kube-linter
```

### Using nix-shell

```
nix-shell -p kube-linter
```

### Using docker

```
docker pull stackrox/kube-linter:latest
```


## Development

### Prerequisites
- Make sure that you have [installed Go](https://golang.org/doc/install) prior to building from source.

### Building KubeLinter

Installing KubeLinter from source is as simple as following these steps:

1. First, clone the KubeLinter repository.

   ```bash
   git clone git@github.com:stackrox/kube-linter.git
   ```

1. Then, compile the source code. This will create the kube-linter binary files for each platform and places them in the `.gobin` folder.

   ```bash
   make build
   ```

1. Finally, you are ready to start using KubeLinter. Verify your version to ensure you've successfully installed KubeLinter.

   ```bash
   .gobin/kube-linter version
   ```

### Testing KubeLinter
There are several layers of testing. Each layer is expected to pass.

1. `go` unit tests:

   ```bash
   make test
   ```

2. end-to-end integration tests:

   ```bash
   make e2e-test
   ```

3. and finally, end-to-end integration tests using `bats-core`:

   ```bash
   make e2e-bats
   ```

## Verifying KubeLinter images

KubeLinter images are signed by [cosign](https://github.com/sigstore/cosign).
We recommend verifying the image before using it.

Once you've installed cosign, you can use the [KubeLinter public key](kubelinter-cosign.pub) to verify the KubeLinter image with:

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

## Using KubeLinter

### Local YAML Linting

Running KubeLinter to Lint your YAML files only requires two steps in its most basic form.

1. Locate the YAML file you'd like to test for security and production readiness best practices:
1. Run the following command:

   ```bash
   kube-linter lint /path/to/your/yaml.yaml
   ```

### Example

Consider the following sample pod specification file `pod.yaml`. This file has two production readiness issues and one security issue:

**Security Issue:**
1. The container in this pod is not running as a read only file system, which could allow it to write to the root filesystem.

**Production readiness:**
1. The container's memory limits are not set, which could allow it to consume excessive memory

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

1. Copy the YAML above to pod.yaml and lint this file by running the following command:

   ```bash
   kube-linter lint pod.yaml
   ```
1. KubeLinter runs its default checks and reports recommendations. Below is the output from our previous command.

   ```
   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) The container "sec-ctx-demo" is using an invalid container image, "busybox". Please use images that are not blocked by the `BlockList` criteria : [".*:(latest)$" "^[^:]*$" "(.*/[^:]+)$"] (check: latest-tag, remediation: Use a container image with a specific tag other than latest.)

   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container "sec-ctx-demo" does not have a read-only root file system (check: no-read-only-root-fs, remediation: Set readOnlyRootFilesystem to true in the container securityContext.)

   pod.yaml: (object: <no namespace>/security-context-demo /v1, Kind=Pod) container "sec-ctx-demo" has memory limit 0 (check: unset-memory-requirements, remediation: Set memory limits for your container based on its requirements. Refer to https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/#requests-and-limits for details.)

   Error: found 3 lint errors
   ```
To learn more about using and configuring KubeLinter, visit the [documentation](./docs) page.

## Mentions/Tutorials

The following are tutorials on KubeLinter written by users. If you have one that you would like to add to this list, please send a PR!

* [Ensuring YAML best practices using KubeLinter](https://www.civo.com/learn/yaml-best-practices-using-kubelinter) at civo.com by Saiyam Pathak.
* [Analyze Kubernetes files for errors with KubeLinter](https://opensource.com/article/21/1/kubelinter) at opensource.com by Jessica Cherry.
* [How to add a new check in KubeLinter?](https://www.psaggu.com/upstream-contribution/2021/08/17/notes.html) by Priyanka Saggu.
* [Extending kube-linter To Build A Custom Template](https://github.com/garethahealy/kubelinter-extending-blog) by [Gareth Healy](https://github.com/garethahealy).

## Community

If you would like to engage with the KubeLinter community, including maintainers and other users, you can join the Slack workspace [here](https://join.slack.com/t/kube-linter/shared_invite/zt-kla9qvyo-Tk~wynTSbr9EE3AjHcv4BQ).

To contribute, check out our [contributing guide](./CONTRIBUTING.md).

As a reminder, all participation in the KubeLinter community is governed by our [code of conduct](./CODE_OF_CONDUCT.md).

## WARNING: Alpha release

KubeLinter is at an early stage of development. There may be breaking changes in
the future to the command usage, flags, and configuration file formats. However,
we encourage you to use KubeLinter to test your environment YAML files, see what
breaks, and [contribute](./CONTRIBUTING.md).

## LICENSE

KubeLinter is licensed under the [Apache License 2.0](./LICENSE).

## StackRox

KubeLinter is made with ❤️ by [StackRox](https://stackrox.com/).

If you're interested in KubeLinter, or in any of the other cool things we do, please know that we're hiring!
Check out our [open positions](https://www.stackrox.com/job-board/). We'd love to hear from you!
