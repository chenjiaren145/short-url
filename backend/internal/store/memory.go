package store

import (
	"sync"
)

// MemoryStore 实现了 Store 接口，使用内存中的 map 来存储数据
// 这是一个简单的实现，适合开发和测试，但在生产环境中，重启服务会导致数据丢失。
type MemoryStore struct {
	// mu 是一个读写互斥锁 (RWMutex)
	// 在 Go 中，内置的 map 不是并发安全的。如果多个 goroutine 同时读写 map，会导致 panic。
	// 因此我们需要使用锁来保护 map 的并发访问。
	mu sync.RWMutex

	// data 用于存储短链接到原始 URL 的映射
	data map[string]string

	// visits 用于存储短链接的访问计数
	visits map[string]int64
}

// NewMemoryStore 创建一个新的 MemoryStore 实例
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data:   make(map[string]string),
		visits: make(map[string]int64),
	}
}

// Save 实现 Store 接口的 Save 方法
func (s *MemoryStore) Save(shortCode string, originalURL string) error {
	// 写操作需要加写锁 (Lock)，确保同一时间只有一个 goroutine 能修改 map
	s.mu.Lock()
	// 使用 defer 确保函数退出时自动释放锁，避免死锁
	defer s.mu.Unlock()

	s.data[shortCode] = originalURL
	return nil
}

// Load 实现 Store 接口的 Load 方法
func (s *MemoryStore) Load(shortCode string) (string, error) {
	// 读操作只需要加读锁 (RLock)
	// 读锁允许并发读取，但不允许写入。这样可以提高读多写少场景下的性能。
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, ok := s.data[shortCode]
	if !ok {
		return "", ErrNotFound
	}
	return url, nil
}

// IncrementVisits 实现 Store 接口的 IncrementVisits 方法
func (s *MemoryStore) IncrementVisits(shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.visits[shortCode]++
	return nil
}

// GetVisits 实现 Store 接口的 GetVisits 方法
func (s *MemoryStore) GetVisits(shortCode string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.visits[shortCode], nil
}

// Delete 实现 Store 接口的 Delete 方法
func (s *MemoryStore) Delete(shortCode string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.data[shortCode]; !ok {
		return ErrNotFound
	}

	delete(s.data, shortCode)
	delete(s.visits, shortCode)
	return nil
}
