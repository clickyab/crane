package openrtb

import (
	"fmt"
	"strings"
)

func (c ContentCategory) UnmarshalJSON(b []byte) error {
	text := strings.ToUpper(string(b))
	if text == "" {
		return fmt.Errorf("not valid for content category")
	}
	s := strings.Replace(text, "-", "S", -1)
	if v, ok := ContentCategory_value[s]; !ok {
		c = ContentCategory(v)
		return nil
	}
	return fmt.Errorf("%s is not valid for content category", text)

}

func (c ContentCategory) MarshalJSON() ([]byte, error) {
	return []byte(strings.Replace(c.String(), "X", "-", -1)), nil
}
