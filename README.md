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

## Usage

By default, the tool will recursively find all files that match the glob path `**/*.{yaml,yml}` extension and attempt to format them with the [basic formatter](formatters/basic). To run the tool with all default settings, simply run the command with no arguments:
```bash
yamlfmt
```
You can also run the command with paths to each individual file, or with glob paths:
```bash
yamlfmt x.yaml y.yaml config/**/*.yaml
```
(NOTE: Glob paths are implemented using the [doublestar](https://github.com/bmatcuk/doublestar) package, which is far more flexible than Go's glob implementation. See the doublestar docs for more details.)

## Configuration

The tool can be configured with a `.yamlfmt` configuration file in the working directory. The configuration is specified in yaml.

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


(NOTE: If providing paths as command line arguments, the flags must be specified before any paths)
