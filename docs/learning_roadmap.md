# 短链接服务 - 学习路线图

本项目是一个完整的短链接服务，包含 Go + Gin 后端和 React + Vite 前端。以下按学习顺序排列的练习项目，每个都有具体实现提示。

---

## 练习 1：添加 URL 验证

**学习目标**：Go 标准库 `net/url` 的使用，错误处理设计

### 需求

在 `POST /shorten` 之前验证 URL 格式，只接受 `http://` 和 `https://` 开头，且必须是合法 URL。

### 实现提示

#### 1. 新建验证文件

创建 `backend/internal/validator/validator.go`：

```go
package validator

import "net/url"

// ValidateURL 检查 URL 是否合法
// - 必须以 http:// 或 https:// 开头
// - 必须是有效的 URL（能被 url.Parse 解析）
func ValidateURL(rawURL string) error {
    // 提示：先用 strings.HasPrefix 检查协议
    // 再用 url.Parse 验证格式
    // 参考: net/url 包的 Parse 函数
    u, err := url.Parse(rawURL)
    if err != nil {
        return errors.New("invalid URL format")
    }
    if u.Scheme != "http" && u.Scheme != "https" {
        return errors.New("URL must use http or https scheme")
    }
    if u.Host == "" {
        return errors.New("URL must have a host")
    }
    return nil
}
```

#### 2. 在 Service 层集成

修改 `backend/internal/service/service.go` 的 `Shorten` 方法：

```go
import (
    "short-url/internal/shortener"
    "short-url/internal/store"
    "short-url/internal/validator"  // 新增
)

func (s *shortenerService) Shorten(originalURL string) (string, error) {
    // 新增：先验证 URL
    if err := validator.ValidateURL(originalURL); err != nil {
        return "", err  // 让上层返回 400 错误
    }

    shortCode, err := shortener.GenerateShortCode()
    if err != nil {
        return "", err
    }
    err = s.store.Save(shortCode, originalURL)
    if err != nil {
        return "", err
    }
    return shortCode, nil
}
```

#### 3. 在 Handler 层返回 400

修改 `backend/internal/api/handler.go` 的 `CreateShortURL`：

```go
// 提示：在调用 service.Shorten 之前不需要改
// 但需要区分 400（验证失败）和 500（服务器错误）
shortCode, err := h.service.Shorten(req.OriginalURL)
if err != nil {
    // 技巧：用 strings.Contains 或类型断言判断是否是验证错误
    // 简单的做法：所有非 nil 错误都返回 400，因为验证错误是最常见的
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
}
```

#### 4. 添加测试

创建 `backend/internal/validator/validator_test.go`：

```go
package validator

import "testing"

func TestValidateURL(t *testing.T) {
    tests := []struct {
        name    string
        url     string
        wantErr bool
    }{
        {"valid https", "https://example.com", false},
        {"valid http", "http://example.com", false},
        {"valid with path", "https://example.com/path/to/page", false},
        {"missing scheme", "example.com", true},
        {"ftp scheme", "ftp://example.com", true},
        {"empty string", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateURL(tt.url)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateURL() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### 验证方式

```bash
make test
# 或手动测试
curl -X POST http://localhost:8081/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "not-a-url"}'
# 期望返回 400

curl -X POST http://localhost:8081/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "ftp://example.com"}'
# 期望返回 400

