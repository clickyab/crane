package simplehash

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash/crc32"
	"hash/crc64"
)

// MD5 is the hash generation func for keys, md5 normally
func MD5(k string) string {
	h := md5.New()
	_, _ = h.Write([]byte(k))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// SHA1 hash string in sha1 method
func SHA1(k string) string {
	h := sha1.New()
	h.Write([]byte(k))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// CRC32 hash string in crc32 method
func CRC32(k string) string {
	x := crc32.New(crc32.IEEETable)
	_, _ = x.Write([]byte(k))
	return fmt.Sprint(x.Sum32())
}

// CRC64 hash string in crc64 method
func CRC64(k string) int64 {
	crc := crc64.New(crc64.MakeTable(crc64.ECMA))
	crc.Write([]byte(k))
	return int64(crc.Sum64())
}
