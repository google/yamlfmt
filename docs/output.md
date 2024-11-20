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

## `gitlab`

Generates a [GitLab Code Quality report](https://docs.gitlab.com/ee/ci/testing/code_quality.html#code-quality-report-format).

Example:

```json
[
  {
    "description": "Not formatted correctly, run yamlfmt to resolve.",
    "check_name": "yamlfmt",
    "fingerprint": "c1dddeed9a8423b815cef59434fe3dea90d946016c8f71ecbd7eb46c528c0179",
    "severity": "major",
    "location": {
      "path": ".gitlab-ci.yml"
    }
  },
]
```

To use in a GitLab CI pipeline, first write the Code Quality report to a file, then upload the file as a Code Quality artifact.
Abbreviated example:

```yaml
yamlfmt:
  script:
    - yamlfmt -dry -output_format gitlab . >yamlfmt-report
  artifacts:
    when: always
    reports:
      codequality: yamlfmt-report
```

With `-quiet`, the GitLab format will omit unnecessary whitespace to produce a more compact output.
