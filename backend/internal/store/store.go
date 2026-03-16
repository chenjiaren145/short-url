package store

import "errors"

// Store 定义了存储操作的接口 (Interface)
// 使用接口可以将业务逻辑与具体的存储实现解耦。
// 这样我们可以在不修改上层逻辑的情况下，轻松切换存储后端（例如从内存切换到 MySQL 或 Redis）。
type Store interface {
	// Save 保存短链接代码和原始 URL 的映射关系
	Save(shortCode string, originalURL string) error

	// Load 根据短链接代码获取原始 URL
	Load(shortCode string) (string, error)

	// IncrementVisits 增加链接的访问次数
	IncrementVisits(shortCode string) error

	// GetVisits 用于获取访问记数
	GetVisits(shortCode string) (int64, error)
}

// ErrNotFound 定义了一个统一的错误，当短链接不存在时返回
// 这样上层调用者不需要关心具体的存储实现返回什么错误，只需要判断是否是这个错误即可。
var ErrNotFound = errors.New("short code not found")
