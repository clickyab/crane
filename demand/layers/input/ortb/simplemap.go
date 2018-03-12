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

func (s simpleMap) Float64(k string) float64 {
	d, ok := s[k]
	if !ok {
		return 0
	}
	switch t := d.(type) {
	case float64:
		return t
	case string:
		f, err := strconv.ParseFloat(t, 64)
		if err != nil {
			return 0
		}
		return f
	case []byte:
		x := string(t)
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0
		}
		return f
	default:
		return 0
	}
}
