package helper

import (
	"crypto/sha1"
	"fmt"
)

func Hash(salt string, password string) string {
	s := fmt.Sprintf("_%s+%s_", salt, password)
	hash := sha1.New()
	hash.Write([]byte(s))
	s = fmt.Sprintf("%x", hash.Sum(nil))
	return s
}
