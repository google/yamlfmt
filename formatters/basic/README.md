# Basic Formatter

The basic formatter is a barebones formatter that simply takes the data provided, serializes it with [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) and encodes it again. This provides a consistent output format that is very opinionated and cannot be configured.

## Configuration

| Key                      | Type           | Default | Description |
|:-------------------------|:---------------|:--------|:------------|
| `indent`                 | int            | 2       | The indentation level in spaces to use for the formatted yaml|
| `include_document_start` | bool           | false   | Include `---` at document start |
| `line_ending`            | `lf` or `crlf` | `crlf` on Windows, `lf` otherwise | Parse and write the file with "lf" or "crlf" line endings. This setting will be overwritten by the global `line_ending`. |
| `emoji_support`          | bool           | false   | Support encoding utf-8 emojis |
| `retain_line_breaks`     | bool           | false   | Retain line breaks in formatted yaml |
| `disallow_anchors`       | bool           | false   | If true, reject any YAML anchors or aliases found in the document. |
