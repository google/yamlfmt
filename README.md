# yamlfmt

`yamlfmt` is an extensible command line tool or library to format yaml files. 

## Goals

* Create a command line yaml formatting tool that is easy to distribute (single binary)
* Make it simple to extend with new custom formatters
* Enable alternative use as a library, providing a foundation for users to create a tool that meets specific needs 

## Installation

To download the `yamlfmt` command, you can download the desired binary from releases or install the module directly:
```
go install github.com/google/yamlfmt/cmd/yamlfmt@latest
```

## Usage

By default, the tool will recursively find all files that match the glob path `**/*.{yaml,yml}` extension and attempt to format them with the [basic formatter](formatters/basic). To run the tool with all default settings, simply run the command with no arguments:
```
yamlfmt
```
You can also run the command with paths to each individual file, or with glob paths:
```
yamlfmt x.yaml y.yaml config/**/*.yaml
```
(NOTE: Glob paths are implemented using the [doublestar](https://github.com/bmatcuk/doublestar) package, which is far more flexible than Go's glob implementation. See the doublestar docs for more details.)

If you would like to have a consistent configuration for include and exclude paths, you can also use a configuration file. The tool will attempt to read a configuration file named `.yamlfmt` in the directory the tool is run on. In it, you can configure paths to include and exclude, for example:
```
include:
- config/**/*.{yaml,yml}
exclude:
- excluded/**/*.yaml
```

The CLI supports 3 operation modes:

* Format (default)
    - Format and write the matched files
* Dry run (`-dry` flag)
    - Format the matched files and output the diff to `stdout`
* Lint (`-lint` flag)
    - Format the matched files and output the diff to `stdout`, exits with status 1 if there are any differences

(NOTE: If providing paths as command line arguments, the flags must be specified before any paths)