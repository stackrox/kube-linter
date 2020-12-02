# Configuring KubeLinter

To configure the checks KubeLinter runs or to run your own custom checks, you
can use a `yaml` configuration file. When you run the `lint` command, use the
`--config` option and provide the path to your configration file.

```bash
kube-linter lint pod.yaml --config kubelinter-config.yaml
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

## Run all default checks

To run all built-in checks, set `addAllBuiltIn` to `true`.
```yaml
checks:
  addAllBuiltIn: true
```

> [!NOTE] 
> 
> - If you set both `doNotAutoAddDefaults` and `addAllBuiltIn` to `true`,
>   `addAllBuiltIn` takes precedence.

## Run specific checks

You can use the `include` and `exclue` keys to run only specific checks. For
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

> [!TIP]
> `exclude` always takes precedence, if you include and exclude the same check,
> KubeLinter always skips the check.
