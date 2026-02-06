# Semver

The `semver` package implements logic to work with [Sementic Versioning 2.0.0](http://semver.org/) in Go.
It provides:
- Parser for semantic versions
- Validation of semantic versions
- Sorting of semantic versions
- Parser for semantic version range constraints
- Range constraint matching
    - whether version contained in range
    - whether a range is contained in another range

This library is loosely based on the awesome:
https://github.com/Masterminds/semver

## Parsing Semantic Versions

This library implements strict parsing of semantic versions as outlined in the [Sementic Versioning](http://semver.org/) spec. Shorthand forms or a `v` prefix e.g. 2.0, 1, v2.1.2 are not considered valid semantic versions following the 2.0.0 spec.

> **Parsing Huge Versions**
> Semver does not limit the amount of Major, Minor or Patch version numbers.
Making `99999999999999999999999.999999999999999999.99999999999999999` a valid semver.
> For simplicity this library uses `uint64` as underlying datatype for Major, Minor and Patch, limiting the maximum number to `9999999999999999999`.

When parsing an invalid version errors are annotated with the character column number that an issue was encountered on:

```go
_, err := semver.NewVersion("1.2")
fmt.Println(err)
// Output: col 4: missing patch
```

## Parsing Semantic Version Constraints

Constraints can be used to filter parsed semantic versions. All constraint expressions expand to one or multiple valid semver ranges.

The constraints parser is more permissive to allow shorthand ranges and wildcards.
e.g. `1 - 2` => `1.0.0 - 2.0.0`

Pre release versions that fall within a range will match. If pre release versions should be excluded, they have to be filtered before checking against constraints.

Specifying pre-release ranges is NOT supported.
e.g. `1.0.0-rc.0 - 1.0.0.rc.10`

In the following examples `<max>` is used to denote the max number possible to put into the Major, Minor or Patch section of a semantic version.

### Available Operators:
- `=`: equal, version must match exactly
- `!=`: not equal, excludes a specific version
- `>`: greater than
- `<`: less than
- `>=`: greater than or equal
- `<=`: less than or equal

### Hypen Range

A hyphen range explicitly defines the minimum and maximum version of a range. The min and max values are inclusive.
e.g. A constraint `1.0.0 - 2.0.0` will match `1.0.0` and `2.0.0`, but not `0.9.1` or `2.0.1`.

- `1.2 - 1.4.5` which is equivalent to `>= 1.2.0 <= 1.4.5`
- `2.3.4 - 4.5` which is equivalent to `>= 2.3.4 <= 4.5.0`

### Wildcards

The `x`, `X` and `*` characters can be used as wildcard in version constraints. They will be automatically expanded to either `0` or `<max>`, depending on the operator.

- `1.2.x` is expanded to `1.2.0 - 1.2.<max>`
- `> 1.2.x` is expanded to `1.2.1 - 1.2.<max>`
- `< 2.4.x` is expanded to `0.0.0 - 2.4.<max>`
- `1.2.x - 3.x` is expanded to `1.2.0 - 3.<max>.<max>`

### Tilde Range

The tilde (~) comparison operator is for patch level ranges when a minor version is specified and major level changes when the minor number is missing.

- `~1.2.3` is expanded to `1.2.3 - 1.2.<max>`
- `~2.3` is expanded to `2.3.0 - 2.3.<max>`
- `~1` is expanded to `1.0.0 - 1.<max>.<max>`
- `~1.x` is expanded to `1.0.0 - 1.<max>.<max>`
- `~2.3.x` is expanded to `2.3.0 - 2.3.<max>`

### Caret Range

The caret (^) comparison operator is for major level changes once a stable (`1.0.0`) release has occurred. Prior to a `1.0.0` release the minor versions acts as the API stability level. This is useful when comparisons of API versions as a major change is API breaking.

**stable versions**
- `^1.2.3` is expanded to `1.2.3 - 1.<max>.<max>`
- `^2.3` is expanded to `2.3.0 - 2.<max>.<max>`
- `^2` is expanded to `2.0.0 - 2.<max>.<max>`
- `^1.2.x` is expanded to `1.2.0 - 1.<max>.<max>`
- `^2.x` is expanded to `2.0.0 - 2.<max>.<max>`

**unstable versions**
- `^0.2.3` is expanded to `0.2.3 - 0.2.<max>`
- `^0.2` is expanded to `0.2.0 - 0.2.<max>`
- `^0` is expanded to `0.0.0 - 0.0.<max>`
