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
│   └── main.go           # 程序入口 (Entry Point)。负责初始化组件、配置路由和启动服务。
├── internal              # 内部代码包。Go 编译器强制限制，外部项目无法直接导入 internal 下的包。
│   ├── api
│   │   └── handler.go    # HTTP 处理器 (Handlers)。负责解析请求、调用业务逻辑和返回响应。
│   ├── shortener
│   │   └── shortener.go  # 核心业务逻辑。负责生成短链接代码（例如 Base62 编码）。
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

- Go 1.20 或更高版本

### 运行项目 (Running the Application)

1. **下载依赖**:
   ```bash
   go mod tidy
   ```

2. **启动服务**:
   ```bash
   go run cmd/main.go
   ```
   服务将启动在 `8081` 端口。

### API 使用示例 (API Usage)

#### 1. 创建短链接 (Create Short URL)

**接口**: `POST /shorten`

**请求体**:
```json
{
  "original_url": "https://www.google.com"
}
```

**命令行测试**:
```bash
curl -X POST -H "Content-Type: application/json" -d '{"original_url": "https://www.google.com"}' http://localhost:8081/shorten
```

**响应**:
```json
{
  "original_url": "https://www.google.com",
  "short_url": "http://localhost:8081/AbCdEf"
}
```

#### 2. 访问短链接 (Access Short URL)

**接口**: `GET /:shortCode`

**示例**:
在浏览器中访问上一步生成的 `short_url`，或者使用 `curl` 查看重定向头：
```bash
curl -I http://localhost:8081/AbCdEf
```

**响应**:
服务器将返回 `302 Found` 状态码，并在 `Location` 头中包含原始 URL。

## 🛠 进阶练习 (Advanced Exercises)

如果你想进一步提升，可以尝试以下改进：

1.  **持久化存储**: 实现一个基于 Redis 或 MySQL 的 `Store`，替换掉 `MemoryStore`。
2.  **配置管理**: 使用 `viper` 或 `godotenv` 从配置文件或环境变量中读取端口号和数据库连接串。
3.  **结构化日志**: 使用 `zap` 或 `logrus` 替换标准的 `log` 库，实现更规范的日志记录。
4.  **单元测试**: 为 `shortener` 和 `api` 包编写单元测试 (`go test`)。

## License

MIT
