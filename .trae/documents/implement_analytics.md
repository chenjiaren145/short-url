# 短链接访问统计功能实现计划

为了进一步学习 Go 语言的接口设计、并发处理以及 Redis 的原子操作，下一步我们将为短链接服务添加访问统计功能。

## 目标
记录每个短链接被访问的次数，并提供接口查询统计数据。

## 详细步骤

### 1. 修改存储层接口 (Store Interface)
- **文件**: [store.go](file:///Users/chenjunwen/dev/temp/short-url/internal/store/store.go)
- **动作**:
  - 在 `Store` 接口中新增 `IncrementVisits(shortCode string) error` 方法，用于增加访问计数。
  - 在 `Store` 接口中新增 `GetVisits(shortCode string) (int64, error)` 方法，用于获取访问计数。

### 2. 实现内存存储 (MemoryStore)
- **文件**: [memory.go](file:///Users/chenjunwen/dev/temp/short-url/internal/store/memory.go)
- **动作**:
  - 在 `MemoryStore` 结构体中增加 `visits map[string]int64` 字段。
  - 修改 `NewMemoryStore` 初始化 `visits` map。
  - 实现 `IncrementVisits` 方法：使用写锁 (`Lock`) 安全地增加计数。
  - 实现 `GetVisits` 方法：使用读锁 (`RLock`) 安全地读取计数。

### 3. 实现 Redis 存储 (RedisStore)
- **文件**: [redis.go](file:///Users/chenjunwen/dev/temp/short-url/internal/store/redis.go)
- **动作**:
  - 实现 `IncrementVisits` 方法：使用 Redis 的 `INCR` 命令（原子操作），key 建议使用 `visits:{shortCode}` 格式。
  - 实现 `GetVisits` 方法：使用 Redis 的 `GET` 命令获取计数。

### 4. 更新业务逻辑层 (Service)
- **文件**: [service.go](file:///Users/chenjunwen/dev/temp/short-url/internal/service/service.go)
- **动作**:
  - 修改 `GetOriginalURL` 方法：在成功获取原始 URL 后，**异步**（使用 `go func()`）调用 `store.IncrementVisits`，确保不阻塞重定向请求。
  - 新增 `GetStats(shortCode string) (int64, error)` 方法，调用 `store.GetVisits`。
  - 在 `ShortenerService` 接口中添加 `GetStats` 定义。

### 5. 更新 API 处理层 (Handler)
- **文件**: [handler.go](file:///Users/chenjunwen/dev/temp/short-url/internal/api/handler.go)
- **动作**:
  - 新增 `GetStats(c *gin.Context)` 方法：
    - 从路径参数获取 `shortCode`。
    - 调用 `service.GetStats`。
    - 返回 JSON 格式的统计信息，例如 `{"short_code": "xyz", "visits": 10}`。

### 6. 注册路由 (Main)
- **文件**: [main.go](file:///Users/chenjunwen/dev/temp/short-url/cmd/short-url/main.go)
- **动作**:
  - 在 `main` 函数中注册新的路由：`r.GET("/:shortCode/stats", handler.GetStats)`。

## 验证计划
1. 启动服务（使用 Memory 或 Redis 模式）。
2. 创建一个短链接。
3. 访问该短链接多次（例如 3 次）。
4. 调用统计接口 `GET /:shortCode/stats`，确认返回的 `visits` 字段为 3。
