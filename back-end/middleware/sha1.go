package middleware

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
)

const (
	MinPasswordSize = 6
	MaxPasswordSize = 64
)

// SHA1 SHA1加密
func SHA1(s string) string {
	// 创建一个SHA1哈希对象
	o := sha1.New()
	// 将输入的字符串转化成字节数组，写入o
	o.Write([]byte(s))
	// 计算并返回SHA1哈希值
	return hex.EncodeToString(o.Sum(nil))
}

// SHAMiddleWare SHA1加密用户密码中间件
func SHAMiddleWare(password string) (string, error) {
	if len(password) < MinPasswordSize || len(password) > MaxPasswordSize {
		return "", errors.New("密码长度小于或大于限制")
	}
	return SHA1(password), nil
}
