# 后端升级计划 - 公开短链接服务

## 概述

配合前端升级，将短链接服务从内部工具升级为公开服务，核心改造：用户认证系统 + 数据库持久化 + 链接管理 API。

---

## 优先级说明

- **P0**: 必须实现，核心功能
- **P1**: 重要功能，用户体验关键
- **P2**: 增强功能，锦上添花
- **P3**: 优化改进，可后续迭代

---

## 改造计划

### Phase 1: 数据库层 (P0)

| 序号 | 任务 | 涉及文件 | 学习要点 |
|------|------|----------|----------|
| 1.1 | 选择数据库 | - | PostgreSQL vs MySQL 选型 |
| 1.2 | 添加数据库驱动 | `go.mod` | `github.com/lib/pq` 或 `github.com/go-sql-driver/mysql` |
| 1.3 | 创建 users 表 | `backend/internal/store/schema.sql` | 用户表设计（id, email, password_hash, created_at） |
| 1.4 | 创建 urls 表 | `backend/internal/store/schema.sql` | 链接表设计（id, short_code, original_url, user_id, visits, created_at） |
| 1.5 | 实现 DBStore | `backend/internal/store/db.go` | 数据库存储实现，实现 Store 接口 |
| 1.6 | 数据库连接池 | `backend/cmd/short-url/main.go` | `sql.Open` + 连接池配置 |
| 1.7 | 数据库迁移方案 | `backend/internal/store/migrate.go` | 简单迁移或使用 golang-migrate |

#### 用户表设计

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

#### 链接表设计

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    visits INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_urls_short_code ON urls(short_code);
CREATE INDEX idx_urls_user_id ON urls(user_id);
```

### Phase 2: 用户认证系统 (P0)

| 序号 | 任务 | 涉及文件 | 学习要点 |
|------|------|----------|----------|
| 2.1 | 密码哈希工具 | `backend/internal/auth/password.go` | `golang.org/x/crypto/bcrypt` |
| 2.2 | JWT 工具 | `backend/internal/auth/jwt.go` | `github.com/golang-jwt/jwt/v5` |
| 2.3 | 用户模型 | `backend/internal/models/user.go` | User struct 定义 |
| 2.4 | 用户 Store 接口 | `backend/internal/store/store.go` | 新增 UserStore 接口 |
| 2.5 | 用户 Store 实现 | `backend/internal/store/user_store.go` | CreateUser, GetUserByEmail, GetUserByID |
| 2.6 | 认证 Service | `backend/internal/service/auth.go` | AuthService 接口 + 实现 |
| 2.7 | 认证 Handler | `backend/internal/api/auth_handler.go` | Register, Login, Me 接口 |
| 2.8 | 认证中间件 | `backend/internal/middleware/auth.go` | JWT 验证，注入 user_id 到 context |
| 2.9 | 配置管理 | `backend/cmd/short-url/main.go` | JWT secret, token 过期时间配置 |

#### 密码哈希示例

```go
// backend/internal/auth/password.go
package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

#### JWT 工具示例

```go
// backend/internal/auth/jwt.go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateToken(userID int, secret string) (string, error) {
    claims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}

func ParseToken(tokenString, secret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil {
        return nil, err
    }
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    return nil, jwt.ErrSignatureInvalid
}
```

#### 认证中间件示例

```go
// backend/internal/middleware/auth.go
package middleware

import (
    "net/http"
    "strings"
    "short-url/internal/auth"
    "github.com/gin-gonic/gin"
)

func Auth(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
            return
        }

        claims, err := auth.ParseToken(parts[1], jwtSecret)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }

        c.Set("user_id", claims.UserID)
        c.Next()
    }
}
```

### Phase 3: 链接管理功能 (P1)

