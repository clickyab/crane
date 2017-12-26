package app

import (
	"crypto/md5"
	"fmt"
)

// Hash is the hash generation func for keys, md5 normally
func Hash(k string) string {
	h := md5.New()
	_, _ = h.Write([]byte(k))
	return fmt.Sprintf("%x", h.Sum(nil))
}
