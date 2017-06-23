package routes

import (
	"context"
	"net/http"

	"github.com/clickyab/services/assert"

	"clickyab.com/exchange/octopus/console/user/aaa"
	"github.com/clickyab/services/eav"
	"github.com/clickyab/services/framework"
	"github.com/clickyab/services/framework/controller"
	"github.com/clickyab/services/trans"
)

type dataType string

const dataKey dataType = "__user_data__"
const tokenKey dataType = "__token_data__"

// Authenticate is a middleware used for authorization in exchange console
func Authenticate(next framework.Handler) framework.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("token")
		if token != "" {
			val := eav.NewEavStore(token).SubKey("token")
			if val != "" {
				usr, err := aaa.NewAaaManager().FindUserByToken(val)
				if err == nil {
					ctx = context.WithValue(ctx, dataKey, usr)
					ctx = context.WithValue(ctx, tokenKey, token)
					next(ctx, w, r)
					return
				}
			}
		}
		framework.JSON(w, http.StatusUnauthorized, controller.ErrorResponseSimple{Error: trans.E("Unauthorized")})
	}
}

// GetUser is the helper function to extract user data from context
func GetUser(ctx context.Context) (*aaa.User, bool) {
	rd, ok := ctx.Value(dataKey).(*aaa.User)
	if !ok {
		return nil, false
	}

	return rd, true
}

// MustGetUser try to get user data, or panic if there is no user data
func MustGetUser(ctx context.Context) *aaa.User {
	rd, ok := GetUser(ctx)
	assert.True(ok, "[BUG] no user in context")
	return rd
}

// GetToken is the helper function to extract user data from context
func GetToken(ctx context.Context) (string, bool) {
	rd, ok := ctx.Value(tokenKey).(string)
	if !ok {
		return "", false
	}

	return rd, true
}

// MustGetToken try to get user data, or panic if there is no user data
func MustGetToken(ctx context.Context) string {
	rd, ok := GetToken(ctx)
	assert.True(ok, "[BUG] no user in context")
	return rd
}
