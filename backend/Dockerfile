# 第一阶段：构建
FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# 复制依赖文件并下载
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 编译应用
RUN CGO_ENABLED=0 GOOS=linux go build -o short-url ./cmd/short-url/main.go

# 第二阶段：运行
FROM alpine:latest

WORKDIR /app

# 从构建阶段复制二进制文件
COPY --from=builder /app/short-url .

# 暴露端口
EXPOSE 8081

# 启动命令
CMD ["./short-url"]
