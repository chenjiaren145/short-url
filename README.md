# Short URL Service

前后端分离单仓项目：
- `backend/`: Go + Gin API 服务
- `frontend/`: React + Vite 前端
- 根目录：`docker-compose.yml`、`docker-compose.debug.yml`、`Makefile`

## 目录结构

```text
.
├── backend
│   ├── cmd/short-url/main.go
│   └── internal/{api,service,shortener,store}
├── frontend
│   └── src/{components,services,types}
├── docker-compose.yml
├── docker-compose.debug.yml
└── Makefile
```

## 运行方式

### 1) 后端本地内存模式

```bash
make run-memory
```

### 2) 后端本地 Redis 模式

```bash
make run-redis
```

### 3) Docker 运行（默认不暴露 Redis）

```bash
make compose-up
```

### 4) Docker 运行（调试模式，暴露 Redis 6379）

```bash
make compose-up-debug-redis
```

### 5) 前端本地开发

```bash
cd frontend
npm ci
npm run dev
```

## 健康检查

- `GET /healthz`
- `GET /readyz`

## API 快速验证

```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"original_url":"https://www.google.com"}' \
  http://localhost:8081/shorten
```

## 常用命令

```bash
make test
make build
make compose-down
```