curl -X POST http://localhost:8081/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://example.com"}'
# 期望返回 200
```

---

## 练习 2：添加删除接口

**学习目标**：接口的后向扩展、多层架构协调

### 需求

添加 `DELETE /:shortCode` 接口，允许删除已创建的短链接。

### 实现提示

#### 1. 扩展 Store 接口

修改 `backend/internal/store/store.go`：

```go
// 在 Store 接口中添加：
type Store interface {
    Save(shortCode string, originalURL string) error
    Load(shortCode string) (string, error)
    IncrementVisits(shortCode string) error
    GetVisits(shortCode string) (int64, error)
    Delete(shortCode string) error  // 新增这一行
}
```

#### 2. 在 MemoryStore 中实现

修改 `backend/internal/store/memory.go`：

```go
// 在 NewMemoryStore 后添加：
func (s *MemoryStore) Delete(shortCode string) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // 检查是否存在
    if _, ok := s.data[shortCode]; !ok {
        return ErrNotFound
    }

    // 删除两个 map 中的数据
    delete(s.data, shortCode)
    delete(s.visits, shortCode)
    return nil
}
```

#### 3. 在 RedisStore 中实现

修改 `backend/internal/store/redis.go`：

```go
// 新增：
func (s *RedisStore) Delete(shortCode string) error {
    ctx := context.Background()
    // 需要删除两个 key：shortCode 和 visits:shortCode
    // 提示：用 redis.NewClient().Del(ctx, keys...).Result()
    // 注意：即使 key 不存在，DEL 也会返回成功，这是 Redis 的正常行为
    // 但为了和 MemoryStore 保持一致，可以先检查是否存在
    exists, err := s.client.Exists(ctx, shortCode).Result()
    if err != nil {
        return err
    }
    if exists == 0 {
        return ErrNotFound
    }
    // DEL 返回删除的 key 数量
    _, err = s.client.Del(ctx, shortCode, "visits:"+shortCode).Result()
    return err
}
```

#### 4. 扩展 Service 层

修改 `backend/internal/service/service.go`：

```go
// 在 ShortenerService 接口中添加：
type ShortenerService interface {
    Shorten(originalURL string) (string, error)
    GetOriginalURL(shortCode string) (string, error)
    GetStats(shortCode string) (int64, error)
    Delete(shortCode string) error  // 新增
}

// 添加实现：
func (s *shortenerService) Delete(shortCode string) error {
    return s.store.Delete(shortCode)
}
```

#### 5. 添加 Handler

修改 `backend/internal/api/handler.go`：

```go
// 新增方法：
func (h *Handler) DeleteShortURL(c *gin.Context) {
    shortCode := c.Param("shortCode")

    err := h.service.Delete(shortCode)
    if err != nil {
        if err == store.ErrNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete short URL"})
        }
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Short URL deleted"})
}
```

#### 6. 注册路由

修改 `backend/cmd/short-url/main.go`，在 `r.GET("/:shortCode/stats", h.GetStats)` 后添加：

```go
r.DELETE("/:shortCode", h.DeleteShortURL)
```

### 验证方式

```bash
# 先创建一个短链接（记住返回的 shortCode）
curl -X POST http://localhost:8081/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://example.com"}'

# 删除它
curl -X DELETE http://localhost:8081/shorten/{shortCode}
# 期望返回 200

# 再次访问（应该 404）
curl -I http://localhost:8081/{shortCode}
# 期望返回 404
```

---

## 练习 3：修复短码碰撞重试

**学习目标**：重试模式、边界条件处理

### 需求

当前 `GenerateShortCode` 碰撞时直接返回错误（概率极低但存在），应该自动重试。

### 实现提示

#### 方案一：在 Generator 层重试

修改 `backend/internal/shortener/generator.go`：

```go
// 提示：
// 1. 生成器不知道 store 的存在，所以需要让 GenerateShortCode 接受一个检查函数
// 2. 或者在 service 层调用 generator 生成后，检查 store.Save 是否返回"已存在"错误
// 推荐方案二（更符合分层原则）

// 在 service.go 中：
const maxRetries = 3

func (s *shortenerService) Shorten(originalURL string) (string, error) {
    // 新增：最多重试 maxRetries 次
    var shortCode string
    var err error

    for i := 0; i < maxRetries; i++ {
        shortCode, err = shortener.GenerateShortCode()
        if err != nil {
            return "", err
        }

        // 尝试保存。如果碰撞了（可以定义一个 ErrExists 错误）就重试
        err = s.store.Save(shortCode, originalURL)
        if err == nil {
            return shortCode, nil
        }

        // 如果是其他错误，不重试，直接返回
        // 问题：store 接口目前没有区分"已存在"和"其他错误"
        // 这是设计问题，需要先扩展 store 接口（参见下面的进阶提示）
    }

    return "", errors.New("failed to generate unique short code after max retries")
}
```

#### 进阶：扩展 Store 接口支持碰撞检查

先修改 `backend/internal/store/store.go` 添加 `Exists` 方法：

```go
// 在 Store 接口中新增：
Exists(shortCode string) (bool, error)
```

在 `memory.go` 和 `redis.go` 中实现，然后在 `service.go` 中：

```go
for i := 0; i < maxRetries; i++ {
    shortCode, err = shortener.GenerateShortCode()
    if err != nil {
        return "", err
    }

    // 新增：先检查是否已存在
    exists, err := s.store.Exists(shortCode)
    if err != nil {
        return "", err
    }
    if exists {
        continue // 碰撞了，重试
    }

    err = s.store.Save(shortCode, originalURL)
    if err == nil {
        return shortCode, nil
    }
    return "", err // 保存失败（不是碰撞），不重试
}
```

### 验证方式

使用单元测试 mock 一个总是返回"已存在"的 store，验证重试逻辑。

---

## 练习 4：编写 Store 集成测试

**学习目标**：Go 单元测试、表驱动测试、测试隔离

### 需求

为 `MemoryStore` 和 `RedisStore` 各自写集成测试，覆盖 Save/Load/Delete/IncrementVisits/GetVisits。

### 实现提示

#### MemoryStore 测试

创建 `backend/internal/store/memory_test.go`：

```go
package store

