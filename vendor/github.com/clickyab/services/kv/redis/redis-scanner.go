package redis

import (
	"sync"

	"github.com/clickyab/services/aredis"
	"github.com/clickyab/services/kv"
)

type scanner struct {
	pattern    string
	cursor     uint64
	terminated bool

	lock sync.Mutex
}

func (s *scanner) Next(max int) ([]string, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.terminated {
		return nil, false
	}
	cmd := aredis.Client.Scan(s.cursor, s.pattern, int64(max))
	var keys []string
	keys, s.cursor = cmd.Val()
	s.terminated = s.cursor == 0
	return keys, s.cursor != 0
}

func (s *scanner) Pattern() string {
	return s.pattern
}

func newRedisScanner(pattern string) kv.Scanner {
	return &scanner{pattern: pattern}
}
