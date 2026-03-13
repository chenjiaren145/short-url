package main

import (
	"flag"
	"log"
	"short-url/internal/api"
	"short-url/internal/service"
	"short-url/internal/store"

	"github.com/gin-gonic/gin"
)

func main() {
	// 定义命令行参数
	storeType := flag.String("store", "memory", "Storage type: memory or redis")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis address")
	redisPassword := flag.String("redis-password", "", "Redis password")
	flag.Parse()

	// 创建一个默认的 Gin 路由器
	r := gin.Default()

	// 1. 初始化存储层 (Store)
	var s store.Store
	switch *storeType {
	case "redis":
		log.Printf("Using Redis storage at %s", *redisAddr)
		s = store.NewRedisStore(*redisAddr, *redisPassword, 0)
	case "memory":
		log.Println("Using In-Memory storage")
		s = store.NewMemoryStore()
	default:
		log.Fatalf("Unknown store type: %s", *storeType)
	}

	// 2. 初始化业务逻辑层 (Service)
	svc := service.NewShortenerService(s)

	// 3. 初始化 API 处理层 (Handler)
	h := api.NewHandler(svc)

	// 定义路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Short URL Service API",
		})
	})

	// POST /shorten -> 创建短链接
	r.POST("/shorten", h.CreateShortURL)

	// GET /:shortCode -> 跳转到原始链接
	r.GET("/:shortCode", h.Redirect)
	r.HEAD("/:shortCode", h.Redirect)

	// TODO GET /getStats -> 获取短链接状态（访问次数）
	r.GET("/:shortCode/stats", h.GetStats)

	// 启动 HTTP 服务，监听 8081 端口
	log.Println("Server starting on :8081")
	err := r.Run(":8081")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
