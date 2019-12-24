# golangci-linter

This action lints your Golang Project with GolangCI-Lint simply.

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

Using a prebuild image from docker

```yaml
uses: docker://mushus/golangci-linter:1.0.0
```
