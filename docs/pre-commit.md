# pre-commit

Starting in v0.7.1, `yamlfmt` can be used as a hook for the popular [pre-commit](https://pre-commit.com/) tool. To include a `yamlfmt` hook in your `pre-commit` config, add the following to the `repos` block in your `.pre-commit-config.yaml`:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.18.1
  hooks:
    - id: yamlfmt
```

## Configuration

The default `entry` for the hook is `yamlfmt .`. This is a reasonable default experience if you are not providing `yamlfmt` with a configuration. You can provide configuration either [through a file](./config-file.md) or [through the command line](./command-usage.md). This may require you to override the `entry`. For example, if you have a configuration file with the exact formatting experience you want (all the right files passed in, all the right formatter settings) then you may want to modify the entry to simply run the command with no arguments:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.18.1
  hooks:
    - id: yamlfmt
      entry: yamlfmt
```

You may also wish to provide all your configuration directly through the configuration flags like so:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.18.1
  hooks:
    - id: yamlfmt
      entry: yamlfmt -doublestar true **/*.{yaml,yml}
```

## Use `yamlfmt` installed on the system instead of pre-commit building with Go

If you would prefer to manage your `yamlfmt` installation yourself, you can have the hook use your installed `yamlfmt` binary instead. As long as `yamlfmt` is in your PATH, you can override the `language` setting to `system`.

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.18.1
  hooks:
    - id: yamlfmt
      language: system
```

## Restore old behaviour

In `v0.18.0` and `v0.18.1`, the experience was changed to what is documented here. What is documented here now is the intended experience. However, originally the hook was configured to only run on the `yaml` filetype, and all discovered files would be passed as a list of arguments to the command. This behaviour can be restored like so:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.18.1
  hooks:
    - id: yamlfmt
      entry: yamlfmt
      types: [yaml]
      pass_filenames: true
```