import "testing"

// TestMemoryStore 是一个表驱动测试
// 表驱动测试是 Go 的最佳实践：把测试用例组织成表格，用一个循环执行
func TestMemoryStore(t *testing.T) {
    s := NewMemoryStore()

    t.Run("Load not found", func(t *testing.T) {
        _, err := s.Load("nonexistent")
        if err != ErrNotFound {
            t.Errorf("expected ErrNotFound, got %v", err)
        }
    })

    t.Run("Save and Load", func(t *testing.T) {
        err := s.Save("abc123", "https://example.com")
        if err != nil {
            t.Fatalf("Save failed: %v", err)
        }

        url, err := s.Load("abc123")
        if err != nil {
            t.Fatalf("Load failed: %v", err)
        }
        if url != "https://example.com" {
            t.Errorf("expected URL, got %s", url)
        }
    })

    t.Run("IncrementVisits", func(t *testing.T) {
        // 先保存一个 URL
        s.Save("visit1", "https://example.com")

        // 增加 3 次访问
        for i := 0; i < 3; i++ {
            if err := s.IncrementVisits("visit1"); err != nil {
                t.Fatalf("IncrementVisits failed: %v", err)
            }
        }

        // 检查访问次数
        visits, err := s.GetVisits("visit1")
        if err != nil {
            t.Fatalf("GetVisits failed: %v", err)
        }
        if visits != 3 {
            t.Errorf("expected 3 visits, got %d", visits)
        }
    })

    t.Run("Delete", func(t *testing.T) {
        s.Save("delete1", "https://example.com")
        err := s.Delete("delete1")
        if err != nil {
            t.Fatalf("Delete failed: %v", err)
        }

        _, err = s.Load("delete1")
        if err != ErrNotFound {
            t.Errorf("expected ErrNotFound after delete, got %v", err)
        }
    })

    // 并发测试（进阶）
    t.Run("Concurrent access", func(t *testing.T) {
        var wg sync.WaitGroup
        for i := 0; i < 100; i++ {
            wg.Add(1)
            go func() {
                defer wg.Done()
                s.Save("concurrent", "https://example.com")
                s.Load("concurrent")
                s.IncrementVisits("concurrent")
            }()
        }
        wg.Wait()

        // 验证没有 panic，且数据一致
        visits, _ := s.GetVisits("concurrent")
        if visits != 100 {
            t.Errorf("expected 100 visits, got %d", visits)
        }
    })
}
```

#### RedisStore 测试（需要 Redis 运行）

创建 `backend/internal/store/redis_test.go`：

```go
package store

import (
    "os"
    "testing"
)

// skipIfNoRedis 是 Go 的一种测试技巧
// 在测试开始前检查环境变量或 Redis 是否可用
// 如果不可用就跳过测试（SKIP），而不是失败
func skipIfNoRedis(t *testing.T) {
    if os.Getenv("REDIS_ADDR") == "" {
        t.Skip("skipping test: REDIS_ADDR not set")
    }
}

func TestRedisStore(t *testing.T) {
    skipIfNoRedis(t)

    addr := os.Getenv("REDIS_ADDR")
    password := os.Getenv("REDIS_PASSWORD")
    s := NewRedisStore(addr, password, 1) // 用 DB 1 避免污染生产数据

    // 每个测试前清空
    defer func() {
        ctx := context.Background()
        s.client.FlushDB(ctx)
    }()

    // 测试用例同上 MemoryStore...
    // （复用表驱动测试的思想，把测试逻辑提取到一个辅助函数中）
}
```

### 验证方式

```bash
# 只跑 memory store 测试（不需要 Redis）
go test ./internal/store/ -v -run TestMemoryStore

