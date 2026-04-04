package api

import (
	"fmt"
	"net/http"
	"short-url/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler 结构体持有 HTTP 处理函数所需的依赖
// 现在它依赖 service.ShortenerService，而不是直接依赖 store.Store
type Handler struct {
	service service.ShortenerService
	baseURL string
}

// NewHandler 创建一个新的 Handler 实例
func NewHandler(svc service.ShortenerService) *Handler {
	return &Handler{service: svc}
}

func NewHandlerWithBaseURL(svc service.ShortenerService, baseURL string) *Handler {
	return &Handler{
		service: svc,
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (h *Handler) shortURL(c *gin.Context, shortCode string) string {
	if h.baseURL != "" {
		return h.baseURL + "/" + shortCode
	}
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	if forwardedProto := c.GetHeader("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = strings.TrimSpace(strings.Split(forwardedProto, ",")[0])
	}
	return fmt.Sprintf("%s://%s/%s", scheme, c.Request.Host, shortCode)
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
	// 可以区分 400/500 code
	shortCode, err := h.service.Shorten(req.OriginalURL)
	if err != nil {

		// 简单的做法：所有非 nil 错误都返回 400，因为验证错误是最常见的
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url":    h.shortURL(c, shortCode),
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

// GetStats 获取访问次数
func (h *Handler) GetStats(c *gin.Context) {
	shortCode := c.Param("shortCode")

	// 调用 Service 层获取访问次数
	visits, err := h.service.GetStats(shortCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取数据失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url": h.shortURL(c, shortCode),
		"visits":    visits,
	})
}
