package pkg

import (
	"math/rand"
	"time"
)

var alphabet = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")

func RandString(n int64) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	length := len(alphabet)
	for i := range b {
		b[i] = alphabet[rand.Intn(length)]
	}
	return string(b)
}
