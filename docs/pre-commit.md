# pre-commit

Starting in v0.7.1, `yamlfmt` can be used as a hook for the popular [pre-commit](https://pre-commit.com/) tool. To include a `yamlfmt` hook in your `pre-commit` config, add the following to the `repos` block in your `.pre-commit-config.yaml`:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.10.0
  hooks:
    - id: yamlfmt
```

When running yamlfmt with the `pre-commit` hook, the only way to configure it is through a `.yamlfmt` configuration file in the root of the repo or a system wide config directory (see [Configuration File](./config-file.md) docs). 

## Use `yamlfmt` installed on the system instead of pre-commit building with Go

If you would prefer to manage your `yamlfmt` installation yourself, you can have the hook use your installed `yamlfmt` binary instead. As long as `yamlfmt` is in your PATH, you can override the `language` setting to `system`.

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.10.0
  hooks:
    - id: yamlfmt
      language: system
```

