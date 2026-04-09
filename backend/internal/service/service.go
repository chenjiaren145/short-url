package service

import (
	"errors"
	"short-url/internal/shortener"
	"short-url/internal/store"
	"short-url/internal/validator"
)

// ShortenerService 定义了短链接服务的业务逻辑接口
type ShortenerService interface {
	Shorten(originalURL string) (string, error)
	GetOriginalURL(shortCode string) (string, error)
	GetStats(shortCode string) (int64, error)
	Delete(shortCode string) error
}

// shortenerService 是 ShortenerService 的具体实现
type shortenerService struct {
	store store.Store
}

// NewShortenerService 创建一个新的 ShortenerService 实例
func NewShortenerService(store store.Store) ShortenerService {
	return &shortenerService{
		store: store,
	}
}

// Shorten 生成短链接并保存
func (s *shortenerService) Shorten(originalURL string) (string, error) {

	// 检查URL是否合法
	if err := validator.ValidateURL(originalURL); err != nil {
		return "", err
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

// GetOriginalURL 获取原始链接
func (s *shortenerService) GetOriginalURL(shortCode string) (string, error) {
	s.store.IncrementVisits(shortCode)
	return s.store.Load(shortCode)
}

// GetStats 获取访问次数
func (s *shortenerService) GetStats(shortCode string) (int64, error) {
	return s.store.GetVisits(shortCode)
}

// Delete 删除短链接
func (s *shortenerService) Delete(shrotCode string) error {
	if shrotCode == "" {
		return errors.New("请传递 shrotCode")
	}

	return s.store.Delete(shrotCode)
}
