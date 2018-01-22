package hash

import (
	"crypto/md5"
	"fmt"
)

// Sign is a method to handle the signing of the request
func Sign(mode int, hash, size, t, ua, ip string) string {
	var res string
	if mode == 0 {
		res = hash + size + t + ua + ip
	} else {
		res = hash + size + t + ip
	}
	m := md5.New()
	_, _ = m.Write([]byte(res))
	return fmt.Sprintf("%d%x", mode, m.Sum(nil))
}
