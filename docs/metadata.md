# Metadata

The `yamlfmt` library supports recognizing a limited amount of metadata from a yaml file.

## How to specify

Metadata is specified with a special token, followed by a colon, and then a type. For example, to add `ignore` metadata to a file:
```
# !yamlfmt!:ignore
```
If this string `!yamlfmt!:ignore` is anywhere in the file, the file will be dropped from the paths to format.

The format of `!yamlfmt!:type` is strict; there must be a colon separating the metadata identifier and the type, and there must be no whitespace separating anything within the metadata identifier block. For example either of these will cause an error:
```
# !yamlfmt!: ignore
# !yamlfmt!ignore
```
Metadata errors are considered non-fatal, and `yamlfmt` will attempt to continue despite them.

## Types

| Type   | Example            | Description | 
|:-------|:-------------------|:------------|
| ignore | `!yamlfmt!:ignore` | If found, `yamlfmt` will exclude the file from formatting. |