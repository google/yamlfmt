
# Command

The command package defines the main command engine that `cmd/yamlfmt` uses. It uses the top level configuration that any run of the yamlfmt command will use.

## Configuration

| Key                      | Type           | Default | Description |
|:-------------------------|:---------------|:--------|:------------|
| `line_ending`            | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This global setting will override any formatter `line_ending` options. |
| `doublestar`             | bool           | false   | Use [doublestar](https://github.com/bmatcuk/doublestar) for include and exclude paths. (This was the default before 0.8.0) |
| `include`                | []string       | []      | The paths for the command to include for formatting. See [Specifying Paths](#specifying-paths) for more details. |
| `exclude`                | []string       | []      | The paths for the command to exclude from formatting. See [Specifying Paths](#specifying-paths) for more details. |
| `extensions`             | []string       | []      | The extensions to use for standard mode path collection. See [Specifying Paths](#specifying-paths) for more details. |
| `formatter`              | map[string]any | default basic formatter | Formatter settings. See [Formatter](#formatter) for more details. |

## Specifying paths

### Standard

In standard path mode, the you can specify a file or directory path directly. If specifying a file, it will simply include the file. If specifying a directory, it will include every file with the correct extension (as specified in `extensions`).

### Doublestar

In Doublestar mode, paths are specified using the format explained in the [doublestar](https://github.com/bmatcuk/doublestar) package. It is almost identical to bash and git's style of glob pattern specification.

# Formatter

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

# Basic Formatter

The basic formatter is a barebones formatter that simply takes the data provided, serializes it with a fork of [gopkg.in/yaml.v3](https://www.github.com/braydonk/yaml) and encodes it again. This provides a consistent output format that is very opinionated with some minor tweak options.

## Configuration

| Key                      | Type           | Default | Description |
|:-------------------------|:---------------|:--------|:------------|
| `indent`                 | int            | 2       | The indentation level in spaces to use for the formatted yaml|
| `include_document_start` | bool           | false   | Include `---` at document start |
| `line_ending`            | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This setting will be overwritten by the global `line_ending`. |
| `retain_line_breaks`     | bool           | false   | Retain line breaks in formatted yaml |
| `disallow_anchors`       | bool           | false   | If true, reject any YAML anchors or aliases found in the document. |
| `max_line_length`        | int            | -1      | Set the maximum line length (see notes below) |

## Note on `max_line_length`

It's not perfect; it uses the `best_width` setting from the yaml library. If there's a very long token that extends too far for the line width, it won't split it up properly. I will keep trying to make this work better, but decided to get a version of the feature in that works for a lot of scenarios even if not all of them.