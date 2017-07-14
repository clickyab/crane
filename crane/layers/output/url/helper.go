package url

import (
	"fmt"
	"net/url"

	"github.com/clickyab/services/assert"
)

func keyGen(m, s string) string {
	return fmt.Sprintf("%s_%s_%s", prefix, m, s)
}

func addMeta(s string, m map[string]interface{}) string {
	u, e := url.Parse(s)
	assert.Nil(e)
	q := u.Query()
	for k, v := range m {
		q.Set(k, fmt.Sprint(v))
	}
	u.RawQuery = q.Encode()
	return u.String()
}
