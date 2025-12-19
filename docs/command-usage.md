# Command Usage

The `yamlfmt` command supports 3 primary modes of operation and a number of command line configuration options.

## Basic Usage

The most basic usage of the `yamlfmt` command is to use the default configuration and provide a single path to format.

```bash
yamlfmt x.yaml
```

You can also provide a directory path, which yamlfmt will search recursively through for any file with a `yml` or `yaml` extension.

```bash
yamlfmt .
```

You can make use of the alternative Doublestar path mode for more dynamic patterns.

```bash
yamlfmt -dstar "**/*.yaml"
```
(Technically, providing the same pattern in bash without quotes in standard path mode will do standard bash glob expansion and work similarly.)

See [Specifying Paths](./paths.md) for more details about the path collection behaviour.

## Three Modes of Operation

The command supports three modes of operation: Format, Lint, and Dry Run.

### Format

This is the default operation (no flag specified). It will collect all paths that match the include patterns and run them straight through formatting, rewriting each file with the results.

### Format stdin

This mode will read input from stdin and output formatted result to stdout. By using the paths `-`, `/dev/stdin`, or using the `-in` flag, `yamlfmt` will only read from stdin and write the result to stdout.
```bash
cat x.yaml | yamlfmt -
cat x.yaml | yamlfmt /dev/stdin
cat x.yaml | yamlfmt -in
```
(Despite `/dev/stdin` and `-` being Linux/Bash standards, this convention was manually implemented in a cross-platform way. It will function identically on Windows. Would you want to use the Linux standards like this on Windows? I don't know, but it will work :smile:)

### Dry Run

This mode is enabled through the `-dry` flag. This will collect all paths that match the include patterns and run them through formatting, printing each file that had a formatting diff to stdout. This mode is affected by the `-quiet` flag, where only the paths of the files with diffs will be printed.

### Lint

This mode is enabled through the `-lint` flag. This will collect all paths that match the include patterns and run them through formatting, and will exit with code 1 (fail) if any files have formatting differences, outputting the diffs to stdout. This mode is also affected by the `-quiet` flag, where only the paths of the files with diffs will be printed.

## Flags

All flags must be specified **before** any path arguments.

### Operation Flags

These flags adjust the command's mode of operation. All of these flags are booleans.

| Name          | Flag             | Example                     | Description                                               |
| :------------ | :--------------- | :-------------------------- | :-------------------------------------------------------- |
| Help          | `-help`          | `yamlfmt -help`             | Print the command usage information.                      |
| Print Version | `-version`       | `yamlfmt -version`          | Print the yamlfmt version.                                |
| Print Config  | `-print_conf`    | `yamlfmt -print_conf`       | Print the merged configuration to use.                    |
| Dry Run       | `-dry`           | `yamlfmt -dry .`            | Use [Dry Run](#dry-run) mode                              |
| Lint          | `-lint`          | `yamlfmt -lint .`           | Use [Lint](#lint) mode                                    |
| Read Stdin    | `-in`            | `cat x.yaml \| yamlfmt -in` | Read input from stdin and output result to stdout.        |
| Quiet Mode    | `-quiet`, `-q`   | `yamlfmt -dry -q .`         | Use quiet mode. Only has effect in Dry Run or Lint modes. |
| Verbose Mode  | `-verbose`, `-v` | `yamlfmt -v .`              | Use verbose mode. Only has effect in Format mode.         |

### Configuration Flags

These flags will configure the underlying behaviour of the command.

The string array flags can be a bit confusing. See the [String Array Flags](#string-array-flags) section for more information.

| Name                  | Flag                  | Type              | Example                                                   | Description |
|:----------------------|:----------------------|:------------------|:----------------------------------------------------------|:------------|
| Config File Path      | `-conf`               | string            | `yamlfmt -conf ./config/.yamlfmt`                         | Specify a path to read a [configuration file](./config-file.md) from. |
| Global Config         | `-global_conf`        | bool              | `yamlfmt -global_conf`                                    | Force yamlfmt to use the configuration file from the system config directory. |
| Disable Global Config | `-no_global_conf`     | bool              | `yamlfmt -no_global_conf`                                 | Disable looking for the configuration file from the system config directory. |
| Doublestar            | `-dstar`              | bool              | `yamlfmt -dstar "**/*.yaml"`                              | Enable [Doublestar](./paths.md#doublestar) path collection mode. Note that doublestar patterns should be specified with quotes in bash to prevent shell expansion. |
| Match type            | `-match_type`         | string            | `yamlfmt -match_type standard`                            | Controls how `include` and `exclude` are interpreted. See [Specifying Paths](./paths.md) for more details. |
| Exclude               | `-exclude`            | []string          | `yamlfmt -exclude ./not/,these_paths.yaml`                | Patterns to exclude from path collection. These are in addition to the exclude patterns specified in the [config file](./config-file.md) |
| Gitignore Excludes    | `-gitignore_excludes` | bool              | `yamlfmt -gitignore_excludes`                             | Use a gitignore file to exclude paths. This is in addition to otherwise specified exclude patterns. |
| Gitignore Path        | `-gitignore_path`     | string            | `yamlfmt -gitignore_path .special_gitignore`              | Specify a path to a gitignore file to use. Defaults to `.gitignore` (in working directory). |
| Extensions            | `-extensions`         | []string          | `yamlfmt -extensions yaml,yml`                            | Extensions to use in standard path collection. Has no effect in Doublestar mode. These add to extensions specified in the [config file](./config-file.md)
| KYAML                 | `-kyaml`              | bool              | `yamlfmat -kyaml`                                         | Enable the alternate [KYAML formatter](./config-file.md#kyaml-formatter). Note that using this option will completely override any formatter configuration from detected config file. |
| Formatter Config      | `-formatter`          | []string          | `yamlfmt -formatter indent=2,include_document_start=true` | Provide configuration values for the formatter. See [Formatter Configuration Options](./config-file.md#basic-formatter) for options. Each field is specified as `configkey=value`. |
| Debug Logging         | `-debug`              | []string          | `yamlfmt -debug paths,config`                             | Enable debug logging. See [Debug Logging](#debug-logging) for more information. |
| Output Format         | `-output_format`      | `default`, `line` | `yamlfmt -output_format line`                             | Choose a different output format. Defaults to `default`. See [Output docs](./output.md) for more information. |

#### String Array Flags

String array flags can be provided in two ways. For example with a flag called `-arrFlag`:

* Individual flags
    - `-arrFlag a -arrFlag b -arrFlag c`
    - Result: `arrFlag: [a b c]`
* Comma separated value in single flag
    - `-arrFlag a,b,c`
    - Result: `arrFlag: [a b c]`
* Technically they can be combined but why would you?
    - `-arrFlag a,b -arrFlag c`
    - Result: `arrFlag: [a b c]`

## Debug Logging

Debug logging can be enabled through the `-debug` [array flag](#string-array-flags). The following is the list of supported debug codes:
* `paths`
    - Log the details for the path discovery process. Use it to debug your include/exclude patterns.
* `config`
    - Log the details for the configuration loading process. Use it to figure out which config file yamlfmt uses and why.
* `diffs`
    - Log each diff for all formatted files.
* `all`
    - Enable all available debug codes.