# 跑 Redis 测试
REDIS_ADDR=localhost:6379 go test ./internal/store/ -v -run TestRedisStore
```

---

## 练习 5：速率限制中间件

**学习目标**：Gin 中间件、滑动窗口算法、context 传递

### 需求

为所有 API 添加速率限制��每个 IP 每分钟最多 60 次请求。

### 实现提示

#### 1. 创建限流中间件

创建 `backend/internal/middleware/ratelimit.go`：

```go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

// slidingWindowRateLimiter 实现滑动窗口限流
type slidingWindowRateLimiter struct {
    mu       sync.Mutex
    requests map[string][]time.Time
    limit    int
    window   time.Duration
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *slidingWindowRateLimiter {
    return &slidingWindowRateLimiter{
        requests: make(map[string][]time.Time),
        limit:    limit,
        window:   window,
    }
}

// Allow 检查是否允许请求
func (l *slidingWindowRateLimiter) Allow(ip string) bool {
    l.mu.Lock()
    defer l.mu.Unlock()

    now := time.Now()
    windowStart := now.Add(-l.window)

    // 清理过期的请求记录
    // 提示：过滤掉 windowStart 之前的记录
    var valid []time.Time
    for _, t := range l.requests[ip] {
        if t.After(windowStart) {
            valid = append(valid, t)
        }
    }

    if len(valid) >= l.limit {
        l.requests[ip] = valid
        return false
    }

    l.requests[ip] = append(valid, now)
    return true
}

// RateLimit 返回一个 Gin 中间件
func RateLimit(limiter *slidingWindowRateLimiter) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()

        if !limiter.Allow(ip) {
            // 提示：用 c.AbortWithStatusJSON 返回 429 Too Many Requests
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
                "error": "rate limit exceeded",
            })
            return
        }

        c.Next()
    }
}
```

#### 2. 注册中间件

修改 `backend/cmd/short-url/main.go`：

```go
import (
    // ... existing imports
    "short-url/internal/middleware"
)

func main() {
    // ... existing setup

    // 创建限流器：每分钟 60 次
    limiter := middleware.NewSlidingWindowLimiter(60, time.Minute)

    // 注册中间件：对所有路由生效
    r.Use(middleware.RateLimit(limiter))

    // ... rest of the code
}
```

### 验证方式

```bash
# 用 Apache Bench 或 Hey 测试
brew install ab  # macOS
ab -n 100 -c 10 http://localhost:8081/healthz
# 观察：前 60 个成功，后面的返回 429
```

---

## 练习 6：访问分析功能（Phase 1）

**学习目标**：数据模型演进、API 扩展、前后端协作

这是 [analytics_feature_plan.md](./analytics_feature_plan.md) 中规划的功能的 Phase 1 实现。

### 需求

记录每次访问的详细信息（时间戳、Referer），并提供分析接口。

### 实现提示

#### 1. 新建访问记录模型

创建 `backend/internal/models/visit.go`：

```go
package models

import "time"

// Visit 记录一次访问
type Visit struct {
    Timestamp time.Time `json:"timestamp"`
    Referer   string    `json:"referer"`
    UserAgent string    `json:"user_agent"`
    IP        string    `json:"ip"`
}
```

#### 2. 扩展 Store 接口

修改 `backend/internal/store/store.go`，添加：

```go
type Store interface {
    // ... existing methods

    // RecordVisit 记录一次访问
    RecordVisit(shortCode string, visit Visit) error

    // GetVisits 获取访问记录列表
    GetVisits(shortCode string) ([]Visit, error)
}
```

#### 3. 在 MemoryStore 中实现

修改 `backend/internal/store/memory.go`：

```go
type MemoryStore struct {
    mu     sync.RWMutex
    data   map[string]string
    visits map[string]int64
    records map[string][]Visit  // 新增：存储访问记录
}

// 提示：在 NewMemoryStore 中初始化 records map
func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        data:    make(map[string]string),
        visits:  make(map[string]int64),
        records: make(map[string][]Visit),  // 新增
    }
}

