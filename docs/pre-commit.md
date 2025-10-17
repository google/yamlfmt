# pre-commit

NOTE: https://github.com/google/yamlfmt/discussions/278 

Starting in v0.7.1, `yamlfmt` can be used as a hook for the popular [pre-commit](https://pre-commit.com/) tool. To include a `yamlfmt` hook in your `pre-commit` config, add the following to the `repos` block in your `.pre-commit-config.yaml`:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.19.0
  hooks:
    - id: yamlfmt
```

When running yamlfmt with the `pre-commit` hook, the only way to configure it is through a `.yamlfmt` configuration file in the root of the repo or a system wide config directory (see [Configuration File](./config-file.md) docs). 

## Use `yamlfmt` installed on the system instead of pre-commit building with Go

If you would prefer to manage your `yamlfmt` installation yourself, you can have the hook use your installed `yamlfmt` binary instead. As long as `yamlfmt` is in your PATH, you can override the `language` setting to `system`.

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.19.0
  hooks:
    - id: yamlfmt
      language: system
```

## Run `yamlfmt` on other filetypes

By default, `yamlfmt` will run on all staged `.yaml` files. If you want to run on other filetypes, you can override the `types` configuration:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.19.0
  hooks:
    - id: yamlfmt
      types: [file]
      files: <filepath regex>
```

You can read more on file filtering in [the pre-commit docs for filtering files with `types`](https://pre-commit.com/#filtering-files-with-types).

## Run `yamlfmt` with configuration

If you are providing your own `yamlfmt` configuration file, the default hook experience is going to make `yamlfmt` behave in adverse ways. The experience the hook gives should be fine without a config file, or if your config file only provides `formatter` configuration, but if you provide path configuration there can be strange behaviour. One way around this is to modify the hook not to pass filenames as arguments:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.19.0
  hooks:
    - id: yamlfmt
      pass_filenames: false
```

This will make the entry simply `yamlfmt`, running the tool in the hook the same as the standard running pattern. It will cause the path configuration from the config file to be used.

NOTE: `pre-commit` may create a `.cache` directory that could have `yaml` files in it. If using this mode, you will have to exclude the `.cache` directory using [the standard exclusion method based on your configured match type](./paths.md#include-and-exclude).

## v0.18.0 series

In v0.18.0, I attempted to make a breaking change to run the hook in a different way than I originally had it. This backfired and was broken in many ways. If you used this version, I am sorry for the confusion. v0.19.0 onward will have the original behaviour restored.
