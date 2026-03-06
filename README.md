# Short URL Service (短链接服务)

这是一个使用 Go 语言和 [Gin](https://github.com/gin-gonic/gin) 框架构建的简易短链接生成服务。
本项目旨在作为学习 Go Web 开发的入门示例，涵盖了项目结构设计、依赖注入、接口抽象、并发安全等核心概念。

## 📚 学习指南 (Learning Guide)

本项目展示了以下 Web 开发的关键方面：

### 1. 项目结构 (Project Structure)
遵循 Go 社区推荐的 [Standard Go Project Layout](https://github.com/golang-standards/project-layout)：

```
.
├── cmd
│   └── short-url
│       └── main.go           # 程序入口 (Entry Point)。负责初始化组件、配置路由和启动服务。
├── internal              # 内部代码包。Go 编译器强制限制，外部项目无法直接导入 internal 下的包。
│   ├── api
│   │   └── handler.go    # HTTP 处理器 (Handlers)。负责解析请求、调用业务逻辑和返回响应。
│   ├── shortener
│   │   └── generator.go  # 核心业务逻辑。负责生成短链接代码（例如 Base62 编码）。
│   ├── service
│   │   └── service.go    # 业务服务层。负责编排存储和生成逻辑。
│   └── store
│       ├── store.go      # 接口定义 (Interface)。定义了存储层的行为规范。
│       └── memory.go     # 接口实现 (Implementation)。基于内存的具体存储实现。
├── go.mod                # 依赖管理文件。
└── README.md             # 项目文档。
```

### 2. 核心概念 (Core Concepts)

*   **依赖注入 (Dependency Injection)**:
    在 `cmd/main.go` 中，我们创建了 `store` 实例，并将其传递给 `api.NewHandler`。这种方式使得 `Handler` 不直接依赖具体的存储实现（如 Redis 或 MySQL），而是依赖 `store.Store` 接口。这提高了代码的可测试性和灵活性。

*   **接口抽象 (Interface Abstraction)**:
    `internal/store/store.go` 定义了 `Store` 接口。任何实现了 `Save` 和 `Load` 方法的结构体都可以作为存储层。这体现了面向接口编程的思想。

*   **并发安全 (Concurrency Safety)**:
    由于 Web 服务器是并发处理请求的，`internal/store/memory.go` 使用了 `sync.RWMutex`（读写锁）来保护共享的 `map`，防止多个 goroutine 同时读写导致的数据竞争 (Data Race)。

*   **Web 框架 (Web Framework)**:
    使用 **Gin** 框架处理 HTTP 请求。学习点包括：
    *   路由注册 (`r.POST`, `r.GET`)
    *   请求参数绑定 (`ShouldBindJSON`)
    *   JSON 响应 (`c.JSON`)
    *   路径参数获取 (`c.Param`)
    *   HTTP 状态码的使用 (`http.StatusOK`, `http.StatusNotFound` 等)

## 🚀 快速开始 (Getting Started)

### 前置要求 (Prerequisites)

- Go 1.20 或更高版本 (本地运行时需要)
- Docker 和 Docker Compose (Docker 运行时需要)

### 运行项目 (Running the Application)

你可以选择**本地内存模式**（快速开发）或**Docker 完整模式**（模拟生产）。

#### 方式一：本地内存模式 (无需依赖)

直接编译并运行 Go 程序，使用内存作为临时存储。重启后数据丢失。

1. **下载依赖**:
   ```bash
   go mod tidy
   ```

2. **启动服务**:
   ```bash
   go run cmd/short-url/main.go
   ```
   服务将启动在 `8081` 端口。

#### 方式二：Docker 完整模式 (App + Redis)

构建并启动所有服务（App + Redis），数据持久化存储在 Docker Volume 中。

1. **一键启动**:
   ```bash
   docker-compose up --build
   ```
   **命令参数解释**:
   *   `up`: 读取 `docker-compose.yml` 配置，创建网络、卷，并启动所有服务容器。
   *   `--build`: 启动前强制重新构建镜像。如果不加此参数，Docker 会尝试使用已存在的镜像，导致最新的代码修改无法生效。

   **常用维护命令**:
   *   `docker-compose down`: 停止并移除容器、网络（清理环境）。
   *   `docker-compose logs -f`: 查看所有服务的实时日志输出。

### API 验证 (Verification)

服务启动后，可以使用 `curl` 命令进行验证。

#### 1. 创建短链接 (Create Short URL)

**命令**:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"original_url": "https://www.google.com"}' http://localhost:8081/shorten
```

**命令解释**:
*   `-X POST`: 指定 HTTP 请求方法为 `POST`。
*   `-H "Content-Type: application/json"`: 设置请求头，告知服务器发送的数据格式为 JSON。
*   `-d '{...}'`: 发送的请求体数据（JSON 格式），包含需要缩短的 `original_url`。

**响应示例**:
```json
{
  "original_url": "https://www.google.com",
  "short_url": "http://localhost:8081/AbCdEf"
}
```

#### 2. 访问短链接 (Access Short URL)

**命令**:
```bash
# 请将 <short_code> 替换为上一步生成的代码 (例如 AbCdEf)
curl -I http://localhost:8081/<short_code>
```

**命令解释**:
*   `-I` (或 `--head`): 仅获取 HTTP 响应头，不下载响应体。这对于检查重定向（302 跳转）非常有用。
*   你将在输出中看到 `HTTP/1.1 302 Found` 和 `Location: https://www.google.com`。

## 🛠 进阶练习 (Advanced Exercises)

如果你想进一步提升，可以尝试以下改进：

1.  **持久化存储**: 实现一个基于 Redis 或 MySQL 的 `Store`，替换掉 `MemoryStore`。
2.  **配置管理**: 使用 `viper` 或 `godotenv` 从配置文件或环境变量中读取端口号和数据库连接串。
3.  **结构化日志**: 使用 `zap` 或 `logrus` 替换标准的 `log` 库，实现更规范的日志记录。
4.  **单元测试**: 为 `shortener` 和 `api` 包编写单元测试 (`go test`)。

## License

MIT