func (s *MemoryStore) RecordVisit(shortCode string, visit Visit) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    s.records[shortCode] = append(s.records[shortCode], visit)
    s.visits[shortCode]++
    return nil
}

func (s *MemoryStore) GetVisits(shortCode string) ([]Visit, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return s.records[shortCode], nil
}
```

#### 4. 在 RedisStore 中实现

修改 `backend/internal/store/redis.go`：

```go
// 提示：Redis 中可以用 LIST 存储访问记录
// RPUSH visits:{shortCode} '{"timestamp":"...","referer":"..."}'
// LRANGE visits:{shortCode} 0 -1

func (s *RedisStore) RecordVisit(shortCode string, visit Visit) error {
    ctx := context.Background()

    data, err := json.Marshal(visit)
    if err != nil {
        return err
    }

    pipe := s.client.Pipeline()
    pipe.RPush(ctx, "visits:"+shortCode, data)
    pipe.Incr(ctx, "count:"+shortCode) // 用单独的计数 key
    _, err = pipe.Exec(ctx)
    return err
}
```

#### 5. 扩展 Service 层

修改 `backend/internal/service/service.go`：

```go
// 新增方法
func (s *shortenerService) RecordVisit(shortCode string, req *VisitRequest) error {
    visit := Visit{
        Timestamp: time.Now(),
        Referer:   req.Referer,
        UserAgent: req.UserAgent,
        IP:        req.IP,
    }
    return s.store.RecordVisit(shortCode, visit)
}

func (s *shortenerService) GetVisitRecords(shortCode string) ([]Visit, error) {
    return s.store.GetVisits(shortCode)
}
```

#### 6. 新增 Handler

修改 `backend/internal/api/handler.go`：

```go
type VisitRequest struct {
    Referer   string
    UserAgent string
    IP        string
}

// 在 Redirect 方法中调用 RecordVisit：
func (h *Handler) Redirect(c *gin.Context) {
    shortCode := c.Param("shortCode")

    // 记录访问
    h.service.RecordVisit(shortCode, &VisitRequest{
        Referer:   c.GetHeader("Referer"),
        UserAgent: c.GetHeader("User-Agent"),
        IP:        c.ClientIP(),
    })

    originalURL, err := h.service.GetOriginalURL(shortCode)
    // ... rest
}

// 新增分析接口：
func (h *Handler) GetAnalytics(c *gin.Context) {
    shortCode := c.Param("shortCode")

    visits, err := h.service.GetVisitRecords(shortCode)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get visits"})
        return
    }

    // 返回结构化数据
    c.JSON(http.StatusOK, gin.H{
        "total":    len(visits),
        "visits":   visits,
    })
}
```

#### 7. 注册路由

在 `main.go` 中添加：

```go
r.GET("/:shortCode/analytics", h.GetAnalytics)
```

### 验证方式

```bash
# 模拟访问
curl -I -H "Referer: https://google.com" http://localhost:8081/abc123

# 查看分析
curl http://localhost:8081/abc123/analytics
```

---

## 练习 7：QR Code 生成

**学习目标**：图片处理、API 响应二进制数据

### 需求

添加 `GET /:shortCode/qrcode` 接口，返回 QR Code 图片。

### 实现提示

#### 方案一：后端生成（推荐）

安装 QR 库：

```bash
go get github.com/skip2/go-qrcode
```

创建 `backend/internal/api/qrcode.go`：

```go
package api

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/skip2/go-qrcode"
)

// GenerateQRCode 生成 QR Code
func (h *Handler) GenerateQRCode(c *gin.Context) {
    shortCode := c.Param("shortCode")
    shortURL := h.shortURL(c, shortCode)

    // 生成 PNG
    png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
        return
    }

    // 直接返回图片二进制，设置正确的 Content-Type
    c.Data(http.StatusOK, "image/png", png)
}
```

在 `main.go` 中注册：

```go
r.GET("/:shortCode/qrcode", h.GenerateQRCode)
```

#### 方案二：前端生成

安装 `qrcode.react`：

```bash
cd frontend
npm install qrcode.react
```

在 `URLList.tsx` 中：

```tsx
import { QRCodeSVG } from 'qrcode.react'

// 在 URLList 中添加 QR 展示：
<div className="flex items-center gap-4">
    <QRCodeSVG value={url.shortUrl} size={100} />
