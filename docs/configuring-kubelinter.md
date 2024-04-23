# Configuring KubeLinter

To configure the checks KubeLinter runs or to run your own custom checks, you
can use a `yaml` configuration file. When you run the `lint` command, use the
`--config` option and provide the path to your configuration file.

If a config file is not explicitly provided to the command,
KubeLinter will look for a configuration file in the current
working directory (by order of preference):

1. `.kube-linter.yaml`
1. `.kube-linter.yml`

Finally, if none is found, the default config is used.

```bash
# specific config file
kube-linter lint pod.yaml --config kubelinter-config.yaml

# will search for config based on the above order or will load defaults
kube-linter lint pod.yaml
```

The configuration file has two sections:

1. `customChecks` for configuring custom checks, and
2. `checks` for configuring default checks.

To view a list of all built-in checks, see [KubeLinter checks](generated/checks.md).

## Disable all default checks

To disable all built in checks, set `doNotAutoAddDefaults` to `true`.

```yaml
checks:
  doNotAutoAddDefaults: true
```

> Equivalent CLI flag is `--do-not-auto-add-defaults`

## Run all default checks

To run all built-in checks, set `addAllBuiltIn` to `true`.

```yaml
checks:
  addAllBuiltIn: true
```

## Ignore paths

To ignore checks on files or directories under certain paths, add ignored paths to `ignorePaths`.
Ignore path uses [`**` match syntax](https://pkg.go.dev/github.com/bmatcuk/doublestar#Match)
```yaml
checks:
  ignorePaths:
    - ~/foo/bar/**
    - /**/*/foo/**
    - ../baz/**
    - /tmp/*.yaml
```
> Equivalent CLI flag is `--ignore-paths`


> [!NOTE]
>
> - If you set both `doNotAutoAddDefaults` and `addAllBuiltIn` to `true`,
>   `addAllBuiltIn` takes precedence.

> Equivalent CLI flag is `--add-all-built-in`

## Run specific checks

You can use the `include` and `exclude` keys to run only specific checks. For
example,

- To disable majority of checks and only run few specific checks,
  use `doNotAutoAddDefaults` along with `include`.
  ```yaml
  checks:
    doNotAutoAddDefaults: true
    include:
      - "privileged-container"
      - "run-as-non-root"
  ```
- To run majority of checks and only exclude few specific checks,
  use `addAllBuiltIn` along with `exclude`.
  ```yaml
  checks:
    addAllBuiltIn: true
    exclude:
      - "unset-cpu-requirements"
      - "unset-memory-requirements"
  ```

> Equivalent CLI flags are `--include` and `--exclude` respectively

> [!TIP] > `exclude` always takes precedence, if you include and exclude the same check,
> KubeLinter always skips the check.

## Ignoring violations for specific cases

To ignore violations for specific objects, users can add an annotation with the key
`ignore-check.kube-linter.io/<check-name>`. We strongly encourage adding an explanation as the value for the annotation.
For example, to ignore a check named "privileged-container" for a specific deployment, you can add an annotation like that:

```yaml
metadata:
  annotations:
    ignore-check.kube-linter.io/privileged-container: "This deployment needs to run as privileged because it needs kernel access"
```

To ignore _all_ checks for a specific object, you can use the special annotation key `kube-linter.io/ignore-all`.

## Run custom checks

You can write custom checks based on existing [templates](generated/templates.md). Every template description includes details about the parameters (`params`) you can use along with that template.

For example,

- To make sure that an annotation exists, you can use the [`required-annotation`](generated/templates?id=required-annotation) template:

  ```yaml
  customChecks:
    - name: required-annotation-responsible
      template: required-annotation
      params:
        key: company.io/responsible
  ```

- To make sure that a specific label exists, you can use the [`required-label`](generated/templates?id=required-label) template:
  ```yaml
  customChecks:
    - name: required-label-release
      template: required-label
      params:
        key: company.io/release
  ```

### Extend custom checks

With custom checks, you can control the checks to run only on specific Kubernetes object types (such as services or deployments). You can also modify the remediation message you get when your custom check fails.

For example,

- To make sure that a specific annotation or label exists on all deployments, you can use the respective template e.g. [`required-annotation`](generated/templates?id=required-annotation) and specify a `scope`

  ```yaml
  customChecks:
    - name: required-annotation-responsible
      template: required-annotation
      params:
        key: company.io/responsible
      scope:
        objectKinds:
          - DeploymentLike
  ```

  For details about `objectKinds` that KubeLinter support, see https://github.com/stackrox/kube-linter/tree/main/pkg/objectkinds.

- Use `remediation` to include a remediation message that users get when your custom check fails:
  ```yaml
  customChecks:
    - name: required-annotation-responsible
      template: required-annotation
      params:
        key: company.io/responsible
      remediation: Please set the annotation 'company.io/responsible'. This will be parsed by xy to generate some docs.
  ```

### Custom `objectKinds`

If you want to add a check for an objectKind that isn't built into kube-linter, you can also register your own custom objectKind alongside your custom check. This is especially useful for resources from CRDs, where the objectKind may not exist in every cluster, and isn't a good candidate for upstream support.

`custom_resource_template_test.go` contains an example of a check that looks for excessively long certificate lifetimes in CNCF `cert-manager` Certificate resources.
