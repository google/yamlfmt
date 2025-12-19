# Configuration File

## Config File Discovery

The config file is a file is a YAML file that contains a valid yamlfmt configuration. yamlfmt will automatically search for files with the following names:

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

In the `-conf` flag, the config file can be named anything. As long as it's valid YAML, yamlfmt will read it as a config file. This can be useful for applying unique configs to different directories in a project. The automatic discovery paths do need to use one of the known names.

In the `-print_conf` flag, merged config values will be printed.

## Command

The command package defines the main command engine that `cmd/yamlfmt` uses. It uses the top level configuration that any run of the yamlfmt command will use.

### Configuration

| Key                      | Type           | Default      | Description |
|:-------------------------|:---------------|:-------------|:------------|
| `line_ending`            | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This global setting will override any formatter `line_ending` options. |
| `doublestar`             | bool                | false         | Use [doublestar](https://github.com/bmatcuk/doublestar) for include and exclude paths. (This was the default before 0.7.0) |
| `continue_on_error`      | bool                | false         | Continue formatting and don't exit with code 1 when there is an invalid YAML file found. |
| `match_type`             | string              | `standard`    | Controls how `include` and `exclude` are interpreted. See [Specifying Paths][] for more details. |
| `include`                | []string            | []            | The paths for the command to include for formatting. See [Specifying Paths][] for more details. |
| `exclude`                | []string            | []            | The paths for the command to exclude from formatting. See [Specifying Paths][] for more details. |
| `gitignore_excludes`     | bool                | false         | Use gitignore files for exclude paths. This is in addition to the patterns from the `exclude` option. |
| `gitignore_path`         | string              | `.gitignore`  | The path to the gitignore file to use. |
| `regex_exclude`          | []string            | []            | Regex patterns to match file contents for, if the file content matches the regex the file will be excluded. Use [Go regexes](https://regex101.com/). |
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
To change to a different formatter, you can change the `type`:
```yaml
formatter:
  type: kyaml
```
Keep in mind that each formatter has its own distinct set of configuration, and attempting to provide options for one formatter to another will cause `yamlfmt` to fail.

## Basic Formatter

The `basic` formatter is the default formatter that takes the data provided, serializes it with a fork of [gopkg.in/yaml.v3][1] and encodes it again. This provides a consistent output format that is very opinionated with some minor tweak options.

### Configuration

| Key                         | Type           | Default | Description |
|:----------------------------|:---------------|:--------|:------------|
| `indent`                    | int            | 2       | The indentation level in spaces to use for the formatted YAML. |
| `include_document_start`    | bool           | false   | Include `---` at document start. |
| `line_ending`               | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This setting will be overwritten by the global `line_ending`. |
| `retain_line_breaks`        | bool           | false   | Retain line breaks in formatted YAML. |
| `retain_line_breaks_single` | bool           | false   | (NOTE: Takes precedence over `retain_line_breaks`) Retain line breaks in formatted YAML, but only keep a single line in groups of many blank lines. |
| `disallow_anchors`          | bool           | false   | If true, reject any YAML anchors or aliases found in the document. |
| `max_line_length`           | int            | 0       | Set the maximum line length ([see note below](#max_line_length)). if not set, defaults to 0 which means no limit. |
| `scan_folded_as_literal`    | bool           | false   | Option that will preserve newlines in folded block scalars (blocks that start with `>`). |
| `indentless_arrays`         | bool           | false   | Render `-` array items (block sequence items) without an increased indent. |
| `drop_merge_tag`            | bool           | false   | Assume that any well formed merge using just a `<<` token will be a merge, and drop the `!!merge` tag from the formatted result. |
| `pad_line_comments`         | int            | 1       | The number of padding spaces to insert before line comments. |
| `trim_trailing_whitespace`  | bool           | false   | Trim trailing whitespace from lines. |
| `eof_newline`               | bool           | false   | Always add a newline at end of file. Useful in the scenario where `retain_line_breaks` is disabled but the trailing newline is still needed. |
| `strip_directives`          | bool           | false   | [YAML Directives](https://yaml.org/spec/1.2.2/#3234-directives) are not supported by this formatter. This feature will attempt to strip the directives before formatting and put them back. [Use this feature at your own risk.](#strip_directives) |
| `array_indent`              | int            | = indent | Set a different indentation level for block sequences specifically. |
| `indent_root_array`         | bool           | false   | Tells the formatter to indent an array that is at the lowest indentation level of the document. |
| `disable_alias_key_correction` | bool        | false   | Disables functionality to fix alias nodes being used as keys. See #247 for details. |
| `force_array_style`         | `flow`, `block`, or empty | empty   | If set, forces arrays to be output in a particular style, either `flow` (`[]`) or `block` (`- x`). If unset, the style from the original document is used. |

### Additional Notes

#### `max_line_length`

It's not perfect; it uses the `best_width` setting from the [gopkg.in/yaml.v3][1] library. If there's a very long token that extends too far for the line width, it won't split it up properly. I will keep trying to make this work better, but decided to get a version of the feature in that works for a lot of scenarios even if not all of them.

#### `strip_directives`

TL;DR:
* If you only have directives at the top of the file this feature will work just fine, otherwise make sure you test it first.
* Please note that directives are completely tossed and ignored by the formatter

This hotfix is flaky. It is very hard to reconstruct data like this without parsing or knowing what may have changed about the structure of the document. What it attempts to do is find the directives in the document before formatting, keep track of them, and put them back where they "belong". However, the only mechanism it has to decide where it "belongs" is the line it was at originally. This can easily change based on what the formatter ended up changing. This means that the only way this fix predictably and reliably works is for directives that are at the top of the document before the document actually starts (i.e. where the `%YAML` directive goes).

In addition, while with this feature the `%YAML` directive may work, the formatter very specifically supports only the [YAML 1.2 spec](https://yaml.org/spec/1.2.2/). So the `%YAML:1.0` directive won't have the desired effect when passing a file through `yamlfmt`, and if you have 1.0-only syntax in your document the formatter may end up failing in other ways that will be unfixable.

## KYAML Formatter

The `kyaml` formatter will read any YAML file and output it in [KYAML format](https://kubernetes.io/docs/reference/encodings/kyaml/). This formatter can read any valid YAML document and output it in KYAML format (this includes KYAML documents, which are themselves valid YAML documents).

As of writing, this formatter takes no configuration; as it is a completely separate formatter, it does not accept any of the [Basic Formatter config options](#configuration) and `yamlfmt` will fail if you attempt to provide those options to the `kyaml` formatter.

[1]: ../pkg/yaml/
[Specifying Paths]: ./paths.md
