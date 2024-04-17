# Output

The tool supports different output formats when writing to stdout/stderr. This can be configured through either a command line flag `-output_format` or through the `output_format` configuration field in the [config file](./config-file.md).

The following values are supported:

## `default`

Example:
```
The following formatting differences were found:

y.yaml:
  a:     a:
-  b: 1    b: 1
+
z.yaml:
  a:     a:
-  b: 1    b: 1
+
x.yaml:
  a:     a:
-  b: 1    b: 1
+
```

## `line`

Example:
```
x.yaml: formatting difference found
y.yaml: formatting difference found
z.yaml: formatting difference found
```