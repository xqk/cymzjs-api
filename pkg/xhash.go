package pkg

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Sum(s string) string {
	m := md5.New()
	_, err := m.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(m.Sum(nil))
}
