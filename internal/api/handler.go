package api

import (
	"net/http"
	"short-url/internal/service"

	"github.com/gin-gonic/gin"
)

// Handler 结构体持有 HTTP 处理函数所需的依赖
// 现在它依赖 service.ShortenerService，而不是直接依赖 store.Store
type Handler struct {
	service service.ShortenerService
}

// NewHandler 创建一个新的 Handler 实例
func NewHandler(svc service.ShortenerService) *Handler {
	return &Handler{service: svc}
}

// CreateShortURLRequest 定义了创建短链接请求的 JSON 结构
type CreateShortURLRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
}

// CreateShortURL 处理创建短链接的请求
func (h *Handler) CreateShortURL(c *gin.Context) {
	var req CreateShortURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 调用 Service 层处理业务逻辑
	shortCode, err := h.service.Shorten(req.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url":    "http://localhost:8081/" + shortCode,
		"original_url": req.OriginalURL,
	})
}

// Redirect 处理短链接跳转
func (h *Handler) Redirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// 调用 Service 层获取原始链接
	originalURL, err := h.service.GetOriginalURL(shortCode)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
