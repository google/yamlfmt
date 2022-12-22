# Basic Formatter

The basic formatter is a barebones formatter that simply takes the data provided, serializes it with a fork of [gopkg.in/yaml.v3](https://www.github.com/braydonk/yaml) and encodes it again. This provides a consistent output format that is very opinionated and cannot be configured.

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