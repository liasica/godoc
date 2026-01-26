# godoc

A small CLI tool that generates OpenAPI/Swagger documentation from a Go project. It is built on `swaggo/swag` and performs post-processing on the generated output (for example, converting enum fields in Swagger into a more user-friendly oneOf/enum representation).

> Note: this repository contains a runnable CLI program whose entry point is `cmd/godoc`. The tool parses `go.mod` (it requires reading the `GOMODCACHE` environment variable) to locate the source code of external dependency modules so types and comments from those modules can be included in the generated documentation.

## Features

- Uses `swaggo/swag` to generate OpenAPI 3 documentation (outputs YAML).
- Supports resolving external module source paths via `replace` in `go.mod` or the module cache, and adds them to the scanning paths.
- Post-processes the generated `swagger.yaml`: beautifies enum descriptions and adjusts `enum` fields based on `x-enum-varnames` / `x-enum-comments` metadata.

## Requirements

- Go 1.20+ (this project uses the `maps` package which requires Go 1.20 or later).
- Dependencies declared in `go.mod` (the repository uses `swaggo/swag`, `gopkg.in/yaml.v3`, `golang.org/x/mod`, etc.).
- The `GOMODCACHE` environment variable must be readable (the program uses it to locate the module cache).

## Installation

### Option 1: Using the install script (recommended)

Use the provided `install.sh` script to automatically download and install the latest version:

```bash
# Download and run the install script
curl -fsSL https://raw.githubusercontent.com/liasica/godoc/main/install.sh | bash
```

Or clone the repository and run locally:

```bash
# Clone the repository
git clone https://github.com/liasica/godoc.git
cd godoc

# Run the install script
bash install.sh
```

The install script will:
- Automatically detect your OS and architecture (supports Linux, macOS, Windows on amd64 and arm64)
- Download the corresponding pre-built binary from GitHub Releases
- Install the binary to `$GOPATH/bin`
- Check the installed version and upgrade automatically if a new version is available
- Verify the installation was successful

> Note: The install script requires a Go environment (to get `GOPATH`) and the `curl` command. After installation, make sure `$GOPATH/bin` is in your `PATH`.

### Option 2: Build from source

If you want to build from source:

```bash
# Clone the repository
git clone https://github.com/liasica/godoc.git
cd godoc

# Build the binary
go build -o godoc ./cmd/godoc

# (Optional) Move the binary to a directory in your PATH
mv godoc $GOPATH/bin/
```

## Quick start

1. Create a default configuration file with the `config init` command:

```bash
# Create .godoc.yaml in the current directory
godoc config init
```

2. Edit the `.godoc.yaml` configuration file according to your project needs.

3. Run the `generate` subcommand to generate documentation:

```bash
# Run generation with the config file
godoc generate -c .godoc.yaml
```

Or run directly with `go run` (for development/debugging):

```bash
# Run generate with an explicit config path
go run ./cmd/godoc generate -c .godoc.yaml

# If you need to download dependencies automatically, use -mod=mod
go run -mod=mod ./cmd/godoc generate -c .godoc.yaml
```

> Note: The `-mod=mod` flag tells Go to automatically download and update dependencies as needed. This is useful when running the tool for the first time or when dependencies are missing.

Generated documentation (YAML) is written by default to:

```
./assets/docs/swagger.yaml
```

## Configuration (`.godoc.yaml`)

`externalPaths`, `paths`, `mainFile`, and `output` can be configured via a YAML file. The default configuration filename is `.godoc.yaml`; you can also specify a different path with the `-config` (or `--config`) command-line flag.

Example configuration:

```yaml
# .godoc.yaml
externalPaths:
  "nexis.run/nexa": "kit/rest"

paths:
  - "./internal/app/rest/route"
  - "./internal/app/rest/controller"
  - "./internal/infrastructure/model"
  - "./internal/infrastructure/vo"
  - "./internal/presentation/pagination"
  - "./internal/presentation/entity"

mainFile: "route.go"
output: "./assets/docs/"
```

