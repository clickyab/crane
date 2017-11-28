package builder

import (
	"clickyab.com/gad/builder/cytid"
	"github.com/clickyab/services/config"
	"clickyab.com/gad/utils"
	"fmt"
)

var copLen = config.RegisterInt("clickyab.cop_len", 10, "cop key len")

// SetWebTID try to create unique user tid, must set it at the end!
func SetWebTID(tid string) ShowOptionSetter {
	return func(o *Context) (*Context, error) {
		if o.common.UserAgent == "" || o.common.IP == nil {
			return nil, fmt.Errorf("call this at the end after setting ip and ua")
		}
		if len(tid) < copLen.Int() {
			tid = utils.CreateHash(copLen.Int(), []byte(o.common.UserAgent), []byte(o.common.IP))
		}
		o.common.TID = tid
		o.common.CopID = cytid.GetCookieProfileID(o.common.TID)
		return o, nil
	}
}
