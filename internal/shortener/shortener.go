package shortener

import (
	"math/rand"
	"time"
)

// charset 定义了用于生成短链接的字符集（Base62：a-z, A-Z, 0-9）
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// seededRand 是一个带种子的随机数生成器
// 注意：在 Go 1.20+ 中，math/rand 的全局函数已经自动播种，但为了演示 clear code，这里显式创建了一个局部生成器。
// 在高并发场景下，共享的 rand.Rand 可能存在锁竞争问题，生产环境通常不需要这样手动加锁，或者使用 crypto/rand。
var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateShortCode 生成一个随机的 6 位短链接代码
// 这是一个简单的实现，可能会有冲突（Collision）。
// 在生产系统中，通常会结合数据库自增 ID 或分布式 ID 生成算法（如 Snowflake）来避免冲突。
func GenerateShortCode() string {
	b := make([]byte, 6)
	for i := range b {
		// 从字符集中随机选择一个字符
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