Field descriptions:
- `externalPaths`: map[string]string where the key is `module/path` and the value is the subpath inside the module (the tool resolves the real path using `require`/`replace` in `go.mod` and `GOMODCACHE`, then adds it to the scan list).
- `paths`: a list of local paths to scan (relative to the repository root or absolute); these are joined and passed to `swag` as `SearchDir`.
- `mainFile`: the filename that `swag` looks for as the main entry (for example `route.go`).
- `output`: the output directory where `swag` will write `swagger.yaml`.

> Note: when the configuration file is missing, the program will fall back to built-in defaults (see `internal.LoadConfig` for the default configuration).

## Usage notes

- The CLI entry `cmd/godoc/main.go` now exposes distinct subcommands:
  - `generate` — generate documentation. This command accepts an optional `-config` (or `--config`) flag which points to the YAML configuration file; if omitted the tool will use built-in defaults or `.godoc.yaml` in the current directory.
  - `config` — configuration helpers. Use `godoc config init` to create a default `.godoc.yaml` in the current directory.

- Common commands and flags:
  - `godoc generate` — run generation using defaults or `.godoc.yaml` if present.
  - `godoc generate -c path/to/config.yaml` — run generation using an explicit config file.
  - `godoc config init` — create a default `.godoc.yaml` in the current directory.
  - `godoc --version` — print the CLI version (the binary sets the version from `version.go`).
  - `godoc help` or `godoc <command> --help` — show command help and flags.

- Typical generation flow:
  1. Create or edit `.godoc.yaml` (or use `godoc config init` to create a default file).
  2. Run `godoc generate` (or `godoc generate -c .godoc.yaml` / `go run ./cmd/godoc generate -c .godoc.yaml`).
  3. The tool resolves external module paths (using `go.mod` and `GOMODCACHE`), runs `swaggo/swag` format and gen, then post-processes the generated `swagger.yaml`.

- If the program exits with an error mentioning `failed to read GOMODCACHE environment variable`, ensure `GOMODCACHE` is set before running, for example:

```bash
export GOMODCACHE=$(go env GOMODCACHE)
```

- Enum post-processing: after `swag` generates `swagger.yaml`, the program calls `internal.ConvertEnum2OneOf` to adjust the file. This step relies on `enum`, `x-enum-varnames`, and `x-enum-comments` produced by the generator to build friendlier `description` fields and to filter special placeholder values (for example `-`).

## Project structure (brief)

- `cmd/godoc/main.go` — CLI entry, loads config, assembles scan paths, runs `swag` format/gen, and performs post-processing.
- `internal/config.go` — parses `.godoc.yaml` and provides default values when the config file is absent.
- `module.go` — reads `go.mod`, parses `require` and `replace`, and uses `GOMODCACHE` to locate dependency source paths.
- `enum.go` — parses and converts enum-related fields in `swagger.yaml` to produce friendlier descriptions and adjusted `enum` values.
- Other files (like `cache.go`) provide project-level caching/config support; see source code for details.

> Note: to add additional output formats (for example JSON or Go) or change `swag` `OutputTypes`, edit `gen.Config` in `cmd/godoc/main.go`.

## Development

- Quick local build:

```bash
go build ./...
```

- Run unit tests (if any):

```bash
go test ./...
```

- It's recommended to run `gofmt` and `go vet` after changes to keep code style and quality consistent.

## Troubleshooting

- Missing module path or permission issues:
  - Example error message: `module path does not exist or is inaccessible` — check the value of `go env GOMODCACHE` and the `require`/`replace` entries in `go.mod`.
- Enum conversion failures:
  - Ensure `x-enum-varnames` and `x-enum-comments` exist in `swagger.yaml` and that their lengths match the `enum` array.
- If `swag` generation fails, run `swag` manually (or install the required `swaggo` tool with `go get`) and verify `SearchDir` and `MainFile` settings.

## Contributing & contact

Contributions via Pull Request or Issue are welcome. If you want me to translate this README to English, add usage examples (actual route/controller snippets), or add CI configuration, tell me which license and details you prefer.

## License

This project is licensed under the MIT License — see the LICENSE file for details.
