# godoc

一个用于从 Go 项目生成 OpenAPI/Swagger 文档的小工具（CLI），基于 `swaggo/swag`，并对生成结果做特定的后处理（如将 swagger 中的枚举字段转换为 oneOf/enum 的更友好形式）。

> 说明：本仓库包含一个可运行的 CLI 程序，入口位于 `cmd/godoc`。项目通过解析 `go.mod`（需要读取 `GOMODCACHE` 环境变量）来定位外部依赖模块的源代码，以便把外部模块中的类型/注释一同纳入文档生成。

## 特性

- 使用 `swaggo/swag` 生成 OpenAPI3 文档（输出 yaml）。
- 支持通过 `go.mod` 中的 `replace` 或 module cache 定位外部依赖的源代码并将其纳入扫描路径。
- 对生成的 `swagger.yaml` 做后处理：根据 `x-enum-varnames` / `x-enum-comments` 信息美化枚举描述并调整 `enum` 字段。

## 要求

- Go 1.20+（本项目使用了 `maps` 包，需要 1.20 及以上）。
- 已在 `go.mod` 中声明的依赖（仓库内已使用 `swaggo/swag`、`gopkg.in/yaml.v3`、`golang.org/x/mod` 等）。
- 环境变量 `GOMODCACHE` 必须可读（程序会使用它定位 module cache）。

## 快速开始

1. 克隆/进入仓库

2. 使用默认配置（当目录下没有 `.godoc.yaml` 时，程序会使用内置默认值）：

```bash
go run ./cmd/godoc
```

或构建二进制并运行：

```bash
go build -o godoc ./cmd/godoc
./godoc
```

3. 使用自定义配置文件（推荐，在多模块或需要自定义扫描路径时）：

```bash
# 指定配置文件路径
go run ./cmd/godoc -config .godoc.yaml
# 或
./godoc -config .godoc.yaml
```

生成的文档（yaml）默认输出到：

```
./assets/docs/swagger.yaml
```

## 配置（`.godoc.yaml`）

现在 `externalPaths`、`paths`、`mainFile`、`output` 都可以通过一个 YAML 配置文件来配置。默认配置文件名为 `.godoc.yaml`，也可以通过 `-config`（或 `--config`）命令行参数指定其他路径。

示例配置：

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

字段说明：
- `externalPaths`：map[string]string，键为 `module/path`，值为模块内的子路径（会通过 `go.mod` 的 `require`/`replace` 和 `GOMODCACHE` 来解析实际路径并加入扫描列表）。
- `paths`：要扫描的本地路径数组（相对于仓库根目录或绝对路径），会被 join 为 `SearchDir` 传给 `swag`。
- `mainFile`：`swag` 查找主入口的文件名（例如 `route.go`）。
- `output`：生成文档的输出目录（目录路径，`swag` 会在该目录写入 `swagger.yaml`）。

> 注意：当配置文件不存在时，程序会使用默认值（可参见 `internal.LoadConfig` 的默认配置）。

## 常用用法说明

- 主入口 `cmd/godoc/main.go` 中会读取配置文件（默认 `.godoc.yaml`），然后：
  - 将 `paths` 与通过 `externalPaths` 解析得到的外部模块路径整合为最终的扫描目录集合。
  - 使用 `swaggo/swag` 的 `format` 与 `gen` 生成 `swagger.yaml`（当前默认 `OutputTypes: ["yaml"]`）。

- 如果程序报错提示 `GOMODCACHE 环境变量读取失败`，请在运行前确保 `GOMODCACHE` 已设置，通常可以通过：

```bash
export GOMODCACHE=$(go env GOMODCACHE)
```

- 枚举后处理：程序运行完 `swag` 生成 `swagger.yaml` 后，会调用 `internal.ConvertEnum2OneOf` 对 `swagger.yaml` 做调整，主要依赖生成器输出的 `enum`、`x-enum-varnames` 与 `x-enum-comments` 字段来生成更友好的 description 并过滤特殊占位值（比如 `-`）。

## 代码结构（简述）

- `cmd/godoc/main.go` — CLI 入口，负责加载配置、组装扫描路径、调用 `swag` format/gen 并执行后处理。
- `internal/config.go` — 解析 `.godoc.yaml` 并提供默认值（如果没有配置文件则加载内置默认配置）。
- `module.go` — 负责读取 `go.mod`、解析 `require` 与 `replace`，并结合 `GOMODCACHE` 定位依赖代码路径。
- `enum.go` — 解析并转换 `swagger.yaml` 中的枚举相关字段，生成更友好的描述并调整 `enum` 内容。
- 其它文件（如 `cache.go`）用于项目内部的缓存/配置支持（具体实现请参阅源码）。

> 注：若要添加对更多输出格式（如 json、go）或更改 `swag` 的 `OutputTypes`，可以在 `cmd/godoc/main.go` 中修改 `gen.Config` 的 `OutputTypes`。

## 开发

- 本地快速构建：

```bash
go build ./...
```

- 运行单元测试（若有）：

```bash
go test ./...
```

- 建议在修改后运行 `gofmt`、`go vet` 来保持代码风格与质量。

## 故障排查

- 找不到模块路径或权限问题：
  - 错误信息示例：`模块路径不存在或无法访问`，请检查 `go env GOMODCACHE` 的值及 `go.mod` 中 `require`/`replace` 是否正确。
- 枚举转换失败：
  - 确认 `swagger.yaml` 中存在 `x-enum-varnames` 与 `x-enum-comments`，且它们的长度与 `enum` 一致。
- 如果 `swag` 生成失败，先单独运行 `swag`（或使用 `go get` 安装所需的 `swaggo` 工具），并检查 `SearchDir` / `MainFile` 设置是否正确。

## 贡献与联系方式

欢迎通过 Pull Request 或 Issue 的方式贡献代码。若希望我帮你把 README 翻译为英文、添加使用示例（实际项目的 route/controller 片段）或添加 CI 配置，请告诉我你希望的 License 与更多细节。

## 许可证

仓库中未指定许可证 —— 若需发布或共享，请在合并前明确选择合适的开源许可证（例如 MIT、Apache-2.0 等）。

