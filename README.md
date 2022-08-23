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

### Formatter

In your `.yamlfmt` file you can also specify configuration for the formatter if that formatter supports it. To change the indentation level of the basic formatter for example:
```yaml
formatter:
  type: basic
  indentation: 4
```
If the type is not specified, the default formatter will be used. In the tool included in this repo, the default is the [basic formatter](formatters/basic).

## Flags

The CLI supports 3 operation modes:

* Format (default, no flags)
    - Format and write the matched files
* Dry run (`-dry` flag)
    - Format the matched files and output the diff to `stdout`
* Lint (`-lint` flag)
    - Format the matched files and output the diff to `stdout`, exits with status 1 if there are any differences

(NOTE: If providing paths as command line arguments, the flags must be specified before any paths)
