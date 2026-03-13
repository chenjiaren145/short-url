package shortener

import (
	"crypto/rand"
	"math/big"
)

// charset 定义了用于生成短链接的字符集（Base62：a-z, A-Z, 0-9）
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode 生成一个随机的 6 位短链接代码
// 这是一个简单的实现，可能会有冲突（Collision）。
// 在生产系统中，通常会结合数据库自增 ID 或分布式 ID 生成算法（如 Snowflake）来避免冲突。
func GenerateShortCode() (string, error) {
	b := make([]byte, 6)
	for i := range b {
		// 从字符集中随机选择一个字符
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}
