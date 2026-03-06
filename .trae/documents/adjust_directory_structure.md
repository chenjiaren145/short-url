# 目录结构调整计划 (修订版)

根据反馈，我们将重命名部分文件以避免命名冲突，并优化目录结构使其更清晰。

## 目标结构

```
.
├── cmd
│   └── short-url
│       └── main.go           # 应用程序入口
├── internal
│   ├── api                   # HTTP API 层 (原 internal/handler)
│   │   └── handler.go        # HTTP 处理器，定义路由处理函数
│   ├── shortener             # 核心业务逻辑/工具 (原 pkg/shortener.go)
│   │   └── generator.go      # 短链接生成算法 (重命名自 shortener.go)
│   ├── service               # 业务服务层
│   │   └── service.go        # 服务接口定义与实现 (重命名自 shortener.go)
│   └── store                 # 数据存储层
│       ├── memory.go
│       ├── redis.go
│       └── store.go
├── go.mod
├── go.sum
├── docker-compose.yml
└── README.md
```

## 执行步骤

### 1. 重构 `pkg` -> `internal/shortener`
- 创建 `internal/shortener` 目录。
- 将 `pkg/shortener.go` 移动并重命名为 `internal/shortener/generator.go`。
- **目的**: 将核心生成逻辑移入 `internal`，避免外部直接依赖，并重命名文件以反映其具体功能（生成器）。

### 2. 重构 `internal/service`
- 将 `internal/service/shortener.go` 重命名为 `internal/service/service.go`。
- **目的**: 避免与 `internal/shortener` 包名混淆，明确这是服务层的入口文件。

### 3. 重构 `internal/handler` -> `internal/api`
- 创建 `internal/api` 目录。
- 将 `internal/handler/handler.go` 移动到 `internal/api/handler.go`。
- 修改 `internal/api/handler.go` 的包名为 `package api`。
- **目的**: 遵循标准布局，将 API 处理逻辑集中在 `api` 目录。

### 4. 更新引用与依赖
- **`cmd/short-url/main.go`**:
    - 更新导入路径：`short-url/internal/handler` -> `short-url/internal/api`。
    - 更新代码引用：`handler.NewHandler` -> `api.NewHandler`。
- **`internal/service/service.go`**:
    - 确认导入路径指向 `short-url/internal/shortener` (无需更改代码，只需确保目录结构正确)。

### 5. 清理与验证
- 删除空的 `pkg` 和 `internal/handler` 目录。
- 运行 `go mod tidy`。
- 运行 `go build ./cmd/short-url` 验证构建。