</div>
```

### 验证方式

```bash
# 后端方案
curl -o qr.png http://localhost:8081/abc123/qrcode
open qr.png  # macOS 打开图片查看

# 或浏览器直接访问 http://localhost:8081/abc123/qrcode
```

---

## 练习 8：前端测试

**学习目标**：Vitest + React Testing Library、组件测试

### 需求

为 URLForm 和 URLList 组件添加单元测试。

### 实现提示

#### 1. 安装测试依赖

```bash
cd frontend
npm install -D vitest @testing-library/react @testing-library/user-event jsdom
```

配置 `vite.config.js`：

```js
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    environment: 'jsdom',
    globals: true,
  },
})
```

#### 2. 测试 URLForm

创建 `frontend/src/components/URLForm.test.tsx`：

```tsx
import { render, screen, fireEvent } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import URLForm from './URLForm'

// describe 是 Vitest 的测试分组语法
describe('URLForm', () => {
    it('renders input and button', () => {
        render(<URLForm onSubmit={vi.fn()} loading={false} />)
        expect(screen.getByRole('textbox')).toBeInTheDocument()
        expect(screen.getByRole('button', { name: '生成短链接' })).toBeInTheDocument()
    })

    it('calls onSubmit with URL when submitted', async () => {
        const user = userEvent.setup()
        const onSubmit = vi.fn()
        render(<URLForm onSubmit={onSubmit} loading={false} />)

        await user.type(screen.getByRole('textbox'), 'https://example.com')
        await user.click(screen.getByRole('button', { name: '生成短链接' }))

        expect(onSubmit).toHaveBeenCalledWith('https://example.com')
    })

    it('disables button when loading', () => {
        render(<URLForm onSubmit={vi.fn()} loading={true} />)
        expect(screen.getByRole('button', { name: '生成短链接' })).toBeDisabled()
    })
})
```

#### 3. 测试 URLList

创建 `frontend/src/components/URLList.test.tsx`：

```tsx
import { render, screen } from '@testing-library/react'
import URLList from './URLList'
import type { URLItem } from '../types/url'

const mockUrls: URLItem[] = [
    {
        shortCode: 'abc123',
        shortUrl: 'http://localhost:8081/abc123',
        originalUrl: 'https://example.com',
        visits: 42,
        createdAt: '2024-01-01T00:00:00Z',
    },
]

describe('URLList', () => {
    it('shows empty state when no URLs', () => {
        render(<URLList urls={[]} onRefreshStats={vi.fn()} />)
        expect(screen.getByText('还没有创建任何短链接')).toBeInTheDocument()
    })

    it('renders URL items', () => {
        render(<URLList urls={mockUrls} onRefreshStats={vi.fn()} />)
        expect(screen.getByText('http://localhost:8081/abc123')).toBeInTheDocument()
        expect(screen.getByText('42')).toBeInTheDocument()
    })
})
```

### 验证方式

```bash
cd frontend
npm test
# 或监听模式
npm test -- --watch
```

---

## 练习 9：结构化日志 + Metrics

**学习目标**：Go `slog`、Prometheus 指标、 observability

### 需求

用 `slog` 替代 `log`，输出 JSON 格式日志，并暴露 Prometheus 指标。

### 实现提示

#### 1. 使用 slog

修改 `backend/cmd/short-url/main.go`：

```go
import (
    "log/slog"
    "os"
)

func main() {
    // 使用 JSON 格式的结构化日志
    slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

    slog.Info("server starting", "port", *port, "store", *storeType)

    // 替换所有 log.Printf 为 slog.Info / slog.Error
    // 例如：
    // log.Printf("Using Redis storage at %s", *redisAddr)
    // 改为：
    // slog.Info("using redis storage", "addr", *redisAddr)
}
```

#### 2. 添加请求日志中间件

创建 `backend/internal/middleware/logger.go`：

```go
package middleware

import (
    "log/slog"
    "time"

    "github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path

        c.Next()

        slog.Info("request",
            "method", c.Request.Method,
            "path", path,
            "status", c.Writer.Status(),
            "latency", time.Since(start),
            "ip", c.ClientIP(),
        )
    }
}
```

#### 3. 添加 Prometheus 指标

安装 Prometheus 库：

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promhttp
```

