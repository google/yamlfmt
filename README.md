# yamlfmt

`yamlfmt` is an extensible command line tool or library to format yaml files. 

## Goals

* Create a command line yaml formatting tool that is easy to distribute (single binary)
* Make it simple to extend with new custom formatters
* Enable alternative use as a library, providing a foundation for users to create a tool that meets specific needs 

## Maintainers

This tool is not yet officially supported by Google. It is currenlty maintained solely by @braydonk.

## Installation

To download the `yamlfmt` command, you can download the desired binary from releases or install the module directly:
```
go install github.com/google/yamlfmt/cmd/yamlfmt@latest
```
NOTE: Recommended setup if this is your first time installing Go would be in [this DigitalOcean blog post](https://www.digitalocean.com/community/tutorials/how-to-build-and-install-go-programs).

You can also simply download the binary you want from releases. The binary is self-sufficient with no dependencies, and can simply be put somewhere on your PATH and run with the command `yamlfmt`.

## Usage

To run the tool with all default settings, simply run the command with a path argument:
```bash
yamlfmt x.yaml y.yaml <...>
```
You can specify as many paths as you want. You can also specify a directory which will be searched recursively for any files with the extension `.yaml` or `.yml`.
```bash
yamlfmt .
```

You can also use an alternate mode that will search paths with doublestar globs by supplying the `-dstar` flag. 
```bash
yamlfmt -dstar **/*.{yaml,yml}
```
See the [doublestar](https://github.com/bmatcuk/doublestar) package for more information on this format.

## Flags

The CLI supports the following flags/arguments:

* Format (default, no flags)
	- Format and write the matched files
* Dry run (`-dry` flag)
	- Format the matched files and output the diff to `stdout`
* Lint (`-lint` flag)
	- Format the matched files and output the diff to `stdout`, exits with status 1 if there are any differences
* Stdin (just `-` or `/dev/stdin` argument, or `-in` flag)
	- Format the yaml data from `stdin` and output the result to `stdout`
* Custom config path (`-conf` flag)
	- If you would like to use a config not stored at `.yamlfmt` in the working directory, you can pass a relative or absolute path to a separate configuration file
* Doublestar path collection (`-dstar` flag)
	- If you would like to use 

(NOTE: If providing paths as command line arguments, the flags must be specified before any paths)

# Configuration File

The `yamlfmt` command can be configured through a yaml configuration file. The tool looks for the config file in the following order:

1. Specified in the `--conf` flag (if this is an invalid path or doesn't exist, the tool will fail)
2. A `.yamlfmt` file in the current working directory
3. A `yamlfmt` folder with a `.yamlfmt` file in the system config directory (`$XDG_CONFIG_HOME`, `$HOME/.config`, `%LOCALAPPDATA%`) e.g. `$HOME/.config/yamlfmt/.yamlfmt`

If none of these are found, the tool's default configuration will be used.

### Include/Exclude

If you would like to have a consistent configuration for include and exclude paths, you can also use a configuration file. The tool will attempt to read a configuration file named `.yamlfmt` in the directory the tool is run on. In it, you can configure paths to include and exclude, for example:
```yaml
include:
  - config/**/*.{yaml,yml}
exclude:
  - excluded/**/*.yaml
```

### Line Ending

The default line ending is `lf` (Unix style, Mac/Linux). The line ending can be changed to `crlf` (Windows style) with the `line_ending` setting:
```yaml
line_ending: crlf
```
This setting will be sent to any formatter as a config field called `line_ending`. If a `line_ending` is specified in the formatter, this will overwrite it. New formatters are free to ignore this setting if they don't need it, but any formatter provided by this repo will handle it accordingly.

### Formatter

In your `.yamlfmt` file you can also specify configuration for the formatter if that formatter supports it. To change the indentation level of the basic formatter for example:
```yaml
formatter:
  type: basic
  indent: 4
```
If the type is not specified, the default formatter will be used. In the tool included in this repo, the default is the [basic formatter](formatters/basic).

For in-depth configuration documentation see the [config docs](docs/config.md).

# pre-commit

Starting in v0.7.1, `yamlfmt` can be used as a hook for the popular [pre-commit](https://pre-commit.com/) tool. To include a `yamlfmt` hook in your `pre-commit` config, add the following to the `repos` block in your `.pre-commit-config.yaml`:

```yaml
- repo: https://github.com/google/yamlfmt
  rev: v0.7.1
  hooks:
    - id: yamlfmt
```

When running yamlfmt with the `pre-commit` hook, the only way to configure it is through a `.yamlfmt` configuration file in the root of the repo or a system wide config directory (see [Configuration File](#configuration-file)). 