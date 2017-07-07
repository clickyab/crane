package url

import (
	"fmt"
)

func keyGen(m, s string) string {
	return fmt.Sprintf("%s_%s_%s", prefix, m, s)
}
