# Paths

`yamlfmt` can collect paths in two modes: Standard, or Doublestar.

## Standard (default)

In standard path mode, the you can specify a file or directory path directly. If specifying a file, it will simply include the file. If specifying a directory, it will include every file with the correct extension (as specified in `extensions`, default is `yml` and `yaml`).

## Doublestar

In Doublestar mode, paths are specified using the format explained in the [doublestar](https://github.com/bmatcuk/doublestar) package. It is almost identical to bash and git's style of glob pattern specification.
