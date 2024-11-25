# Paths

`yamlfmt` can collect paths in two modes: Standard, or Doublestar.

## Standard (default)

In standard path mode, you can specify a file or directory path directly. If specifying a file, it will simply include the file. If specifying a directory, it will include every file with the correct extension (as specified in `extensions`, default is `yml` and `yaml`).

This mode does *not* support wildcards, aka. globbing. That means with `*.yaml` yamlfmt will look for a file named asterisk dot yaml. If you require globbing, use the [Doublestar mode](#doublestar) instead.

## Doublestar

In Doublestar mode, paths are specified using wildcard patterns explained in the [doublestar](https://github.com/bmatcuk/doublestar) package. It is almost identical to bash and git's style of glob pattern specification.

To enable the doublestar mode, set `doublestar: true` in the config file or use the `-dstar` command line flag.

## Include and Exclude

In both modes, `yamlfmt` will allow you to configure include and exclude paths. These can be paths to files in Standard or Doublestar modes, paths to directories in Standard mode, and valid doublestar patterns in Doublestar mode. These paths should be specified **relative to the working directory of `yamlfmt`**. They will work as absolute paths if both the includes and excludes are specified as absolute paths or if both are relative paths, however it will not work as expected if they are mixed together. It usually easier to reason about includes and excludes when always specifying both as relative paths from the directory `yamlfmt` is going to be run in.

Exclude paths can be specified on the command line using the `-exclude` flag.
Paths excluded from the command line are **added* to excluded paths from the config file.

Include paths can be specified on the command line via the positional arguments, i.e. there is no flag for it.
Paths from the command line take precedence over and **replace** any paths configured in the config file.

yamlfmt will build a list of all files to format using the include list, then discard any files matching the exclude list.

## Extensions

*Only in standard mode*

By default, yamlfmt formats all files ending in `.yaml` and `.yml`.
You can modify this behavior using the config file and command line flags.

The config file **sets** the list of extensions.
For example, with `extensions: ["foo"]`, yamlfmt will only match files ending in `.foo` and will *not* match files ending in `.yaml` or `yml`.
An empty list triggers the default behavior.

The `-extensions` command line flag **adds** to the list of extensions from the config file.
For example, `-extensions yaml.gotmpl` will match files ending in `.yaml.gotmpl` *in addition to* files ending in `.yaml` and `.yml`.
