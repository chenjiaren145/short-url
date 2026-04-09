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
	storeType := flag.String("store", "memory", "Storage type: memory or redis")
	redisAddr := flag.String("redis-addr", "localhost:6379", "Redis address")
	redisPassword := flag.String("redis-password", "", "Redis password")
	port := flag.String("port", "8081", "HTTP server port")
	baseURL := flag.String("base-url", "", "Public base URL for short links, e.g. https://sho.rt")
	flag.Parse()

	r := gin.Default()

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

	svc := service.NewShortenerService(s)

	h := api.NewHandlerWithBaseURL(svc, *baseURL)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the Short URL Service API",
		})
	})
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.GET("/readyz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ready"})
	})

	r.POST("/shorten", h.CreateShortURL)

	r.GET("/:shortCode", h.Redirect)
	r.HEAD("/:shortCode", h.Redirect)

	r.GET("/:shortCode/stats", h.GetStats)

	r.DELETE("/:shortCode", h.Delete)

	addr := ":" + *port
	log.Printf("Server starting on %s", addr)
	err := r.Run(addr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
