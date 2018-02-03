package hash

import (
	"fmt"

	"github.com/clickyab/services/simplehash"
)

// Sign is a method to handle the signing of the request
func Sign(mode int, hash, size, t, ua, ip string) string {
	var res string
	if mode == 0 {
		res = hash + size + t + ua + ip
	} else {
		res = hash + size + t + ip
	}
	return simplehash.MD5(fmt.Sprintf("%d%s", mode, res))
}
