# Basic Formatter

The basic formatter is a barebones formatter that simply takes the data provided, serializes it with [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) and encodes it again. This provides a consistent output format that is very opinionated and cannot be configured.

## Configuration

| Key                      | Default | Description |
|:-------------------------|:--------|:------------|
| `indent`                 | 2       | The indentation level in spaces to use for the formatted yaml|
| `include_document_start` | false   | Include `---` at document start |
