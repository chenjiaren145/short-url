# Short URL Service

前后端分离短链接服务：Go + Gin API，React + Vite 前端

## 快速开始

```bash
# 安装前端依赖
make frontend-install

# 终端 1: 启动后端
make run-memory

# 终端 2: 启动前端
make frontend-dev
```

访问 http://localhost:5173

## 常用命令

```bash
make run-memory      # 后端（内存模式）
make run-redis       # 后端（Redis 模式）
make frontend-dev    # 前端开发服务器
make test            # 后端测试
make compose-up      # Docker 启动完整服务
make compose-down    # Docker 停止服务
```

## 目录结构

```
.
├── backend/          # Go + Gin API
│   ├── cmd/short-url/main.go
│   └── internal/{api,service,shortener,store}
├── frontend/         # React + Vite
│   └── src/{components,services,types}
└── docs/             # 文档
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/shorten` | 创建短链接 |
| GET | `/:code` | 重定向 |
| GET | `/:code/stats` | 访问统计 |
| GET | `/healthz` | 健康检查 |

```bash
curl -X POST http://localhost:8081/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://www.google.com"}'
```
