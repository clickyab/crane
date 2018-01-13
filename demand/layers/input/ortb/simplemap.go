package ortb

import "strconv"

type simpleMap map[string]interface{}

func (s simpleMap) Bool(k string) bool {
	d, ok := s[k]
	if !ok {
		return false
	}
	switch t := d.(type) {
	case float64:
		return t != 0
	case string:
		b, _ := strconv.ParseBool(t)
		return b
	case bool:
		return t
	default:
		return false
	}
}

func (s simpleMap) String(k string) string {
	d, ok := s[k]
	if !ok {
		return ""
	}
	switch t := d.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return ""
	}
}
