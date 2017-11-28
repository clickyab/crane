package cytid

import (
	"fmt"
	"strconv"
	"time"

	"clickyab.com/gad/utils"
	"github.com/clickyab/services/assert"
	"github.com/clickyab/services/config"
	"github.com/clickyab/services/kv"
)

var (
	expire = config.RegisterDuration("clickyab.cookie_profile.expire", 3*24*time.Hour, "expirey of the cookie profile")
)

func randInt64() int64 {
	x := <-utils.ID
	i, err := strconv.ParseInt(x[:8], 16, 64)
	assert.Nil(err)
	return i
}

// GetCookieProfileID return the cookie profile id for current user
// this is simplest implementation. we need to add many other data for
// business logic
func GetCookieProfileID(cop string) string {
	ot := kv.NewOneTimeSetter(cop, expire.Duration())
	return ot.Set(fmt.Sprint(randInt64()))
}
