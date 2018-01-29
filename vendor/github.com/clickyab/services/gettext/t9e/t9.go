package t9e

import (
	"fmt"

	"github.com/clickyab/services/gettext/internal/t9"
)

// T9 is translation error
type T9 struct {
	t9.Base
}

func (t9 T9) Error() string {
	return fmt.Sprintf(t9.Text, t9.Params...)
}

// G for basic translation (roughly equivalent to gettext() or _())
func G(translationID string, args ...interface{}) T9 {
	return T9{
		t9.Base{

			Text:   translationID,
			Params: args,
		},
	}
}

// PG for translation with context (pgettext())
func PG(s string, p ...interface{}) T9 {
	return G(s, p...)
}

// NG for translation with quantities (ngettext())
func NG(s string, p ...interface{}) T9 {
	return G(s, p...)
}

// NPG for translation with both quantities and context (npgettext())
func NPG(s string, p ...interface{}) T9 {
	return G(s, p...)
}