| 序号 | 任务 | 涉及文件 | 学习要点 |
|------|------|----------|----------|
| 3.1 | 扩展 URL 模型 | `backend/internal/models/url.go` | URL struct 添加 UserID 字段 |
| 3.2 | 扩展 Store 接口 | `backend/internal/store/store.go` | 新增 GetByUserID, GetByID, Update, Delete |
| 3.3 | 实现 Store 方法 | `backend/internal/store/db.go` | SQL 查询实现 |
| 3.4 | 扩展 Service | `backend/internal/service/service.go` | 链接 CRUD 业务逻辑 |
| 3.5 | 链接 Handler | `backend/internal/api/url_handler.go` | GetMyURLs, UpdateURL, DeleteURL |
| 3.6 | 权限校验 | `backend/internal/service/service.go` | 确保用户只能操作自己的链接 |
| 3.7 | 分页支持 | `backend/internal/store/db.go` | LIMIT/OFFSET 分页 |

#### 链接列表接口

```go
// GET /api/urls?page=1&page_size=20
func (h *Handler) GetMyURLs(c *gin.Context) {
    userID := c.GetInt("user_id")
    page := c.DefaultQuery("page", "1")
    pageSize := c.DefaultQuery("page_size", "20")

    urls, total, err := h.service.GetURLsByUserID(userID, page, pageSize)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get urls"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "urls": urls,
        "total": total,
        "page": page,
        "page_size": pageSize,
    })
}
```

#### 权限校验示例

```go
func (s *shortenerService) DeleteURL(userID int, shortCode string) error {
    url, err := s.store.GetByShortCode(shortCode)
    if err != nil {
        return err
    }

    if url.UserID != userID {
        return ErrUnauthorized
    }

    return s.store.Delete(shortCode)
}
```

### Phase 4: 安全与防滥用 (P2)

| 序号 | 任务 | 涉及文件 | 学习要点 |
|------|------|----------|----------|
| 4.1 | 用户配额限制 | `backend/internal/models/quota.go` | 每用户创建链接数限制 |
| 4.2 | 全局速率限制 | `backend/internal/middleware/ratelimit.go` | 滑动窗口/令牌桶算法 |
| 4.3 | 用户级速率限制 | `backend/internal/middleware/ratelimit.go` | 认证用户单独限流 |
| 4.4 | 输入验证增强 | `backend/internal/validator/` | 邮箱格式、URL 安全检查 |
| 4.5 | CORS 配置 | `backend/cmd/short-url/main.go` | gin-cors 中间件 |
| 4.6 | 请求日志 | `backend/internal/middleware/logger.go` | 记录用户操作审计日志 |

### Phase 5: 高级功能 (P3)

| 序号 | 任务 | 涉及文件 | 学习要点 |
|------|------|----------|----------|
| 5.1 | 自定义短码 | `backend/internal/service/service.go` | 允许用户指定短码（唯一性检查） |
| 5.2 | 链接过期 | `backend/internal/store/db.go` | expired_at 字段 + 定时清理 |
| 5.3 | 批量创建 | `backend/internal/api/url_handler.go` | 批量导入链接 |
| 5.4 | 导出功能 | `backend/internal/api/url_handler.go` | 导出用户所有链接 CSV/JSON |
| 5.5 | 访问分析增强 | `backend/internal/models/visit.go` | 详细访问记录（IP、UA、Referer） |
| 5.6 | 统计报表 | `backend/internal/api/stats_handler.go` | 按时间段统计访问趋势 |

---

## API 设计

### 认证接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | 无 |
| POST | `/api/auth/login` | 用户登录 | 无 |
| GET | `/api/auth/me` | 获取当前用户 | 需要 |

### 链接接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/shorten` | 创建短链接 | 需要 |
| GET | `/api/urls` | 获取用户链接列表 | 需要 |
| GET | `/api/urls/:id` | 获取单个链接详情 | 需要 |
| PUT | `/api/urls/:id` | 更新链接 | 需要 |
| DELETE | `/api/urls/:id` | 删除链接 | 需要 |
| GET | `/:shortCode` | 短链接跳转 | 无 |
| GET | `/:shortCode/stats` | 访问统计 | 无（或需要） |

---

## 目录结构预览

