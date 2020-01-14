# golangci-linter

This action lints your Golang Project with GolangCI-Lint simply.

## Feature

* Output lint result as error annotations of GitHub Actions.

## Demo

[Demo](https://github.com/Mushus/golangci-linter/commit/2c532c7486a7dbb12096082d74c5f94273596325#annotation_9255967392559673)

## Inputs

### `config`

**Optional** Path of the GolangCI-Lint config file. Default `none`.

### `basePath`

**Optional** Path of the Golang project root directory. Default `"."`

## Outputs

nothing.

## Example usage

Same as `golangci-lint run`

```yaml
uses: Mushus/golangci-linter@v1
```

Same as `golangci-lint run .golangci.yml`

```yaml
uses: Mushus/golangci-linter@v1
with:
  config: .golangci.yml
```

Using a prebuild image from DockerHub

```yaml
uses: docker://mushus/golangci-linter:1.1.1
```
