# Paths

`yamlfmt` can collect paths in two modes: Standard, or Doublestar.

## Standard (default)

In standard path mode, you can specify a file or directory path directly. If specifying a file, it will simply include the file. If specifying a directory, it will include every file with the correct extension (as specified in `extensions`, default is `yml` and `yaml`).

## Doublestar

In Doublestar mode, paths are specified using the format explained in the [doublestar](https://github.com/bmatcuk/doublestar) package. It is almost identical to bash and git's style of glob pattern specification.

## Include and Exclude

In both modes, `yamlfmt` will allow you to configure include and exclude paths. These can be paths to files in Standard or Doublestar modes, paths to directories in Standard mode, and valid doublestar patterns in Doublestar mode. These paths should be specified **relative to the working directory of `yamlfmt`**. They will work as absolute paths if both the includes and excludes are specified as absolute paths or if both are relative paths, however it will not work as expected if they are mixed together. It usually easier to reason about includes and excludes when always specifying both as relative paths from the directory `yamlfmt` is going to be run in.
