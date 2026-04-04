package validator

import (
	"errors"
	"net/url"
)

// ValidateURL 检查 URL 是否合法
// - 必须以 http:// 或 https:// 开头
// - 必须是有效的 URL（能被 url.Parse 解析）
func ValidateURL(rawUrl string) error {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return errors.New("invalid URL fromat")
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return errors.New("invalid URL format")
	}
	if u.Host == "" {
		return errors.New("URL must have a host")
	}
	return nil
}