```
backend/
├── cmd/short-url/main.go
├── internal/
│   ├── api/
│   │   ├── handler.go          # 已有
│   │   ├── auth_handler.go     # P2.7
│   │   └── url_handler.go      # P3.5
│   ├── auth/
│   │   ├── password.go         # P2.1
│   │   └── jwt.go              # P2.2
│   ├── middleware/
│   │   ├── auth.go             # P2.8
│   │   ├── ratelimit.go        # P4.2
│   │   └── logger.go           # P4.6
│   ├── models/
│   │   ├── user.go             # P2.3
│   │   └── url.go              # P3.1
│   ├── service/
│   │   ├── service.go          # 已有，需改造
│   │   └── auth.go             # P2.6
│   ├── store/
│   │   ├── store.go            # 已有，需扩展
│   │   ├── db.go               # P1.5
│   │   ├── user_store.go       # P2.5
│   │   ├── schema.sql          # P1.3-1.4
│   │   └── migrate.go          # P1.7
│   └── validator/
│       └── validator.go        # P4.4
└── go.mod
```

---

## 数据库选型建议

| 数据库 | 优点 | 缺点 | 推荐场景 |
|--------|------|------|----------|
| PostgreSQL | 功能强大、JSON 支持、扩展性好 | 配置稍复杂 | 生产环境首选 |
| MySQL | 流行度高、资料丰富、简单易用 | 功能相对较少 | 学习项目可选 |
| SQLite | 无需安装、零配置、单文件 | 并发性能差 | 本地开发/测试 |

**推荐**: PostgreSQL（学习企业级实践）或 MySQL（学习基础 SQL）

---

## 环境变量配置

```bash
# .env 或启动参数
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=shorturl
DB_USER=postgres
DB_PASSWORD=postgres

JWT_SECRET=your-super-secret-key
JWT_EXPIRE=24h

# 可选
QUOTA_MAX_URLS=100
RATE_LIMIT_RPM=60
```

---

## 学习要点总结

### Phase 1-2 核心学习点

1. **数据库操作**: `database/sql` + 驱动，预处理语句，事务
2. **密码安全**: bcrypt 哈希，防止明文存储
3. **JWT 认证**: Token 生成、验证、刷新机制
4. **中间件模式**: Gin 中间件链，context 数据传递
5. **接口设计**: RESTful 规范，错误响应统一格式

### Phase 3-4 核心学习点

1. **权限控制**: 资源归属校验，避免越权访问
2. **SQL 注入防护**: 预处理语句，参数化查询
3. **分页查询**: LIMIT/OFFSET 性能优化
4. **限流算法**: 滑动窗口/令牌桶实现

### Phase 5 核心学习点

1. **定时任务**: 后台清理过期数据
2. **批量操作**: 事务处理批量插入
3. **数据导出**: CSV/JSON 格式化输出

---

## 测试策略

| 层级 | 测试内容 | 工具 |
|------|----------|------|
| 单元测试 | 密码哈希、JWT 生成解析 | `testing` |
| 集成测试 | Store 层 CRUD | `testing` + test DB |
| API 测试 | Handler 层 HTTP 接口 | `httptest` |
| E2E 测试 | 完整用户流程 | 手动/curl |

---

## 实施建议

1. **先完成后端 Phase 1-2**，前端才能开始认证开发
2. **每个 Phase 完成后写测试**，确保稳定性
3. **使用分支开发**，每个 Phase 一个分支
4. **保持向后兼容**，旧接口（无认证）可暂时保留

---

## 命令速查

```bash
# 安装依赖
go get golang.org/x/crypto/bcrypt
go get github.com/golang-jwt/jwt/v5
go get github.com/lib/pq  # PostgreSQL

# 运行服务
make run-memory  # 内存模式（开发）
make run-db      # 数据库模式（新增）

# 运行测试
make test
go test ./internal/store/... -v
go test ./internal/auth/... -v

# 数据库迁移
psql -d shorturl -f internal/store/schema.sql
```

---

## 配合前端计划

| 后端 Phase | 前端依赖 |
|------------|----------|
| Phase 1 (数据库) | 无直接依赖 |
| Phase 2 (认证) | 前端 Phase 1 (认证系统) |
| Phase 3 (链接管理) | 前端 Phase 2 (链接管理) |
| Phase 4 (安全) | 无直接依赖 |
| Phase 5 (高级) | 前端 Phase 4 (高级功能) |

**建议顺序**: 后端先行一个 Phase，前端跟随实现。
