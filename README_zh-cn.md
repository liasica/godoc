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

## 安装

### 方式一：使用安装脚本（推荐）

使用提供的 `install.sh` 脚本自动下载并安装最新版本：

```bash
# 下载并运行安装脚本
curl -fsSL https://raw.githubusercontent.com/liasica/godoc/master/install.sh | bash
```

或者克隆仓库后本地运行：

```bash
# 克隆仓库
git clone https://github.com/liasica/godoc.git
cd godoc

# 运行安装脚本
bash install.sh
```

安装脚本会：
- 自动检测您的操作系统和架构（支持 Linux、macOS、Windows 的 amd64 和 arm64 架构）
- 从 GitHub Releases 下载对应平台的预编译二进制文件
- 将二进制文件安装到 `$GOPATH/bin` 目录
- 检查已安装版本，如果有新版本则自动升级
- 验证安装是否成功

> 注意：安装脚本需要 Go 环境（用于获取 `GOPATH`）和 `curl` 命令。安装完成后，请确保 `$GOPATH/bin` 已添加到您的 `PATH` 环境变量中。

### 方式二：从源码构建

如果您想从源码构建，可以按照以下步骤：

```bash
# 克隆仓库
git clone https://github.com/liasica/godoc.git
cd godoc

# 构建二进制文件
go build -o godoc ./cmd/godoc

# （可选）将二进制文件移动到 PATH 中的目录
mv godoc $GOPATH/bin/
```

## 快速开始

1. 使用 `config init` 命令创建默认配置文件：

```bash
# 在当前目录创建 .godoc.yaml
godoc config init
```

2. 根据您的项目需求编辑 `.godoc.yaml` 配置文件。

3. 运行 `generate` 子命令生成文档：

```bash
# 使用配置文件运行生成
godoc generate -c .godoc.yaml
```

或者直接使用 `go run`（适用于开发调试）：

```bash
# 使用显式配置路径运行 generate
go run ./cmd/godoc generate -c .godoc.yaml

# 如果需要自动下载依赖，可以使用 -mod=mod
go run -mod=mod ./cmd/godoc generate -c .godoc.yaml
```

> 注意：`-mod=mod` 标志会告诉 Go 在需要时自动下载和更新依赖。这在首次运行工具或依赖缺失时很有用。

生成的文档（YAML）默认写入：

```
./assets/docs/swagger.yaml
```

## 配置（`.godoc.yaml`）

`externalPaths`、`paths`、`mainFile`、`output` 都可以通过 YAML 配置文件配置。默认配置文件名为 `.godoc.yaml`；也可以使用 `-config`（或 `--config`）命令行标志指定其他路径。

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
- `externalPaths`：map[string]string，键为 `module/path`，值为模块内的子路径（工具会通过 `go.mod` 的 `require`/`replace` 和 `GOMODCACHE` 来解析实际路径并加入扫描列表）。
- `paths`：要扫描的本地路径数组（相对于仓库根目录或绝对路径），会被 join 为 `SearchDir` 传给 `swag`。
- `mainFile`：`swag` 查找主入口的文件名（例如 `route.go`）。
- `output`：生成文档的输出目录（目录路径，`swag` 会在该目录写入 `swagger.yaml`）。

> 注意：当配置文件不存在时，程序会使用默认值（参见 `internal.LoadConfig` 的默认配置）。

## 使用说明

- CLI 主入口 `cmd/godoc/main.go` 现在暴露了两个子命令：
  - `generate` — 生成文档。此命令接受可选的 `-config`（或 `--config`）标志，指向 YAML 配置文件；如果省略，则工具会使用内置默认值或当前目录下的 `.godoc.yaml`（若存在）。
  - `config` — 配置相关工具。使用 `godoc config init` 在当前目录创建默认 `.godoc.yaml`。

- 常用命令及标志：
  - `godoc generate` — 使用默认配置或当前目录的 `.godoc.yaml` 运行生成。
  - `godoc generate -c path/to/config.yaml` — 使用显式配置文件运行生成。
  - `godoc config init` — 在当前目录创建默认 `.godoc.yaml`。
  - `godoc --version` — 打印 CLI 版本（版本由 `version.go` 提供）。
  - `godoc help` 或 `godoc <command> --help` — 显示命令帮助与标志说明。

- 典型生成流程：
  1. 创建或编辑 `.godoc.yaml`（也可以使用 `godoc config init` 创建默认配置）。
  2. 运行 `godoc generate`（或 `godoc generate -c .godoc.yaml` / `go run ./cmd/godoc generate -c .godoc.yaml`）。
  3. 工具会解析外部模块路径（使用 `go.mod` 和 `GOMODCACHE`），运行 `swaggo/swag` 的 format 与 gen，然后对生成的 `swagger.yaml` 做后处理。

- 如果程序报错提示 `failed to read GOMODCACHE environment variable`，请确保运行前已设置 `GOMODCACHE`，例如：

```bash
export GOMODCACHE=$(go env GOMODCACHE)
```

- 枚举后处理：`swag` 生成 `swagger.yaml` 后，程序会调用 `internal.ConvertEnum2OneOf` 对文件进行调整。此步骤依赖生成器输出的 `enum`、`x-enum-varnames` 和 `x-enum-comments` 字段，用于生成更友好的 description 并过滤特殊占位值（例如 `-`）。

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

本项目使用 MIT 许可证 —— 详情请参阅仓库根目录中的 LICENSE 文件。
