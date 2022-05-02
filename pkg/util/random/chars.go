package random

import (
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"math/rand"
	"time"
)

// 使用私有的随机生成器
var irand = rand.New(rand.NewSource(time.Now().UnixNano()))

const (
	defaultChars = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func UUID() string {
	return hex.EncodeToString(uuid.NewV4().Bytes())
}

func Chars(n int) string {
	return GenString(defaultChars, n)
}

func GenString(chars string, n int) string {
	charLen := len(chars)
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = chars[irand.Intn(charLen)]
	}
	return string(b)
}
