package ortb

import (
	"strconv"
)

type simpleMap map[string]interface{}

func (s simpleMap) Bool(k string) (bool, bool) {
	d, ok := s[k]
	if !ok {
		return false, false
	}
	switch t := d.(type) {
	case float64:
		return t != 0, true
	case string:
		b, _ := strconv.ParseBool(t)
		return b, true
	case bool:
		return t, true
	default:
		return false, true
	}
}

func (s simpleMap) String(k string) (string, bool) {
	d, ok := s[k]
	if !ok {
		return "", false
	}
	switch t := d.(type) {
	case string:
		return t, true
	case []byte:
		return string(t), true
	default:
		return "", true
	}
}
