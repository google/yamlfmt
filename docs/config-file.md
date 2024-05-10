# Configuration File

## Config File Discovery

The config file is a file is a yaml file that contains a valid yamlfmt configuration. yamlfmt will automatically search for files with the following names:

  - `.yamlfmt`
  - `yamlfmt.yml`
  - `yamlfmt.yaml`
  - `.yamlfmt.yaml`
  - `.yamlfmt.yml`

You can also pass a config file to yamlfmt using the `-conf` flag. When using the `-conf` flag, the config file can be named anything.
When not using one of the config file flags, it will be automatically discovered in the following priority order:

1. Specified in the `-conf` flag (if this is an invalid path or doesn't exist, the tool will fail)
1. A config file in the current working directory
1. The first config file found up the tree step by step from the current working directory
1. A `yamlfmt` folder with a config file in the system config directory (`$XDG_CONFIG_HOME`, `$HOME/.config`, `%LOCALAPPDATA%`) e.g. `$HOME/.config/yamlfmt/.yamlfmt`

If none of these are found, the tool's default configuration will be used.

### Config File Discovery Caveats

If the flag `-global_conf` is passed, all other steps will be circumvented and the config file will be discovered from the system config directory. See [the command line flag docs](./command-usage.md#configuration-flags).

In the `-conf` flag, the config file can be named anything. As long as it's valid yaml, yamlfmt will read it as a config file. This can be useful for applying unique configs to different directories in a project. The automatic discovery paths do need to use one of the known names.

In the `-resolved_conf` flag, merged config values will be printed.

## Command

The command package defines the main command engine that `cmd/yamlfmt` uses. It uses the top level configuration that any run of the yamlfmt command will use.

### Configuration

| Key                      | Type           | Default      | Description |
|:-------------------------|:---------------|:-------------|:------------|
| `line_ending`            | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This global setting will override any formatter `line_ending` options. |
| `doublestar`             | bool                | false         | Use [doublestar](https://github.com/bmatcuk/doublestar) for include and exclude paths. (This was the default before 0.7.0) |
| `continue_on_error`      | bool                | false         | Continue formatting and don't exit with code 1 when there is an invalid yaml file found. |
| `include`                | []string            | []            | The paths for the command to include for formatting. See [Specifying Paths][] for more details. |
| `exclude`                | []string            | []            | The paths for the command to exclude from formatting. See [Specifying Paths][] for more details. |
| `gitignore_excludes`     | bool                | false         | Use gitignore files for exclude paths. This is in addition to the patterns from the `exclude` option. |
| `gitignore_path`         | string              | `.gitignore`  | The path to the gitignore file to use. |
| `regex_exclude`          | []string            | []            | Regex patterns to match file contents for, if the file content matches the regex the file will be excluded. Use [Golang regexes](https://regex101.com/). |
| `extensions`             | []string            | []            | The extensions to use for standard mode path collection. See [Specifying Paths][] for more details. |
| `formatter`              | map[string]any      | `type: basic` | Formatter settings. See [Formatter](#formatter) for more details. |
| `output_format`          | `default` or `line` | `default`     | The output format to use. See [Output docs](./output.md) for more details. |

## Formatter

Formatter settings are specified by giving a formatter type in the `type` field, and specifying the rest of the formatter settings in the same block. For example, to get a default `basic` formatter, use the following configuration:
```yaml
formatter:
  type: basic
```
To include other settings for the basic formatter, include them in the same `formatter` block.
```yaml
formatter:
  type: basic
  include_document_start: true
```
Currently, there is only a `basic` formatter, however there is full support for making your own formatter in a fork or the potential for a new formatter to exist in the future.

## Basic Formatter

The basic formatter is a barebones formatter that simply takes the data provided, serializes it with a fork of [gopkg.in/yaml.v3][1] and encodes it again. This provides a consistent output format that is very opinionated with some minor tweak options.

### Configuration

| Key                         | Type           | Default | Description |
|:----------------------------|:---------------|:--------|:------------|
| `indent`                    | int            | 2       | The indentation level in spaces to use for the formatted yaml. |
| `include_document_start`    | bool           | false   | Include `---` at document start. |
| `line_ending`               | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This setting will be overwritten by the global `line_ending`. |
| `retain_line_breaks`        | bool           | false   | Retain line breaks in formatted yaml. |
| `retain_line_breaks_single` | bool           | false   | (NOTE: Takes precedence over `retain_line_breaks`) Retain line breaks in formatted yaml, but only keep a single line in groups of many blank lines. |
| `disallow_anchors`          | bool           | false   | If true, reject any YAML anchors or aliases found in the document. |
| `max_line_length`           | int            | 0       | Set the maximum line length (see notes below). if not set, defaults to 0 which means no limit. |
| `scan_folded_as_literal`    | bool           | false   | Option that will preserve newlines in folded block scalars (blocks that start with `>`). |
| `indentless_arrays`         | bool           | false   | Render `-` array items (block sequence items) without an increased indent. |
| `drop_merge_tag`            | bool           | false   | Assume that any well formed merge using just a `<<` token will be a merge, and drop the `!!merge` tag from the formatted result. |
| `pad_line_comments`         | int            | 1       | The number of padding spaces to insert before line comments. |

### Note on `max_line_length`

It's not perfect; it uses the `best_width` setting from the [gopkg.in/yaml.v3][1] library. If there's a very long token that extends too far for the line width, it won't split it up properly. I will keep trying to make this work better, but decided to get a version of the feature in that works for a lot of scenarios even if not all of them.

[1]: https://www.github.com/braydonk/yaml
[Specifying Paths]: ./paths.md