创建 `backend/internal/metrics/metrics.go`：

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    HTTPRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "path", "status"},
    )

    HTTPRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "path"},
    )
)
```

在 `main.go` 中暴露 `/metrics` 端：

```go
import (
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

r.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### 验证方式

```bash
# 查看日志输出
make run-memory
# 观察 JSON 格式的日志

# 查看 metrics
curl http://localhost:8081/metrics | grep http_requests_total
```

---

## 练习 10：Graceful Shutdown

**学习目标**：信号处理、HTTP Server 生命周期

### 需求

服务收到 SIGTERM/SIGINT 时，优雅地关闭：停止接收新请求，等待现有请求处理完毕。

### 实现提示

修改 `backend/cmd/short-url/main.go`：

```go
import (
    "context"
    "flag"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    // ... 现有设置不变 ...

    // 创建 HTTP Server（替代直接用 r.Run）
    srv := &http.Server{
        Addr:    ":" + *port,
        Handler: r,
    }

    // 在 goroutine 中启动服务器
    go func() {
        slog.Info("server starting", "addr", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            slog.Error("server error", "error", err)
            os.Exit(1)
        }
    }()

    // 等待中断信号
    quit := make(chan os.Signal, 1)
    // 通知进程收到 SIGTERM 或 SIGINT 信号
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    slog.Info("shutting down server...")

    // 创建有超时的 context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // 优雅关闭：停止接收新请求，等待现有请求完成
    if err := srv.Shutdown(ctx); err != nil {
        slog.Error("server forced to shutdown", "error", err)
    }

    slog.Info("server stopped")
}
```

### 验证方式

```bash
make run-memory
# 在另一个终端
kill -TERM $(pgrep short-url)
# 观察：服务器收到信号后等待，没有立即退出
```

---

## 练习 11：TTL 过期机制

**学习目标**：Redis TTL、后台定时清理

### 需求

支持创建带过期时间的短链接，过期后自动删除。

### 实现提示

#### 1. 修改请求结构

在 `backend/internal/api/handler.go` 中：

```go
type CreateShortURLRequest struct {
    OriginalURL string `json:"original_url" binding:"required"`
    TTLSeconds  int    `json:"ttl_seconds"` // 可选，0 表示永不过期
}
```

#### 2. 修改 Store 接口

在 `store.go` 中新增 `SaveWithExpiry` 方法：

```go
SaveWithExpiry(shortCode string, originalURL string, ttl time.Duration) error
```

#### 3. Redis 实现（天然支持）

```go
func (s *RedisStore) SaveWithExpiry(shortCode string, originalURL string, ttl time.Duration) error {
    ctx := context.Background()
    return s.client.Set(ctx, shortCode, originalURL, ttl).Err()
}
```

#### 4. MemoryStore 实现（需要定时清理）

创建 `backend/internal/store/expiry.go`：

```go
package store

// StartCleanup 启动后台清理 goroutine
// 定期删除过期的条目（适用于 MemoryStore）
func (s *MemoryStore) StartCleanup(interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        for range ticker.C {
            s.mu.Lock()
            // 清理逻辑：如果存储了过期时间
            for key, expiry := range s.expiry {
                if time.Now().After(expiry) {
                    delete(s.data, key)
                    delete(s.visits, key)
                    delete(s.expiry, key)
                }
            }
            s.mu.Unlock()
        }
    }()
}
```

### 验证方式

```bash
# 创建带 10 秒过期的链接
curl -X POST http://localhost:8081/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://example.com", "ttl_seconds": 10}'

# 10 秒后访问，应该返回 404
```

---

## 文档目录

- [练习 1：URL 验证](#练习-1添加-url-验证)
- [练习 2：删除接口](#练习-2添加删除接口)
- [练习 3：碰撞重试](#练习-3修复短码碰撞重试)
- [练习 4：Store 集成测试](#练习-4编写-store-集成测试)
- [练习 5：速率限制](#练习-5速率限制中间件)
- [练习 6：访问分析 Phase 1](#练习-6访问分析功能phase-1)
- [练习 7：QR Code](#练习-7qr-code-生成)
- [练习 8：前端测试](#练习-8前端测试)
- [练习 9：结构化日志 + Metrics](#练习-9结构化日志--metrics)
- [练习 10：Graceful Shutdown](#练习-10graceful-shutdown)
- [练习 11：TTL 过期机制](#练习-11ttl-过期机制)
