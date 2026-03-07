.PHONY: run-memory run-redis test build clean

# 默认使用内存存储启动（数据重启丢失）
run-memory:
	go run cmd/short-url/main.go -store memory

# 使用 Redis 存储启动（需先启动 Docker 和 Redis）
run-redis:
	@echo "Ensuring Redis is running..."
	docker-compose up -d redis
	go run cmd/short-url/main.go -store redis -redis-addr localhost:6379

# 运行测试
test:
	go test ./...

# 编译二进制文件
build:
	go build -o short-url cmd/short-url/main.go

# 清理构建文件
clean:
	rm -f short-url